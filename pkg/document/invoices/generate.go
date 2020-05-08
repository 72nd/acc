package invoices

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"os"
	"path"
)

func GenerateAllInvoices(invoices schema.Invoices, dstFolder string, doOverwrite bool) {
	nFiles := len(invoices)
	for i := range invoices {
		fileName := fmt.Sprintf("%s.pdf", invoices[i].FileString())
		filePath := path.Join(dstFolder, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) && !doOverwrite {
			logrus.Infof("(%d/%d) File %s exists, skipping", i+i, nFiles, fileName)
			continue
		}
		logrus.Infof("(%d/%d) Generate %s...", i+1, nFiles, fileName)
		GenerateInvoice(invoices[i], filePath)
	}
}

func GenerateInvoice(invoice schema.Invoice, dstPath string) {
	pdf := initPdf()
	pdf = page(pdf, invoice)
	save(pdf, dstPath)
}
