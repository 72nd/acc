package query

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
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

func (q Queryables) QueryAcc(a schema.Acc, termsInput, dateInput string, mode OutputMode) {
	var ele ElementGroup
	for i := range q {
		ele = append(ele, accElementsFromQueryable(a, q[i])...)
	}
	if termsInput != "" {
		terms := searchTermsFromUserInput(termsInput)
		ele = ele.Match(terms)
	}
	if dateInput != "" {
		ranges := dateTermsFromUserInput(dateInput)
		ele = ele.DateMatch(ranges)
	}
	OutputsFromElements(ele).PPKeyValue(mode)
}
