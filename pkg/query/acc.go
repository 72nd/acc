package query

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
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
	case "transaction":
		return NewElements(s.Statement.Transactions)
	default:
		logrus.Fatalf("no acc elements for \"%s\" found", q.Name)
	}
	return []Element{}
}

func (q Queryables) QueryAcc(s schema.Schema, termsInput, dateInput, selectInput string, mode OutputMode, render, caseSensitive bool) {
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
	OutputsFromElements(s, ele).PPKeyValue(s, mode, render)
}
