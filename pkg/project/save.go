package project

import (
	"sync"

	"gitlab.com/72th/acc/pkg/schema"
)

// Save saves a schema as a folder structure to the given folder (path).
// This function is only directly called when converting a project to folder mode.
func Save(s schema.Schema, path string) {
	var wg *sync.WaitGroup
	cst := customersToSave(s)

	wg.Add(1)
	saveCustomers(path, cst, wg)
	wg.Wait()
}

// SaveWithEnv does the same as Save() but uses the `ACC_FOLDER` env variable.
// This is the default use case.
func SaveWithEnv(s schema.Schema) {
	Save(s, repositoryPath())
}

// customersToSave transforms the schema into the optimized structure to save the customers
// and their projects in folder mode.
func customersToSave(s schema.Schema) CustomersToSave {
	var wg *sync.WaitGroup
	cstChan := make(chan CustomerToSave)

	for i := range s.Parties.Customers {
		wg.Add(1)
		go customerToSave(s, s.Parties.Customers[i], cstChan, wg)
	}
	wg.Wait()

	rsl := make(CustomersToSave, len(cstChan))
	i := 0
	for c := range cstChan {
		rsl[i] = c
	}
	return rsl
}

func saveCustomers(path string, cst CustomersToSave, wg *sync.WaitGroup) {
	wg.Done()
}

// customerToSave builds the CustomerToSave structure for a given customer Party and
// adds it to the channel.
func customerToSave(s schema.Schema, cst schema.Party, cstChan chan CustomerToSave, wg *sync.WaitGroup) {
	var prjWg *sync.WaitGroup
	prjChan := make(chan ProjectFile)

	for i := range s.Projects {
		prjWg.Add(1)
		go projectFile(s, s.Projects[i], prjChan, prjWg)
	}
	prjWg.Wait()

	prjFiles := make(ProjectFiles, len(prjChan))
	i := 0
	for p := range prjChan {
		prjFiles[i] = p
	}
	cstChan <- CustomerToSave{
		Customer:     cst,
		ProjectFiles: prjFiles,
	}
	wg.Done()
}

// projectFile builds the ProjectFile structure for a given schema.Project and
// adds it to the channel.
func projectFile(s schema.Schema, prj schema.Project, prjChan chan ProjectFile, wg *sync.WaitGroup) {
	var exp schema.Expenses
	var inv schema.Invoices

	for i := range s.Expenses {
		if s.Expenses[i].ProjectName == prj.Name { // TODO change this to Id
			exp = append(exp, s.Expenses[i])
		}
		// TODO do this also for invoices.
	}

	prjChan <- ProjectFile{
		Project:  prj,
		Expenses: exp,
		Invoices: inv,
	}
	wg.Done()
}
