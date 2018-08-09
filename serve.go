package sesr

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const (
	listenPort = 8080
)

var (
	w http.ResponseWriter
	r *http.Request
)

func Serve(awsRegion string, awsAccessKeyId string, awsSecretAccessKey string) {
	log.Printf("initialize sesr")
	log.Debug("debug logging enabled")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqHandler(w, r, awsRegion, awsAccessKeyId, awsSecretAccessKey)
	})
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/health/", healthCheckHandler)

	log.Info("listening for requests on :" + fmt.Sprintf("%v", listenPort))
	if err := http.ListenAndServe(":"+fmt.Sprintf("%v", listenPort), nil); err != nil {
		log.Fatalf("http; %v", err)
	}
}
