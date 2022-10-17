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
	errSize     = errors.New("size must be positive")
	errBufClose = errors.New("ring buffer is closed")
	errBufNil   = errors.New("ring buffer is nil")
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

// boxHeader - header representation
type boxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

// getHeaderBoxInfo - getting header information
func getHeaderBoxInfo(data []byte) (boxHeader boxHeader) {
	buf := bytes.NewBuffer(data)

	binary.Read(buf, binary.BigEndian, &boxHeader)

	return
}

// getFourccType - header type
func getFourccType(boxHeader boxHeader) string {
	return string(boxHeader.FourccType[:])
}

// New creates a new buffer of the specified size. The size must be greater than 0
func New(size int) (*Buffer, error) {
	if size <= 0 {
		return nil, errSize
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
		return errBufNil
	}

	b.isClosed = true

	return nil
}

// Write writes headers and up to N atoms to the inner ring, overwriting older data if necessary
func (b *Buffer) Write(buf []byte) (int, error) {
	if b.isClosed { // to unsubscribe, you need to return an error
		return 0, errBufClose
	}

	// Total number of bytes written
	var n int

	// Separately store the header, which is easy to determine - each atom has a name. And the name of the header atoms is ftyp and moov. It is enough to watch the first 4 bytes of the stream, if there is an ftp, then this is the header and you need to save it
	bHead := getHeaderBoxInfo(buf)
	fourccType := getFourccType(bHead)

	switch fourccType {
	case `ftyp`:
		n = copy(b.ftyp, buf[:bHead.Size])
	case `moov`:
		b.moov = make([]byte, len(buf))
		n = copy(b.moov, buf[:bHead.Size])
	default:
		tmp := make([]byte, len(buf))
		n = copy(tmp, buf[:bHead.Size])
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
