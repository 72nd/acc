package providers

import "gitlab.com/72th/acc/pkg/schema"

type AccProvider interface {
	GetExpenses() schema.Expenses
	GetInvoices() schema.Invoices
	GetStatement() schema.Statement
}
