package sesr

import (
	"io"
	"net/http"

	//log "github.com/Sirupsen/logrus"
)

func sendResponse(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, body)
}
