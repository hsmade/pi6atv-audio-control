package api

import (
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/config"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/i2cMultiplexer"
	"github.com/hsmade/pi6atv-audio-control/backend/pkg/ic2IOExpander"
	"github.com/pkg/errors"
	"sync"
)

type Control struct {
	ioExpander  *ic2IOExpander.PCA9671
	multiplexer *i2cMultiplexer.TCA9548a
	lock        sync.Locker
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
		lock:        &sync.Mutex{},
	}

	// it could be that when we last exited, the i2c to the multiplexer was disabled, so make sure it's enabled
	_ = ioExpander.Set(0, true)
	_ = ioExpander.Set(17, true)
	return &c, nil
}
