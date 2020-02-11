package document

import (
	"fmt"
	"gitlab.com/72th/acc/pkg/schema"
)

func GenerateExpenses(expenses schema.Expenses, dstPath string) {
	for i := range expenses {
		GenerateExpense(expenses[i], dstPath)
	}
}

func GenerateExpense(expense schema.Expense, dstPath string) {
	props := Properties{
		Type:       "Expense",
		Identifier: expense.Identifier,
		DstName:    expense.Identifier,
		Line1:      fmt.Sprintf("id: %s // name: %s // amount: %.2f", expense.Id, expense.Name, expense.Amount),
		Line2:      fmt.Sprintf("accrual at: %s // advanced by 3th: %t // settlement at: %s", expense.DateOfAccrual, expense.AdvancedByThirdParty, expense.DateOfAccrual),
		Line3:      fmt.Sprintf("3th party: %s // customer: %s", "not-implemented", "not-implemented"),
	}
	pdf := NewPdf(expense.Path, dstPath)
	pdf.Generate(props)
}

func GenerateInvoices(invoices schema.Invoices, dstPath string) {
	for i := range invoices {
		GenerateInvoice(invoices[i], dstPath)
	}

}

func GenerateInvoice(invoice schema.Invoice, dstPath string) {

}
