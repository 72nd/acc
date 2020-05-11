package camt

import (
	"encoding/xml"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"io/ioutil"
	"os"
)

type BankToCustomerStatement struct {
	CamtPath string
}

func NewBankToCustomerStatement(path string) BankToCustomerStatement {
	return BankToCustomerStatement{
		CamtPath: path,
	}
}

func (s BankToCustomerStatement) Transactions(interactive bool) []schema.Transaction {
	file, err := os.Open(s.CamtPath)
	if err != nil {
		logrus.Fatalf("error reading %s: %s", s.CamtPath, err)
	}
	raw, _ := ioutil.ReadAll(file)
	var doc Document
	if err := xml.Unmarshal(raw, &doc); err != nil {
		logrus.Fatalf("error unmarshalling %s: %s", s.CamtPath, err)
	}
	return []schema.Transaction{}
}
