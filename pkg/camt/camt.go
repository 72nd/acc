package camt

import (
	"encoding/xml"
	"fmt"
	"gitlab.com/72th/acc/pkg/schema"
)

const DateLayout = "2006-01-02"

type Document struct {
	XMLName xml.Name `xml:"Document"`
	Entries []Entry  `xml:"BkToCstmrStmt>Stmt>Ntry"`
}

type Entry struct {
	XMLName xml.Name `xml:"Ntry"`
	// Amount of transaction.
	Amount float64 `xml:"Amt"`
	// Booking is a reversal, should be checked.
	ReversalIndicator bool `xml:"RvslInd"`
	// `BOOK` or `PDNG`, in camt.053 only BOOK entries should be apparent.
	Status string `xml:"Sts"`
	// Unique number from the bank, identifies transactions.
	AccountServicerReference string        `xml:"AcctSvcrRef"`
	BookingData              string        `xml:"BookgDt>Dt"`
	ValueData                string        `xml:"ValDt>Dt"`
	Transactions             []Transaction `xml:"NtryDtls>TxDtls"`
}

func (e Entry) AccTransactions() []schema.Transaction {
	result := make([]schema.Transaction, len(e.Transactions))
	for i := range e.Transactions {
		result[i] = e.Transactions[i].AccTransaction(e.BookingData)
	}
	return result
}

type Transaction struct {
	XMLName              xml.Name `xml:"TxDtls"`
	Amount               float64  `xml:"Amt"`
	Description          string   `xml:"RmtInf>Ustrd"`
	CreditDebitIndicator string   `xml:"CdtDbtInd"` // `CRDT` or `DBIT`.
	Creditor             Party    `xml:"RltdPties>Cdtr"`
	Debitor              Party    `xml:"RltdPties>Dbtr"`
	Iban                 string   `xml:"RltdPties>CdtrAcct>Id>IBAN"`
	AccountCode          string   `xml:"RltdPties>CdtrAcct>Id>Othr>Id"`
	BankName             string   `xml:"RltdAgts>CdtrAgt>FinInstnId>Nm"`
}

func (t Transaction) AccTransaction(date string) schema.Transaction {
	trnType := schema.CreditTransaction
	thirdParty := t.Debitor.String()
	if t.CreditDebitIndicator == "DBIT" {
		trnType = schema.DebitTransaction
		thirdParty = t.Creditor.String()
	}
	trn := schema.Transaction{
		Description:     t.String(),
		TransactionType: trnType,
		ThirdParty:      thirdParty,
		ThirdPartyIdent: "",
		Date:            date,
		Amount:          t.Amount,
	}
	trn.SetId()
	return trn
}

func (t Transaction) String() string {
	typeStr := fmt.Sprintf("Received %.2f.- from %s", t.Amount, t.Debitor)
	if t.CreditDebitIndicator == "DBIT" {
		typeStr = fmt.Sprintf("Payed %.2f.- to %s", t.Amount, t.Creditor)
	}

	var description string
	if t.Description != "" {
		description = fmt.Sprintf(" with description: %s", t.Description)
	}
	return fmt.Sprintf("%s%s", typeStr, description)
}

type Party struct {
	Name           string `xml:"Nm"`
	AddressLine    string `xml:"PstlAdr>AdrLine"`
	StreetName     string `xml:"PstlAdr>StrtNm"`
	BuildingNumber int    `xml:"PstlAdr>BldgNb"`
	PostalCode     int    `xml:"PstlAdr>PstCd"`
	TownName       string `xml:"PstlAdr>TwnNm"`
	Country        string `xml:"PstlAdr>Ctry"`
}

func (p Party) String() string {
	result := p.Name
	var address string
	if p.StreetName != "" && p.BuildingNumber != 0 {
		address = fmt.Sprintf("%s %d", p.StreetName, p.BuildingNumber)
	}
	if p.PostalCode != 0 && p.TownName != "" {
		address = fmt.Sprintf("%s, %d %s", address, p.PostalCode, p.TownName)
	}
	if p.AddressLine != "" {
		address = p.AddressLine
	}
	if address != "" {
		result = fmt.Sprintf("%s (%s)", result, address)
	}
	return result
}
