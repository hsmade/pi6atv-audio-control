package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Error struct {
	Message string
	Error   error
	Code    int
}

func (e Error) Send(w http.ResponseWriter) {
	if e.Code == 0 {
		e.Code = 500
	}
	logrus.WithError(e.Error).Errorf("sending error: %v", e.Message)
	w.WriteHeader(e.Code)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(e)
}
