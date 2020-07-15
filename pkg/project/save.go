package project

import (
	"fmt"
	"path/filepath"
	"sync"

	"gitlab.com/72th/acc/pkg/schema"
)

type SaveContainer struct {
	wg     sync.WaitGroup
	cst    CustomersToSave
	cstMux sync.Mutex
	prj    ProjectFiles
	prjMux sync.Mutex
}

func (c *SaveContainer) Wait() {
	c.wg.Wait()
}

func (c *SaveContainer) AddCst(cst CustomerToSave) {
	c.wg.Add(1)
	go func() {
		c.cstMux.Lock()
		c.cst = append(c.cst, cst)
		c.cstMux.Unlock()
		c.wg.Done()
	}()
}

func (c *SaveContainer) AddPrj(prj ProjectFile) {
	c.wg.Add(1)
	go func() {
		c.prjMux.Lock()
		c.prj = append(c.prj, prj)
		c.prjMux.Unlock()
		c.wg.Done()
	}()
}

// Save saves a schema as a folder structure to the given folder (path).
// This function is only directly called when converting a project to folder mode.
func Save(s schema.Schema, path string) {
	var wg sync.WaitGroup
	cst := customersToSave(s, filepath.Join(path, projectFolderName))

	wg.Add(1)
	go saveCustomers(path, cst, s.FileHashes, &wg)
	wg.Add(1)
	go saveEmployees(path, s.Parties.Employees, s.FileHashes, &wg)
	wg.Add(1)
	go saveInternalExpenses(path, s.Expenses, s.FileHashes, &wg)
	wg.Wait()
}

// customersToSave transforms the schema into the optimized structure to save the customers
// and their projects in folder mode.
func customersToSave(s schema.Schema, path string) CustomersToSave {
	var wg sync.WaitGroup
	cnt := &SaveContainer{}

	for i := range s.Parties.Customers {
		wg.Add(1)
		cstFolder := filepath.Join(path, folderName(s.Parties.Customers[i].Name))
		go customerToSave(s, s.Parties.Customers[i], cnt, cstFolder, &wg)
	}
	wg.Wait()
	cnt.Wait()

	return cnt.cst
}

// customerToSave builds the CustomerToSave structure for a given customer Party and
// adds it to the channel.
func customerToSave(s schema.Schema, cst schema.Party, cnt *SaveContainer, path string, wg *sync.WaitGroup) {
	var prjWg sync.WaitGroup
	prjCnt := &SaveContainer{}

	for i := range s.Projects {
		prjWg.Add(1)
		prjFolder := filepath.Join(path, folderName(s.Projects[i].Name))
		go projectFile(s, s.Projects[i], prjCnt, prjFolder, &prjWg)
	}
	prjWg.Wait()
	prjCnt.Wait()

	cnt.AddCst(CustomerToSave{
		Customer:     cst,
		ProjectFiles: prjCnt.prj,
	})
	wg.Done()
}

// projectFile builds the ProjectFile structure for a given schema.Project and
// adds it to the channel.
func projectFile(s schema.Schema, prj schema.Project, cnt *SaveContainer, path string, wg *sync.WaitGroup) {
	var exp schema.Expenses
	var inv schema.Invoices

	for i := range s.Expenses {
		fmt.Println(s.Expenses[i])
		if s.Expenses[i].ProjectId == prj.Id {
			e := s.Expenses[i]
			e.Path = relativeAssetPath(path, e.Path)
			fmt.Println(e.Path)
			exp = append(exp, e)
		}
	}
	for i := range s.Invoices {
		if s.Invoices[i].ProjectId == prj.Id {
			inv = append(inv, s.Invoices[i])
		}
	}

	cnt.AddPrj(ProjectFile{
		Project:  prj,
		Expenses: exp,
		Invoices: inv,
	})
	wg.Done()
}

func saveCustomers(path string, cst CustomersToSave, hashes map[string]string, wg *sync.WaitGroup) {
	prjFolder := filepath.Join(path, projectFolderName)
	createNonExistingDir(prjFolder)
	for i := range cst {
		wg.Add(1)
		go saveCustomer(prjFolder, cst[i], hashes, wg)
	}
	wg.Done()
}

func saveCustomer(path string, cst CustomerToSave, hashes map[string]string, wg *sync.WaitGroup) {
	cstFolder := filepath.Join(path, folderName(cst.Customer.Name))
	createNonExistingDir(cstFolder)

	cstFile := filepath.Join(cstFolder, customerFileName)
	schema.SaveYamlOnChange(cst.Customer, cstFile, "customer", hashes[cstFile])

	for i := range cst.ProjectFiles {
		wg.Add(1)
		go saveProject(cstFolder, cst.ProjectFiles[i], hashes, wg)
	}
	wg.Done()
}

func saveProject(path string, prj ProjectFile, hashes map[string]string, wg *sync.WaitGroup) {
	prjFolder := filepath.Join(path, folderName(prj.Project.Name))
	createNonExistingDir(prjFolder)

	prjFile := filepath.Join(prjFolder, projectFileName)
	schema.SaveYamlOnChange(prj, prjFile, "project", hashes[prjFolder])
	wg.Done()
}

func saveInternalExpenses(path string, exp schema.Expenses, hashes map[string]string, wg *sync.WaitGroup) {
	createNonExistingDir(filepath.Join(path, internalFolderName))
	var intExp schema.Expenses
	for i := range exp {
		if exp[i].Internal {
			intExp = append(intExp, exp[i])
		}
	}
	sorted := intExp.SortByYear()

	var expWg sync.WaitGroup
	for k, v := range sorted {
		expWg.Add(1)
		go saveInternalYearExpenses(filepath.Join(path, internalFolderName), k, v, hashes, &expWg)
	}
	expWg.Wait()
	wg.Done()
}

func saveInternalYearExpenses(path string, year int, exp schema.Expenses, hashes map[string]string, wg *sync.WaitGroup) {
	filename := fmt.Sprintf("expenses-%d.yaml", year)
	if year == 0 {
		filename = "expenses-other.yaml"
	}
	expPath := filepath.Join(path, filename)
	schema.SaveYamlOnChange(exp, expPath, "internal expenses", hashes[expPath])
	wg.Done()
}

func saveEmployees(path string, emp []schema.Party, hashes map[string]string, wg *sync.WaitGroup) {
	empPath := filepath.Join(path, employeesFileName)
	schema.SaveYamlOnChange(emp, empPath, "employees", hashes[empPath])
	wg.Done()
}
