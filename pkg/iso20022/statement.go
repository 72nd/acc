package iso20022

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/72nd/acc/pkg/schema"
)

// BankToCustomerStatement contains the path to a ISO 20022 camt XML file.
type BankToCustomerStatement struct {
	CamtPath string
}

// NewBankToCustomerStatement returns a new BankToCustomerStatement with the given path.
func NewBankToCustomerStatement(path string) BankToCustomerStatement {
	return BankToCustomerStatement{
		CamtPath: path,
	}
}

// Transactions reads the file for a given statement and returns the Transactions in the Acc data format.
func (s BankToCustomerStatement) Transactions() []schema.Transaction {
	file, err := os.Open(s.CamtPath)
	if err != nil {
		logrus.Fatalf("error reading %s: %s", s.CamtPath, err)
	}
	raw, _ := ioutil.ReadAll(file)
	var doc Document
	if err := xml.Unmarshal(raw, &doc); err != nil {
		logrus.Fatalf("error unmarshalling %s: %s", s.CamtPath, err)
	}
	trn := doc.AccTransactions()
	return trn
}
