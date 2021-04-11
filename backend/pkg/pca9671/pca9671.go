package pca9671

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"sync"
)

type I2CWriter interface {
	Write(b []byte) (int, error)
	Tx(w, r []byte) error
}

// PCA9671 describes the PCA9671 IC
// This is an I2C port multiplexer. It has 16 ports, named P00 - P07 and P10 - P17
type PCA9671 struct {
	state    [2]byte
	device   I2CWriter
	bus      i2c.BusCloser
	address  uint16
	logger   *logrus.Entry
	lock     sync.Locker
	filename string
	metric   *prometheus.GaugeVec
}

// NewPCA9671 creates a new PCA9671 object using address as I2C address
func NewPCA9671(address uint16, filename string) (*PCA9671, error) {
	p := PCA9671{
		address: address,
		state:   [2]byte{0x00, 0x00}, // P07 P06 P05 P04 P03 P02 P01 P00, P17 P16 P15 P14 P13 P12 P11 P10
		logger: logrus.WithFields(logrus.Fields{
			"address": fmt.Sprintf("%x", address),
			"package": "pca9671",
		}),
		lock:     &sync.Mutex{},
		filename: filename,
	}

	p.metric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "control_port_status",
			Help: "The state of the control port",
		},
		[]string{"port"},
		)
	prometheus.MustRegister(p.metric)

	if _, err := host.Init(); err != nil {
		p.logger.WithError(err).Fatal("registering I2C")
	}

	b, err := i2creg.Open("")
	if err != nil {
		p.logger.WithError(err).Fatal("opening I2C bus")
	}
	p.bus = b
	p.device = &i2c.Dev{Addr: address, Bus: b}

	err = p.restoreDump() // restore the previous known state
	return &p, err        // FIXME: return p.Check()
}

func (p *PCA9671) Close() error {
	return p.bus.Close()
}

// setMetrics updates the prometheus metrics
func (p *PCA9671) setMetrics() {
	if p.metric == nil {
		return // for running the tests
	}
	for port, state := range p.GetAll() {
		value := 0
		if state {
			value = 1
		}
		p.metric.WithLabelValues(fmt.Sprintf("%d", port)).Set(float64(value))
	}
}

// restoreDump reads the json file containing the state, written by storeDump
func (p *PCA9671) restoreDump() error {
	p.logger.WithField("func", "restoreDump").Debugf("restoring state from '%v'", p.filename)

	file, err := ioutil.ReadFile(p.filename)
	if os.IsNotExist(err) {
		// file doesn't exist, so we don't attempt to read it
		p.logger.WithField("func", "restoreDump").Debugf("file '%v' does not exist, not restoring anything", p.filename)
		return p.writeState()
	}

	if err != nil {
		return errors.Wrap(err, "opening data dump")
	}

	var stateMap map[int]bool
	err = json.Unmarshal(file, &stateMap)
	if err != nil {
		return errors.Wrap(err, "parsing data dump")
	}

	err = p.SetAll(stateMap)
	p.setMetrics()
	return errors.Wrap(err, "setting state from data dump")
}

// storeDump dumps the current state to the json file
func (p *PCA9671) storeDump() error {
	p.logger.WithField("func", "storeDump").Debugf("dumping state to '%v'", p.filename)
	p.setMetrics()
	file, err := json.MarshalIndent(p.GetAll(), "", " ")
	if err != nil {
		return errors.Wrap(err, "dumping state to json")
	}

	p.logger.WithField("func", "storeDump").Debugf("dumping state '%s'", file)
	err = ioutil.WriteFile(p.filename, file, 0644)
	if err != nil {
		return errors.Wrap(err, "writing json file")
	}
	return nil
}

// Check polls the device to see that it's connected
func (p *PCA9671) Check() error {
	// addr 1111 1000, addr-device+0
	//  Re-START
	// addr 1111 1001, read
	// NACK
	device := &i2c.Dev{Addr: 248, Bus: p.bus}
	data := make([]byte, 3)
	err := device.Tx([]byte{byte(p.address)}, data)
	if err != nil {
		return errors.Wrap(err, "Opening reading device ID")
	}
	p.logger.Debugf("received ID byte: %#b", data)
	// FIXME: check response
	// m m m m m m m m  c c c c c c c f  f p p p p r r r
	// manufacturer(8), category(7), feature(2+4), revision(3)

	return errors.New("Not implemented yet")
}

// GetAll gets the state of all ports
func (p *PCA9671) GetAll() map[int]bool {
	result := make(map[int]bool, 16)
	p.lock.Lock()
	defer p.lock.Unlock()
	for i := 0; i < 16; i++ {
		port := i
		if port > 7 {
			port += 2 // Ports go from 0-7 and 10-17
		}
		result[port] = getBit(p.state, port)
	}
	return result
}

// SetAll sets the state of all ports
func (p *PCA9671) SetAll(state map[int]bool) error {
	p.lock.Lock()
	for i, isSet := range state {
		p.state = setBit(p.state, i, isSet)
	}
	p.lock.Unlock()
	return p.writeState()
}

// writeState sends the port states over I2C
func (p *PCA9671) writeState() error {
	p.logger.WithField("func", "writeState").Debugf("Writing state: %#b", p.state)
	first := p.state[0]
	second := p.state[1]
	err := p.device.Tx([]byte{first, second}, nil)
	if err != nil {
		return errors.Wrap(err, "writing to device")
	}

	return p.storeDump()
}

func (p *PCA9671) ReadState() error {
	data := make([]byte, 2)
	err := p.device.Tx(nil, data)
	p.logger.WithField("func", "ReadState").Debugf("data: %#b err: %v", data, err)
	result := [2]byte{data[0], data[1]}
	if err == nil {
		p.state = result
	}
	return errors.Wrap(err, "reading from device")
}

// Get gets the state of the requested port
func (p *PCA9671) Get(port int) (bool, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	err := p.ReadState()
	return getBit(p.state, port), errors.Wrap(err, "reading state")
}

// Set sets the state of the requested port
func (p *PCA9671) Set(port int, state bool) error {
	p.lock.Lock()
	p.state = setBit(p.state, port, state)
	p.lock.Unlock()
	return p.writeState()
}

// setBit returns the a copy of the store with the port bit set
func setBit(store [2]byte, port int, state bool) [2]byte {
	if port < 8 {
		if state {
			store[0] |= 1 << port
		} else {
			store[0] &^= 1 << port
		}
	} else {
		port -= 10
		if state {
			store[1] |= 1 << port
		} else {
			store[1] &^= 1 << port
		}
	}
	return store
}

// getBit returns the state of the port bit from the store
func getBit(store [2]byte, port int) bool {
	logrus.Debugf("getBit: store='%#b', port:%d", store, port)
	if port < 8 {
		return store[0]>>port%2 == 1
	} else {
		port -= 10
		return store[1]>>port%2 == 1
	}
}
