/*
Cyclic buffer for fragmented mp4 stream in order to save memory when pre-recording an event
*/
package mp4ring

import (
	"container/ring"
	"reflect"
	"testing"
)

func Test_getHeaderBoxInfo(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name          string
		args          args
		wantBoxHeader boxHeader
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBoxHeader, err := getHeaderBoxInfo(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHeaderBoxInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBoxHeader, tt.wantBoxHeader) {
				t.Errorf("getHeaderBoxInfo() = %v, want %v", gotBoxHeader, tt.wantBoxHeader)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    *Buffer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_Close(t *testing.T) {
	type fields struct {
		ftyp     []byte
		moov     []byte
		r        *ring.Ring
		size     int64
		isClosed bool
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
			b := &Buffer{
				ftyp:     tt.fields.ftyp,
				moov:     tt.fields.moov,
				r:        tt.fields.r,
				size:     tt.fields.size,
				isClosed: tt.fields.isClosed,
			}
			if err := b.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Buffer.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuffer_Write(t *testing.T) {
	type fields struct {
		ftyp     []byte
		moov     []byte
		r        *ring.Ring
		size     int64
		isClosed bool
	}
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				ftyp:     tt.fields.ftyp,
				moov:     tt.fields.moov,
				r:        tt.fields.r,
				size:     tt.fields.size,
				isClosed: tt.fields.isClosed,
			}
			got, err := b.Write(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Buffer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Buffer.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_Bytes(t *testing.T) {
	type fields struct {
		ftyp     []byte
		moov     []byte
		r        *ring.Ring
		size     int64
		isClosed bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				ftyp:     tt.fields.ftyp,
				moov:     tt.fields.moov,
				r:        tt.fields.r,
				size:     tt.fields.size,
				isClosed: tt.fields.isClosed,
			}
			if got := b.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Buffer.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_Size(t *testing.T) {
	type fields struct {
		ftyp     []byte
		moov     []byte
		r        *ring.Ring
		size     int64
		isClosed bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				ftyp:     tt.fields.ftyp,
				moov:     tt.fields.moov,
				r:        tt.fields.r,
				size:     tt.fields.size,
				isClosed: tt.fields.isClosed,
			}
			if got := b.Size(); got != tt.want {
				t.Errorf("Buffer.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
