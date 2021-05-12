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
	router.HandleFunc("/control/mpx/", control.GetMultiplexer).Methods("GET")
	router.HandleFunc("/control/mpx/{port}", control.SetMultiplexer).Methods("POST")
	router.HandleFunc("/control/io/", control.GetIOExpander).Methods("GET")
	router.HandleFunc("/control/io/{port}", control.GetIOExpanderPort).Methods("GET")
	router.HandleFunc("/control/io/{port}/{state}", control.SetIOExpanderPort).Methods("POST")
	router.Handle("/metrics", promhttp.Handler())

	return router, nil
}
