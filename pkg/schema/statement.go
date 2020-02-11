package schema

import (
	"encoding/json"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const DefaultBankStatementFile = "bank-statement.json"

// BankStatement represents a bank statement.
type BankStatement struct {
	Id           string        `json:"id" default:"-"`
	Identifier   string        `json:"identifier" default:"e-19-01"`
	Period       string        `json:"period" default:"2019"`
	Transactions []Transaction `json:"transactions" default:"[]"`
}

// NewBankStatement returns a new BankStatement struct with the one Expense in it.
func NewBankStatement() BankStatement {
	stm := BankStatement{}
	if err := defaults.Set(&stm); err != nil {
		logrus.Fatal(err)
	}
	stm.Transactions = []Transaction{NewTransaction()}
	return stm
}

// OpenBankStatement opens a BankStatement struct saved in the json file given by the path.
func OpenBankStatement(path string) BankStatement {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	stm := BankStatement{}
	if err := json.Unmarshal(raw, &stm); err != nil {
		logrus.Fatal(err)
	}
	return stm
}

// Save writes the element as a json to the given path.
func (s BankStatement) Save(path string) {
	raw, err := json.Marshal(s)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := ioutil.WriteFile(path, raw, 0644); err != nil {
		logrus.Fatal(err)
	}
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

func NewTransaction() Transaction {
	trn := Transaction{}
	if err := defaults.Set(&trn); err != nil {
		logrus.Fatal(err)
	}
	return trn
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
