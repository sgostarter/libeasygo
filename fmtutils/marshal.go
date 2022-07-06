package fmtutils

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Marshal ...
func Marshal(v interface{}) string {
	dat, err := yaml.Marshal(v)
	if err == nil {
		return string(dat)
	}

	return fmt.Sprintf("%#v", v)
}

func JSONMarshal(v interface{}) string {
	dat, err := json.Marshal(v)
	if err == nil {
		return string(dat)
	}

	return fmt.Sprintf("%#v", v)
}
