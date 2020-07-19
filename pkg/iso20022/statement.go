package iso20022

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/72nd/acc/pkg/schema"
)

type BankToCustomerStatement struct {
	CamtPath string
}

func NewBankToCustomerStatement(path string) BankToCustomerStatement {
	return BankToCustomerStatement{
		CamtPath: path,
	}
}

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
