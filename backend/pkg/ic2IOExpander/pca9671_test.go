package ic2IOExpander

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"periph.io/x/conn/v3/i2c"
	"reflect"
	"sync"
	"testing"
)

type fakeI2CDevice struct {
	WrittenData []byte
	ReadData    []byte
}

func (f *fakeI2CDevice) Write(b []byte) (int, error) {
	f.WrittenData = make([]byte, len(b))
	copy(f.WrittenData, b)
	return len(b), nil
}

func (f *fakeI2CDevice) Tx(w []byte, r []byte) error {
	f.WrittenData = make([]byte, len(w))
	copy(f.WrittenData, w)
	if len(f.ReadData) > 0 {
		logrus.Debugf("fakeI2CDevice: returning: %#b", f.ReadData)
		copy(r, f.ReadData)
	}
	return nil
}

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	code := m.Run()
	os.Exit(code)
}

// TODO
func TestNewPCA9671(t *testing.T) {
	tests := []struct {
		name         string
		wantState    [2]byte
		wantErr      bool
		haveFile     bool
		fileContents []byte
	}{}
	for _, tt := range tests {
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			t.Error(err)
		}

		if tt.fileContents != nil || tt.haveFile {
			_, err := tmpfile.Write(tt.fileContents)
			if err != nil {
				t.Error(err)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPCA9671(0x00, tmpfile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPCA9671() error = %v, wantErr %v", err, tt.wantErr)
				_ = os.Remove(tmpfile.Name())
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewPCA9671() got = %v, want %v", got, tt.want)
			//}
			if !reflect.DeepEqual(got.state, tt.wantState) {
				t.Errorf("NewPCA9671() got state = %v, want %v", got, tt.wantState)
			}
			_ = os.Remove(tmpfile.Name())
		})
	}
}

