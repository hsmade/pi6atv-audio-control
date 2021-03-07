package ina260

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/io/i2c"
)

type Ina260 struct {
	device *i2c.Device
	logger *logrus.Entry
}

func NewIna260(dev int) *Ina260 {
	I := &Ina260{
		logger: logrus.WithFields(logrus.Fields{
			"address": fmt.Sprintf("%x", dev),
			"package": "ina260",
		}),
	}

	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, dev)
	if err != nil {
		I.logger.WithError(err).Fatal("failed to open i2c device")
	}
	I.device = d

	if !I.isIna260() {
		I.logger.Fatal("This is not an INA260")
	}
	return I
}

func (I *Ina260) Begin() error {
	return I.reset()
}

func (I *Ina260) isIna260() bool {
	buffer := make([]byte, 2)
	err := I.device.ReadReg(0xFE, buffer)
	return err == nil && len(buffer) == 2 && buffer[0] == 0x54 && buffer[1] == 0x49
}

func (I *Ina260) reset() error {
	buffer := make([]byte, 2)
	err := I.device.ReadReg(0x00, buffer)
	if err != nil {
		return err
	}
	buffer[1] |= 1 << 7
	err = I.device.WriteReg(0x00, buffer)
	return err
}

func (I *Ina260) readValue(address byte) (float64, error) {
	buffer := make([]byte, 2)
	err := I.device.ReadReg(address, buffer)
	if err != nil {
		return 0, err
	}

	value := float64(int(buffer[0])*256 + int(buffer[1]))
	I.logger.Debugf("read from register %v: %v -> %f", address, buffer, value)
	return value, nil
}

func (I *Ina260) ReadVoltage() (float64, error) {
	v, err := I.readValue(0x02)
	return v * 1.25 / 1000, err
}

func (I *Ina260) ReadCurrent() (float64, error) {
	v, err := I.readValue(0x01)
	return float64(int16(v)) * 1.25, err
}

func (I *Ina260) ReadPower() (float64, error) {
	v, err := I.readValue(0x03)
	return v * 10, err
}
