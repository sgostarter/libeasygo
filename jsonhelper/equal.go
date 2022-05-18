package jsonhelper

import (
	"encoding/json"
	"reflect"
)

func Equal(d1, d2 []byte) bool {
	ok, _ := EqualE(d1, d2)

	return ok
}

func EqualE(d1, d2 []byte) (ok bool, err error) {
	var o1, o2 interface{}

	err = json.Unmarshal(d1, &o1)
	if err != nil {
		return
	}

	err = json.Unmarshal(d2, &o2)
	if err != nil {
		return
	}

	ok = reflect.DeepEqual(o1, o2)

	return
}
