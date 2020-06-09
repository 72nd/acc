package ledger

import (
	"fmt"

	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
)

// OCCURRENCE ENTRIES

// EntriesForExpense returns the journal entries for a given schema.Expense.
// Depending on the nature of the expense the appropriate function will be called.
func EntriesForExpense(a schema.Acc, exp schema.Expense) []Entry {
	if exp.AdvancedByThirdParty {
		return entriesForEmployeeAdvancedExpense(a, exp)
	}
	return entriesForCompanyPaidExpenses(a, exp)
}

// entriesForEmployeeAdvancedExpense returns the journal entries for expenses advanced
// by employees.
func entriesForEmployeeAdvancedExpense(a schema.Acc, exp schema.Expense) []Entry {
	cmt := NewComment("employee advanced expense", exp.String())
	cat, err := a.JournalConfig.ExpenseCategories.CategoryByName(exp.ExpenseCategory)
	cmt.add(err)

	acc1 := "no account found"
	if err == nil {
		acc1 = cat.Account
	}

	emp, err := a.Parties.EmployeeById(exp.AdvancedThirdPartyId)
	cmt.add(err)

	desc := "no employee found"
	if err == nil {
		data := map[string]string{
			"Identifier": exp.Identifier,
			"Party":      fmt.Sprintf("%s (%s)", emp.Name, emp.Identifier),
		}
		desc = util.ApplyTemplate(
			"expense advanced by employee description",
			a.JournalConfig.ExpenseAdvancedByEmployeeDescription,
			data)
	}

	return []Entry{
		{
			Date:        exp.AccrualDateTime(),
			Status:      UnmarkedStatus,
			Code:        exp.Identifier,
			Description: desc,
			Comment:     cmt,
			Account1:    acc1,
			Account2:    fmt.Sprintf("%s:%s", a.JournalConfig.EmployeeLiabilitiesAccount, emp.Name),
			Amount:      exp.Amount,
		}}
}

// entriesForCompanyPaidExpenses returns the journal entries for expenses paid by the company itself.
func entriesForCompanyPaidExpenses(a schema.Acc, exp schema.Expense) []Entry {
	cmt := NewComment("company paid expense", exp.String())
	cmt.DoManual = true

	cat, err := a.JournalConfig.ExpenseCategories.CategoryByName(exp.ExpenseCategory)
	cmt.add(err)

	acc1 := "no account found"
	if err == nil {
		acc1 = cat.Account
	}

	var desc string
	if exp.Internal {
		data := map[string]string{
			"Name":       exp.Name,
			"Identifier": exp.Identifier,
		}
		desc = util.ApplyTemplate(
			"internal expense occurence description",
			a.JournalConfig.InternalExpenseOccurenceDescription,
			data)
	} else {
		data := map[string]string{
			"Name":       exp.Name,
			"Identifier": exp.Identifier,
			"Project":    exp.ProjectName,
		}
		desc = util.ApplyTemplate(
			"production expense occurence description",
			a.JournalConfig.ProductionExpenseOccurenceDescription,
			data)
	}
	var acc2 string
	if exp.PayedWithDebit {
		acc2 = a.JournalConfig.BankAccount
	} else {
		acc2 = a.JournalConfig.PayableAccount
	}

	return []Entry{
		{
			Date:        exp.AccrualDateTime(),
			Status:      UnmarkedStatus,
			Code:        exp.Identifier,
			Description: desc,
			Comment:     cmt,
			Account1:    acc1,
			Account2:    acc2,
			Amount:      exp.Amount,
		}}
}

// SETTLEMENT ENTRIES

// SettlementEntriesForExpense takes the related schema.Transaction and schema.Expense and returns
// the settlement entries.
func SettlementEntriesForExpense(a schema.Acc, trn schema.Transaction, exp schema.Expense) []Entry {
	if trn.TransactionType == util.DebitTransaction && exp.AdvancedByThirdParty {
		return settlementEntriesForAdvancedSettlement(a, trn, exp)
	}
	return settlementEntriesForAdvancedSettlement(a, trn, exp)
}

// settlementEntriesForAdvancedSettlement returns the entries for the settlement of an employee advance.
func settlementEntriesForAdvancedSettlement(a schema.Acc, trn schema.Transaction, exp schema.Expense) []Entry {
	cmt := NewComment("settlement of employee advancement", trn.String())
	cmt.add(compareAmounts(trn.Amount, exp.Amount))

	emp, err := a.Parties.EmployeeById(exp.AdvancedThirdPartyId)
	cmt.add(err)

	desc := "no employee found"
	if err == nil {
		data := map[string]string{
			"Identifier": exp.Identifier,
			"Party":      fmt.Sprintf("%s (%s)", emp.Name, emp.Identifier),
		}
		desc = util.ApplyTemplate(
			"employee advanced expense settlement",
			a.JournalConfig.AdvancedExpenseSettlementDescription,
			data)
	}

	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: desc,
			Comment:     cmt,
			Account1:    fmt.Sprintf("%s:%s", a.JournalConfig.EmployeeLiabilitiesAccount, emp.Name),
			Account2:    a.JournalConfig.BankAccount,
			Amount:      trn.Amount,
		}}
}

// settlementEntriesForCompanyPaidExpenses returns the journal entries for the settlement of company
// paid expenses.
func settlementEntriesForCompanyPaidExpenses(a schema.Acc, trn schema.Transaction, exp schema.Expense) []Entry {
	cmt := NewComment("settlement of employee advancement", trn.String())
	cmt.add(compareAmounts(trn.Amount, exp.Amount))

	data := map[string]string{
		"Identifier": exp.Identifier,
	}
	desc := util.ApplyTemplate(
		"company paid expense settlement",
		a.JournalConfig.CompanyPaidExpenseSettlementDescription,
		data)

	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: desc,
			Comment:     cmt,
			Account1:    a.JournalConfig.PayableAccount,
			Account2:    a.JournalConfig.BankAccount,
			Amount:      trn.Amount,
		}}
}
