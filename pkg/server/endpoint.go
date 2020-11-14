// REST API Server
package server

import (
	"fmt"
	"sync"

	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/server/api"
	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/neko-neko/echo-logrus/v2/log"
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
	/*
		swagger, err := GetSwagger()
		if err != nil {
			logrus.Fatalf("error loading OpenAPI spec: %s", err)
		}
		swagger.Servers = nil
	*/
	echo := echo.New()
	echo.Logger = log.Logger()
	echo.Use(middleware.Logger())

	api.RegisterHandlers(echo, e)

	if port == 0 {
		port = 8000
	}
	echo.Logger.Fatal(echo.Start(fmt.Sprintf("0.0.0.0:%d", port)))
}

// Get all customers
// (GET /customers)
func (e *Endpoint) GetCustomers(ctx echo.Context, params api.GetCustomersParams) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()



	return nil
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
func (e *Endpoint) GetEmployees(ctx echo.Context, params api.GetEmployeesParams) error {
	return nil
}

// Get all Expenses
// (GET /expenses)
func (e *Endpoint) GetExpenses(ctx echo.Context, params api.GetExpensesParams) error {
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
func (e *Endpoint) GetInvoices(ctx echo.Context, params api.GetInvoicesParams) error {
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
func (e *Endpoint) GetMiscRecords(ctx echo.Context, params api.GetMiscRecordsParams) error {
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
func (e *Endpoint) GetProjects(ctx echo.Context, params api.GetProjectsParams) error {
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
