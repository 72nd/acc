package ledger

import (
	"fmt"

	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
)

// EntriesForTransaction returns the journal entries for a given schema.Transaction.
func EntriesForTransaction(s schema.Schema, trn schema.Transaction) []Entry {
	if trn.AssociatedDocumentId != "" {
		return entriesForTransactionWithDocument(s, trn)
	}
	return entrieForDefaultTransaction(s, trn, nil)
}

// entriesForTransactionWithDocument returns the entries for transactions with an associated
// document.
func entriesForTransactionWithDocument(s schema.Schema, trn schema.Transaction) []Entry {
	exp, err := s.Expenses.ExpenseById(trn.AssociatedDocumentId)
	if err == nil {
		return SettlementEntriesForExpense(s, trn, *exp)
	}
	inv, err := s.Invoices.InvoiceById(trn.AssociatedDocumentId)
	if err == nil {
		return SettlementEntriesForInvoice(s, trn, *inv)
	}
	return entrieForDefaultTransaction(s, trn, fmt.Errorf("no expense/invoice for id \"%s\" found", trn.AssociatedDocumentId))
}

// entrieForDefaultTransaction is the fallback function. It is possible to give an additional
// error as parameter. This error will be appended to the transaction comment.
func entrieForDefaultTransaction(s schema.Schema, trn schema.Transaction, err error) []Entry {
	cmt := NewManualComment("default", trn.String())
	cmt.add(err)

	var acc1, acc2 string
	if trn.TransactionType == util.CreditTransaction {
		// Incoming transaction
		acc1 = s.JournalConfig.BankAccount
		acc2 = defaultAccount
	} else {
		// Outgoing transaction
		acc1 = defaultAccount
		acc2 = s.JournalConfig.BankAccount
	}
	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: fmt.Sprintf("some help: %s", trn.String()),
			Comment:     cmt,
			Account1:    acc1,
			Account2:    acc2,
			Amount:      trn.Amount,
		}}
}

/*

func (t Transaction) defaultEntries(a Acc) []Entry {
	var account1, account2 string
	// Incoming transaction
	if t.TransactionType == util.CreditTransaction {
		account1 = a.JournalConfig.BankAccount
		account2 = defaultAccount
	} else {
		account1 = defaultAccount
		account2 = a.JournalConfig.BankAccount
	}
	return []Entry{
		{
			Date:        t.DateTime(),
			Status:      UnmarkedStatus,
			Code:        t.Identifier,
			Description: t.JournalDescription(a),
			Comment:     NewManualComment("default", t.String()),
			Account1:    account1,
			Account2:    account2,
			Amount:      t.Amount,
		}}
}
*/
