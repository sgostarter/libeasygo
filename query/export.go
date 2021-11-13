package query

import "net/url"

func ParseQuery(query string) (url.Values, error) {
	return url.ParseQuery(query)
}

func EncodeValues(values url.Values) string {
	return values.Encode()
}

func Values2Orm(values url.Values, s interface{}) error {
	return Unmarshal(values, s)
}

func Orm2Values(s interface{}, values url.Values) error {
	return Marshal(s, values)
}
