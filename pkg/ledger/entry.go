package ledger

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

type EntryStatus int

const (
	UnmarkedStatus EntryStatus = iota
	PendingStatus
	ClearedStatus
)

func (s EntryStatus) String() string {
	switch s {
	case UnmarkedStatus:
		return ""
	case PendingStatus:
		return "!"
	case ClearedStatus:
		return "*"
	}
	logrus.Fatalf("no journal element string found for %d", int(s))
	return "UNDEFINED"
}

type Entry struct {
	Date        time.Time
	Status      EntryStatus
	Description string
	Account1    string
	Account2    string
	Amount      float64
}

func (e Entry) Transaction() string {
	data := struct {
	} {
	}
	tpl, err := template.New("transaction").Parse(``)
	if err != nil {
		logrus.Fatal("error while parsing transaction template: ", err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		logrus.Fatal("error while applying data to transaction template: ", err)
	}
	return buf.String()
}