func TestPCA9671_Get(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	type fields struct {
		state   [2]byte
		device  I2CWriter
		address uint16
		logger  *logrus.Entry
		lock    sync.Locker
	}
	type args struct {
		port int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Get P03, set",
			fields:  fields{state: [2]byte{0b00001000, 0b00000000}, lock: &sync.Mutex{}, address: 0x20},
			args:    args{port: 3},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Get P15, cleared",
			fields:  fields{state: [2]byte{0b11111111, 0b11011111}, lock: &sync.Mutex{}},
			args:    args{port: 15},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := fakeI2CDevice{ReadData: []byte{tt.fields.state[0], tt.fields.state[1]}}
			p := &PCA9671{
				state:   tt.fields.state,
				device:  &device,
				address: tt.fields.address,
				logger:  logrus.WithField("foo", "bar"),
				lock:    tt.fields.lock,
			}
			got, err := p.Get(tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPCA9671_GetAll(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	type fields struct {
		state   [2]byte
		device  I2CWriter
		address uint16
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
		address uint16
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
			tmpfile, err := ioutil.TempFile("", "example")
			if err != nil {
				t.Error(err)
			}

			p := &PCA9671{
				state:    tt.fields.state,
				device:   tt.fields.device,
				address:  tt.fields.address,
				logger:   logrus.WithField("foo", "bar"),
				lock:     tt.fields.lock,
				filename: tmpfile.Name(),
			}
			if err := p.Set(tt.args.port, tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.fields.device.(*fakeI2CDevice).WrittenData
			if bytes.Compare(tt.want, got) != 0 {
				t.Logf("got : %v, %d, %#b", got, len(got), got)
				t.Logf("want: %v, %d, %#b", tt.want, len(tt.want), tt.want)
				t.Errorf("written state=%#b, want=%#b", got, tt.want)
			}
			_ = os.Remove(tmpfile.Name())
		})
	}
}

func TestPCA9671_SetAll(t *testing.T) {
	device := fakeI2CDevice{}

	type fields struct {
		state   [2]byte
		device  I2CWriter
		address uint16
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
			tmpfile, err := ioutil.TempFile("", "example")
			if err != nil {
				t.Error(err)
			}

			p := &PCA9671{
				state:    tt.fields.state,
				device:   tt.fields.device,
				address:  tt.fields.address,
				logger:   logrus.WithField("foo", "bar"),
				lock:     tt.fields.lock,
				filename: tmpfile.Name(),
			}
			if err := p.SetAll(tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("SetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.fields.device.(*fakeI2CDevice).WrittenData
			if bytes.Compare(tt.want, got) != 0 {
				t.Logf("got : %v, %d, %#b", got, len(got), got)
				t.Logf("want: %v, %d, %#b", tt.want, len(tt.want), tt.want)
				t.Errorf("written state=%#b, want=%#b", got, tt.want)
			}
			_ = os.Remove(tmpfile.Name())
		})
	}
}

func TestPCA9671_writeState(t *testing.T) {
	device := fakeI2CDevice{}

	type fields struct {
		state   [2]byte
		device  I2CWriter
		address uint16
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
			tmpfile, err := ioutil.TempFile("", "example")
			if err != nil {
				t.Error(err)
			}

			p := &PCA9671{
				state:    tt.fields.state,
				device:   tt.fields.device,
				address:  tt.fields.address,
				logger:   logrus.WithField("foo", "bar"),
				lock:     &sync.Mutex{},
				filename: tmpfile.Name(),
			}
			if err := p.writeState(); (err != nil) != tt.wantErr {
				t.Errorf("writeState() error = %v, wantErr %v", err, tt.wantErr)
				got := tt.fields.device.(*fakeI2CDevice).WrittenData
				if bytes.Compare(tt.want, got) != 0 {
					t.Logf("got : %v, %d, %#b", got, len(got), got)
					t.Logf("want: %v, %d, %#b", tt.want, len(tt.want), tt.want)
					t.Errorf("written state=%#b, want=%#b", got, tt.want)
				}
			}
			_ = os.Remove(tmpfile.Name())
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

func TestPCA9671_restoreDump(t *testing.T) {
	type fields struct {
		state    [2]byte
		device   I2CWriter
		bus      i2c.BusCloser
		address  uint16
		logger   *logrus.Entry
		lock     sync.Locker
		filename string
	}
	tests := []struct {
		name         string
		fields       fields
		wantErr      bool
		haveFile     bool
		fileContents []byte
		filename     string
		wantState    [2]byte
	}{
		{
			name:         "empty file",
			wantErr:      true,
			wantState:    [2]byte{0b00000000, 0b00000000},
			fileContents: nil,
			haveFile:     true,
		},
		{
			name:      "nonexistent file",
			wantErr:   false,
			wantState: [2]byte{0b00000000, 0b00000000},
			haveFile:  false,
			filename:  "non-existent",
			fields:    fields{device: &fakeI2CDevice{}},
		},
		{
			name:      "valid json file",
			wantErr:   false,
			wantState: [2]byte{0b01010001, 0b11010101},
			fileContents: []byte(`
				{
				  "0": true,
				  "1": false,
				  "2": false,
				  "3": false,
				  "4": true,
				  "5": false,
				  "6": true,
				  "7": false,
				  "10": true,
				  "11": false,
				  "12": true,
				  "13": false,
				  "14": true,
				  "15": false,
				  "16": true,
				  "17": true
						}
			`),
			haveFile: true,
			fields:   fields{device: &fakeI2CDevice{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.haveFile {
				tmpfile, err := ioutil.TempFile("", "example")
				if err != nil {
					t.Error(err)
				}
				tt.filename = tmpfile.Name()

				if tt.fileContents != nil {
					_, _ = tmpfile.Write(tt.fileContents)
				}
			}

			p := &PCA9671{
				state:    tt.fields.state,
				device:   tt.fields.device,
				bus:      tt.fields.bus,
				address:  tt.fields.address,
				logger:   logrus.WithField("foo", "bar"),
				lock:     &sync.Mutex{},
				filename: tt.filename,
			}
			if err := p.restoreDump(); (err != nil) != tt.wantErr {
				t.Errorf("restoreDump() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(p.state, tt.wantState) {
				t.Errorf("restoreDump() state = %#b, want %#b", p.state, tt.wantState)
			}
			_ = os.Remove(tt.filename)
		})
	}
}

func TestPCA9671_storeDump(t *testing.T) {
	type fields struct {
		state    [2]byte
		device   I2CWriter
		bus      i2c.BusCloser
		address  uint16
		logger   *logrus.Entry
		lock     sync.Locker
		filename string
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		wantFile map[int]bool
	}{
		{
			name: "valid data",
			fields: fields{
				state:  [2]byte{0b01010111, 0b11010101},
				logger: logrus.WithField("foo", "bar"),
				lock:   &sync.Mutex{},
			},
			wantErr: false,
			wantFile: map[int]bool{0: true, 1: true, 2: true, 3: false, 4: true, 5: false, 6: true, 7: false,
				10: true, 11: false, 12: true, 13: false, 14: true, 15: false, 16: true, 17: true},
		},
		{
			name: "invalid file path",
			fields: fields{
				state:    [2]byte{0b01010111, 0b11010101},
				logger:   logrus.WithField("foo", "bar"),
				lock:     &sync.Mutex{},
				filename: "/dev/no/such/file",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				state:    tt.fields.state,
				device:   tt.fields.device,
				bus:      tt.fields.bus,
				address:  tt.fields.address,
				logger:   tt.fields.logger,
				lock:     tt.fields.lock,
				filename: tt.fields.filename,
			}

			if p.filename == "" {
				tmpfile, err := ioutil.TempFile("", "example")
				if err != nil {
					t.Error(err)
				}
				p.filename = tmpfile.Name()
			}
			if err := p.storeDump(); (err != nil) != tt.wantErr {
				t.Errorf("storeDump() error = %v, wantErr %v", err, tt.wantErr)
			}
			contents, err := ioutil.ReadFile(p.filename)
			if err != nil {
				t.Logf("failed to read tempfile %s: %v", p.filename, err)
			}
			var data map[int]bool
			_ = json.Unmarshal(contents, &data)
			if !reflect.DeepEqual(data, tt.wantFile) {
				t.Errorf("storeDump() file contents = \n\t%v, want \n\t%v", data, tt.wantFile)
			}
			_ = os.Remove(p.filename)
		})
	}
}
