package query

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrStruct = errors.New("unmarshal() expects struct input")
)

func Unmarshal(values url.Values, s interface{}) error {
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ErrStruct
		}

		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrStruct
	}

	return reflectValueFromTag(values, val)
}

// nolint: gocognit
func reflectSliceValueFromTag(tag string, val reflect.Value, vs []string) error {
	if val.IsNil() {
		val.Set(reflect.MakeSlice(val.Type(), 0, len(vs)))
	}

	switch t := val.Type().Elem().Kind(); t {
	case reflect.String:
		for _, v := range vs {
			val.Set(reflect.Append(val, reflect.ValueOf(v)))
		}
	case reflect.Bool:
		for _, v := range vs {
			f, err := strconv.ParseBool(v)
			if err != nil {
				return newParseError(err, "cast bool", tag, t.String(), v, tag)
			}

			val.Set(reflect.Append(val, reflect.ValueOf(f)))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		for _, v := range vs {
			rv := reflect.New(val.Type().Elem()).Elem()

			if v == "" {
				v = "0"
			}

			n, err := strconv.ParseUint(v, 0, 64)

			if err != nil && rv.OverflowUint(n) {
				return newParseError(err, "cast uint", tag, t.String(), v, tag)
			}

			rv.SetUint(n)
			val.Set(reflect.Append(val, rv))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for _, v := range vs {
			rv := reflect.New(val.Type().Elem()).Elem()

			if v == "" {
				v = "0"
			}

			n, err := strconv.ParseInt(v, 0, 64)

			if err != nil && rv.OverflowInt(n) {
				return newParseError(err, "cast int", tag, t.String(), v, tag)
			}

			rv.SetInt(n)

			val.Set(reflect.Append(val, rv))
		}
	case reflect.Float32, reflect.Float64:
		for _, v := range vs {
			rv := reflect.New(val.Type().Elem()).Elem()

			if v == "" {
				v = "0"
			}

			n, err := strconv.ParseFloat(v, 64)

			if err != nil && rv.OverflowFloat(n) {
				return newParseError(err, "cast float", tag, t.String(), v, tag)
			}

			rv.SetFloat(n)

			val.Set(reflect.Append(val, rv))
		}
	}

	return nil
}

// nolint: gocognit
func reflectValueFromTag(values url.Values, val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		kt := typ.Field(i)
		tag := kt.Tag.Get("query")

		if tag == "-" {
			continue
		}

		if tag == "" {
			tag = selectTagByURLValues(values, kt.Name)
		}

		sv := val.Field(i)
		vs, fv := getVals(values, tag)

		switch sv.Kind() {
		case reflect.String:
			sv.SetString(fv)
		case reflect.Bool:
			b, err := strconv.ParseBool(fv)
			if err != nil {
				return newParseError(err, "cast bool", tag, sv.Type().String(), fv, tag)
			}

			sv.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if fv == "" {
				fv = "0"
			}

			n, err := strconv.ParseUint(fv, 0, 64)

			if err != nil || sv.OverflowUint(n) {
				return newParseError(err, "cast uint", tag, sv.Type().String(), fv, tag)
			}

			sv.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fv == "" {
				fv = "0"
			}

			n, err := strconv.ParseInt(fv, 0, 64)

			if err != nil || sv.OverflowInt(n) {
				return newParseError(err, "cast int", tag, sv.Type().String(), fv, tag)
			}

			sv.SetInt(n)
		case reflect.Float32, reflect.Float64:
			if fv == "" {
				fv = "0"
			}

			n, err := strconv.ParseFloat(fv, sv.Type().Bits())

			if err != nil || sv.OverflowFloat(n) {
				return newParseError(err, "cast float", tag, sv.Type().String(), fv, tag)
			}

			sv.SetFloat(n)
		case reflect.Struct:
			err := reflectValueFromTag(values, sv)
			if err != nil {
				return newParseError(err, "cast struct", tag, sv.Type().String(), fv, tag)
			}
		case reflect.Ptr:
			if sv.IsNil() {
				sv.Set(reflect.New(sv.Type().Elem()))
				sv = sv.Elem()
			}

			err := reflectValueFromTag(values, sv)

			if err != nil {
				return err
			}
		case reflect.Slice:
			err := reflectSliceValueFromTag(tag, sv, vs)
			if err != nil {
				return err
			}
		default:
		}
	}

	return nil
}

//get val, if absent get from tag default val
func getVals(values url.Values, tag string) (vs []string, fv string) {
	name, opts := parseTag(tag)
	vs = values[name]

	if len(vs) > 0 {
		fv = vs[0]
	}

	if fv == "" && len(vs) <= 1 && len(opts) > 0 {
		vs = opts
		if len(vs) > 0 {
			fv = vs[0]
		}
	}

	return
}

type tagOptions []string

func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")

	return s[0], s[1:]
}
