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

func (err Error) Send(w http.ResponseWriter) {
	if err.Code == 0 {
		err.Code = 500
	}
	logrus.WithError(err.Error).Errorf("responding with error: %v", err.Message)
	w.WriteHeader(err.Code)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(err)
}
