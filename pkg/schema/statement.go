package schema

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

const DefaultBankStatementFile = "bank-statement.yaml"
const DefaultTransactionPrefix = "t-"

// BankStatement represents a bank statement.
type BankStatement struct {
	Id           string        `yaml:"id" default:"-"`
	Identifier   string        `yaml:"identifier" default:"e-19-01"`
	Period       string        `yaml:"period" default:"2019"`
	Transactions []Transaction `yaml:"transactions" default:"[]"`
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
	if err := yaml.Unmarshal(raw, &stm); err != nil {
		logrus.Fatal(err)
	}
	return stm
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (s BankStatement) Save(path string) {
	SaveToYaml(s, path)
}

// SetId sets a unique id to all elements in the slice.
func (s BankStatement) SetId() {
	for i := range s.Transactions {
		s.Transactions[i].SetId()
	}
}

func (s BankStatement) GetIdentifiables() []Identifiable {
	trn := make([]Identifiable, len(s.Transactions))
	for i := range s.Transactions {
		trn[i] = s.Transactions[i]
	}
	return trn
}

// Type returns a string with the type name of the element.
func (s BankStatement) Type() string {
	return ""
}

// String returns a human readable representation of the element.
func (s BankStatement) String() string {
	return fmt.Sprintf("")
}

// Conditions returns the validation conditions.
func (s BankStatement) Conditions() util.Conditions {
	return util.Conditions{

	}
}

// Validate the element and return the result.
func (s BankStatement) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(s)}
}

func (s BankStatement) TransactionSearchItems() util.SearchItems {
	result := make(util.SearchItems, len(s.Transactions))
	for i := range s.Transactions {
		result[i] = s.Transactions[i].SearchItem()
	}
	return result
}

type TransactionType int

const (
	IncomingTransaction TransactionType = iota
	OutgoingTransaction
)

// Transaction represents a single transaction of a bank statement.
type Transaction struct {
	Id              string          `yaml:"id" default:""`
	Identifier      string          `yaml:"identifier" default:""`
	Description     string          `yaml:"description" default:""`
	TransactionType TransactionType `yaml:"transactionType" default:"0"`
	Date            string          `yaml:"date" default:""`
	Amount          float64         `yaml:"amount" default:"10.00"`
}

func NewTransaction() Transaction {
	trn := Transaction{}
	if err := defaults.Set(&trn); err != nil {
		logrus.Fatal(err)
	}
	return trn
}

func NewTransactionWithUuid() Transaction {
	trn := NewTransaction()
	trn.Id = GetUuid()
	return trn
}

func InteractiveNewTransaction(s BankStatement) Transaction {
	trn := NewTransactionWithUuid()
	trn.Identifier = util.AskString(
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(s.GetIdentifiables(), DefaultTransactionPrefix),
	)
	trn.Description = util.AskString(
		"Description",
		"Some information about the transaction",
		"",
	)
	trn.TransactionType = TransactionType(util.AskIntFromListSearch(
		"Transaction Type",
		"",
		util.SearchItems{
			util.SearchItem{
				Name:       "Incoming transaction",
				Identifier: "0",
				Value:      "1 Incoming Transaction",
			},
			util.SearchItem{
				Name:       "Outgoing transaction",
				Identifier: "2",
				Value:      "2 Outgoing Transaction",
			},
		}))
	trn.Date = util.AskDate(
		"Date",
		"Transaction date",
		time.Now(),
	)
	trn.Amount = util.AskFloat(
		"Amount",
		"",
		23.42,
	)
	return trn
}

// GetId returns the unique id of the element.
func (t Transaction) GetId() string {
	return t.Id
}

func (t Transaction) GetIdentifier() string {
	return t.Identifier
}

// SetId generates a unique id for the element if there isn't already one defined.
func (t *Transaction) SetId() {
	if t.Id != "" {
		return
	}
	t.Id = uuid.Must(uuid.NewRandom()).String()
}

// Type returns a string with the type name of the element.
func (Transaction) Type() string {
	return ""
}

// String returns a human readable representation of the element.
func (Transaction) String() string {
	return fmt.Sprintf("")
}

// Conditions returns the validation conditions.
func (Transaction) Conditions() util.Conditions {
	return util.Conditions{

	}
}

// Validate the element and return the result.
func (t Transaction) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(t)}
}

func (t Transaction) SearchItem() util.SearchItem {
	return util.SearchItem{
		Name:       t.Description,
		Identifier: t.Id,
		Value:      t.Description,
	}
}
