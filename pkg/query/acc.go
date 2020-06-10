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

func accElementsFromQueryable(a schema.Acc, q Queryable) []Element {
	switch q.Name {
	case "customer":
		return NewElements(a.Parties.Customers)
	case "employee":
		return NewElements(a.Parties.Employees)
	case "expense":
		return NewElements(a.Expenses)
	case "invoice":
		return NewElements(a.Invoices)
	case "transaction":
		return NewElements(a.BankStatement.Transactions)
	default:
		logrus.Fatalf("no acc elements for \"%s\" found", q.Name)
	}
	return []Element{}
}

func (q Queryables) QueryAcc(a schema.Acc, termsInput, dateInput, selectInput string, mode OutputMode, noRender, caseSensitive bool) {
	var ele ElementGroup
	for i := range q {
		ele = append(ele, accElementsFromQueryable(a, q[i])...)
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
	OutputsFromElements(a, ele).PPKeyValue(a, mode)
}
