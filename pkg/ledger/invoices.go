package ledger

import (
	"fmt"

	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
)

// INVOICING ENTRIES

// EntriesForInvoicing returns the journal entries for invoicing of the given
// schema.Invoice.
func EntriesForInvoicing(a schema.Acc, inv schema.Invoice) []Entry {
	if inv.Revoked {
		return []Entry{}
	}
	cmt := NewComment("invoice sent", inv.String())

	cmp, err := a.Parties.CustomerById(inv.CustomerId)
	cmt.add(err)

	desc := "no customer found"
	if err == nil {
		data := map[string]string{
			"Identifier": inv.Identifier,
			"Party":      fmt.Sprintf("%s (%s)", cmp.Name, cmp.Identifier),
		}
		desc = util.ApplyTemplate(
			"invoice sent transaction description",
			a.JournalConfig.InvoicingTransactionDescription,
			data)
	} 

	return []Entry{
		{
			Date:        inv.SendDateTime(),
			Status:      UnmarkedStatus,
			Description: desc,
			Comment:     cmt,
			Account1:    a.JournalConfig.ReceivableAccount,
			Account2:    a.JournalConfig.RevenueAccount,
			Amount:      inv.Amount,
		}}
}

// SETTLEMENT

// SettlementEntriesForInvoice returns the entries for the settlement (aka receiving the
// money from the customer) of the related invoice.
func SettlementEntriesForInvoice(a schema.Acc, trn schema.Transaction, inv schema.Invoice) []Entry {
	cmt := NewComment("invoice settlement", trn.String())
	cmt.add(compareAmounts(trn.Amount, inv.Amount))

	cmp, err := a.Parties.CustomerById(inv.CustomerId)
	cmt.add(err)

	desc := "TODO no customer found"
	if err == nil {
		data := map[string]string{
			"Identifier": inv.Identifier,
			"Party":      fmt.Sprintf("%s (%s)", cmp.Name, cmp.Identifier),
		}
		desc = util.ApplyTemplate(
			"invoice settlement transaction description",
			a.JournalConfig.InvoiceSettlementTransactionDescription,
			data)
	} 

	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: desc,
			Comment:     cmt,
			Account1:    a.JournalConfig.BankAccount,
			Account2:    a.JournalConfig.ReceivableAccount,
			Amount:      trn.Amount,
		}}
}

