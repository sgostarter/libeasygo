package frame

import (
	"io"

	"github.com/sgostarter/libeasygo/helper"
)

func NewSimpleFrame(rw io.ReadWriteCloser) Frame {
	return &simpleFrame{
		rw: rw,
	}
}

type simpleFrame struct {
	rw io.ReadWriteCloser
}

func (s *simpleFrame) ReadFrame() ([]byte, error) {
	buf := make([]byte, 4)

	_, err := io.ReadFull(s.rw, buf)
	if err != nil {
		return nil, err
	}

	buf = make([]byte, helper.BytesToInt(buf))

	_, err = io.ReadFull(s.rw, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *simpleFrame) WriteFrame(p []byte) error {
	bs := helper.IntToBytes(len(p))

	_, err := s.rw.Write(bs)
	if err != nil {
		return err
	}

	_, err = s.rw.Write(p)

	return err
}

func (s *simpleFrame) Close() error {
	return s.rw.Close()
}
