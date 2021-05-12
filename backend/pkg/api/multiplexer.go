package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// GetMultiplexer fetches the state of the multiplexer and returns it
func (c *Control) GetMultiplexer(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("GetMultiplexer called with %v", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	selectedPort, err := c.multiplexer.Get()
	if err != nil {
		msg := "GetMultiplexer returned error"
		logrus.WithError(err).Warn(msg)
		Error{Message: msg, Error: err}.Send(w)
		return
	}
	logrus.Debugf("GetMultiplexer: result: %d", selectedPort)

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(struct {
		Port int
	}{selectedPort})
}

// SetMultiplexer sets the port to be selected on the multiplexer.
func (c *Control) SetMultiplexer(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("SetMultiplexer called with %v", r.URL.Path)
	c.lock.Lock()
	defer c.lock.Unlock()
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["port"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("SetMultiplexer: failed to parse port from '%s'", vars["port"]), Error: err}.Send(w)
		return
	}

	err = c.multiplexer.Set(port)
	if err != nil {
		msg := "SetMultiplexer returned error"
		logrus.WithError(err).Warn(msg)
		Error{Message: msg, Error: err}.Send(w)
		return
	}
	w.WriteHeader(200)
	_,_ = w.Write([]byte("{}")) // empty json
}
