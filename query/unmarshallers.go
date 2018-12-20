package query

import (
	"strconv"
	"strings"
)

// IntList implements encoding.TextUnmarshaler in order to extract query
// of form ?key=1,2,3,4 onto an int slice
type IntList []int

// UnmarshalText implements encoding.TextUnmarshaler
func (l *IntList) UnmarshalText(text []byte) error {
	s := string(text)
	items := strings.Split(s, ",")
	for _, item := range items {
		i, err := strconv.Atoi(item)
		if err != nil {
			return err
		}
		*l = append(*l, i)
	}
	return nil
}

// StringList implements encoding.TextUnmarshaler in order to extract query
// of form ?key=foo,bar onto a string slice
type StringList []string

// UnmarshalText implements encoding.TextUnmarshaler
func (l *StringList) UnmarshalText(text []byte) error {
	*l = StringList(strings.Split(string(text), ","))
	return nil
}
