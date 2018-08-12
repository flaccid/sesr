package sesr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type SendPayload struct {
	Sender     string   `json:"sender"`
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
}

func reqHandler(w http.ResponseWriter, r *http.Request, awsRegion string, awsAccessKeyId string, awsSecretAccessKey string) {
	logRequest(r)

	// we only want to process real requests i.e. reject robots, favicon etc.
	if r.URL.Path != "/" {
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
		return
	}

	var jsonResponse string

	switch r.Method {
	case "GET":
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	case "POST":
		var bodyBytes []byte
		var payload SendPayload

		// read the body content and unmarshal the expected json
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jsonResponse = `{"error": "internal server error"}`
			sendResponse(w, http.StatusInternalServerError, jsonResponse)
			log.Error("error reading request body: ", err)
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		err = json.Unmarshal(bodyBytes, &payload)
		if err != nil {
			log.Error("error unmarshalling json: ", err)
		}

		if len(payload.Sender) < 1 {
			sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters: no sender provided"}`)
		} else if len(payload.Recipients) < 1 {
			sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters: no recipients provided"}`)
		} else if len(payload.Subject) < 1 {
			sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters: no subject provided"}`)
		} else if len(payload.Body) < 1 {
			sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters: no body provided"}`)
		} else {
			err := Send(awsRegion, awsAccessKeyId, awsSecretAccessKey, payload.Sender, payload.Recipients, payload.Subject, payload.Body, "UTF-8")
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("failure sending email")
			} else {
				log.WithFields(log.Fields{
					"sender":     payload.Sender,
					"recipients": payload.Recipients,
					"subject":    payload.Subject,
					"charset":    "UTF-8",
					"body":       "[redacted]",
				}).Info("email(s) sent")
			}
		}
	default:
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	}
}
