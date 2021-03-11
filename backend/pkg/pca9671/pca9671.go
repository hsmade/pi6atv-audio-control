package pca9671

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/periph/host"
	"sync"
)

type I2CWriter interface {
	Write(b []byte) (int, error)
	Tx(w, r []byte) error
}

// PCA9671 describes the PCA9671 IC
// This is an I2C port multiplexer. It has 16 ports, named P00 - P07 and P10 - P17
type PCA9671 struct {
	state   [2]byte
	device  I2CWriter
	bus  	i2c.BusCloser
	address uint16
	logger  *logrus.Entry
	lock    sync.Locker
}

// NewPCA9671 creates a new PCA9671 object using address as I2C address
func NewPCA9671(address uint16) (*PCA9671, error) {
	p := PCA9671{
		address: address,
		state:   [2]byte{0x00, 0x00}, // P07 P06 P05 P04 P03 P02 P01 P00, P17 P16 P15 P14 P13 P12 P11 P10
		logger: logrus.WithFields(logrus.Fields{
			"address": fmt.Sprintf("%x", address),
			"package": "pca9671",
		}),
		lock: &sync.Mutex{},
	}

	if _, err := host.Init(); err != nil {
		p.logger.WithError(err).Fatal("registering I2C")
	}

	b, err := i2creg.Open("")
	if err != nil {
		p.logger.WithError(err).Fatal("opening I2C bus")
	}
	p.bus = b
	p.device = &i2c.Dev{Addr: address, Bus: b}
	return &p, nil
	//return &p, p.Check()
}

// Check polls the device to see that it's connected
func (p *PCA9671) Check() error {
	// addr 1111 1000, addr-device+0
	//  Re-START
	// addr 1111 1001, read
	// NACK
	device := &i2c.Dev{Addr: 0xF8, Bus: p.bus}
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
	defer p.lock.Unlock()
	for i, isSet := range state {
		p.state = setBit(p.state, i, isSet)
	}
	return p.writeState()
}

// writeState sends the port states over I2C
func (p *PCA9671) writeState() error {
	first := p.state[0]
	second := p.state[1]
	size, err := p.device.Write([]byte{first, second})
	if err != nil {
		return errors.Wrap(err, "writing to device")
	}
	if size != 2 {
		return errors.New(fmt.Sprintf("write wrote %d bytes, instead of 2", size))
	}
	return nil
}

// Get gets the state of the requested port
func (p *PCA9671) Get(port int) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return getBit(p.state, port)
}

// Set sets the state of the requested port
func (p *PCA9671) Set(port int, state bool) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.state = setBit(p.state, port, state)
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
