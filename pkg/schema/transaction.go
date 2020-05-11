package schema

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"regexp"
	"time"
)

type TransactionType int

const (
	CreditTransaction TransactionType = iota // Incoming transaction
	DebitTransaction                         // Outgoing transaction
)

// Transaction represents a single transaction of a bank statement.
type Transaction struct {
	Id                   string          `yaml:"id" default:""`
	Identifier           string          `yaml:"identifier" default:""`
	Description          string          `yaml:"description" default:""`
	TransactionType      TransactionType `yaml:"transactionType" default:"0"`
	AssociatedPartyId    string          `yaml:"associatedPartyId" default:""`
	AssociatedDocumentId string          `yaml:"associatedDocumentId" default:""`
	Date                 string          `yaml:"date" default:""`
	Amount               float64         `yaml:"amount" default:"10.00"`
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

func (t *Transaction) AssistedCompletion(a Acc) {
	t.SetId()
	fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Optimize transaction:")), aurora.BrightMagenta(t.Description))
	identifier := SuggestNextIdentifier(a.BankStatement.GetIdentifiables(), DefaultTransactionPrefix)
	if t.Identifier == "" && identifier != "" {
		t.Identifier = util.AskString(
			"Identifier",
			"Unique human readable identifier",
			identifier)
	}
	t.Description = util.AskString(
		"Description",
		"Description of the transaction",
		t.Description)
	parties := append(a.Parties.CustomersSearchItems(), a.Parties.EmployeesSearchItems()...)
	t.AssociatedPartyId = util.AskStringFromSearch(
		"Associated Party",
		"customer/employee which is originator/recipient of the transaction",
		parties)
	document, err := t.parseAssociatedDocument(a.Expenses, a.Invoices)
	if err == nil && util.AskForConformation(fmt.Sprintf("Use «%s» as associated document?", document.String())) {
		t.AssociatedDocumentId = document.GetId()
	}
}

func (t Transaction) parseAssociatedDocument(expenses Expenses, invoices Invoices) (Identifiable, error) {
	r := regexp.MustCompile(`([ei]-(|.*-)(\d+))(\s|$)`)
	matches := r.FindAllStringSubmatch(t.Description, -1)
	if len(matches) != 1 || len(matches[0]) != 5 {
		return nil, fmt.Errorf("no document identifier found")
	}
	ident := matches[0][1]
	expense, err1 := expenses.ExpenseByIdent(ident)
	invoice, err2 := invoices.InvoiceByIdent(ident)
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("no expense or invoice for identifier %s found", ident)
	}
	if invoice == nil {
		return expense, nil
	}
	return invoice, nil
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
