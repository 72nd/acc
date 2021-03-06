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

	if err := validateGetRequest(ctx, params); err != nil {
		return err
	}
	if params.Identifier != nil {
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

// Add a customer
// (POST /customers)
func (e *Endpoint) PostCustomers(ctx echo.Context) error {
	/*
	var cst Customer
	if err := ctx.Bind(&cst); err != nil {
	}
	*/
	return nil
}

// Remove a customer
// (DELETE /customers/{id})
func (e *Endpoint) DeleteCustomersId(ctx echo.Context, id string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	temp := e.schema.Parties.Customers
	for i := range temp {
		if temp[i].Id != id {
			continue
		}
		temp[i] = temp[len(temp)-1]
		temp = temp[:len(temp)-1]
	}
	e.schema.Parties.Customers = temp
	e.schema.Save()
	return onDeleteSuccess(ctx, id, "customer")
}

// Get a customer by ID
// (GET /customers/{id})
func (e *Endpoint) GetCustomersId(ctx echo.Context, id string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	rsl, err := e.schema.Parties.CustomerByRef(schema.NewRef(id))
	if err != nil {
		return onIdNotFound(ctx, id, "customer")
	}

	return ctx.JSON(http.StatusOK, fromAccParty(*rsl))
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

// Add a employee
// (POST /employees)
func (e *Endpoint) PostEmployees(ctx echo.Context) error {
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

// Add a Expense
// (POST /expenses)
func (e *Endpoint) PostExpenses(ctx echo.Context) error {
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

// Add a Invoices
// (POST /invoices)
func (e *Endpoint) PostInvoices(ctx echo.Context) error {
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

// Add a Miscellaneous Record
// (POST /misc_records)
func (e *Endpoint) PostMiscRecords(ctx echo.Context) error {
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

// Add a Projects
// (POST /projects)
func (e *Endpoint) PostProjects(ctx echo.Context) error {
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

// Update a Project
// (PUT /projects/{id})
func (e *Endpoint) PutProjectsId(ctx echo.Context, id string) error {
	return nil
}

// validateGetRequest checks if the user has set both queries (identifier and query) at the
// same time. If so an error message gets logged and also returned to the caller. Otherwise
// the function will return nil.
func validateGetRequest(ctx echo.Context, params GetCustomersParams) error {
	if params.Identifier != nil && params.Query != nil {
		ctx.Logger().Error("using the query and identifier parameters at the same time is forbidden")
		return ctx.String(http.StatusBadRequest, "using the query and identifier parameters at the same time is forbidden")
	}
	return nil
}

// onIdNotFound handles the event when there was no element for an given id. The incident is
// logged and the appropriate HTTP response is given to the callee.
func onIdNotFound(ctx echo.Context, id, typeName string) error {
	msg := fmt.Sprintf("no %s for id %s found", typeName, id)
	ctx.Logger().Error(msg)
	return ctx.String(http.StatusNotFound, msg)
}

// onDeleteSuccess handles the result and logging when a element was removed.
func onDeleteSuccess(ctx echo.Context, id, typeName string) error {
	msg := fmt.Sprintf("%s with id %s successfully deleted", typeName, id)
	ctx.Logger().Debug(msg)
	return ctx.String(http.StatusOK, msg)
}
