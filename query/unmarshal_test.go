package query

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AI interface {
}

type A struct {
	N          int
	N2         int
	ND         int `query:"nd,99"`
	S          string
	Ss         []string
	Bs         []bool
	Int8s      []int8
	Int64s     []int64
	AI         AI
	M          map[string]int
	MyJSONOk   float64
	MyJSANFall float32
}

type B struct {
	A *A `param:"A"`
}

func TestMarshal1(t *testing.T) {
	s := url.Values{}
	s.Add("N", "10")
	s.Add("Ss", "1")
	s.Add("Ss", "2")
	s.Add("Bs", "true")
	s.Add("int_8_s", "1")
	s.Add("int_8_s", "2")
	s.Add("Int64s", "3")
	s.Add("Int64s", "4")

	b := &B{}
	err := Unmarshal(s, b)
	assert.Nil(t, err)

	d, _ := json.Marshal(b)
	t.Log(string(d))

	vs2 := url.Values{}
	err = Marshal(b, vs2)
	assert.Nil(t, err)

	t.Log(vs2.Encode())
}
