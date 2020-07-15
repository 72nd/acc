package schema

import (
	"fmt"
	"time"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gitlab.com/72nd/acc/pkg/util"
)

const DefaultStatementFile = "bank-statement.yaml"
const DefaultTransactionPrefix = "t-"

// BankStatement represents a bank statement.
type Statement struct {
	Name         string        `yaml:"name" default:"e-19-01"`
	Period       string        `yaml:"period" default:"2019"`
	Transactions []Transaction `yaml:"transactions" default:"[]"`
}

// NewBankStatement returns a new BankStatement struct with the one Expense in it.
func NewBankStatement(useDefaults bool) Statement {
	stm := Statement{}
	if useDefaults {
		if err := defaults.Set(&stm); err != nil {
			logrus.Fatal("error setting defaults: ", err)
		}
		stm.Transactions = []Transaction{NewTransaction()}
	}
	return stm
}

// OpenBankStatement opens a BankStatement struct saved in the json file given by the path.
func OpenBankStatement(path string) Statement {
	var stm Statement
	util.OpenYaml(&stm, path, "statement")
	return stm
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (t Statement) Save(path string) {
	util.SaveToYaml(t, path, "statement")
}

func (t *Statement) AddTransaction(trn []Transaction) {
	t.Transactions = append(t.Transactions, trn...)
}

// SetId sets a unique id to all elements in the slice.
func (t Statement) SetId() {
	for i := range t.Transactions {
		t.Transactions[i].SetId()
	}
}

func (t Statement) GetIdentifiables() []Identifiable {
	trn := make([]Identifiable, len(t.Transactions))
	for i := range t.Transactions {
		trn[i] = t.Transactions[i]
	}
	return trn
}

func (t Statement) TransactionById(id string) (*Transaction, error) {
	for i := range t.Transactions {
		if t.Transactions[i].Id == id {
			return &t.Transactions[i], nil
		}
	}
	return nil, fmt.Errorf("no transaction for id \"%s\" found", id)
}

// Type returns a string with the type name of the element.
func (t Statement) Type() string {
	return "Bank-Statement"
}

// String returns a human readable representation of the element.
func (t Statement) String() string {
	return fmt.Sprintf("%s, for %s", t.Name, t.Period)
}

// Conditions returns the validation conditions.
func (t Statement) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: t.Name == "",
			Message:   "name is not set (Name is empty)",
		},
		{
			Condition: t.Period == "",
			Message:   "period is not set (Period is empty)",
		},
	}
}

// Validate the element and return the result.
func (t Statement) Validate() util.ValidateResults {
	results := util.ValidateResults{util.Check(t)}
	for i := range t.Transactions {
		results = append(results, util.Check(t.Transactions[i]))
	}
	return results
}

func (t Statement) TransactionSearchItems() util.SearchItems {
	result := make(util.SearchItems, len(t.Transactions))
	for i := range t.Transactions {
		result[i] = t.Transactions[i].SearchItem()
	}
	return result
}

func (t *Statement) AssistedCompletion(s Schema, doAll, autoSave, autoMode, askSkip, documentsOnly bool) {
	first := true
	for i := range t.Transactions {
		if !first {
			fmt.Println()
		} else {
			first = false
		}
		t.Transactions[i] = t.Transactions[i].AssistedCompletion(s, doAll, autoMode, askSkip, documentsOnly)
		if autoSave {
			s.Save()
		}
	}
}

func (t Statement) TransactionForDocument(id string) (*Transaction, error) {
	for i := range t.Transactions {
		if t.Transactions[i].AssociatedDocumentId == id {
			return &t.Transactions[i], nil
		}
	}
	return nil, fmt.Errorf("no transaction with associated document \"%s\" found", id)
}

func (t Statement) FilterTransactions(from *time.Time, to *time.Time) ([]Transaction, error) {
	var result []Transaction
	for i := range t.Transactions {
		date, err := time.Parse(util.DateFormat, t.Transactions[i].Date)
		rsl := false
		if err != nil {
			return nil, fmt.Errorf("transaction \"%s\": %s", t.Transactions[i].String(), err)
		}
		if from != nil && to == nil && (date.After(*from) || date.Equal(*from)) {
			rsl = true
		}
		if to != nil && from == nil && (date.Before(*to) || date.Equal(*to)) {
			rsl = true
		}
		if to != nil && from != nil && (date.After(*from) || date.Equal(*from)) && (date.Before(*to) || date.Equal(*to)) {
			rsl = true
		}
		if rsl {
			result = append(result, t.Transactions[i])
		}
	}
	return result, nil
}
