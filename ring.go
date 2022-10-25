/*
Cyclic buffer for fragmented mp4 stream in order to save memory when pre-recording an event
*/
package mp4ring

import (
	"bytes"
	"container/ring"
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	//ErrSize - size must be positive
	ErrSize = errors.New("size must be positive")

	// ErrBufClose - ring buffer is closed
	ErrBufClose = errors.New("buffer is closed")

	// ErrBufNil - ring buffer is nil
	ErrBufNil = errors.New("buffer is nil")
)

// Buffer implements a cyclic buffer. Has a fixed size,
// and the new data overwrites the old, so that for a buffer
// of size N, for any number of write operations, only the last N mp4 atoms are saved.
// At the same time, the ftyp and moov headers are stored separately.
type Buffer struct {
	ftyp     []byte     // ftyp header
	moov     []byte     // moov header
	r        *ring.Ring // atoms
	size     int64
	isClosed bool // closing flag
}

// header representation (https://docs.fileformat.com/video/mp4/#structure-of-mp4-files)
type boxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

// getting header information
func getHeaderBoxInfo(data []byte) (boxHeader, error) {
	buf := bytes.NewBuffer(data)

	var box boxHeader

	err := binary.Read(buf, binary.BigEndian, &box)

	return box, err
}

// New creates a new buffer of the specified size. The size must be greater than 0
func New(size int) (*Buffer, error) {
	if size <= 0 {
		return nil, ErrSize
	}

	return &Buffer{
		r:        ring.New(size),
		ftyp:     make([]byte, 32),
		isClosed: false,
	}, nil
}

// Close closing the buffer
func (b *Buffer) Close() error {
	if b == nil {
		return ErrBufNil
	}

	b.isClosed = true

	return nil
}

// Write writes headers and up to N atoms to the inner ring, overwriting older data if necessary
func (b *Buffer) Write(buf []byte) (int, error) {
	if b.isClosed { // to unsubscribe, you need to return an error
		return 0, ErrBufClose
	}

	// Total number of bytes written
	var n int

	// Separately store the header, which is easy to determine - each atom has a name. And the name of the header atoms is ftyp and moov. It is enough to watch the first 4 bytes of the stream, if there is an ftp, then this is the header and you need to save it
	bHead, err := getHeaderBoxInfo(buf)
	if err != nil {
		return 0, err
	}

	// get header type
	fourccType := string(bHead.FourccType[:])

	switch fourccType {
	case `ftyp`:
		b.ftyp = make([]byte, len(buf))
		n = copy(b.ftyp, buf)
	case `moov`:
		b.moov = make([]byte, len(buf))
		n = copy(b.moov, buf)
	default:
		tmp := make([]byte, len(buf))
		n = copy(tmp, buf)
		b.r.Value = tmp
		b.r = b.r.Next()
	}

	return n, nil
}

// Bytes provides all recorded mp4 headers and atoms
func (b *Buffer) Bytes() []byte {
	var buf bytes.Buffer

	buf.Write(b.ftyp)
	buf.Write(b.moov)

	b.r.Do(func(p interface{}) {
		atom, ok := p.([]uint8)
		if ok {
			buf.Write(atom)
		}
	})

	return buf.Bytes()
}

// Size provides the size of all recorded data in kilobytes
func (b *Buffer) Size() string {
	b.size = 0

	b.size += int64(len(b.ftyp) + len(b.moov))

	b.r.Do(func(p interface{}) {
		atom, ok := p.([]uint8)
		if ok {
			b.size += int64(len(atom))
		}
	})

	return fmt.Sprintf("%.2f", float64(b.size)/(1<<10))
}
