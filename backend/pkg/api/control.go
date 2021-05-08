package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/config"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/i2cMultiplexer"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/ic2IOExpander"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Control struct {
	ioExpander  *ic2IOExpander.PCA9671
	multiplexer *i2cMultiplexer.TCA9548a
}

// NewControl takes a config and creates a new Control object
func NewControl(config *config.Config) (*Control, error) {
	ioExpander, err := ic2IOExpander.NewPCA9671(config.Backend.Pca.Address, config.Backend.Pca.Filename)
	if err != nil {
		return nil, errors.Wrap(err, "initialising new PCA object")
	}
	multiplexer, err := i2cMultiplexer.NewTCA9548a(config.Backend.Tca.Address)
	if err != nil {
		return nil, errors.Wrap(err, "initialising new TCA object")
	}
	c := Control{
		ioExpander:  ioExpander,
		multiplexer: multiplexer,
	}

	// it could be that when we last exited, the i2c to the multiplexer was disabled, so make sure it's enabled
	_ = ioExpander.Set(0, true)
	_ = ioExpander.Set(17, true)
	return &c, nil
}

// delete below

func (c *Control) GetAll(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("GetAll called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	err := c.ioExpander.ReadState()
	if err != nil {
		logrus.WithError(err).Warn("PCA ReadState returned error")
		Error{Message: "PCA ReadState returned error", Error: err}.Send(w)
		return
	}
	tcaPort, err := c.multiplexer.Get()
	if err != nil {
		logrus.WithError(err).Warn("TCA Get returned error")
		Error{Message: "TCA Get returned error", Error: err}.Send(w)
		return
	} else {
		w.WriteHeader(200)
	}
	state := c.ioExpander.GetAll()

	_ = json.NewEncoder(w).Encode(struct {
		Pca map[int]bool
		Tca int
	}{
		Pca: state,
		Tca: tcaPort,
	})
}

func (c *Control) CarrierGet(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("CarrierGet called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["relay"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse port from '%s'", vars["relay"]), Error: err}.Send(w)
		return
	}
	state, err := c.ioExpander.Get(port)
	if err != nil {
		Error{Message: fmt.Sprint("failed to read from device"), Error: err}.Send(w)
		return
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(state)
}

func (c *Control) CarrierSet(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("CarrierSet called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["relay"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse port from '%s'", vars["relay"]), Error: err}.Send(w)
		return
	}
	state, err := strconv.ParseBool(vars["state"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse state from '%s'", vars["relay"]), Error: err}.Send(w)
		return
	}
	err = c.ioExpander.Set(port, state)

	if err != nil {
		Error{Message: "failed to set PCA state", Error: err}.Send(w)
	} else {
		w.WriteHeader(200)
		_,_ = w.Write([]byte("{}"))
	}
}

func (c *Control) ProgrammerSet(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ProgrammerSet called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["port"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse port from '%s'", vars["port"]), Error: err}.Send(w)
		return
	}
	err = c.multiplexer.Set(port)

	if err != nil {
		Error{Message: "failed to set PCA state", Error: err}.Send(w)
		return
	}
	w.WriteHeader(200)
	_,_ = w.Write([]byte("{}"))
}

func (c *Control) ControlCheck(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ControlCheck called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	err := c.ioExpander.Check()
	if err != nil {
		Error{Message: "error checking PCA", Error: err}.Send(w)
		return
	}
	w.WriteHeader(200)
}
