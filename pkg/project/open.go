package project

import (
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
)

/*
 * TODO's
 * - Project structure done
 * - Referencing to projects is done via id
 * - Generate ProjectFiles from Project
 * - Save all the stuff
 * - Import from flat file
 */

type StrTuple []string

// Open loads the schema for the project mode.
func Open(path string, cmp schema.Company, jfg schema.JournalConfig, saveFunc func(schema.Schema)) schema.Schema {
	var wg sync.WaitGroup
	cstChan := make(chan schema.Party)
	empChan := make(chan schema.Party)
	expChan := make(chan schema.Expense)
	prjChan := make(chan ProjectFile)
	fileChan := make(chan StrTuple)
	wg.Add(1)
	go openCustomersProjects(path, cstChan, prjChan, fileChan, &wg)
	wg.Add(1)
	go openInternalExpenses(path, expChan, fileChan, &wg)
	wg.Add(1)
	go openEmployeeFile(path, empChan, fileChan, &wg)
	wg.Wait()

	close(cstChan)
	close(empChan)
	close(expChan)
	close(prjChan)
	close(fileChan)

	cst := make([]schema.Party, len(cstChan))
	i := 0
	for c := range cstChan {
		cst[i] = c
		i++
	}
	emp := make([]schema.Party, len(empChan))
	i = 0
	for e := range empChan {
		emp[i] = e
		i++
	}
	exp := make(schema.Expenses, len(expChan))
	i = 0
	for e := range expChan {
		exp[i] = e
		i++
	}
	prj := make(ProjectFiles, len(prjChan))
	i = 0
	for p := range prjChan {
		prj[i] = p
		i++
	}
	var files map[string]string
	for f := range fileChan {
		files[f[0]] = f[1]
	}

	return schema.Schema{
		Company:       cmp,
		Expenses:      append(exp, prj.Expenses()...),
		Invoices:      prj.Invoices(),
		JournalConfig: jfg,
		Parties: schema.Parties{
			Customers: cst,
			Employees: emp,
		},
		Projects:   prj.Projects(),
		FileHashes: files,
		SaveFunc:   saveFunc,
	}
}

// openCustomersProjects walks the given projects folder (path) and returns all found customers and projects.
func openCustomersProjects(path string, cstChan chan schema.Party, prjChan chan ProjectFile, fileChan chan StrTuple, wg *sync.WaitGroup) {
	path = filepath.Join(path, "projects")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logrus.Fatalf("projects folder in acc repository doesn't exist, expected path: %s", path)
	}
	folders := getFoldersInPath(path)

	for i := range folders {
		wg.Add(1)
		go customerWalk(folders[i], cstChan, prjChan, fileChan, wg)
	}
	wg.Done()
}

// customerWalk goes trough one customer folder and puts the customer and all found projects into channels.
func customerWalk(path string, cstChan chan schema.Party, prjChan chan ProjectFile, fileChan chan StrTuple, wg *sync.WaitGroup) {
	wg.Add(1)
	go openCustomerFile(path, cstChan, fileChan, wg)

	folders := getFoldersInPath(path)
	for i := range folders {
		wg.Add(1)
		go openProjectFile(folders[i], prjChan, fileChan, wg)
	}

	wg.Done()
}

// openCustomerFile tries to open the `customer.yaml` file in the given folder path.
// If the file exists it will be parsed and the customer get added to the customer channel.
func openCustomerFile(path string, cstChan chan schema.Party, fileChan chan StrTuple, wg *sync.WaitGroup) {
	cstFile := filepath.Join(path, customerFileName)
	if _, err := os.Stat(cstFile); os.IsNotExist(err) {
		logrus.Errorf("the %s file does not exist in %s", customerFileName, path)
	} else {
		var cst schema.Party
		hash := schema.OpenYamlHashed(&cst, cstFile, "customer file")
		fileChan <- StrTuple{cstFile, hash}
		cstChan <- cst
	}
	wg.Done()
}

// openProjectFile tries to open the `project.yaml` file in the given folder path.
// If the file exists it will be parsed and the project get added to the project channel.
func openProjectFile(path string, prjChan chan ProjectFile, fileChan chan StrTuple, wg *sync.WaitGroup) {
	prjFile := filepath.Join(path, projectFileName)
	if _, err := os.Stat(prjFile); os.IsNotExist(err) {
		logrus.Errorf("the %s file does not exist in %s", projectFileName, path)
	} else {
		var prj ProjectFile
		hash := schema.OpenYamlHashed(&prj, prjFile, "project file")
		fileChan <- StrTuple{prjFile, hash}
		prjChan <- prj.AbsolutePaths(path)
	}
	wg.Done()
}

// openInternalExpenses opens the internal expenses in the `internal` folder.
func openInternalExpenses(path string, expChan chan schema.Expense, fileChan chan StrTuple, wg *sync.WaitGroup) {
	intFolder := filepath.Join(path, internalFolderName)
	if _, err := os.Stat(intFolder); os.IsNotExist(err) {
		logrus.Errorf("the %s folder does not exist in %s", internalFolderName, path)
		wg.Done()
		return
	}
	files := getMatchingFilesInPath(intFolder, regexp.MustCompile(`expenses-2\d\d\d\.yaml`))
	for i := range files {
		wg.Add(1)
		go openExpenseFile(files[i], expChan, fileChan, wg)
	}
	wg.Done()
}

// openExpenseFile opens an expense file by the given path and adds the expenses into to channel.
func openExpenseFile(path string, expChan chan schema.Expense, fileChan chan StrTuple, wg *sync.WaitGroup) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logrus.Errorf("the internal expense file \"%s\" does not exists", path)
		wg.Done()
		return
	}
	var exp schema.Expenses
	util.OpenYaml(&exp, path, "internal expenses file")
	for i := range exp {
		expChan <- exp[i]
	}
	wg.Done()
}

// openExpenseFile opens the employee file by the given path and adds the employees into the channel.
func openEmployeeFile(path string, empChan chan schema.Party, fileChan chan StrTuple, wg *sync.WaitGroup) {
	empPath := filepath.Join(path, employeesFileName)
	if _, err := os.Stat(empPath); os.IsNotExist(err) {
		logrus.Errorf("the %s file does not exist in %s", employeesFileName, path)
		wg.Done()
		return
	}
	var emp []schema.Party
	_ = schema.OpenYamlHashed(&emp, empPath, "employee file")
	// fileChan <- StrTuple{empPath, hash}
	if len(emp) == 0 {
		wg.Done()
		return
	}
	for i := range emp {
		empChan <- emp[i]
	}
	wg.Done()
}
