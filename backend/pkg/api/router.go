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

	sensors, err := NewSensors()
	if err != nil {
		return nil, errors.Wrap(err, "initialising sensors")
	}

	router := mux.NewRouter()
	router.HandleFunc("/control/", control.ControlGetAll)
	router.HandleFunc("/control/{relay}", control.ControlGetRelay).Methods("GET")
	router.HandleFunc("/control/{relay}", control.ControlSetRelay).Methods("POST")
	router.HandleFunc("/sensors", sensors.SensorsGetAll).Methods("GET")

	return router, nil
}
