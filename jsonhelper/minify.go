package jsonhelper

import (
	"bytes"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"
)

func Minify(d []byte) (md []byte, err error) {
	w := &bytes.Buffer{}

	err = json.Minify(minify.New(), w, bytes.NewReader(d), nil)
	if err != nil {
		return
	}

	md = w.Bytes()

	return
}
