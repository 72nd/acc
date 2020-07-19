package iso20022

import (
	"encoding/xml"
	"fmt"

	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
	"github.com/sirupsen/logrus"
)

// DateLayout states the default date layout used by the ISO 20022 standard.
const DateLayout = "2006-01-02"

// Document is the root node  of a bank statement.
type Document struct {
	XMLName xml.Name `xml:"Document"`
	Entries []Entry  `xml:"BkToCstmrStmt>Stmt>Ntry"`
}

// AccTransactions pareses the Transactions of a given file and returns it as Transaction structs.
func (d Document) AccTransactions(currency string) []schema.Transaction {
	var result []schema.Transaction
	for i := range d.Entries {
		result = append(result, d.Entries[i].AccTransactions(currency)...)
	}
	return result
}

// Entry is a ISO 20022 entry.
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

// AccTransactions returns the transactions of a given entry.
func (e Entry) AccTransactions(currency string) []schema.Transaction {
	result := make([]schema.Transaction, len(e.Transactions))
	for i := range e.Transactions {
		result[i] = e.Transactions[i].AccTransaction(e.BookingData, currency)
	}
	return result
}

// Transaction reassembles a ISO 20022 transaction.
type Transaction struct {
	XMLName              xml.Name `xml:"TxDtls"`
	Amount               string   `xml:"Amt"`
	Description          string   `xml:"RmtInf>Ustrd"`
	CreditDebitIndicator string   `xml:"CdtDbtInd"` // `CRDT` or `DBIT`.
	Creditor             Party    `xml:"RltdPties>Cdtr"`
	Debitor              Party    `xml:"RltdPties>Dbtr"`
	Iban                 string   `xml:"RltdPties>CdtrAcct>Id>IBAN"`
	AccountCode          string   `xml:"RltdPties>CdtrAcct>Id>Othr>Id"`
	BankName             string   `xml:"RltdAgts>CdtrAgt>FinInstnId>Nm"`
}

// AccTransaction converts an ISO 20022 transaction into a Acc bank account transaction.
func (t Transaction) AccTransaction(date, currency string) schema.Transaction {
	trnType := util.CreditTransaction
	if t.CreditDebitIndicator == "DBIT" {
		trnType = util.DebitTransaction
	}
	amount, err := util.NewMonyFromDotNotation(t.Amount, currency)
	if err != nil {
		logrus.Fatal(err)
	}
	trn := schema.Transaction{
		Description:       t.String(),
		TransactionType:   trnType,
		AssociatedPartyId: "",
		Date:              date,
		Amount:            amount,
	}
	trn.SetId()
	return trn
}

// String returns a human readable string of a given Transaction.
func (t Transaction) String() string {
	typeStr := fmt.Sprintf("Received %.2f.- from %s", t.Amount, t.Debitor)
	if t.CreditDebitIndicator == "DBIT" {
		typeStr = fmt.Sprintf("Paid %.2f.- to %s", t.Amount, t.Creditor)
	}

	var description string
	if t.Description != "" {
		description = fmt.Sprintf(" with description: %s", t.Description)
	}
	return fmt.Sprintf("%s%s", typeStr, description)
}

// Party reassembles a ISO 20022 party.
type Party struct {
	Name           string `xml:"Nm"`
	AddressLine    string `xml:"PstlAdr>AdrLine"`
	StreetName     string `xml:"PstlAdr>StrtNm"`
	BuildingNumber int    `xml:"PstlAdr>BldgNb"`
	PostalCode     int    `xml:"PstlAdr>PstCd"`
	TownName       string `xml:"PstlAdr>TwnNm"`
	Country        string `xml:"PstlAdr>Ctry"`
}

// String returns a human readable string of a given party.
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
