/*
Cyclic buffer for fragmented mp4 stream in order to save memory when pre-recording an event
*/
package mp4ring

import (
	"container/ring"
	"reflect"
	"testing"
)

// examples of headers
const (
	ft = `    ftypisom    isomiso2avc1mp41`
	mv = `    moov   lmvhd              _�                                                    @                                  �trak   \tkhd                                                                        @        p      mdia    mdhd              _�    U�     -hdlr        vide            VideoHandler    *minf    vmhd               $dinf    dref            url        �stbl   �stsd           �avc1                           p H   H                                          ��   8avcC � (��  g� (�ٲ�� O��-@@@P         (� '   h��D�    stts            stsc            stsz                stco           (mvex    trex                        `
	r  = `   lmoof    mfhd       4   Ttraf    tfhd   8       �   �        tfdt          �|    trun           t   �   �`
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
		{name: "ftyp", args: args{[]byte(ft + mv + r)}, wantBoxHeader: boxHeader{Size: 538976288, FourccType: [4]byte{102, 116, 121, 112}}, wantErr: false},
		{name: "moov", args: args{[]byte(mv + ft + r)}, wantBoxHeader: boxHeader{Size: 538976288, FourccType: [4]byte{109, 111, 111, 118}}, wantErr: false},
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
	var (
		b   *Buffer
		err error
	)

	b, err = New(10)

	if b == nil || err != nil {
		t.Fatal(`failed to create ring buffer`)
	}

	b, err = New(-12) //nolint:ineffassign,staticcheck

	if err == nil {
		t.Fatal(`failed check size buffer`)
	}
}

func TestBuffer_Close(t *testing.T) {
	tests := []struct {
		name    string
		buf     *Buffer
		wantErr bool
	}{
		{name: "ok", buf: &Buffer{ftyp: []byte(ft), moov: []byte(mv), r: ring.New(1), size: 100, isClosed: false}, wantErr: false},
		{name: "nil", buf: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.buf.Close(); (err != nil) != tt.wantErr {
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
