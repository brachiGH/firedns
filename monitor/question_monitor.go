package monitor

import (
	"strings"
	"sync"
	"time"

	"github.com/brachiGH/firedns/internal/utils"
	"github.com/brachiGH/firedns/monitor/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	var db database.Analytics_DB
	db.Connect()
	defer db.Disconnect()
	tick := time.Tick(30 * time.Second)
	for range tick {
		updates := []mongo.WriteModel{}

		for _, p := range passed {
			updates = append(updates, createMongoModel(p.ip, p.lables, "passed"))
		}
		for _, d := range dorped {
			updates = append(updates, createMongoModel(d.ip, d.lables, "dorped"))
		}

		if len(updates) != 0 {
			db.UpdateMany(updates)
		}

		QuestionMonitorWG.Add(1)
		passed = []*questionPair{}
		dorped = []*questionPair{}
		QuestionMonitorWG.Done()
	}
}
