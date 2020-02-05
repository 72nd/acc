package schema

import "github.com/google/uuid"

// Expenses is a slice of multiple expenses.
type Expenses []Expense

// SetId sets a unique id to all elements in the slice.
func (e Expenses) SetId() {
	for i := range e {
		e[i].SetId()
	}
}

// Expense represents a payment done by the company or a third party to assure the ongoing of the business.
type Expense struct {
	// Id is the internal unique identifier of the Expense.
	Id string `json:"id" default:"1"`
	// Identifier is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `json:"identifier" default:"e-19-1"`
	// Name describes meaningful the kind of the Expense.
	Name string `json:"name" default:"Expense Name"`
	// Amount states the amount of the Expense.
	Amount float64 `json:"amount" default:"10.00"`
	// Path is the full path to the voucher pdf.
	Path string `json:"path" default:"/path/to/file.pdf"`
	// DateOfAccrual represents the day the obligation emerged.
	DateOfAccrual string `json:"dateOfAccrual" default:"2019-12-20"`
	// Billable states if the costs for the Expense will be forwarded to the customer.
	Billable bool `json:"billable" default:"false"`
	// ObligedCustomerId refers to the customer which have to pay the Expense.
	ObligedCustomerId int `json:"obligedCustomerId" default:"0"`
	// AdvancedByThirdParty states if a third party (employee, etc.) advanced the payment of this expense for the company.
	AdvancedByThirdParty bool `json:"advancedByThirdParty" default:"false"`
	// AdvancePartyId refers to the third party which advanced the payment.
	AdvancedThirdPartyId int `json:"advancedThirdPartyId" default:"0"`
	// DateOfSettlement states the date of the settlement of the expense (the company has not to take further actions).
	DateOfSettlement string `json:"dateOfSettlement" default:"2019-12-25"`
	// SettlementTransactionId refers to a possible bank transaction which settled the Expense for the company.
	SettlementTransactionId string `json:"settlementTransactionId" default:""`
}

// GetId returns the unique id of the element.
func (e Expense) GetId() string {
	return e.Id
}

// SetId generates a unique id for the element if there isn't already one defined.
func (e *Expense) SetId() {
	if e.Id != "" {
		return
	}
	e.Id = uuid.Must(uuid.NewRandom()).String()
}
