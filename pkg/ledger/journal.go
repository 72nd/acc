package ledger

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
)

const HLedgerDateFormat = "2006-01-02"

type Journal []Entry

func JournalFromStatement(acc schema.Acc) Journal {
	var result Journal
	for i := range acc.BankStatement.Transactions {
		result = append(result, journalFromTransaction(acc, acc.BankStatement.Transactions[i])...)
	}
	return result
}

func (j Journal) SaveHLedgerFile(path string) {
	ledger := j.HLedger()
	if err := ioutil.WriteFile(path, []byte(ledger), 0644); err != nil {
		logrus.Fatalf("error writing file %s: %s", path, err)
	}
}

func (j Journal) HLedger() string {
	var result string
	for i := range j {
		result = fmt.Sprintf("%s\n\n%s", result, j[i].Transaction())
	}
	return result
}

func journalFromTransaction(acc schema.Acc, trn schema.Transaction) Journal {
	if trn.AssociatedDocumentId != "" {
		exp, err := acc.Expenses.ExpenseById(trn.AssociatedDocumentId)
		if err == nil {
			return journalFromExpense(*exp, trn)
		}
		inv, err := acc.Invoices.InvoiceById(trn.AssociatedDocumentId)
		if err == nil {
			return journalFromInvoice(*inv, trn)
		}
	}
	return defaultJournal(acc, trn)
}

func defaultJournal(acc schema.Acc, trn schema.Transaction) Journal {
	return Journal{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Description: trn.JournalDescription(acc),
			Comment:     "TODO: manual correction needed",
			Account1:    acc.BankAccount,
			Account2:    "other:unknow",
			Amount:      trn.Amount,
		}}
}

func journalFromExpense(exp schema.Expense, trn schema.Transaction) Journal {
	return Journal{}
}

func journalFromInvoice(inv schema.Invoice, trn schema.Transaction) Journal {
	return Journal{}
}
