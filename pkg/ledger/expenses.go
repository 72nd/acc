package ledger

import (
	"fmt"

	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
)

// INITIAL ENTRIES

// EntriesForExpense returns the journal entries for a given schema.Expense.
// Depending on the nature of the expense the appropriate function will be called.
func EntriesForExpense(a schema.Acc, exp schema.Expense) []Entry {
	if exp.AdvancedByThirdParty {
		return entriesForEmployeeAdvancedExpense(a, exp)
	} else if exp.Internal {
		return entriesForInternalExpense(a, exp)
	}
	return entriesForProductionExpense(a, exp)

	ele := "default expense"
	if e.AdvancedByThirdParty {
		ele = "employee advanced expense"
	}

	cmt := NewComment(ele, e.String())
	acc1, err := expenseAccount(a)
	cmt.add(err)
	acc2 := a.JournalConfig.PayableAccount
	var desc string
	if e.AdvancedByThirdParty {
		acc2, err = e.employeeLiabilityAccount(a)
		cmt.add(err)
		desc = e.employeeAdvancedDescription(a)
	} else {
		desc = fmt.Sprintf("Einkauf von %s (%s)", e.Name, e.Identifier)
		cmt.DoManual = true
	}

	return []Entry{
		{
			Date:        e.AccrualDateTime(),
			Status:      UnmarkedStatus,
			Description: desc,
			Comment:     cmt,
			Account1:    acc1,
			Account2:    acc2,
			Amount:      e.Amount,
		}}
}

// entriesForEmployeeAdvancedExpense returns the journal entries for expenses advanced
// by employees.
func entriesForEmployeeAdvancedExpense(a schema.Acc, exp schema.Expense) []Entry {
	cmt := NewComment("employee advanced expense", exp.String())
	cat, err := a.JournalConfig.ExpenseCategories.CategoryByName(exp.ExpenseCategory)
	cmt.add(err)

	emp, err := a.Parties.EmployeeById(exp.AdvancedThirdPartyId)
	cmt.add(err)

	desc := "no employee found"
	if err == nil {
		data := map[string]string{
			"Identifier": exp.Identifier,
			"Party":      fmt.Sprintf("%s (%s)", emp.Name, emp.Identifier),
		}
		desc = util.ApplyTemplate("expense advanced by employee description", a.JournalConfig.ExpenseAdvancedByEmployeeDescription, data)
	}

	return []Entry{
		{
			Date:        exp.AccrualDateTime(),
			Status:      UnmarkedStatus,
			Code:        exp.Identifier,
			Description: desc,
			Comment:     cmt,
			Account1:    cat.Account,
			Account2:    fmt.Sprintf("%s:%s", a.JournalConfig.EmployeeLiabilitiesAccount, emp.Name),
			Amount:      exp.Amount,
		}}
}

// entriesForInternalExpense returns the journal entries for internal expenses. Example:
// buying a printer for the companies office.
func entriesForInternalExpense(a schema.Acc, exp schema.Expense) []Entry {
}

// entriesForProductionExpense returns the journal entries for expenses made
// for customer projects.
func entriesForProductionExpense(a schema.Acc, exp schema.Expense) []Entry {
}

// SETTLEMENT ENTRIES

// SettlementEntriesFromExpense takes the related schema.Transaction and schema.Expense and returns
// the settlement entries.
func SettlementEntriesFromExpense(a schema.Acc, trn schema.Transaction, exp schema.Expense) []Entry {
}

// settlementEntriesForAdvancedSettlement returns the entries for the settlement of an employee advance.
func settlementEntriesForAdvancedSettlement(a schema.Acc, trn schema.Transaction, exp schema.Expense) []Entry {
}

// settlementEntriesForInternalExpense returns the journal entries for the settlement of an internal
// expense (aka paying the bill).
func settlementEntriesForInternalExpense(a schema.Acc, trn schema.Transaction, exp schema.Expense) []Entry {
}

// settlementEntriesForProductionExpense returns the entries for the settlement of an production expense.
// Aka. for paying the bill.
func settlementEntriesForProductionExpense() []Entry {
}

// TODO: Settlement für nicht AdvancedByThirdParty transactions
func (e *Expense) SettlementJournal(a Acc, trn Transaction, update bool) []Entry {
	cmt := NewComment("advanced expense settlement", trn.String())
	acc1, err := e.expenseAccount(a)
	cmt.add(err)
	acc2, err := e.employeeLiabilityAccount(a)
	cmt.add(err)
	if trn.Amount != e.Amount {
		cmt.add(fmt.Errorf("amount of transaction (%.2f) doesn't match amount of colligated expense %s", trn.Amount, e.String()))
	}

	if update {
		e.DateOfSettlement = trn.Date
		e.SettlementTransactionId = trn.Id
	}

	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: "TODO expense employee booking description",
			Comment:     cmt,
			Account1:    acc1,
			Account2:    acc2,
			Amount:      trn.Amount,
		}}

}

func (e Expense) InternalSettlementEntries(a Acc, trn Transaction) []Entry {
	cmt := NewComment("internal expense settlement", trn.String())
	if trn.Amount != e.Amount {
		cmt.add(fmt.Errorf("amount of transaction (%.2f) doesn't match amount of colligated expense %s", trn.Amount, e.String()))
	}
	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: e.internalSettlementDescription(a),
			Comment:     cmt,
			Account1:    a.JournalConfig.PayableAccount,
			Account2:    a.JournalConfig.BankAccount,
			Amount:      trn.Amount,
		}}
}

func (e Expense) internalSettlementDescription(a Acc) string {
	data := map[string]string{
		"Identifier": e.Identifier,
	}
	return util.ApplyTemplate("internal expense settlement description", a.JournalConfig.InternalExpenseTransactionDescription, data)
}

// HELPERS
