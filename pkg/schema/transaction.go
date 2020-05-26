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

type JournalMode int

const (
	UnknownJournalMode JournalMode = iota
	ManualJournalMode
	AutoJournalMode
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
	JournalMode          JournalMode     `yaml:"journalMode" default:"0"`
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
		"Value",
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
				Name:        "Incoming transaction",
				Value:       "0",
				SearchValue: "1 Incoming Transaction",
			},
			util.SearchItem{
				Name:        "Outgoing transaction",
				Value:       "2",
				SearchValue: "2 Outgoing Transaction",
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

func (t Transaction) AssistedCompletion(a Acc, doAll bool) Transaction {
	tmp := t
	if !doAll && t.Id != "" && t.Identifier != "" {
		fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Skip transaction:")), aurora.BrightMagenta(t.Description))
		return t
	}
	fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Optimize transaction:")), aurora.BrightMagenta(t.Description))
	identifier := SuggestNextIdentifier(a.BankStatement.GetIdentifiables(), DefaultTransactionPrefix)
	if t.Id == "" {
		t.SetId()
	}
	if t.Identifier == "" {
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
	} else {
		docs := append(a.Expenses.SearchItems(), a.Invoices.SearchItems()...)
		t.AssociatedDocumentId = util.AskStringFromSearch(
			"Associated Document",
			"couldn't find associated document, manual search",
			docs)
	}
	t.JournalMode = JournalMode(util.AskIntFromList(
		"Journal Mode",
		"choose how journal entry will be generated for this transaction",
		util.SearchItems{
			{
				Name: "Manual Mode",
				Value: ManualJournalMode,
			},
			{
				Name: "Auto Mode",
				Value: AutoJournalMode,
			},
		}))

	ok := util.AskForConformation("Were your entries correct?")
	if !ok {
		for {
			strategy := util.AskIntFromList(
				"Strategy",
				"how do you want to resolve this situation?",
				util.SearchItems{
					{
						Name:  "Redo",
						Value: 1,
					},
					{
						Name:  "Skip",
						Value: 2,
					},
				})
			switch strategy {
			case 1:
				t.AssistedCompletion(a, doAll)
			case 2:
				return tmp
			default:
				logrus.Error("invalid input, try again")
			}
		}
	}
	return t
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
	return "Transaction"
}

// String returns a human readable representation of the element.
func (t Transaction) String() string {
	if t.TransactionType == CreditTransaction {
		return fmt.Sprintf("%s: received %.2f at %s", t.Identifier, t.Amount, t.Date)
	}
	return fmt.Sprintf("%s: payed %.2f at %s", t.Identifier, t.Amount, t.Date)
}

// Conditions returns the validation conditions.
func (t Transaction) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: t.Id == "",
			Message:   "unique identifier not set (Id is empty)",
			Level:     util.BeforeExportFlaw,
		},
		{
			Condition: t.Identifier == "",
			Message:   "human readable identifier not set (Identifier is empty)",
			Level:     util.BeforeExportFlaw,
		},
		{
			Condition: t.Description == "",
			Message:   "no description set",
			Level:     util.BeforeExportFlaw,
		},
		{
			Condition: t.TransactionType < 0 || t.TransactionType > 1,
			Message:   "transaction type not valid",
			Level:     util.BeforeExportFlaw,
		},
		{
			Condition: t.AssociatedPartyId == "" && t.JournalMode == AutoJournalMode,
			Message:   "no associated party set although auto journal mode is set",
			Level:     util.BeforeMergeFlaw,
		},
		{
			Condition: t.AssociatedDocumentId == "" && t.JournalMode == AutoJournalMode,
			Message:   "no associated document set although auto journal mode is set",
			Level:     util.BeforeMergeFlaw,
		},
		{
			Condition: t.Date == "",
			Message:   "date not set",
			Level:     util.BeforeMergeFlaw,
		},
		{
			Condition: t.Amount <= 0,
			Message:   "amount is not set",
			Level:     util.BeforeMergeFlaw,
		},
	}
}

// Validate the element and return the result.
func (t Transaction) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(t)}
}

func (t Transaction) SearchItem() util.SearchItem {
	return util.SearchItem{
		Name:        t.Description,
		Value:       t.Id,
		SearchValue: t.Description,
	}
}
