package camt

import (
	"encoding/xml"
	"fmt"
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
	fmt.Println(doc)
	fmt.Println("++++++")
	for _, entry := range doc.Entries {
		fmt.Printf("%.2f %s %s\n", entry.Amount, entry.BookingData, entry.ValueData)
		for _, trans := range entry.Transactions {
			fmt.Println(trans)
		}
		fmt.Println("-------------")
	}
	return []schema.Transaction{}
}