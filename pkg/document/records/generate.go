package records

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"os"
	"path"
)

func GenerateExpensesRec(s schema.Schema, dstFolder string, doOverwrite, downConvert bool) {
	nFiles := len(s.Expenses)
	for i := range s.Expenses {
		fileName := fmt.Sprintf("%s.pdf", s.Expenses[i].FileString())
		filePath := path.Join(dstFolder, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		GenerateExpenseRec(s, s.Expenses[i], filePath, downConvert)
	}
}

func GenerateExpenseRec(s schema.Schema, exp schema.Expense, dstPath string, downConvert bool) {
	emp := "no 3rd party"
	if exp.AdvancedByThirdParty {
		emp = s.Parties.EmployeeStringById(exp.AdvancedThirdPartyId)
	}
	props := Properties{
		Type:       "Expense",
		Identifier: exp.Identifier,
		DstName:    exp.Identifier,
		Line1:      fmt.Sprintf("id: %s", exp.Id),
		Line2:      fmt.Sprintf("name: %s // amount: %.2f", exp.Name, exp.Amount),
		Line3:      fmt.Sprintf("accrual at: %s // advanced by 3th: %t // settlement at: %s", exp.DateOfAccrual, exp.AdvancedByThirdParty, exp.DateOfAccrual),
		Line4:      fmt.Sprintf("3rd party: %s // customer: %s", emp, s.Parties.CustomerStringById(exp.ObligedCustomerId)),
	}
	pdf := NewPdf(exp.Path, dstPath)
	pdf.Generate(props, downConvert)
}

func GenerateInvoicesRec(s schema.Schema, dstPath string, doOverwrite, downConvert bool) {
	nFiles := len(s.Invoices)
	for i := range s.Invoices {
		fileName := fmt.Sprintf("%s.pdf", s.Invoices[i].FileString())
		filePath := path.Join(dstPath, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		GenerateInvoiceRec(s, s.Invoices[i], filePath, downConvert)
	}

}

func GenerateInvoiceRec(s schema.Schema, inv schema.Invoice, dstPath string, downConvert bool) {
	props := Properties{
		Type:       "Invoice",
		Identifier: inv.Identifier,
		DstName:    inv.Identifier,
		Line1:      fmt.Sprintf("id %s", inv.Id),
		Line2:      fmt.Sprintf("name: %s // amount: %.2f", inv.Name, inv.Amount),
		Line3: fmt.Sprintf("send at: %s // settlement at %s", inv.SendDate, inv.DateOfSettlement),
		Line4: fmt.Sprintf("customer: %s", s.Parties.CustomerStringById(inv.CustomerId)),
	}
	pdf := NewPdf(inv.Path, dstPath)
	pdf.Generate(props, downConvert)
}
