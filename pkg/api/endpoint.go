// REST API Server
package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/72nd/acc/pkg/schema"
	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

// Defines the endpoint for the REST interface.
type Endpoint struct {
	// The Schema to operate on.
	schema *schema.Schema
	// Locks the endpoint thus only one request can change the state at any given moment.
	mutex sync.Mutex
}

// NewEndpoint returns a new endpoint. Takes a Schema as parameter. The request
// received by the server will be applied on this given
// data.
func NewEndpoint(s *schema.Schema) Endpoint {
	return Endpoint{
		schema: s,
	}
}

// Serve Runs the REST endpoint on the given port.
func (e *Endpoint) Serve(port int) {
	echo := echo.New()
	echo.Logger = log.Logger()
	echo.Use(middleware.Logger())

	RegisterHandlers(echo, e)

	if port == 0 {
		port = 8000
	}
	logrus.Warn("Do NOT expose this API to any network or use it with mutliple users")
	echo.Logger.Fatal(echo.Start(fmt.Sprintf("0.0.0.0:%d", port)))
}

// Get all customers
// (GET /customers)
func (e *Endpoint) GetCustomers(ctx echo.Context, params GetCustomersParams) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	var rsl Parties

	if params.Identifier != nil && params.Query != nil {
		ctx.Logger().Error("using the query and identifier parameters at the same time is forbidden")
		return ctx.String(http.StatusBadRequest, "using the query and identifier parameters at the same time is forbidden")
	} else if params.Identifier != nil {
		cst, err := e.schema.Parties.CustomerByIdentifier(*params.Identifier)
		if err != nil {
			msg := fmt.Sprintf("found multiple customers for given ident %sb, please fix this duplication first", *params.Identifier)
			ctx.Logger().Error(msg)
			return ctx.String(
				http.StatusInternalServerError, msg)
		}
		rsl = Parties{fromAccParty(*cst)}
	} else if params.Query != nil {
		items := e.schema.Parties.CustomersSearchItems().Match(*params.Query)
		rsl = make(Parties, len(items))
		for i := range items {
			rsl[i] = fromAccParty(items[i].Element.(schema.Party))
		}
	} else {
		rsl = fromAccParties(e.schema.Parties.Customers)
	}
	return ctx.JSON(http.StatusOK, rsl)
}

// Remove a customer
// (DELETE /customers/{id})
func (e *Endpoint) DeleteCustomersId(ctx echo.Context, id string) error {
	return nil
}

// Get a customer by ID
// (GET /customers/{id})
func (e *Endpoint) GetCustomersId(ctx echo.Context, id string) error {
	return nil
}

// Add a customer
// (POST /customers/{id})
func (e *Endpoint) PostCustomersId(ctx echo.Context, id string) error {
	return nil
}

// Update a customer
// (PUT /customers/{id})
func (e *Endpoint) PutCustomersId(ctx echo.Context, id string) error {
	return nil
}

// Get all employees
// (GET /employees)
func (e *Endpoint) GetEmployees(ctx echo.Context, params GetEmployeesParams) error {
	return nil
}

// Remove a employee
// (DELETE /employees/{id})
func (e *Endpoint) DeleteEmployeesId(ctx echo.Context, id string) error {
	return nil
}

// Get a employee by ID
// (GET /employees/{id})
func (e *Endpoint) GetEmployeesId(ctx echo.Context, id string) error {
	return nil
}

// Add a employee
// (POST /employees/{id})
func (e *Endpoint) PostEmployeesId(ctx echo.Context, id string) error {
	return nil
}

// Update a employee
// (PUT /employees/{id})
func (e *Endpoint) PutEmployeesId(ctx echo.Context, id string) error {
	return nil
}

// Get all Expenses
// (GET /expenses)
func (e *Endpoint) GetExpenses(ctx echo.Context, params GetExpensesParams) error {
	return nil
}

// Remove a Expense
// (DELETE /expenses/{id})
func (e *Endpoint) DeleteExpensesId(ctx echo.Context, id string) error {
	return nil
}

// Get a Expense by ID
// (GET /expenses/{id})
func (e *Endpoint) GetExpensesId(ctx echo.Context, id string) error {
	return nil
}

// Add a Expense
// (POST /expenses/{id})
func (e *Endpoint) PostExpensesId(ctx echo.Context, id string) error {
	return nil
}

// Update a Expense
// (PUT /expenses/{id})
func (e *Endpoint) PutExpensesId(ctx echo.Context, id string) error {
	return nil
}

// Get all invoices
// (GET /invoices)
func (e *Endpoint) GetInvoices(ctx echo.Context, params GetInvoicesParams) error {
	return nil
}

// Remove a Invoice
// (DELETE /invoices/{id})
func (e *Endpoint) DeleteInvoicesId(ctx echo.Context, id string) error {
	return nil
}

// Get a Invoice by ID
// (GET /invoices/{id})
func (e *Endpoint) GetInvoicesId(ctx echo.Context, id string) error {
	return nil
}

// Add a Invoice
// (POST /invoices/{id})
func (e *Endpoint) PostInvoicesId(ctx echo.Context, id string) error {
	return nil
}

// Update a Invoice
// (PUT /invoices/{id})
func (e *Endpoint) PutInvoicesId(ctx echo.Context, id string) error {
	return nil
}

// Get all Miscellaneous Records
// (GET /misc_records)
func (e *Endpoint) GetMiscRecords(ctx echo.Context, params GetMiscRecordsParams) error {
	return nil
}

// Remove a Miscellaneous Record
// (DELETE /misc_records/{id})
func (e *Endpoint) DeleteMiscRecordsId(ctx echo.Context, id string) error {
	return nil
}

// Get a Miscellaneous Record by ID
// (GET /misc_records/{id})
func (e *Endpoint) GetMiscRecordsId(ctx echo.Context, id string) error {
	return nil
}

// Add a Miscellaneous Record
// (POST /misc_records/{id})
func (e *Endpoint) PostMiscRecordsId(ctx echo.Context, id string) error {
	return nil
}

// Update a Miscellaneous Record
// (PUT /misc_records/{id})
func (e *Endpoint) PutMiscRecordsId(ctx echo.Context, id string) error {
	return nil
}

// Get all Projects
// (GET /projects)
func (e *Endpoint) GetProjects(ctx echo.Context, params GetProjectsParams) error {
	return nil
}

// Remove a Project
// (DELETE /projects/{id})
func (e *Endpoint) DeleteProjectsId(ctx echo.Context, id string) error {
	return nil
}

// Get a Project by ID
// (GET /projects/{id})
func (e *Endpoint) GetProjectsId(ctx echo.Context, id string) error {
	return nil
}

// Add a Project
// (POST /projects/{id})
func (e *Endpoint) PostProjectsId(ctx echo.Context, id string) error {
	return nil
}

// Update a Project
// (PUT /projects/{id})
func (e *Endpoint) PutProjectsId(ctx echo.Context, id string) error {
	return nil
}
