package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

type cumulativeArg []string

func (a *cumulativeArg) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func (a *cumulativeArg) String() string {
	return strings.Join(*a, " ")
}

func (a *cumulativeArg) IsCumulative() bool {
	return true
}

// CumulativeArg is a special type of Arg which greedily consumes all args
func CumulativeArg(s kingpin.Settings) *[]string {
	target := new([]string)
	s.SetValue((*cumulativeArg)(target))
	return target
}
