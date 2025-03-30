package transport

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/brachiGH/firedns/internal/server"
	"github.com/brachiGH/firedns/internal/utils/logger"
	"go.uber.org/zap"
)

func StartTLSServer(port string) {
	log := logger.NewLogger()

	pubkey := os.Getenv("CertFile")
	prvkey := os.Getenv("KeyFile")
	cert, err := tls.LoadX509KeyPair(pubkey, prvkey)
	if err != nil {
		log.Fatal("Error loading certificates:", zap.Error(err))
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the ID from the URL path (assuming the ID is provided after the leading '/')
		id := r.URL.Path[1:]

		log.Debug("Received request for ID", zap.String("id", id))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Error("Invalid remote address", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		clientIP := net.ParseIP(host)
		data, err := server.HandleDnsMessage(body, clientIP, id)
		if err != nil {
			log.Error("Fail to handle request: ", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/dns-message")
		_, err = w.Write(data)
		if err != nil {
			log.Error("Fail to write request: ", zap.Error(err))
			return
		}
	})

	server := &http.Server{
		Addr:      ":" + port,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal("Server failed:", zap.Error(err))
	}
}
