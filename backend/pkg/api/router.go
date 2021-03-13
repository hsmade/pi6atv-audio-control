package api

import (
	"github.com/gorilla/mux"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/config"
	"github.com/pkg/errors"
)

// NewRouter initialises the control and sensors objects and creates the http router to reach them
func NewRouter(config config.Config) (*mux.Router, error) {
	control, err := NewControl()
	if err != nil {
		return nil, errors.Wrap(err, "initialising control")
	}

	router := mux.NewRouter()
	router.HandleFunc("/control/", control.ControlGetAll)
	router.HandleFunc("/control/check", control.ControlCheck)
	router.HandleFunc("/control/{relay}", control.ControlGetRelay).Methods("GET")
	router.HandleFunc("/control/{relay}/{state}", control.ControlSetRelay).Methods("POST")

	return router, nil
}
