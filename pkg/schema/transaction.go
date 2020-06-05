package schema

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
)

type JournalMode int

const (
	UnknownJournalMode JournalMode = iota
	ManualJournalMode
	AutoJournalMode
)

// Transaction represents a single transaction of a bank statement.
type Transaction struct {
	Id                   string               `yaml:"id" default:""`
	Identifier           string               `yaml:"identifier" default:""`
	Description          string               `yaml:"description" default:""`
	TransactionType      util.TransactionType `yaml:"transactionType" default:"0"`
	AssociatedPartyId    string               `yaml:"associatedPartyId" default:""`
	AssociatedDocumentId string               `yaml:"associatedDocumentId" default:""`
	Date                 string               `yaml:"date" default:""`
	Amount               float64              `yaml:"amount" default:"10.00"`
	JournalMode          JournalMode          `yaml:"journalMode" default:"0"`
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
	trn.TransactionType = util.TransactionType(util.AskIntFromListSearch(
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

func (t Transaction) AssistedCompletion(a Acc, doAll, autoMode, askSkip bool) Transaction {
	tmp := t
	if autoMode {
		t.JournalMode = AutoJournalMode
	}
	if !doAll && util.Check(t).Valid() {
		fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Skip transaction:")), aurora.BrightMagenta(t.String()))
		return t
	}
	fmt.Printf("%s %s %s\n", aurora.BrightMagenta(aurora.Bold("Optimize transaction:")), aurora.BrightMagenta(t.String()), aurora.BrightMagenta(t.Description))
	if !doAll && askSkip {
		skip := util.AskBool(
			"Skip",
			"Skip this non valid entry?",
			false)
		if skip {
			return t
		}
	}
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
	t.JournalMode = JournalMode(util.AskIntFromList(
		"Journal Mode",
		"choose how journal entry will be generated for this transaction",
		util.SearchItems{
			{
				Name:  "Manual Mode",
				Value: int(ManualJournalMode),
			},
			{
				Name:  "Auto Mode",
				Value: int(AutoJournalMode),
			},
		}))
	fmt.Println(t.JournalMode)
	if t.AssociatedPartyId == "" && t.JournalMode == AutoJournalMode {
		parties := append(a.Parties.CustomersSearchItems(), a.Parties.EmployeesSearchItems()...)
		suggestion, err := t.parseAssociatedParty(t.Description, a.Parties)
		if err == nil && util.AskForConformation(fmt.Sprintf("Use \"%s\" as associeted third party?", suggestion.String())) {
			t.AssociatedPartyId = suggestion.GetId()
		} else {
			var pty interface{}
			t.AssociatedPartyId, pty = util.AskStringFromSearchWithNew(
				"Associated Party",
				"customer/employee which is originator/recipient of the transaction",
				parties,
				InteractiveNewGenericParty,
				a)
			if pty != nil {
				value, ok := pty.(Party)
				if !ok {
					logrus.Fatal("returned party has invalid type")
				}
				if value.PartyType == CustomerType {
					a.Parties.Customers = append(a.Parties.Customers, value)
				} else if value.PartyType == EmployeeType {
					a.Parties.Employees = append(a.Parties.Employees, value)
				}
				t.AssociatedPartyId = value.Id
			}
		}
	}

	document, err := t.parseAssociatedDocument(a.Expenses, a.Invoices)
	if err == nil && util.AskForConformation(fmt.Sprintf("Use \"%s\" as associated document?", document.String())) {
		t.AssociatedDocumentId = document.GetId()
	} else {
		docs := append(a.Expenses.SearchItems(), a.Invoices.SearchItems(a)...)
		t.AssociatedDocumentId = util.AskStringFromSearch(
			"Associated Document",
			"couldn't find associated document, manual search",
			docs)
	}

	strategy := util.AskForStategy()
	switch strategy {
	case util.RedoStrategy:
		t.AssistedCompletion(a, doAll, autoMode, askSkip)
	case util.SkipStrategy:
		return tmp
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

func (t Transaction) parseAssociatedParty(desc string, parties Parties) (Identifiable, error) {
	pty := append(parties.Customers, parties.Employees...)
	for i := range pty {
		if strings.Contains(desc, pty[i].Name) {
			return pty[i], nil
		}
	}
	return nil, fmt.Errorf("no party for description found")
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

func (t Transaction) DateTime() time.Time {
	result, err := time.Parse(util.DateFormat, t.Date)
	if err != nil {
		logrus.Fatalf("could not parse «%s» as date with YYYY-MM-DD: %s", t.Date, err)
	}
	return result
}

// Type returns a string with the type name of the element.
func (Transaction) Type() string {
	return "Transaction"
}

// String returns a human readable representation of the element.
func (t Transaction) String() string {
	if t.TransactionType == util.CreditTransaction {
		return fmt.Sprintf("%s: received %.2f at %s", t.Identifier, t.Amount, t.Date)
	}
	return fmt.Sprintf("%s: payed %.2f at %s", t.Identifier, t.Amount, t.Date)
}

func (t Transaction) JournalDescription(a Acc) string {
	return fmt.Sprintf("TODO %s", t.Description)
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

func (t Transaction) Journal(a Acc, update bool) Journal {
	if t.AssociatedDocumentId != "" {
		exp, err := a.Expenses.ExpenseById(t.AssociatedDocumentId)
		if err == nil {
			return exp.SettlementJournal(a, t, update)
		}
		inv, err := a.Invoices.InvoiceById(t.AssociatedDocumentId)
		if err == nil {
			return inv.SettlementJournal(a, t, update)
		}
	}
	return t.defaultJournal(a)
}

func (t Transaction) defaultJournal(a Acc) Journal {
	var account1, account2 string
	// Incoming transaction
	if t.TransactionType == util.CreditTransaction {
		account1 = a.JournalConfig.BankAccount
		account2 = defaultAccount
	} else {
		account1 = defaultAccount
		account2 = a.JournalConfig.BankAccount
	}
	return Journal{
		{
			Date:        t.DateTime(),
			Status:      UnmarkedStatus,
			Description: t.JournalDescription(a),
			Comment:     NewManualComment("default", t.String()),
			Account1:    account1,
			Account2:    account2,
			Amount:      t.Amount,
		}}
}
