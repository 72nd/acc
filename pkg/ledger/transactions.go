package ledger

/*
func (t Transaction) JournalEntries(a Acc) []Entry {
	if t.AssociatedDocumentId != "" {
		exp, err := a.Expenses.ExpenseById(t.AssociatedDocumentId)
		if err == nil {
			if exp.Internal {
				return exp.InternalSettlementEntries(a, t)
			}
			return exp.SettlementJournal(a, t, update)
		}
		inv, err := a.Invoices.InvoiceById(t.AssociatedDocumentId)
		if err == nil {
			return inv.SettlementJournal(a, t, update)
		}
	}
	return t.defaultEntries(a)
}

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
