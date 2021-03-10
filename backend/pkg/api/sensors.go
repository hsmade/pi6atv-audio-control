package api

import "net/http"

type Sensors struct {
}

func NewSensors() (*Sensors, error) {
	return nil, nil
}

func (s *Sensors) SensorsGetAll(w http.ResponseWriter, r *http.Request) {

}
