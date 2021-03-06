package distributed

import (
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
	"github.com/sirupsen/logrus"
)

/*
 * TODO's
 * - Import from flat file
 * - Internal Expenses
 */

/*
type sync.WaitGroup struct {
	Wg      *sync.WaitGroup
	Name    string
	Counter int
	Open    []string
}

func Newsync.WaitGroup(name string) sync.WaitGroup {
	return sync.WaitGroup{
		Wg:   &sync.WaitGroup{},
		Name: name}
}

func (w *sync.WaitGroup) Add(n int, location string) {
	w.Wg.Add(n)
	w.Counter++
	w.Open = append(w.Open, location)
	logrus.Debugf("+ wg %s %s -> %d", w.Name, location, w.Counter)
}

func (w *sync.WaitGroup) Done(location string) {
	logrus.Debugf("- wg %s %s -> %d", w.Name, location, w.Counter)
	w.Counter--
	for i := range w.Open {
		if w.Open[i] == location {
			w.Open = append(w.Open[:i], w.Open[i+1:]...)
			break
		}
	}
	w.Wg.Done()
}

func (w sync.WaitGroup) Wait(location string) {
	logrus.Debugf("w wg %s %s -> %d", w.Name, location, w.Counter)
	w.Wg.Wait()
	logrus.Debugf("wait done wg %s %s -> %d", w.Name, location, w.Counter)
	w.PrintOpen()
}

func (w sync.WaitGroup) PrintOpen() {
	rsl := ""
	for _, o := range w.Open {
		rsl = fmt.Sprintf("%s, %s", rsl, o)
	}
	logrus.Debug(rsl)
}
*/

// StrTuple represents a tuple of two strings.
type StrTuple []string

// OpenContainer encapsulates all data slices which have to be open in a concurrency-safe manner.
// Multiple go-routines can add data to the same collection without having to wait.
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

// Wait has to be called when the data of the container should be read. The function will wait
// until all pending writes are done.
func (c *OpenContainer) Wait() {
	c.wg.Wait()
}

// AddCst adds a customer to the customer slice.
func (c *OpenContainer) AddCst(cst schema.Party) {
	c.wg.Add(1)
	go func() {
		c.cstMux.Lock()
		c.cst = append(c.cst, cst)
		c.cstMux.Unlock()
		c.wg.Done()
	}()
}

// AddEmp adds a employee to the employee slice.
func (c *OpenContainer) AddEmp(emp []schema.Party) {
	c.wg.Add(1)
	go func() {
		c.empMux.Lock()
		c.emp = append(c.emp, emp...)
		c.empMux.Unlock()
		c.wg.Done()
	}()
}

// AddExp adds a expense to the expense slice.
func (c *OpenContainer) AddExp(exp schema.Expenses) {
	c.wg.Add(1)
	go func() {
		c.expMux.Lock()
		c.exp = append(c.exp, exp...)
		c.expMux.Unlock()
		c.wg.Done()
	}()
}

// AddPrj adds a project file to the project file slice.
func (c *OpenContainer) AddPrj(prj ProjectFiles) {
	c.wg.Add(1)
	go func() {
		c.prjMux.Lock()
		c.prj = append(c.prj, prj...)
		c.prjMux.Unlock()
		c.wg.Done()
	}()
}

// AddFile adds a Tuple containing the path and the hash of it's content to the files map.
func (c *OpenContainer) AddFile(file StrTuple) {
	c.wg.Add(1)
	go func() {
		c.filesMux.Lock()
		c.files[file[0]] = file[1]
		c.filesMux.Unlock()
		c.wg.Done()
	}()
}

// Open loads the schema for the distributed mode.
func Open(path string, cmp schema.Company, jfg schema.JournalConfig, saveFunc func(schema.Schema), currency string) schema.Schema {
	cnt := &OpenContainer{}
	cnt.files = make(map[string]string)

	// wg := Newsync.WaitGroup("open")
	var wg sync.WaitGroup
	wg.Add(1)
	// wg.Add(1, "openCustomersProjects")
	go openCustomersProjects(path, cnt, &wg)
	wg.Add(1)
	// wg.Add(1, "openInternalExpenses")
	go openInternalExpenses(path, cnt, &wg)
	wg.Add(1)
	// wg.Add(1, "openEmployeeFile")
	go openEmployeeFile(path, cnt, &wg)
	wg.Wait()
	// wg.Wait("open")
	cnt.Wait()

	return schema.Schema{
		Currency:      currency,
		Company:       cmp,
		Expenses:      append(cnt.exp, cnt.prj.Expenses()...),
		Invoices:      cnt.prj.Invoices(),
		JournalConfig: jfg,
		Parties: schema.PartiesCollection{
			Customers: cnt.cst,
			Employees: cnt.emp,
		},
		Projects:   cnt.prj.Projects(),
		FileHashes: cnt.files,
		SaveFunc:   saveFunc,
		BaseFolder: path,
	}
}

