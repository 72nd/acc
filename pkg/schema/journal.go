package schema

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
)

const HLedgerDateFormat = "2006-01-02"
const defaultAccount = "other:unknown"

type Journal []Entry

func JournalFromStatement(a Acc, update bool) Journal {
	var result Journal
	for i := range a.Expenses {
		result = append(result, a.Expenses[i].Journal(a)...)
	}
	for i := range a.Invoices {
		result = append(result, a.Invoices[i].Journal(a)...)
	}
	for i := range a.BankStatement.Transactions {
		result = append(result, a.BankStatement.Transactions[i].Journal(a, update)...)
	}
	return result
}

func (j Journal) SaveHLedgerFile(path string) {
	ledger := j.HLedger()
	if err := ioutil.WriteFile(path, []byte(ledger), 0644); err != nil {
		logrus.Fatalf("error writing file %s: %s", path, err)
	}
}

func (j Journal) HLedger() string {
	var result string
	sort.Sort(j)
	for i := range j {
		result = fmt.Sprintf("%s\n\n%s", result, j[i].Transaction())
	}
	return result
}

func (j Journal) Len() int {
	return len(j)
}

func (j Journal) Swap(i, k int) {
	j[i], j[k] = j[k], j[i]
}

func (j Journal) Less(i, k int) bool {
	return j[i].Date.Before(j[k].Date)
}

type Comment struct {
	Mode    string
	Element string
	DoManual bool
	Errors  []error
}

func NewComment(mode, element string) Comment {
	return Comment{
		Mode:    mode,
		Element: element,
		DoManual: true,
		Errors:  []error{},
	}
}

func NewManualComment(mode, element string) Comment {
	cmt := NewComment(mode, element)
	cmt.DoManual = false
	logrus.Warnf("journal entry of «%s» needs manual correction", element)
	return cmt
}

func (c *Comment) add(err error) {
	if err == nil {
		return
	}
	logrus.Warnf("error while converting «%s» to journal entries: %s", c.Element, err)
	c.Errors = append(c.Errors, err)
}

func (c Comment) String() string {
	if c.DoManual {
		return "TODO: manual correction needed"
	}
	if len(c.Errors) == 0 {
		return fmt.Sprint("parsed as ", c.Mode)
	}
	result := "TODO:"
	for i := range c.Errors {
		sep := ", "
		if i == 0 {
			sep = " "
		}
		result = fmt.Sprintf("%s%s%s", result, sep, c.Errors[i].Error())
	}
	return result
}

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
	TransactionType util.TransactionType
	Date            time.Time
	Status          EntryStatus
	Description     string
	Comment         Comment
	Account1        string
	Account2        string
	Amount          float64
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
		Comment:     e.Comment.String(),
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