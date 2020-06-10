package providers

import "gitlab.com/72th/acc/pkg/schema"

// AccProvider represents a layer between the data and the way it's saved.
// This is done so a flat-file (all YAML files are in one folder) and a
// customer-project folder-tree is both possible.
type AccProvider interface {
	GetExpenses() schema.Expenses
	GetInvoices() schema.Invoices
	GetMiscRecords() schema.MiscRecords
	GetParties() schema.Parties
	GetProjects() schema.Projects
	GetStatements() schema.Statement

	AddCustomer(cst schema.Party)
	AddEmployee(emp schema.Party)
	AddExpense(exp schema.Expense)
	AddInvoice(inv schema.Invoice)
	AddMiscRecord(mrc schema.MiscRecord)
	AddProject(prj schema.Project)
	AddTransaction(trn schema.Transaction)

	UpdateCustomer(id string, cst schema.Party)
	UpdateEmployee(id string, emp schema.Party)
	UpdateExpense(id string, exp schema.Expense)
	UpdateInvoice(id string, inv schema.Invoice)
	UpdateMiscRecord(id string, mrc schema.MiscRecord)
	UpdateProject(id string, prj schema.Project)
	UpdateTransaction(id string, trn schema.Transaction)
}
