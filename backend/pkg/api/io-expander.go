package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// GetIOExpanderPort gets the state of a port on the IO expander
func (c *Control) GetIOExpanderPort(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("GetIOExpanderPort called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["port"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse port from '%s'", vars["port"]), Error: err}.Send(w)
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

// GetIOExpander gets the state of the IO expander
func (c *Control) GetIOExpander(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("GetIOExpander called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	err := c.ioExpander.ReadState()
	if err != nil {
		Error{Message: fmt.Sprint("failed to read from device"), Error: err}.Send(w)
		return
	}
	state := c.ioExpander.GetAll()

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(state)
}

// SetIOExpanderPort sets the state of a port on the IO expander
func (c *Control) SetIOExpanderPort(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("SetIOExpanderPort called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["port"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse port from '%s'", vars["port"]), Error: err}.Send(w)
		return
	}

	state, err := strconv.ParseBool(vars["state"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse state from '%s'", vars["state"]), Error: err}.Send(w)
		return
	}

	err = c.ioExpander.Set(port, state)
	if err != nil {
		Error{Message: "failed to set IO expander state", Error: err}.Send(w)
		return
	}

	w.WriteHeader(200)
	_,_ = w.Write([]byte("{}"))
}
