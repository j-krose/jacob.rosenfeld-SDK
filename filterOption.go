package jrosenfeldLotrSdk

import (
	"fmt"

	"github.com/j-krose/jrosenfeldLotrSdk/rest"
)

var _ rest.UrlParameter = (*FilterOption)(nil)

type FilterOption struct {
	field   string
	values  []string
	include bool
}

func Matches(field string, value string) rest.UrlParameter {
	return FilterOption{field, []string{value}, true}
}

func DoesntMatch(field string, value string) rest.UrlParameter {
	return FilterOption{field, []string{value}, false}
}

func Includes(field string, values []string) rest.UrlParameter {
	return FilterOption{field, values, true}
}

func Excludes(field string, values []string) rest.UrlParameter {
	return FilterOption{field, values, false}
}

func (fo FilterOption) GetUrlParameter() string {
	// Negations should be phrased as "?name!=Frodo"
	field := fo.field
	if !fo.include {
		field = (field + "!")
	}

	// Organize the list of values into a comma separated string with no spaces
	if len(fo.values) == 0 {
		// Ideally we would do something better than log here, but it will do for now
		fmt.Println("Unexpected blank value for filter option")
	}
	commaSeparatedValues := ""
	first := true
	for _, value := range fo.values {
		if !first {
			commaSeparatedValues += ","
		} else {
			first = false
		}
		commaSeparatedValues += value
	}

	return rest.BuildUrlParameter(field, commaSeparatedValues)
}
