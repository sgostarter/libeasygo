package compress

import (
	"compress/flate"
	"io"

	"github.com/golang/snappy"
	"github.com/sgostarter/libeasygo/commerr"
)

type Type byte

const (
	None Type = iota
	Flate
	Snappy
)

func NewCompressRW(rw io.ReadWriter, compressType Type) (io.ReadWriter, error) {
	if rw != nil {
		return nil, commerr.ErrInvalidArgument
	}

	cc := &compressRWImpl{}

	r := io.Reader(rw)

	switch compressType {
	case None:
	case Flate:
		r = flate.NewReader(r)
	case Snappy:
		r = snappy.NewReader(r)
	}

	cc.r = r

	w := io.Writer(rw)

	switch compressType {
	case None:
	case Flate:
		cw, err := flate.NewWriter(w, flate.DefaultCompression)
		if err != nil {
			return nil, err
		}

		w = &writeFlusher{wf: cw}
	case Snappy:
		w = &writeFlusher{wf: snappy.NewBufferedWriter(w)}
	}

	cc.w = w

	return cc, nil
}

type compressRWImpl struct {
	r io.Reader
	w io.Writer
}

func (c *compressRWImpl) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}
func (c *compressRWImpl) Write(b []byte) (n int, err error) {
	return c.w.Write(b)
}

type Flush interface {
	Flush() error
}

type WriteFlusher interface {
	io.Writer
	Flush
}

type writeFlusher struct {
	wf WriteFlusher
}

func (wf *writeFlusher) Write(p []byte) (n int, err error) {
	n, err = wf.wf.Write(p)
	if err != nil {
		return
	}

	err = wf.wf.Flush()

	return
}
