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

type OpenContainer struct {
	wg       sync.WaitGroup
	cst      []schema.Party
	cstMux   sync.Mutex
	emp      []schema.Party
	empMux   sync.Mutex
	exp      []schema.Expense
	expMux   sync.Mutex
	prj      ProjectFiles
	prjMux   sync.Mutex
	files    map[string]string
	filesMux sync.Mutex
}

func (c *OpenContainer) Wait() {
	c.wg.Wait()
}

func (c *OpenContainer) AddCst(cst schema.Party) {
	c.wg.Add(1)
	go func() {
		c.cstMux.Lock()
		c.cst = append(c.cst, cst)
		c.cstMux.Unlock()
		c.wg.Done()
	}()
}

func (c *OpenContainer) AddEmp(emp []schema.Party) {
	c.wg.Add(1)
	go func() {
		c.empMux.Lock()
		c.emp = append(c.emp, emp...)
		c.empMux.Unlock()
		c.wg.Done()
	}()
}

func (c *OpenContainer) AddExp(exp schema.Expenses) {
	c.wg.Add(1)
	go func() {
		c.expMux.Lock()
		c.exp = append(c.exp, exp...)
		c.expMux.Unlock()
		c.wg.Done()
	}()
}

func (c *OpenContainer) AddPrj(prj ProjectFiles) {
	c.wg.Add(1)
	go func() {
		c.prjMux.Lock()
		c.prj = append(c.prj, prj...)
		c.prjMux.Unlock()
		c.wg.Done()
	}()
}

func (c *OpenContainer) AddFile(file StrTuple) {
	c.wg.Add(1)
	go func() {
		c.filesMux.Lock()
		c.files[file[0]] = file[1]
		c.filesMux.Unlock()
		c.wg.Done()
	}()
}

// Open loads the schema for the project mode.
func Open(path string, cmp schema.Company, jfg schema.JournalConfig, saveFunc func(schema.Schema)) schema.Schema {
	cnt := &OpenContainer{}
	cnt.files = make(map[string]string)

	var wg sync.WaitGroup
	wg.Add(1)
	go openCustomersProjects(path, cnt, &wg)
	wg.Add(1)
	go openInternalExpenses(path, cnt, &wg)
	wg.Add(1)
	go openEmployeeFile(path, cnt, &wg)
	cnt.Wait()
	wg.Wait()

	return schema.Schema{
		Company:       cmp,
		Expenses:      append(cnt.exp, cnt.prj.Expenses()...),
		Invoices:      cnt.prj.Invoices(),
		JournalConfig: jfg,
		Parties: schema.Parties{
			Customers: cnt.cst,
			Employees: cnt.emp,
		},
		Projects:   cnt.prj.Projects(),
		FileHashes: cnt.files,
		SaveFunc:   saveFunc,
	}
}

// openCustomersProjects walks the given projects folder (path) and returns all found customers and projects.
func openCustomersProjects(path string, cnt *OpenContainer, wg *sync.WaitGroup) {
	path = filepath.Join(path, "projects")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logrus.Fatalf("projects folder in acc repository doesn't exist, expected path: %s", path)
	}
	folders := getFoldersInPath(path)

	for i := range folders {
		wg.Add(1)
		go customerWalk(folders[i], cnt, wg)
	}
	wg.Done()
}

// customerWalk goes trough one customer folder and puts the customer and all found projects into channels.
func customerWalk(path string, cnt *OpenContainer, wg *sync.WaitGroup) {
	wg.Add(1)
	go openCustomerFile(path, cnt, wg)

	folders := getFoldersInPath(path)
	for i := range folders {
		wg.Add(1)
		go openProjectFile(folders[i], cnt, wg)
	}

	wg.Done()
}

// openCustomerFile tries to open the `customer.yaml` file in the given folder path.
// If the file exists it will be parsed and the customer get added to the customer channel.
func openCustomerFile(path string, cnt *OpenContainer, wg *sync.WaitGroup) {
	cstFile := filepath.Join(path, customerFileName)
	if _, err := os.Stat(cstFile); os.IsNotExist(err) {
		logrus.Errorf("the %s file does not exist in %s", customerFileName, path)
	} else {
		var cst schema.Party
		hash := schema.OpenYamlHashed(&cst, cstFile, "customer file")
		cnt.AddFile(StrTuple{cstFile, hash})
		/*
		raw, err := ioutil.ReadFile(cstFile)
		if err != nil {
			logrus.Fatal(err)
		}
		if err := yaml.Unmarshal(raw, &cst); err != nil {
			logrus.Fatal(err)
		}
		*/
		cnt.AddCst(cst)
	}
	wg.Done()
}

// openProjectFile tries to open the `project.yaml` file in the given folder path.
// If the file exists it will be parsed and the project get added to the project channel.
func openProjectFile(path string, cnt *OpenContainer, wg *sync.WaitGroup) {
	prjFile := filepath.Join(path, projectFileName)
	if _, err := os.Stat(prjFile); os.IsNotExist(err) {
		logrus.Errorf("the %s file does not exist in %s", projectFileName, path)
	} else {
		var prj ProjectFile
		hash := schema.OpenYamlHashed(&prj, prjFile, "project file")
		cnt.AddFile(StrTuple{prjFile, hash})
		cnt.AddPrj(ProjectFiles{prj.AbsolutePaths(path)})
	}
	wg.Done()
}

// openInternalExpenses opens the internal expenses in the `internal` folder.
func openInternalExpenses(path string, container *OpenContainer, wg *sync.WaitGroup) {
	intFolder := filepath.Join(path, internalFolderName)
	if _, err := os.Stat(intFolder); os.IsNotExist(err) {
		logrus.Errorf("the %s folder does not exist in %s", internalFolderName, path)
		wg.Done()
		return
	}
	files := getMatchingFilesInPath(intFolder, regexp.MustCompile(`expenses-2\d\d\d\.yaml`))

	otherExpPath := filepath.Join(intFolder, "expenses-other.yaml")
	if _, err := os.Stat(otherExpPath); !os.IsNotExist(err) {
		files = append(files, otherExpPath)
	}
	for i := range files {
		wg.Add(1)
		go openExpenseFile(files[i], container, wg)
	}
	wg.Done()
}

// openExpenseFile opens an expense file by the given path and adds the expenses into to channel.
func openExpenseFile(path string, container *OpenContainer, wg *sync.WaitGroup) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logrus.Errorf("the internal expense file \"%s\" does not exists", path)
		wg.Done()
		return
	}
	var exp schema.Expenses
	util.OpenYaml(&exp, path, "internal expenses file")
	container.AddExp(exp)
	wg.Done()
}

// openExpenseFile opens the employee file by the given path and adds the employees into the channel.
func openEmployeeFile(path string, cnt *OpenContainer, wg *sync.WaitGroup) {
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
	cnt.AddEmp(emp)
	wg.Done()
}
