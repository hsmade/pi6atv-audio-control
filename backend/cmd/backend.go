package main

import (
	"flag"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/api"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/config"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	configFile = flag.String("config", "", "the path to the config file")
	verbose = flag.Bool("verbose", false, "enable verbose mode")
)
func main() {
	flag.Parse()

	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	c, err := config.NewConfig(*configFile)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("config: %v", c)
	router, err := api.NewRouter(c)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Fatal(http.ListenAndServe(":8001", router))
}
