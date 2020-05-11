package schema

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"time"
)

type TransactionType int

const (
	CreditTransaction TransactionType = iota // Incoming transaction
	DebitTransaction                         // Outgoing transaction
)

// Transaction represents a single transaction of a bank statement.
type Transaction struct {
	Id              string          `yaml:"id" default:""`
	Identifier      string          `yaml:"identifier" default:""`
	Description     string          `yaml:"description" default:""`
	TransactionType TransactionType `yaml:"transactionType" default:"0"`
	ThirdPartyIdent string          `yaml:"thirdPartyIdent" default:""`
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

func (t *Transaction) AssistedCompletion(thirdParty string) {
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
func (t Transaction) String() string {
	if t.TransactionType == CreditTransaction {
		return fmt.Sprintf("%s: received %.2f at %s", t.Identifier, t.Amount, t.Date)
	}
	return fmt.Sprintf("%s: payed %.2f at %s", t.Identifier, t.Amount, t.Date)
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
