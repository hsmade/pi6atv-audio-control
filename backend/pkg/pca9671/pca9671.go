package pca9671

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/io/i2c"
	"sync"
)

type I2CWriter interface {
	Write(b []byte) error
}

// PCA9671 describes the PCA9671 IC
// This is an I2C port multiplexer. It has 16 ports, named P00 - P07 and P10 - P17
type PCA9671 struct {
	state   [2]byte
	device  I2CWriter
	address int
	logger  *logrus.Entry
	lock    sync.Locker
}

// NewPCA9671 creates a new PCA9671 object using address as I2C address
func NewPCA9671(address int) (*PCA9671, error) {
	p := PCA9671{
		address: address,
		state:   [2]byte{0x00, 0x00}, // P07 P06 P05 P04 P03 P02 P01 P00, P17 P16 P15 P14 P13 P12 P11 P10
		logger: logrus.WithFields(logrus.Fields{
			"address": fmt.Sprintf("%x", address),
			"package": "ina260",
		}),
		lock: &sync.Mutex{},
	}

	device, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, address)
	if err != nil {
		p.logger.WithError(err).Fatal("failed to open i2c device")
	}
	p.device = device
	return &p, p.Check()
}

// Check polls the device to see that it's connected
func (p *PCA9671) Check() error {
	// addr 1111 1000, addr-device+0
	//  Re-START
	// addr 1111 1001, read
	// NACK
	device, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0xF8)
	if err != nil {
		return errors.Wrap(err, "Opening device for ID")
	}
	data := make([]byte, 3)
	err = device.WriteReg(byte(p.address), data)
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
	first := 0xFF & ^(p.state[1])
	second := 0xFF & ^(p.state[0])
	return p.device.Write([]byte{first, second})
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
