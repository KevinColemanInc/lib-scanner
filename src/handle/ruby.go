package handle

import (
	"io/ioutil"
	"regexp"
	"github.com/KevinColemanInc/lib-crawl/src/report"

)

func RubyScan(done <-chan struct{}, path string, c chan<- report.Warning) {
	// Check file names
	for _, problem := range regexPathProblems {
		matches := problem.Re.FindAllStringSubmatch(path, -1)
		for _, match := range matches {
			for _, line := range match {
				select {
				case c <- report.Warning{GemName: path, Filepath: path, ProblemType: problem, Line: line}:
				case <-done:
					return
				}

			}
		}
	}

	// Check file content
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		select {
		case c <- report.Warning{GemName: path, Filepath: path, Err: err}:
		case <-done:
			return
		}
	}
	for _, problem := range regexProblems {
		matches := problem.Re.FindAllStringSubmatch(string(dat), -1)
		for _, match := range matches {
			select {
			case c <- report.Warning{GemName: path, Filepath: path, ProblemType: problem, Line: match[0]}:
			case <-done:
				return
			}

		}
	}
}

var regexProblems = []report.Problem{
	report.Problem{
		Severity:    5,
		Name:        "#send",
		Description: "Uses the send method which can be used for RCE",
		Re:          regexp.MustCompile(`^([^#\n]*;? *send[\( ].*)`),
	},
	report.Problem{
		Severity:    5,
		Name:        "#eval",
		Description: "Uses the send method which can be used for RCE",
		Re:          regexp.MustCompile(`^([^#\n]*;? *eval[\( ].*)`),
	},
	report.Problem{
		Severity:    1,
		Name:        "has http",
		Description: "HTTP may signal that they are trying to share send a network request",
		Re:          regexp.MustCompile(`(?i)^([^#\n]*;? *http.*)`),
	},
	report.Problem{
		Severity:    1,
		Name:        "has tcp",
		Description: "tcp may signal that they are trying to share send a network request",
		Re:          regexp.MustCompile(`(?i)^([^#\n]*;? *tcp.*)`),
	},
	report.Problem{
		Severity:    1,
		Name:        "has udp",
		Description: "udp may signal that they are trying to share send a network request",
		Re:          regexp.MustCompile(`(?i)^([^#\n]*;? *udp.*)`),
	},
	report.Problem{
		Severity:    5,
		Name:        "#exec",
		Description: "Uses the exec method which can be used for RCE",
		Re:          regexp.MustCompile(`^([^#\n]*;? *exec[\( ].*)`),
	},
	report.Problem{
		Severity:    5,
		Name:        "#system",
		Description: "Uses the system method which can be used for RCE",
		Re:          regexp.MustCompile(`^([^#\n]*;? *system[\( ].*)`),
	},
}

var regexPathProblems = []report.Problem{
	report.Problem{
		Severity:    5,
		Name:        "#ext_config",
		Description: "Runs code on install; Can be used to run or install malware",
		Re:          regexp.MustCompile(`extconf.rb`),
	},
}
