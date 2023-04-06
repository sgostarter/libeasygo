package mwf

import "encoding/json"

type JSONSerial struct {
}

func (serial *JSONSerial) Marshal(t any) ([]byte, error) {
	return json.Marshal(t)
}

func (serial *JSONSerial) Unmarshal(d []byte, t any) error {
	return json.Unmarshal(d, t)
}
