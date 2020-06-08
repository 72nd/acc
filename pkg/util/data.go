package util

import (
	"bytes"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
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
		logrus.Fatalf("acc doesn't support years before the common era (given: \"%s\")", year)
	}
	from = time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)
	return from, to
}
