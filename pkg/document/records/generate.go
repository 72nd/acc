package records

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"os"
	"path"
)

func GenerateExpensesRec(expenses schema.Expenses, dstFolder string, doOverwrite bool) {
	nFiles := len(expenses)
	for i := range expenses {
		fileName := fmt.Sprintf("%s.utils", expenses[i].FileString())
		filePath := path.Join(dstFolder, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		GenerateExpenseRec(expenses[i], filePath)
	}
}

func GenerateExpenseRec(expense schema.Expense, dstPath string) {
	props := Properties{
		Type:       "Expense",
		Identifier: expense.Identifier,
		DstName:    expense.Identifier,
		Line1:      fmt.Sprintf("id: %s", expense.Id),
		Line2:      fmt.Sprintf("name: %s // amount: %.2f", expense.Name, expense.Amount),
		Line3:      fmt.Sprintf("accrual at: %s // advanced by 3th: %t // settlement at: %s", expense.DateOfAccrual, expense.AdvancedByThirdParty, expense.DateOfAccrual),
		Line4:      fmt.Sprintf("3th party: %s // customer: %s", "not-implemented", "not-implemented"),
	}
	pdf := NewPdf(expense.Path, dstPath)
	pdf.Generate(props)
}

func GenerateInvoicesRec(invoices schema.Invoices, dstPath string) {
	for i := range invoices {
		GenerateInvoiceRec(invoices[i], dstPath)
	}

}

func GenerateInvoiceRec(invoice schema.Invoice, dstPath string) {

}
