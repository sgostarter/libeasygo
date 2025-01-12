package compressex

import (
	"bytes"
	"io"
	"sync"

	"compress/gzip"
)

var (
	gzipPool = sync.Pool{
		New: func() interface{} {
			z, _ := gzip.NewWriterLevel(nil, gzip.BestCompression)

			return z //return gzip.NewWriter(nil)
		},
	}

	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

func GzipWithBuildOnPool(data []byte) ([]byte, error) {
	w, _ := gzipPool.Get().(*gzip.Writer)
	defer gzipPool.Put(w)

	buf, _ := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)

	buf.Reset()
	w.Reset(buf)

	return GzipEx(w, buf, data)
}

func Gzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)

	return GzipEx(w, &buf, data)
}

func GzipEx(w *gzip.Writer, buf *bytes.Buffer, data []byte) ([]byte, error) {
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return append([]byte{}, buf.Bytes()...), nil
}

func UnGzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	result, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = r.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
