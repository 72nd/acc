package schema

import "github.com/google/uuid"

// BankStatement represents a bank statement.
type BankStatement struct {
	Id           string        `json:"id" default:"-"`
	Identifier   string        `json:"identifier" default:"e-19-01"`
	Period       string        `json:"period" default:"2019"`
	Transactions []Transaction `json:"transactions" default:"[]"`
}

// SetId sets a unique id to all elements in the slice.
func (s BankStatement) SetId() {
	for i := range s.Transactions {
		s.Transactions[i].SetId()
	}
}

// Transaction represents a single transaction of a bank statement.
type Transaction struct {
	Id     string  `json:"id" default:""`
	Amount float64 `json:"amount" default:"10.00"`
}

// GetId returns the unique id of the element.
func (t Transaction) GetId() string {
	return t.Id
}

// SetId generates a unique id for the element if there isn't already one defined.
func (t *Transaction) SetId() {
	if t.Id != "" {
		return
	}
	t.Id = uuid.Must(uuid.NewRandom()).String()
}
