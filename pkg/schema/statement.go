package schema

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const DefaultBankStatementFile = "bank-statement.yaml"
const DefaultTransactionPrefix = "t-"

// BankStatement represents a bank statement.
type BankStatement struct {
	Name         string        `yaml:"name" default:"e-19-01"`
	Period       string        `yaml:"period" default:"2019"`
	Transactions []Transaction `yaml:"transactions" default:"[]"`
}

// NewBankStatement returns a new BankStatement struct with the one Expense in it.
func NewBankStatement(useDefaults bool) BankStatement {
	stm := BankStatement{}
	if useDefaults {
		if err := defaults.Set(&stm); err != nil {
			logrus.Fatal("error setting defaults: ", err)
		}
		stm.Transactions = []Transaction{NewTransaction()}
	}
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

func (s *BankStatement) AddTransaction(trn []Transaction) {
	s.Transactions = append(s.Transactions, trn...)
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
	return "Bank-Statement"
}

// String returns a human readable representation of the element.
func (s BankStatement) String() string {
	return fmt.Sprintf("%s, for %s", s.Name, s.Period)
}

// Conditions returns the validation conditions.
func (s BankStatement) Conditions() util.Conditions {
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
func (s BankStatement) Validate() util.ValidateResults {
	results := util.ValidateResults{util.Check(s)}
	for i := range s.Transactions {
		results = append(results, util.Check(s.Transactions[i]))
	}
	return results
}

func (s BankStatement) TransactionSearchItems() util.SearchItems {
	result := make(util.SearchItems, len(s.Transactions))
	for i := range s.Transactions {
		result[i] = s.Transactions[i].SearchItem()
	}
	return result
}

func (s *BankStatement) AssistedCompletion(a Acc, doAll, autoSave, autoMode, askSkip bool) {
	first := true
	for i := range s.Transactions {
		if !first {
			fmt.Println()
		} else {
			first = false
		}
		s.Transactions[i] = s.Transactions[i].AssistedCompletion(a, doAll, autoMode, askSkip)
		if autoSave {
			a.SaveProject()
		}
	}
}
