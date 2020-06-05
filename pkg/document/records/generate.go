package records

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"os"
	"path"
)

func GenerateExpensesRec(a schema.Acc, dstFolder string, doOverwrite, downConvert bool) {
	nFiles := len(a.Expenses)
	for i := range a.Expenses {
		fileName := fmt.Sprintf("%s.pdf", a.Expenses[i].FileString())
		filePath := path.Join(dstFolder, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		GenerateExpenseRec(a, a.Expenses[i], filePath, downConvert)
	}
}

func GenerateExpenseRec(a schema.Acc, exp schema.Expense, dstPath string, downConvert bool) {
	emp := "no 3rd party"
	if exp.AdvancedByThirdParty {
		emp = a.Parties.EmployeeStringById(exp.AdvancedThirdPartyId)
	}
	props := Properties{
		Type:       "Expense",
		Identifier: exp.Identifier,
		DstName:    exp.Identifier,
		Line1:      fmt.Sprintf("id: %s", exp.Id),
		Line2:      fmt.Sprintf("name: %s // amount: %.2f", exp.Name, exp.Amount),
		Line3:      fmt.Sprintf("accrual at: %s // advanced by 3th: %t // settlement at: %s", exp.DateOfAccrual, exp.AdvancedByThirdParty, exp.DateOfAccrual),
		Line4:      fmt.Sprintf("3rd party: %s // customer: %s", emp, a.Parties.CustomerStringById(exp.ObligedCustomerId)),
	}
	pdf := NewPdf(exp.Path, dstPath)
	pdf.Generate(props, downConvert)
}

func GenerateInvoicesRec(a schema.Acc, dstPath string, doOverwrite, downConvert bool) {
	nFiles := len(a.Invoices)
	for i := range a.Invoices {
		fileName := fmt.Sprintf("%s.pdf", a.Invoices[i].FileString())
		filePath := path.Join(dstPath, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		GenerateInvoiceRec(a, a.Invoices[i], filePath, downConvert)
	}

}

func GenerateInvoiceRec(a schema.Acc, inv schema.Invoice, dstPath string, downConvert bool) {
	props := Properties{
		Type:       "Invoice",
		Identifier: inv.Identifier,
		DstName:    inv.Identifier,
		Line1:      fmt.Sprintf("id %s", inv.Id),
		Line2:      fmt.Sprintf("name: %s // amount: %.2f", inv.Name, inv.Amount),
		Line3: fmt.Sprintf("send at: %s // settlement at %s", inv.SendDate, inv.DateOfSettlement),
		Line4: fmt.Sprintf("customer: %s", a.Parties.CustomerStringById(inv.CustomerId)),
	}
	pdf := NewPdf(inv.Path, dstPath)
	pdf.Generate(props, downConvert)
}
