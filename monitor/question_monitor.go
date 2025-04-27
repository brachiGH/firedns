package monitor

import (
	"strings"
	"sync"
	"time"

	"github.com/brachiGH/firedns/internal/utils"
	"github.com/brachiGH/firedns/internal/utils/config"
	"github.com/brachiGH/firedns/internal/utils/logger"
	"github.com/brachiGH/firedns/monitor/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type questionPair struct {
	ip     utils.IP
	lables []utils.Lable
}

var passed []*questionPair
var dorped []*questionPair

var QuestionMonitorWG sync.WaitGroup

func Passed(ip utils.IP, lables []utils.Lable) {
	QuestionMonitorWG.Wait()

	passed = append(passed, &questionPair{ip: ip, lables: lables})
}

func Droped(ip utils.IP, lables []utils.Lable) {
	QuestionMonitorWG.Wait()

	dorped = append(dorped, &questionPair{ip: ip, lables: lables})
}

func createMongoModel(ip utils.IP, lables []utils.Lable, modelField string) mongo.WriteModel {
	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"ip": ip}).
		SetUpdate(bson.M{"$push": bson.M{modelField: bson.A{strings.Join(utils.LablesToStrings(lables), "."), time.Now()}}}).SetUpsert(true)
}

func UpdateQuestions_Routine() {
	log := logger.NewLogger()

	db, err := database.GetAnalyticsDB()
	if err != nil {
		log.Error("Failded to connect to the db: %w", zap.Error(err))
	}

	tick := time.Tick(config.UpdateQuestions__TickDuration)
	for range tick {
		updates := []mongo.WriteModel{}

		for _, p := range passed {
			updates = append(updates, createMongoModel(p.ip, p.lables, "passed"))
		}
		for _, d := range dorped {
			updates = append(updates, createMongoModel(d.ip, d.lables, "dorped"))
		}

		if len(updates) != 0 {
			err := db.UpdateMany(updates)
			if err != nil {
				log.Error("update question routine failed to update database", zap.Error(err))
			} else {
				log.Debug("dns question", zap.Any("updates", updates))
			}
		}

		QuestionMonitorWG.Add(1)
		passed = []*questionPair{}
		dorped = []*questionPair{}
		QuestionMonitorWG.Done()
	}
}
