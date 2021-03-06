package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type TransactionType int

const (
	CreditTransaction TransactionType = iota // Incoming transaction
	DebitTransaction                         // Outgoing transaction
)

const DateFormat = "2006-01-02"

func Contains(list []string, key string) bool {
	for i := range list {
		if list[i] == key {
			return true
		}
	}
	return false
}

func ApplyTemplate(name, tpl string, data interface{}) string {
	t, err := template.New(name).Parse(tpl)
	if err != nil {
		logrus.Fatalf("error while parsing %s template: %s", name, err)
	}
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		logrus.Fatalf("error while apling data to %s template: %s", name, err)
	}
	return b.String()
}

func DateRangeFromYear(year int) (from, to time.Time) {
	if year < 1 {
		logrus.Fatalf("acc doesn't support years before the common era (given: \"%d\")", year)
	}
	from = time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)
	return from, to
}

// CompareFloats rounds both numbers to their third decimal place and compares them.
func CompareFloats(a float64, b float64) bool {
	return math.Floor(a*1000)/1000 == math.Floor(b*1000)/1000
}

// EscapedSplit separates a string with a given separator while ignoring separators which are escaped with a backslash (ex.: "\:" is ignored when splitting by ":" ).
func EscapedSplit(input, sep string) []string {
	const esc = "ESCAPE"
	if strings.Contains(input, esc) {
		logrus.Fatalf("match string may not contain \"%s\"", esc)
	}
	input = strings.Replace(input, fmt.Sprintf("\\%s", sep), esc, -1)
	ele := strings.Split(input, sep)
	for i := range ele {
		ele[i] = strings.Replace(ele[i], esc, sep, -1)
	}
	return ele
}

// OpenYaml opens a file and tries to marshal its content to the given interface.
// The elementType parameter is used in error messages.
func OpenYaml(data interface{}, path, dataType string) {
	if path == "" {
		logrus.Fatalf("error reading %s file: given path is empty", dataType)
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatalf("error reading %s file \"%s\": %s", dataType, path, err)
	}
	if err := yaml.Unmarshal(raw, data); err != nil {
		logrus.Fatalf("error converting (unmarshalling) %s data of file \"%s\" %s", dataType, path, err)
	}
}

// SaveToYaml writes the element (utils) as a json file to the given path.
// The elementType parameter is used in error messages.
func SaveToYaml(data interface{}, path, dataType string) {
	var raw []byte
	var err error
	raw, err = yaml.Marshal(data)
	if err != nil {
		logrus.Fatalf("error converting (marshalling) %s data for file \"%s\" to YAML: %s", dataType, path, err)
	}
	if err := ioutil.WriteFile(path, raw, 0644); err != nil {
		logrus.Fatalf("error writing %s file \"%s\" %s", dataType, path, err)
	}
}
