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
	pca, err := pca9671.NewPCA9671(20) // FIXME
	if err != nil {
		return nil, errors.Wrap(err, "initialising new Control object")
	}
	c := Control{
		pca: pca,
	}
	return &c, nil
}

func (c *Control) ControlGetAll(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ControlGetAll called with %v", r.URL.RawQuery)
	w.Header().Set("Content-Type", "application/json")

	state := c.pca.GetAll()

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(state)
}

func (c *Control) ControlGetRelay(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("ControlGetAll called with %v", r.URL.RawQuery)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	port, err := strconv.Atoi(vars["relay"]) // FIXME: how to make sure the param is there?
	if err != nil {
		Error{Message: fmt.Sprintf("failed to parse port from '%s'", vars["relay"]), Error: err}.Send(w)
	}
	state := c.pca.Get(port)

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(state)
}

func (c *Control) ControlSetRelay(w http.ResponseWriter, r *http.Request) {

}
