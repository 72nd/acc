package invoices

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"os"
	"path"
)

func GenerateAllInvoices(acc schema.Acc, dstFolder, place string, doOverwrite bool) {
	nFiles := len(acc.Invoices)
	for i := range acc.Invoices {
		fileName := fmt.Sprintf("%s.pdf", acc.Invoices[i].FileString())
		filePath := path.Join(dstFolder, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		customer, err := acc.Parties.CustomerById(acc.Invoices[i].CustomerId)
		if err != nil {
			logrus.Errorf("found for invoice %s no customer (given: %s): %s", acc.Invoices[i].Id, acc.Invoices[i].CustomerId, err)
			continue
		}
		GenerateInvoice(acc.Company, acc.Invoices[i], *customer, place, filePath)
	}
}

func GenerateInvoice(company schema.Company, invoice schema.Invoice, customer schema.Party, place, dstPath string) {
	doc := NewInvoiceDocument(12, place)
	save(doc.Generate(company, invoice, customer), dstPath)
}