// openCustomersProjects walks the given projects folder (path) and returns all found customers and projects.
func openCustomersProjects(path string, cnt *OpenContainer, wg *sync.WaitGroup) {
	path = filepath.Join(path, projectFolderName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logrus.Fatalf("projects folder in acc repository doesn't exist, expected path: %s", path)
	}
	folders := getFoldersInPath(path)

	for i := range folders {
		wg.Add(1)
		// wg.Add(1, "customerWalk")
		go customerWalk(folders[i], cnt, wg)
	}
	wg.Done()
	// wg.Done("openCustomersProjects")
}

// customerWalk goes trough one customer folder and puts the customer and all found projects into channels.
func customerWalk(path string, cnt *OpenContainer, wg *sync.WaitGroup) {
	wg.Add(1)
	// wg.Add(1, "openCustomerFile")
	go openCustomerFile(path, cnt, wg)

	folders := getFoldersInPath(path)
	for i := range folders {
		wg.Add(1)
		// wg.Add(1, "openProjectFile")
		go openProjectFile(folders[i], cnt, wg)
	}
	wg.Done()
	// wg.Done("customerWalk")
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
		cnt.AddCst(cst)
	}
	wg.Done()
	// wg.Done("openCustomerFile")
}

// openProjectFile tries to open the `project.yaml` file in the given folder path.
// If the file exists it will be parsed and the project get added to the container.
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
	// wg.Done("openProjectFile")
}

// openInternalExpenses opens the internal expenses in the `internal` folder.
func openInternalExpenses(path string, container *OpenContainer, wg *sync.WaitGroup) {
	intFolder := filepath.Join(path, internalFolderName)
	if _, err := os.Stat(intFolder); os.IsNotExist(err) {
		logrus.Errorf("the %s folder does not exist in %s", internalFolderName, path)
		wg.Done()
		// wg.Done("openInternalExpenses")
		return
	}
	files := getMatchingFilesInPath(intFolder, regexp.MustCompile(`expenses-2\d\d\d\.yaml`))

	otherExpPath := filepath.Join(intFolder, "expenses-other.yaml")
	if _, err := os.Stat(otherExpPath); !os.IsNotExist(err) {
		files = append(files, otherExpPath)
	}
	for i := range files {
		wg.Add(1)
		// wg.Add(1, "openExpenseFile")
		go openExpenseFile(files[i], container, wg)
	}
	wg.Done()
	// wg.Done("openInternalExpenses")
}

// openExpenseFile opens an expense file by the given path and adds the expenses into to channel.
func openExpenseFile(path string, container *OpenContainer, wg *sync.WaitGroup) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logrus.Errorf("the internal expense file \"%s\" does not exists", path)
		wg.Done()
		// wg.Done("openExpenseFile")
		return
	}
	var exp schema.Expenses
	util.OpenYaml(&exp, path, "internal expenses file")
	container.AddExp(exp)
	wg.Done()
	// wg.Done("openExpenseFile")
}

// openExpenseFile opens the employee file by the given path and adds the employees into the channel.
func openEmployeeFile(path string, cnt *OpenContainer, wg *sync.WaitGroup) {
	empPath := filepath.Join(path, employeesFileName)
	if _, err := os.Stat(empPath); os.IsNotExist(err) {
		logrus.Errorf("the %s file does not exist in %s", employeesFileName, path)
		wg.Done()
		// wg.Done("openEmployeeFile")
		return
	}
	var emp []schema.Party
	hash := schema.OpenYamlHashed(&emp, empPath, "employee file")
	cnt.AddFile(StrTuple{empPath, hash})
	if len(emp) == 0 {
		wg.Done()
		// wg.Done("openEmployeeFile")
		return
	}
	cnt.AddEmp(emp)
	wg.Done()
	// wg.Done("openEmployeeFile")
}
