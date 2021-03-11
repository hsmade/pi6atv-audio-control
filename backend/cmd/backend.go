package main

import (
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/api"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/config"
	"log"
	"net/http"
	"periph.io/x/periph/host"

)

func main() {
	host.Init()
	router, err := api.NewRouter(config.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8001", router))
}
