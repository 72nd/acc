package schema

import "github.com/google/uuid"

// Invoices is a slice of invoices.
type Invoices []Invoice

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
	// ProjectId refers to the project the invoice is associated with.
	ProjectId string `json:"projectId" default:""`
	// SendDate states the date, the invoice was sent to the customer.
	SendDate string `json:"sendDate" default:"2019-12-20"`
	// DateOfSettlement states the date the customer paid the outstanding amount.
	DateOfSettlement string `json:"dateOfSettlement" default:"2019-12-25"`
	// SettlementTransactionId refers to a possible bank transaction which settled the Expense for the company.
	SettlementTransactionId string `json:"settlementTransactionId" default:""`
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
