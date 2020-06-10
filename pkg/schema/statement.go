package schema

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
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
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	stm := Statement{}
	if err := yaml.Unmarshal(raw, &stm); err != nil {
		logrus.Fatal(err)
	}
	return stm
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (s Statement) Save(path string) {
	SaveToYaml(s, path)
}

func (s *Statement) AddTransaction(trn []Transaction) {
	s.Transactions = append(s.Transactions, trn...)
}

// SetId sets a unique id to all elements in the slice.
func (s Statement) SetId() {
	for i := range s.Transactions {
		s.Transactions[i].SetId()
	}
}

func (s Statement) GetIdentifiables() []Identifiable {
	trn := make([]Identifiable, len(s.Transactions))
	for i := range s.Transactions {
		trn[i] = s.Transactions[i]
	}
	return trn
}

func (s Statement) TransactionById(id string) (*Transaction, error) {
	for i := range s.Transactions {
		if s.Transactions[i].Id == id {
			return &s.Transactions[i], nil
		}
	}
	return nil, fmt.Errorf("no transaction for id \"%s\" found", id)
}

// Type returns a string with the type name of the element.
func (s Statement) Type() string {
	return "Bank-Statement"
}

// String returns a human readable representation of the element.
func (s Statement) String() string {
	return fmt.Sprintf("%s, for %s", s.Name, s.Period)
}

// Conditions returns the validation conditions.
func (s Statement) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: s.Name == "",
			Message:   "name is not set (Name is empty)",
		},
		{
			Condition: s.Period == "",
			Message:   "period is not set (Period is empty)",
		},
	}
}

// Validate the element and return the result.
func (s Statement) Validate() util.ValidateResults {
	results := util.ValidateResults{util.Check(s)}
	for i := range s.Transactions {
		results = append(results, util.Check(s.Transactions[i]))
	}
	return results
}

func (s Statement) TransactionSearchItems() util.SearchItems {
	result := make(util.SearchItems, len(s.Transactions))
	for i := range s.Transactions {
		result[i] = s.Transactions[i].SearchItem()
	}
	return result
}

func (s *Statement) AssistedCompletion(a Acc, doAll, autoSave, autoMode, askSkip bool) {
	first := true
	for i := range s.Transactions {
		if !first {
			fmt.Println()
		} else {
			first = false
		}
		s.Transactions[i] = s.Transactions[i].AssistedCompletion(a, doAll, autoMode, askSkip)
		if autoSave {
			a.SaveAccComplex()
		}
	}
}

func (s Statement) TransactionForDocument(id string) (*Transaction, error) {
	for i := range s.Transactions {
		fmt.Println(s.Transactions[i].AssociatedDocumentId)
		if s.Transactions[i].AssociatedDocumentId == id {
			return &s.Transactions[i], nil
		}
	}
	return nil, fmt.Errorf("no transaction with associated document \"%s\" found", id)
}

func (s Statement) FilterTransactions(from *time.Time, to *time.Time) ([]Transaction, error) {
	var result []Transaction
	for i := range s.Transactions {
		date, err := time.Parse(util.DateFormat, s.Transactions[i].Date)
		rsl := false
		if err != nil {
			return nil, fmt.Errorf("transaction \"%s\": %s", s.Transactions[i].String(), err)
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
			result = append(result, s.Transactions[i])
		}
	}
	return result, nil
}
