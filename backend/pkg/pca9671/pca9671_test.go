package pca9671

import (
	"reflect"
	"testing"
)

func TestNewPCA9671(t *testing.T) {
	type args struct {
		address [2]byte
		bus     int
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
			got, err := NewPCA9671(tt.args.address, tt.args.bus)
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

func TestPCA9671_Check(t *testing.T) {
	type fields struct {
		Address [2]byte
		Bus     int
		state   [4]byte
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
				Address: tt.fields.Address,
				Bus:     tt.fields.Bus,
				state:   tt.fields.state,
			}
			if err := p.Check(); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPCA9671_Get(t *testing.T) {
	type fields struct {
		Address [2]byte
		Bus     int
		state   [4]byte
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				Address: tt.fields.Address,
				Bus:     tt.fields.Bus,
				state:   tt.fields.state,
			}
			got, err := p.Get(tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPCA9671_GetAll(t *testing.T) {
	type fields struct {
		Address [2]byte
		Bus     int
		state   [4]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[int]bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				Address: tt.fields.Address,
				Bus:     tt.fields.Bus,
				state:   tt.fields.state,
			}
			got, err := p.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPCA9671_Set(t *testing.T) {
	type fields struct {
		Address [2]byte
		Bus     int
		state   [4]byte
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				Address: tt.fields.Address,
				Bus:     tt.fields.Bus,
				state:   tt.fields.state,
			}
			if err := p.Set(tt.args.port, tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPCA9671_SetAll(t *testing.T) {
	type fields struct {
		Address [2]byte
		Bus     int
		state   [4]byte
	}
	type args struct {
		state map[int]bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PCA9671{
				Address: tt.fields.Address,
				Bus:     tt.fields.Bus,
				state:   tt.fields.state,
			}
			if err := p.SetAll(tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("SetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
