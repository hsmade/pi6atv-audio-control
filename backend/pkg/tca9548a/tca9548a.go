package tca9548a

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type I2CWriter interface {
	Write(b []byte) (int, error)
	Tx(w, r []byte) error
}

// TCA9548a describes the TCA9548a IC
// This is an I2C port multiplexer. It has 7 i2c ports, named P00 - P07
type TCA9548a struct {
	device      I2CWriter
	bus         i2c.BusCloser
	address     uint16
	logger      *logrus.Entry
	portMetric  *prometheus.GaugeVec

}

func NewTCA9548a(address uint16) (*TCA9548a, error) {
	if address < 112 || address > 119 {
		return nil, errors.New("invalid address. Must be between 0x70 and 0x77 or 112 and 119 decimal")
	}

	t := TCA9548a{
		address: address,
		logger: logrus.WithFields(logrus.Fields{
			"address": fmt.Sprintf("%x", address),
			"package": "TCA9548a",
		}),
	}

	t.portMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "programmer_port_status",
			Help: "The state of the programmer port",
		},
		[]string{"port"},
	)

	prometheus.MustRegister(t.portMetric)

	if _, err := host.Init(); err != nil {
		t.logger.WithError(err).Fatal("registering I2C")
	}

	b, err := i2creg.Open("")
	if err != nil {
		t.logger.WithError(err).Fatal("opening I2C bus")
	}
	t.bus = b
	t.device = &i2c.Dev{Addr: address, Bus: b}

	return &t, err
}

func (t *TCA9548a) Close() error {
	return t.bus.Close()
}

// Set links the port specified
// It will only allow one port to be enabled, so it disables all others
func (t *TCA9548a) Set(port int) error {
	t.logger.Infof("setting control to port %d", port)
	data := []byte{0x00}
	if port != 255 {
		data = []byte{1 << port}
	}
	_, err := t.device.Write(data)
	if err != nil {
		t.logger.WithError(err).Warnf("failed to set control to port %d", port)
	} else {
		t.logger.Debugf("control set to port %d", port)
	}
	return errors.Wrap(err, "fail writing control byte")
}

// Get returns the current enabled port, returns 255/error if multiple or when it fails to read
func (t *TCA9548a) Get() (int, error) {
	t.logger.Debug("reading control port")
	data := make([]byte, 1)
	err := t.device.Tx(nil, data)
	if err != nil {
		t.logger.WithError(err).Warn("failed to read control port")
		return 255, err
	}
	t.logger.Debugf("control port read: %#b", data)
	switch uint8(data[0]) {
	case 0b00000000:
		return 255, nil
	case 0b00000001:
		return 0, nil
	case 0b00000010:
		return 1, nil
	case 0b00000100:
		return 2, nil
	case 0b00001000:
		return 3, nil
	case 0b00010000:
		return 4, nil
	case 0b00100000:
		return 5, nil
	case 0b01000000:
		return 6, nil
	case 0b10000000:
		return 7, nil
	default:
		msg := fmt.Sprintf("multiple ports selected: %#b", data)
		t.logger.Warn(msg)
		return 255, errors.New(msg)
	}
}

