package report

import (
	"regexp"
)
type Problem struct {
	Severity int
	Description string
	Name string
	Re *regexp.Regexp
}

func (b Problem) String() string {
	return b.Name
}