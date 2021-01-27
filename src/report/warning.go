package report

import (
	"fmt"
	"strings"
)

type Warning struct {
	Filepath    string
	ProblemType Problem
	Line        string
	GemName     string
	Err         error
}

func (w *Warning) ToCSV() string {
	return strings.Join(w.ToArray(), "≫")
}

func (w *Warning) ToCLI(level string) {
	switch level {
	case "verbose":
		fmt.Println(w.ToCSV())
	default:
		fmt.Println(w.ToArray()[:3])
	}
}

func (w *Warning) ToArray() []string {
	return []string{w.ProblemType.String(), w.GemName, w.Line}
}
