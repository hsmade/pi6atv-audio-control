package pca9671

type PCA9671 struct {
	Address [2]byte
	Bus     int
	state	[4]byte
}

func NewPCA9671(address [2]byte, bus int) (*PCA9671, error) {
	p := PCA9671{
		Address: address,
		Bus:     bus,
		state:   [4]byte{0x00, 0x00},
	}
	return &p, p.Check()
}

// check polls the device to see that it's connected
func (p *PCA9671) Check() error {
	return nil
}

// GetAll gets the state of all ports
func (p *PCA9671) GetAll() (map[int]bool, error){
	return nil, nil
}

// SetAll sets the state of all ports
func (p *PCA9671) SetAll(state map[int]bool) error{
	return nil
}
// Get gets the state of the requested port
func (p *PCA9671) Get(port int) (bool, error){
	return false, nil
}

// Set sets the state of the requested port
func (p *PCA9671) Set(port int, state bool) error{
	return nil
}
