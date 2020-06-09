package ledger

// The ledger package delivers the functionality to generate hledger journals out of a given
// schema.Acc struct.
//
// Design Rationale
//
// Initially the mechanisms for generating journals was part of the schema package. But
// the complexity of generating the transactions made the code quit hard to understand.
// As this generation doesn't alter any data of a given acc project it was decided to move
// the functionality into it's own package.

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
)

// HLedgerDateFormat defines the default date format as required by hledger.
const HLedgerDateFormat = "2006-01-02"

// defaultAccount is used when no account can be chosen and the user has to manually complete the journal entry.
const defaultAccount = "other:unknown"

// Journal is a data structure which can be converted into a hledger journal.
//
// On Aliases
//
// Some standard chart of accounts uses multiple root accounts for the same account type.
// An example is Switzerland where "Personalaufwand" and "Betriebsfremder Aufwand" are
// both root accounts. In hledger both of them have to be children of the "expense" account.
// To accomplish this the use of aliases is possible. Alias declaration for the given example:
//
//	aliases := [][]string{
//		[]string{"Personalaufwand", "expenses"},
//		[]string{"Betriebsfremder Aufwand", "expenses"}}
type Journal struct {
	Aliases [][]string
	Entries []Entry
}

// NewJournalConfig returns a new Journal with the given aliases.
func NewJournal(aliases [][]string) Journal {
	return Journal{
		Aliases: aliases,
	}
}

// AddEntries adds new entries to the journal.
func (j *Journal) AddEntries(entries []Entry) {
	j.Entries = append(j.Entries, entries...)
}

// JournalFromAcc takes an schema.Acc project and converts it into an Journal. This is
// mainly used to export the Journal afterwards into a hledger journal. Optionally the
// year can be filtered, if the given year parameter is > 0, only events happened in
// this year will be converted into transactions.
func JournalFromAcc(a schema.Acc, year int) Journal {
	rsl := NewJournal(parseAliases(a.JournalConfig.AccountAliases))
	a = a.FilterYear(year)

	for i := range a.Expenses {
		rsl.AddEntries(EntriesForExpense(a, a.Expenses[i]))
	}
	for i := range a.Invoices {
		rsl.AddEntries(EntriesForInvoicing(a, a.Invoices[i]))
	}
	/*
		for i := range stn.Transactions {
			result.Entries = append(result.Entries, stn.Transactions[i].JournalEntries(a, update)...)
		}
	*/
	return rsl
}

func parseAliases(input []string) [][]string {
	result := make([][]string, len(input))
	for i := range input {
		ele := strings.Split(input[i], ":")
		if len(ele) != 2 {
			logrus.Fatalf("error while parsing account aliases \"%s\" couldn't be parsed as ALIAS:REPLACE", input[i])
		}
		result[i] = []string{ele[0], ele[1]}
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
	result := j.HLedgerHeader()
	sort.Sort(j)
	for i := range j.Entries {
		result = fmt.Sprintf("%s\n\n%s", result, j.Entries[i].Transaction())
	}
	return result
}

func (j Journal) HLedgerHeader() string {
	var result string
	first := true
	for i := range j.Aliases {
		if first {
			result = fmt.Sprintf("alias %s = %s", j.Aliases[i][0], j.Aliases[i][1])
			first = false
			continue
		}
		result = fmt.Sprintf("%s\nalias %s = %s", result, j.Aliases[i][0], j.Aliases[i][1])
	}
	return fmt.Sprintf("\n%s", result)
}

func (j Journal) Len() int {
	return len(j.Entries)
}

func (j Journal) Swap(i, k int) {
	j.Entries[i], j.Entries[k] = j.Entries[k], j.Entries[i]
}

func (j Journal) Less(i, k int) bool {
	return j.Entries[i].Date.Before(j.Entries[k].Date)
}

type Comment struct {
	Mode     string
	Element  string
	DoManual bool
	Errors   []error
}

func NewComment(mode, element string) Comment {
	return Comment{
		Mode:     mode,
		Element:  element,
		DoManual: false,
		Errors:   []error{},
	}
}

func NewManualComment(mode, element string) Comment {
	cmt := NewComment(mode, element)
	cmt.DoManual = true
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
	Code            string
	Description     string
	Comment         Comment
	Account1        string
	Account2        string
	Amount          float64
}

const trnTpl = `
{{.Date}} {{if .Code }}({{.Code}}) {{end}}{{.Description}} {{if ne .Comment ""}}; {{.Comment}}{{end}}
    {{.Account1}}{{.Space1}}{{.Amount1}}
    {{.Account2}}{{.Space2}}{{.Amount2}}
`

func (e Entry) Transaction() string {
	data := struct {
		Date        string
		Code        string
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
		Code:        e.Code,
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

func compareAmounts(a float64, b float64) error {
	if util.CompareFloats(a, b) {
		return nil
	}
	return fmt.Errorf("the two involved amounts don't match: %.3f vs %.3f", a, b)
}
