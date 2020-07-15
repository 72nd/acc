package invoices

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/72nd/acc/pkg/schema"
	"os"
	"path"
)

func GenerateAllInvoices(s schema.Schema, dstFolder, place string, doOverwrite bool) {
	nFiles := len(s.Invoices)
	for i := range s.Invoices {
		fileName := fmt.Sprintf("%s.pdf", s.Invoices[i].FileString())
		filePath := path.Join(dstFolder, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		customer, err := s.Parties.CustomerById(s.Invoices[i].CustomerId)
		if err != nil {
			logrus.Errorf("found for invoice %s no customer (given: %s): %s", s.Invoices[i].Id, s.Invoices[i].CustomerId, err)
			continue
		}
		GenerateInvoice(s.Company, s.Invoices[i], *customer, place, filePath)
	}
}

func GenerateInvoice(company schema.Company, invoice schema.Invoice, customer schema.Party, place, dstPath string) {
	doc := NewInvoiceDocument(12, place)
	save(doc.Generate(company, invoice, customer), dstPath)
}
