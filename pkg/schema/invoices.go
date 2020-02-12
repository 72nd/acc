package schema

import (
	"encoding/json"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const DefaultInvoicesFile = "invoices.json"

// Invoices is a slice of invoices.
type Invoices []Invoice

// NewInvoices returns a new Expense slice with the one Expense in it.
func NewInvoices() Invoices {
	return []Invoice{NewInvoice()}
}

// OpenInvoices opens a Expenses saved in the json file given by the path.
func OpenInvoices(path string) Invoices {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	inv := Invoices{}
	if err := json.Unmarshal(raw, &inv); err != nil {
		logrus.Fatal(err)
	}
	return inv
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (i *Invoices) Save(path string, indented bool) {
	SaveToJson(i, path, indented)
}

// SetId sets a unique id to all elements in the slice.
func (i Invoices) SetId() {
	for j := range i {
		i[j].SetId()
	}
}

// Invoice represents an invoice sent to a customer for some services.
type Invoice struct {
	// Id is the internal unique identifier of the Expense.
	Id string `json:"id" default:"1"`
	// Identifier is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `json:"identifier" default:"i-19-1"`
	// Name describes meaningful the kind of the Expense.
	Name string `json:"name" default:"Expense Name"`
	// Amount states the amount of the Expense.
	Amount float64 `json:"amount" default:"10.00"`
	// Path is the full path to the voucher pdf.
	Path string `json:"path" default:"/path/to/file.pdf"`
	// CustomerId refers to the customer the invoice was sent to.
	CustomerId string `json:"customerId" default:""`
	// ProjectName refers to the project the invoice is associated with.
	ProjectId string `json:"projectId" default:""`
	// SendDate states the date, the invoice was sent to the customer.
	SendDate string `json:"sendDate" default:"2019-12-20"`
	// DateOfSettlement states the date the customer paid the outstanding amount.
	DateOfSettlement string `json:"dateOfSettlement" default:"2019-12-25"`
	// SettlementTransactionId refers to a possible bank transaction which settled the Expense for the company.
	SettlementTransactionId string `json:"settlementTransactionId" default:""`
}

// NewInvoice returns a new Acc element with the default values.
func NewInvoice() Invoice {
	inv := Invoice{}
	if err := defaults.Set(&inv); err != nil {
		logrus.Fatal(err)
	}
	return inv
}

// GetId returns the unique id of the element.
func (i Invoice) GetId() string {
	return i.Id
}

// SetId generates a unique id for the element if there isn't already one defined.
func (i *Invoice) SetId() {
	if i.Id != "" {
		return
	}
	i.Id = uuid.Must(uuid.NewRandom()).String()
}
