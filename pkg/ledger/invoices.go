package ledger

/*
func (i Invoice) Journal(a Acc) []Entry {
	cmt := NewComment("invoice sent", i.String())
	return []Entry{
		{
			Date:        i.SendDateTime(),
			Status:      UnmarkedStatus,
			Description: i.sendTransactionDescription(a),
			Comment:     cmt,
			Account1:    a.JournalConfig.ReceivableAccount,
			Account2:    a.JournalConfig.RevenueAccount,
			Amount:      i.Amount,
		}}
}

func (i Invoice) sendTransactionDescription(a Acc) string {
	cmp, err := a.Parties.CustomerById(i.CustomerId)
	if err != nil {
		logrus.Error("no customer found: ", err)
	}
	data := map[string]string{
		"Identifier": i.Identifier,
		"Party":      fmt.Sprintf("%s (%s)", cmp.Name, cmp.Identifier),
	}
	return util.ApplyTemplate("invoice transaction description", a.JournalConfig.InvoicingTransactionDescription, data)
}

func (i Invoice) SettlementJournal(a Acc, trn Transaction, update bool) []Entry {
	cmt := NewComment("invoice settlement", trn.String())
	if trn.Amount != i.Amount {
		cmt.add(fmt.Errorf("amount of transaction (%.2f) doesn't match amount of colligated invoice %s", trn.Amount, i.String()))
	}
	if update {
		i.DateOfSettlement = trn.Date
		i.SettlementTransactionId = trn.Id
	}

	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: i.settlementTransactionDescription(a),
			Comment:     cmt,
			Account1:    a.JournalConfig.BankAccount,
			Account2:    a.JournalConfig.ReceivableAccount,
			Amount:      trn.Amount,
		}}
}

func (i Invoice) settlementTransactionDescription(a Acc) string {
	cmp, err := a.Parties.CustomerById(i.CustomerId)
	if err != nil {
		logrus.Error("no customer found: ", err)
	}
	data := map[string]string{
		"Identifier": i.Identifier,
		"Party":      fmt.Sprintf("%s (%s)", cmp.Name, cmp.Identifier),
	}
	return util.ApplyTemplate("invoice settlement transaction description", a.JournalConfig.InvoiceSettlementTransactionDescription, data)
}
*/
