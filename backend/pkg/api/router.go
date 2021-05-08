package api

import (
	"github.com/gorilla/mux"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/config"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter initialises the control and sensors objects and creates the http router to reach them
func NewRouter(config *config.Config) (*mux.Router, error) {
	control, err := NewControl(config)
	if err != nil {
		return nil, errors.Wrap(err, "initialising control")
	}

	router := mux.NewRouter()
	// /control
	//   /io/
	//     / (get states)
	//     /{port}/{state} (set)
	//  /mpx
	//    / (get state)
	//    /{port} (set)
	router.HandleFunc("/control/", control.GetAll)
	router.HandleFunc("/control/programmer/{port}", control.ProgrammerSet).Methods("POST")
	router.HandleFunc("/control/{relay}", control.CarrierGet).Methods("GET")
	router.HandleFunc("/control/{relay}/{state}", control.CarrierSet).Methods("POST")
	router.Handle("/metrics", promhttp.Handler())

	return router, nil
}
