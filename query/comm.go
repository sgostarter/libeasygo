package query

import (
	"net/url"
	"strings"

	"github.com/iancoleman/strcase"
)

func selectTagByURLValues(values url.Values, name string) (tag string) {
	tag = strcase.ToSnake(name)
	if values.Get(tag) != "" {
		return
	}

	tag = name
	if values.Get(tag) != "" {
		return
	}

	tag = strings.ToLower(name)
	if values.Get(tag) != "" {
		return
	}

	return
}
