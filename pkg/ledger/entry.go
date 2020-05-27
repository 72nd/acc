package ledger

import (
	"bytes"
	"fmt"
	"strings"
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

func (s EntryStatus) TrnEle() string {
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
	Comment     string
	Account1    string
	Account2    string
	Amount      float64
}

const trnTpl = `
{{.Date}} {{.Description}} {{if ne .Comment ""}}; {{.Comment}}{{end}}
    {{.Account1}}{{.Space1}}{{.Amount1}}
    {{.Account2}}{{.Space2}}{{.Amount2}}
`

func (e Entry) Transaction() string {
	data := struct {
		Date        string
		Description string
		Comment     string
		Account1    string
		Space1      string
		Amount1     string
		Account2    string
		Space2      string
		Amount2     string
	}{
		Date:        e.trnDate(),
		Description: e.Description,
		Comment:     e.Comment,
		Account1:    e.Account1,
		Space1:      e.trnSpace(e.Account1),
		Amount1:     e.trnAmount(false),
		Account2:    e.Account2,
		Space2:      e.trnSpace(e.Account2),
		Amount2:     e.trnAmount(true),
	}
	tpl, err := template.New("transaction").Parse(trnTpl)
	if err != nil {
		logrus.Fatal("error while parsing transaction template: ", err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		logrus.Fatal("error while applying data to transaction template: ", err)
	}
	return buf.String()
}

func (e Entry) trnDate() string {
	return e.Date.Format(HLedgerDateFormat)
}

func (e Entry) trnSpace(account string) string {
	var max int
	if len(e.Account1) > len(e.Account2) {
		max = len(e.Account1)
	} else {
		max = len(e.Account2)
	}

	spaces := 8
	if len(account) < max {
		spaces += max - len(account) - 2
	}
	return strings.Repeat(" ", spaces)
}

func (e Entry) trnAmount(invers bool) string {
	sign := ""
	if invers {
		sign = "-"
	}
	whole := int64(e.Amount)
	if e.Amount == float64(whole) {
		return fmt.Sprintf("CHF%s%d", sign, whole)
	}
	return fmt.Sprintf("CHF%s%.2f", sign, e.Amount)
}
