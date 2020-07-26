package ledger

import (
	"fmt"

	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
)

// OCCURRENCE ENTRIES

// EntriesForExpense returns the journal entries for a given schema.Expense.
// Depending on the nature of the expense the appropriate function will be called.
func EntriesForExpense(s schema.Schema, exp schema.Expense) []Entry {
	if exp.AdvancedByThirdParty {
		return entriesForEmployeeAdvancedExpense(s, exp)
	}
	return entriesForCompanyPaidExpenses(s, exp)
}

// entriesForEmployeeAdvancedExpense returns the journal entries for expenses advanced
// by employees.
func entriesForEmployeeAdvancedExpense(s schema.Schema, exp schema.Expense) []Entry {
	cmt := NewComment("employee advanced expense", exp.String())
	cat, err := s.JournalConfig.ExpenseCategories.CategoryByName(exp.ExpenseCategory)
	cmt.add(err)

	acc1 := "no account found"
	if err == nil {
		acc1 = cat.Account
	}

	emp, err := s.Parties.EmployeeByRef(exp.AdvancedThirdParty)
	cmt.add(err)

	desc := "no employee found"
	if err == nil {
		data := map[string]string{
			"Identifier": exp.Identifier,
			"Party":      fmt.Sprintf("%s (%s)", emp.Name, emp.Identifier),
		}
		desc = util.ApplyTemplate(
			"expense advanced by employee description",
			s.JournalConfig.ExpenseAdvancedByEmployeeDescription,
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
			Account2:    fmt.Sprintf("%s:%s", s.JournalConfig.EmployeeLiabilitiesAccount, emp.Name),
			Amount:      exp.Amount,
		}}
}

// entriesForCompanyPaidExpenses returns the journal entries for expenses paid by the company itself.
func entriesForCompanyPaidExpenses(s schema.Schema, exp schema.Expense) []Entry {
	cmt := NewComment("company paid expense", exp.String())
	cmt.DoManual = true

	cat, err := s.JournalConfig.ExpenseCategories.CategoryByName(exp.ExpenseCategory)
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
			"internal expense occurrence description",
			s.JournalConfig.InternalExpenseOccurenceDescription,
			data)
	} else {
		data := map[string]string{
			"Name":       exp.Name,
			"Identifier": exp.Identifier,
			"Project":    exp.ProjectName,
		}
		desc = util.ApplyTemplate(
			"production expense occurrence description",
			s.JournalConfig.ProductionExpenseOccurenceDescription,
			data)
	}
	var acc2 string
	if exp.PaidWithDebit {
		acc2 = s.JournalConfig.BankAccount
	} else {
		acc2 = s.JournalConfig.PayableAccount
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
func SettlementEntriesForExpense(s schema.Schema, trn schema.Transaction, exp schema.Expense) []Entry {
	if trn.TransactionType == util.DebitTransaction && exp.AdvancedByThirdParty {
		return settlementEntriesForAdvancedSettlement(s, trn, exp)
	}
	return settlementEntriesForCompanyPaidExpenses(s, trn, exp)
}

// settlementEntriesForAdvancedSettlement returns the entries for the settlement of an employee advance.
func settlementEntriesForAdvancedSettlement(s schema.Schema, trn schema.Transaction, exp schema.Expense) []Entry {
	cmt := NewComment("settlement of employee advancement", trn.String())
	cmt.add(compareAmounts(trn.Amount, exp.Amount))

	emp, err := s.Parties.EmployeeByRef(exp.AdvancedThirdParty)
	cmt.add(err)

	desc := "no employee found"
	if err == nil {
		data := map[string]string{
			"Identifier": exp.Identifier,
			"Party":      fmt.Sprintf("%s (%s)", emp.Name, emp.Identifier),
		}
		desc = util.ApplyTemplate(
			"employee advanced expense settlement",
			s.JournalConfig.AdvancedExpenseSettlementDescription,
			data)
	}

	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: desc,
			Comment:     cmt,
			Account1:    fmt.Sprintf("%s:%s", s.JournalConfig.EmployeeLiabilitiesAccount, emp.Name),
			Account2:    s.JournalConfig.BankAccount,
			Amount:      trn.Amount,
		}}
}

// settlementEntriesForCompanyPaidExpenses returns the journal entries for the settlement of company
// paid expenses.
func settlementEntriesForCompanyPaidExpenses(s schema.Schema, trn schema.Transaction, exp schema.Expense) []Entry {
	cmt := NewComment("settlement of company paid exense", trn.String())
	cmt.add(compareAmounts(trn.Amount, exp.Amount))

	data := map[string]string{
		"Identifier": exp.Identifier,
	}
	desc := util.ApplyTemplate(
		"company paid expense settlement",
		s.JournalConfig.CompanyPaidExpenseSettlementDescription,
		data)

	return []Entry{
		{
			Date:        trn.DateTime(),
			Status:      UnmarkedStatus,
			Code:        trn.Identifier,
			Description: desc,
			Comment:     cmt,
			Account1:    s.JournalConfig.PayableAccount,
			Account2:    s.JournalConfig.BankAccount,
			Amount:      trn.Amount,
		}}
}
