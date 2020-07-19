package query

import (
	"github.com/sirupsen/logrus"
	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
)

var AccQueryables = Queryables{
	{
		Name: "customer",
		Type: schema.Party{},
	},
	{
		Name: "employee",
		Type: schema.Party{},
	},
	{
		Name: "expense",
		Type: schema.Expense{},
	},
	{
		Name: "invoice",
		Type: schema.Invoice{},
	},
	{
		Name: "misc-record",
		Type: schema.MiscRecord{},
	},
	{
		Name: "project",
		Type: schema.Project{},
	},
	{
		Name: "transaction",
		Type: schema.Transaction{},
	},
}

func accElementsFromQueryable(s schema.Schema, q Queryable) []Element {
	switch q.Name {
	case "customer":
		return NewElements(s.Parties.Customers)
	case "employee":
		return NewElements(s.Parties.Employees)
	case "expense":
		return NewElements(s.Expenses)
	case "invoice":
		return NewElements(s.Invoices)
	case "misc-record":
		return NewElements(s.MiscRecords)
	case "project":
		return NewElements(s.Projects)
	case "transaction":
		return NewElements(s.Statement.Transactions)
	default:
		logrus.Fatalf("no acc elements for \"%s\" found", q.Name)
	}
	return []Element{}
}

func (q Queryables) QueryAcc(s schema.Schema, termsInput, dateInput, selectInput string, mode OutputMode, render, caseSensitive, openAttachment bool) {
	var ele ElementGroup
	for i := range q {
		ele = append(ele, accElementsFromQueryable(s, q[i])...)
	}
	if termsInput != "" {
		terms := searchTermsFromUserInput(termsInput, caseSensitive)
		ele = ele.MatchTerm(terms, caseSensitive)
	}
	if dateInput != "" {
		ranges := dateTermsFromUserInput(dateInput)
		ele = ele.DateMatch(ranges)
	}
	if selectInput != "" {
		sel := util.EscapedSplit(selectInput, ",")
		ele = ele.Select(sel, caseSensitive)
	}
	OutputsFromElements(s, ele).PPKeyValue(s, mode, render, openAttachment)
}
