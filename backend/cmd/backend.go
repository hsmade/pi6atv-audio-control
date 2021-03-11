package main

import (
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/api"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"periph.io/x/periph/host"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	host.Init()
	router, err := api.NewRouter(config.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Fatal(http.ListenAndServe(":8001", router))
}
