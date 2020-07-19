package records

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
	"github.com/sirupsen/logrus"
)

// GenerateExpensesRec generates an record document for all expenses in the given schema.
// The PDFs will then be saved into the given destination folder. DownConvert (converting
// all PDF's to a PostScript file and back) should be done always  as gopdf will have
// otherwise problems embedding most PDFs.
func GenerateExpensesRec(s schema.Schema, dstFolder string, doOverwrite, downConvert bool) {
	sort.Sort(s.Expenses)
	nFiles := len(s.Expenses)
	var wg sync.WaitGroup
	for i := range s.Expenses {
		fileName := fmt.Sprintf("e-%03d_%s.pdf", i+1, s.Expenses[i].FileString())
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

// GenerateExpenseRec generates and saves the record PDF for a given expense.
func GenerateExpenseRec(s schema.Schema, exp schema.Expense, dstPath string, downConvert bool, wg *sync.WaitGroup) {
	emp := "no 3rd party"
	if exp.AdvancedByThirdParty {
		emp = s.Parties.EmployeeStringById(exp.AdvancedThirdPartyId)
	}
	props := Properties{
		Type:       "Expense",
		Identifier: exp.Identifier,
		DstName:    exp.Identifier,
		Line1:      fmt.Sprintf("id: %s", na(exp.Id)),
		Line2:      fmt.Sprintf("name: %s // amount: %s", na(exp.Name), sfr(exp.Amount)),
		Line3:      fmt.Sprintf("accrual at: %s // expense category: %s // internal: %t", date(exp.DateOfAccrual), na(exp.ExpenseCategory), exp.Internal),
		Line4:      fmt.Sprintf("settlement at: %s // advanced by 3th: %s // 3rd party: %s", date(exp.DateOfSettlement), na(exp.AdvancedByThirdParty), emp),
		Line5:      fmt.Sprintf("customer: %s // payed with debit: %s", na(s.Parties.CustomerStringById(exp.ObligedCustomerId)), na(exp.PayedWithDebit)),
	}
	pdf := NewPdf(exp.Path, dstPath)
	pdf.Generate(props, downConvert)
	wg.Done()
}

// GenerateInvoicesRec generates an record document for all invoices in the given schema.
// The PDFs will then be saved into the given destination folder. DownConvert (converting
// all PDF's to a PostScript file and back) should be done always  as gopdf will have
// otherwise problems embedding most PDFs.
func GenerateInvoicesRec(s schema.Schema, dstPath string, doOverwrite, downConvert bool) {
	sort.Sort(s.Invoices)
	nFiles := len(s.Invoices)
	var wg sync.WaitGroup
	for i := range s.Invoices {
		fileName := fmt.Sprintf("i-%03d_%s.pdf", i+1, s.Invoices[i].FileString())
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

// GenerateInvoiceRec generates and saves the record PDF for a given invoice.
func GenerateInvoiceRec(s schema.Schema, inv schema.Invoice, dstPath string, downConvert bool, wg *sync.WaitGroup) {
	props := Properties{
		Type:       "Invoice",
		Identifier: inv.Identifier,
		DstName:    inv.Identifier,
		Line1:      fmt.Sprintf("id: %s", na(inv.Id)),
		Line2:      fmt.Sprintf("name: %s // amount: %s", na(inv.Name), sfr(inv.Amount)),
		Line3:      fmt.Sprintf("send at: %s // settlement: at %s", date(inv.SendDate), date(inv.DateOfSettlement)),
		Line4:      fmt.Sprintf("customer: %s", s.Parties.CustomerStringById(inv.CustomerId)),
	}
	pdf := NewPdf(inv.Path, dstPath)
	pdf.Generate(props, downConvert)
	wg.Done()
}

// GenerateMiscsRec generates an record document for all misc-records in the given schema.
// The PDFs will then be saved into the given destination folder. DownConvert (converting
// all PDF's to a PostScript file and back) should be done always  as gopdf will have
// otherwise problems embedding most PDFs.
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

// GenerateMiscRec generates and saves the record PDF for a given misc-record.
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

const NA = "n.a."

// na returns the string of a given data. If data is empty "N/A" will be returned.
// Remark: Do not use this if numbers (int, floats) <= 0 are valid.
func na(data interface{}) string {
	switch data.(type) {
	case string:
		if data == "" {
			return NA
		}
		return data.(string)
	case int:
		if data.(int) <= 0 {
			return NA
		}
		return strconv.FormatInt(data.(int64), 10)
	case bool:
		return strconv.FormatBool(data.(bool))
	case float64:
		if data.(float64) <= 0.0 {
			return NA
		}
		return strconv.FormatFloat(data.(float64), 'E', 2, 64)
	case float32:
		if data.(float32) <= 0.0 {
			return NA
		}
		return strconv.FormatFloat(data.(float64), 'E', 2, 64)
	}
	return "not implemented"
}

func sfr(amount float64) string {
	return fmt.Sprintf("SFr. %.2f", amount)
}

func date(data string) string {
	if data == "" {
		return NA
	}
	date, err := time.Parse(util.DateFormat, data)
	if err != nil {
		return fmt.Sprintf("invalid date: %s", data)
	}
	return date.Format("02.01.2006")
}
