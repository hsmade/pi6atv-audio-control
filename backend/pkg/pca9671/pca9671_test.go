package pca9671

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/io/i2c"
	"reflect"
	"sync"
	"testing"
)

type fakeI2CDevice struct {
	Data []byte
}

func (f *fakeI2CDevice) Write(b []byte) error {
	f.Data = make([]byte, len(b))
	copy(f.Data, b)
	return nil
}

// TODO
func TestNewPCA9671(t *testing.T) {
	type args struct {
		address int
	}
	tests := []struct {
		name    string
		args    args
		want    *PCA9671
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPCA9671(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPCA9671() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPCA9671() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO
func TestPCA9671_Check(t *testing.T) {
	type fields struct {
		state   [2]byte
		device  *i2c.Device
		address int
		logger  *logrus.Entry
		lock    sync.Locker
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				state:   tt.fields.state,
				device:  tt.fields.device,
				address: tt.fields.address,
				logger:  tt.fields.logger,
				lock:    tt.fields.lock,
			}
			if err := p.Check(); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPCA9671_Get(t *testing.T) {
	type fields struct {
		state   [2]byte
		device  *i2c.Device
		address int
		logger  *logrus.Entry
		lock    sync.Locker
	}
	type args struct {
		port int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Get P03, set",
			fields: fields{state: [2]byte{0b00001000, 0b00000000}, lock: &sync.Mutex{}},
			args:   args{port: 3},
			want:   true,
		},
		{
			name:   "Get P15, cleared",
			fields: fields{state: [2]byte{0b11111111, 0b11011111}, lock: &sync.Mutex{}},
			args:   args{port: 15},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				state:   tt.fields.state,
				device:  tt.fields.device,
				address: tt.fields.address,
				logger:  tt.fields.logger,
				lock:    tt.fields.lock,
			}
			if got := p.Get(tt.args.port); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPCA9671_GetAll(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	type fields struct {
		state   [2]byte
		device  *i2c.Device
		address int
		logger  *logrus.Entry
		lock    sync.Locker
	}
	tests := []struct {
		name   string
		fields fields
		want   map[int]bool
	}{
		{
			name:   "get all random 1",
			fields: fields{state: [2]byte{0b10101010, 0b01010101}, lock: &sync.Mutex{}},
			want: map[int]bool{
				0: false, 1: true, 2: false, 3: true, 4: false, 5: true, 6: false, 7: true,
				10: true, 11: false, 12: true, 13: false, 14: true, 15: false, 16: true, 17: false},
		},
		{
			name:   "get all random 1",
			fields: fields{state: [2]byte{0b11001001, 0b00101110}, lock: &sync.Mutex{}},
			want: map[int]bool{
				0: true, 1: false, 2: false, 3: true, 4: false, 5: false, 6: true, 7: true,
				10: false, 11: true, 12: true, 13: true, 14: false, 15: true, 16: false, 17: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				state:   tt.fields.state,
				device:  tt.fields.device,
				address: tt.fields.address,
				logger:  tt.fields.logger,
				lock:    tt.fields.lock,
			}
			if got := p.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPCA9671_Set(t *testing.T) {
	device := fakeI2CDevice{}

	type fields struct {
		state   [2]byte
		device  I2CWriter
		address int
		logger  *logrus.Entry
		lock    sync.Locker
	}
	type args struct {
		port  int
		state bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    []byte
	}{
		{
			name: "set P03",
			fields: fields{
				state:   [2]byte{},
				device:  &device,
				address: 23,
				lock:    &sync.Mutex{},
			},
			args:    args{port: 3, state: true},
			wantErr: false,
			want:    []byte{0b00001000, 0b00000000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				state:   tt.fields.state,
				device:  tt.fields.device,
				address: tt.fields.address,
				logger:  tt.fields.logger,
				lock:    tt.fields.lock,
			}
			if err := p.Set(tt.args.port, tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.fields.device.(*fakeI2CDevice).Data
			if bytes.Compare(tt.want, got) != 0 {
				t.Logf("got : %v, %d, %#b", got, len(got), got)
				t.Logf("want: %v, %d, %#b", tt.want, len(tt.want), tt.want)
				t.Errorf("written state=%#b, want=%#b", got, tt.want)
			}
		})
	}
}

func TestPCA9671_SetAll(t *testing.T) {
	device := fakeI2CDevice{}

	type fields struct {
		state   [2]byte
		device  I2CWriter
		address int
		logger  *logrus.Entry
		lock    sync.Locker
	}
	type args struct {
		state map[int]bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    []byte
	}{
		{
			name: "write 0b10010011 0b11001011",
			fields: fields{
				state:   [2]byte{},
				device:  &device,
				address: 23,
				lock:    &sync.Mutex{},
			},
			args: args{state: map[int]bool{
				0: true, 1: true, 2: false, 3: false, 4: true, 5: false, 6: false, 7: true,
				10: true, 11: true, 12: false, 13: true, 14: false, 15: false, 16: true, 17: true,
			}},
			wantErr: false,
			want:    []byte{0b10010011, 0b11001011},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				state:   tt.fields.state,
				device:  tt.fields.device,
				address: tt.fields.address,
				logger:  tt.fields.logger,
				lock:    tt.fields.lock,
			}
			if err := p.SetAll(tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("SetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.fields.device.(*fakeI2CDevice).Data
			if bytes.Compare(tt.want, got) != 0 {
				t.Logf("got : %v, %d, %#b", got, len(got), got)
				t.Logf("want: %v, %d, %#b", tt.want, len(tt.want), tt.want)
				t.Errorf("written state=%#b, want=%#b", got, tt.want)
			}
		})
	}
}

func TestPCA9671_writeState(t *testing.T) {
	device := fakeI2CDevice{}

	type fields struct {
		state   [2]byte
		device  I2CWriter
		address int
		logger  *logrus.Entry
		lock    sync.Locker
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    []byte
	}{
		{
			name: "write 0b10101100 0b00110001",
			fields: fields{
				state:  [2]byte{0b10101100, 0b00110001},
				device: &device,
			},
			wantErr: false,
			want:    []byte{0b10101100, 0b00110001},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				state:   tt.fields.state,
				device:  tt.fields.device,
				address: tt.fields.address,
				logger:  tt.fields.logger,
				lock:    tt.fields.lock,
			}
			if err := p.writeState(); (err != nil) != tt.wantErr {
				t.Errorf("writeState() error = %v, wantErr %v", err, tt.wantErr)
				got := tt.fields.device.(*fakeI2CDevice).Data
				if bytes.Compare(tt.want, got) != 0 {
					t.Logf("got : %v, %d, %#b", got, len(got), got)
					t.Logf("want: %v, %d, %#b", tt.want, len(tt.want), tt.want)
					t.Errorf("written state=%#b, want=%#b", got, tt.want)
				}
			}
		})
	}
}

func Test_getBit(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	type args struct {
		store [2]byte
		port  int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "get port P00, set",
			args: args{
				store: [2]byte{0b00000001, 0b00000000},
				port:  0,
			},
			want: true,
		},
		{
			name: "get port P03, set",
			args: args{
				store: [2]byte{0b00001000, 0b00000000},
				port:  3,
			},
			want: true,
		},
		{
			name: "get port P00, clear",
			args: args{
				store: [2]byte{0b00000000, 0b00000000},
				port:  0,
			},
			want: false,
		},
		{
			name: "get port P07, set",
			args: args{
				store: [2]byte{0b10000000, 0b00000000},
				port:  7,
			},
			want: true,
		},
		{
			name: "get port P10, set",
			args: args{
				store: [2]byte{0b00000000, 0b00000001},
				port:  10,
			},
			want: true,
		},
		{
			name: "get port P13, set",
			args: args{
				store: [2]byte{0b00000000, 0b00001000},
				port:  13,
			},
			want: true,
		},
		{
			name: "get port P17, set",
			args: args{
				store: [2]byte{0b00000000, 0b10000000},
				port:  17,
			},
			want: true,
		},
		{
			name: "get port P00, clear, rest set",
			args: args{
				store: [2]byte{0b11111110, 0b11111111},
				port:  0,
			},
			want: false,
		},
		{
			name: "get port P10, clear, rest set",
			args: args{
				store: [2]byte{0b11111111, 0b11111110},
				port:  10,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBit(tt.args.store, tt.args.port); got != tt.want {
				t.Errorf("getBit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setBit(t *testing.T) {
	type args struct {
		store [2]byte
		port  int
		state bool
	}
	tests := []struct {
		name string
		args args
		want [2]byte
	}{
		{
			name: "Set P0 on 0X00 store",
			args: args{
				store: [2]byte{0b00000000, 0b00000000},
				port:  0,
				state: true,
			},
			want: [2]byte{0b00000001, 0b00000000},
		},
		{
			name: "clear P0 on 0xFF store",
			args: args{
				store: [2]byte{0b11111111, 0b11111111},
				port:  0,
				state: false,
			},
			want: [2]byte{0b11111110, 0b11111111},
		},
		{
			name: "Set P3 on 0X00 store",
			args: args{
				store: [2]byte{0b00000000, 0b00000000},
				port:  3,
				state: true,
			},
			want: [2]byte{0b00001000, 0b00000000},
		},
		{
			name: "clear P3 on 0xFF store",
			args: args{
				store: [2]byte{0b11111111, 0b11111111},
				port:  3,
				state: false,
			},
			want: [2]byte{0b11110111, 0b11111111},
		},
		{
			name: "Set P7 on 0X00 store",
			args: args{
				store: [2]byte{0b00000000, 0b00000000},
				port:  7,
				state: true,
			},
			want: [2]byte{0b10000000, 0b00000000},
		},
		{
			name: "clear P7 on 0xFF store",
			args: args{
				store: [2]byte{0b11111111, 0b11111111},
				port:  7,
				state: false,
			},
			want: [2]byte{0b01111111, 0b11111111},
		},
		{
			name: "Set P10 on 0X00 store",
			args: args{
				store: [2]byte{0b00000000, 0b00000000},
				port:  10,
				state: true,
			},
			want: [2]byte{0b00000000, 0b00000001},
		},
		{
			name: "clear P10 on 0xFF store",
			args: args{
				store: [2]byte{0b11111111, 0b11111111},
				port:  10,
				state: false,
			},
			want: [2]byte{0b11111111, 0b11111110},
		},
		{
			name: "Set P13 on 0X00 store",
			args: args{
				store: [2]byte{0b00000000, 0b00000000},
				port:  13,
				state: true,
			},
			want: [2]byte{0b00000000, 0b00001000},
		},
		{
			name: "clear P13 on 0xFF store",
			args: args{
				store: [2]byte{0b11111111, 0b11111111},
				port:  13,
				state: false,
			},
			want: [2]byte{0b11111111, 0b11110111},
		},
		{
			name: "Set P17 on 0X00 store",
			args: args{
				store: [2]byte{0b00000000, 0b00000000},
				port:  17,
				state: true,
			},
			want: [2]byte{0b00000000, 0b10000000},
		},
		{
			name: "clear P17 on 0xFF store",
			args: args{
				store: [2]byte{0b11111111, 0b11111111},
				port:  17,
				state: false,
			},
			want: [2]byte{0b11111111, 0b01111111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setBit(tt.args.store, tt.args.port, tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setBit() = %v, want %v", got, tt.want)
			}
		})
	}
}
