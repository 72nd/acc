package records

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
)

func GenerateExpensesRec(s schema.Schema, dstFolder string, doOverwrite, downConvert bool) {
	nFiles := len(s.Expenses)
	var wg sync.WaitGroup
	for i := range s.Expenses {
		fileName := fmt.Sprintf("%s.pdf", s.Expenses[i].FileString())
		filePath := path.Join(dstFolder, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		wg.Add(1)
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		go GenerateExpenseRec(s, s.Expenses[i], filePath, downConvert, &wg)
	}
	wg.Wait()
}

func GenerateExpenseRec(s schema.Schema, exp schema.Expense, dstPath string, downConvert bool, wg *sync.WaitGroup) {
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
		Line3:      fmt.Sprintf("accrual at: %s // advanced by 3th: %t // settlement at: %s", exp.DateOfAccrual, exp.AdvancedByThirdParty, exp.DateOfSettlement),
		Line4:      fmt.Sprintf("3rd party: %s // customer: %s", emp, s.Parties.CustomerStringById(exp.ObligedCustomerId)),
	}
	pdf := NewPdf(exp.Path, dstPath)
	pdf.Generate(props, downConvert)
	wg.Done()
}

func GenerateInvoicesRec(s schema.Schema, dstPath string, doOverwrite, downConvert bool) {
	nFiles := len(s.Invoices)
	var wg sync.WaitGroup

	for i := range s.Invoices {
		fileName := fmt.Sprintf("%s.pdf", s.Invoices[i].FileString())
		filePath := path.Join(dstPath, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		wg.Add(1)
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		go GenerateInvoiceRec(s, s.Invoices[i], filePath, downConvert, &wg)
	}
	wg.Wait()
}

func GenerateInvoiceRec(s schema.Schema, inv schema.Invoice, dstPath string, downConvert bool, wg *sync.WaitGroup) {
	props := Properties{
		Type:       "Invoice",
		Identifier: inv.Identifier,
		DstName:    inv.Identifier,
		Line1:      fmt.Sprintf("id %s", inv.Id),
		Line2:      fmt.Sprintf("name: %s // amount: %.2f", inv.Name, inv.Amount),
		Line3:      fmt.Sprintf("send at: %s // settlement at %s", inv.SendDate, inv.DateOfSettlement),
		Line4:      fmt.Sprintf("customer: %s", s.Parties.CustomerStringById(inv.CustomerId)),
	}
	pdf := NewPdf(inv.Path, dstPath)
	pdf.Generate(props, downConvert)
	wg.Done()
}

func GenerateMiscsRec(s schema.Schema, dstPath string, doOverwrite, downConvert bool) {
	nFiles := len(s.MiscRecords)

	var wg sync.WaitGroup
	for i := range s.MiscRecords {
		fileName := fmt.Sprintf("%s.pdf", s.MiscRecords[i].FileString())
		filePath := path.Join(dstPath, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		wg.Add(1)
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		go GenerateMiscRec(s, s.MiscRecords[i], filePath, downConvert, &wg)
	}
	wg.Wait()
}

func GenerateMiscRec(s schema.Schema, mrc schema.MiscRecord, dstPath string, downConvert bool, wg *sync.WaitGroup) {
	props := Properties{
		Type:       "Miscellaneous Record",
		Identifier: mrc.Identifier,
		DstName:    mrc.Identifier,
		Line1:      fmt.Sprintf("id %s", mrc.Id),
		Line2:      fmt.Sprintf("name: %s", mrc.Name),
		Line3:      fmt.Sprintf("received at: %s", mrc.Date),
	}
	pdf := NewPdf(mrc.Path, dstPath)
	pdf.Generate(props, downConvert)
	wg.Done()
}
