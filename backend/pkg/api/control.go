package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/pca9671"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Control struct {
	pca *pca9671.PCA9671
}

// NewControl takes a config and creates a new Control object
func NewControl() (*Control, error) {
	pca, err := pca9671.NewPCA9671(0x20) // FIXME
	if err != nil {
		return nil, errors.Wrap(err, "initialising new Control object")
	}
	c := Control{
		pca: pca,
	}
	return &c, nil
}

func (c *Control) ControlGetAll(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ControlGetAll called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	state := c.pca.GetAll()

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(state)
}

func (c *Control) ControlGetRelay(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ControlGetRelay called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["relay"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse port from '%s'", vars["relay"]), Error: err}.Send(w)
	}
	state, err := c.pca.Get(port)
	if err != nil {
		Error{Message: fmt.Sprint("failed to read from device"), Error: err}.Send(w)
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(state)
}

func (c *Control) ControlSetRelay(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ControlSetRelay called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["relay"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse port from '%s'", vars["relay"]), Error: err}.Send(w)
	}
	state, err := strconv.ParseBool(vars["state"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse state from '%s'", vars["relay"]), Error: err}.Send(w)
	}
	result := c.pca.Set(port, state)

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(result)
}

func (c *Control) ControlQuit(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ControlQuit called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	err := c.pca.Close()
	if err != nil {
		Error{Message: "error closing PCA", Error: err}.Send(w)
	}
	w.WriteHeader(200)
}

func (c *Control) ControlCheck(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ControlCheck called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	err := c.pca.Check()
	if err != nil {
		Error{Message: "error checking PCA", Error: err}.Send(w)
	}
	w.WriteHeader(200)
}
