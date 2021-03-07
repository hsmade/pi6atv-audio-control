package api

import (
	"github.com/gorilla/mux"
)

func Router() {
	router := mux.NewRouter()
	router.HandleFunc("/control/", ControlGetAll)
	router.HandleFunc("/control/{relay}", ControlGetRelay).Methods("GET")
	router.HandleFunc("/control/{relay}", ControlSetRelay).Methods("POST")
	router.HandleFunc("/sensors", SensorsGetAll).Methods("GET")
}