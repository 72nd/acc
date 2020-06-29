package folder

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
)

// Open loads the schema for the project mode.
func Open() schema.Schema {
	prt, prj := openCustomersProjects(repositoryPath())
	return schema.Schema{
		Parties:  prt,
		Projects: prj,
	}
}

// openCustomersProjects walks the given projects folder (path) and returns all found customers and projects.
func openCustomersProjects(path string) (schema.Parties, schema.Projects) {
	var cst schema.Parties
	var prj schema.Projects

	path = filepath.Join(path, "projects")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logrus.Fatalf("projects folder in acc repository doesn't exist, expected path: %s", path)
	}
	folders := getFoldersInPath(path)

	cstChan := make(chan schema.Party)
	prjCahn := make(chan schema.Project)

	var wg sync.WaitGroup
	for i := range folders {
		wg.Add(1)
		go customerWalk(folders[i], cstChan, prjCahn, &wg)
	}

	return cst, prj
}

// customerWalk goes trough one customer folder and puts the customer and all found projects into channels.
func customerWalk(path string, cstChan chan schema.Party, prjChan chan schema.Project, wg *sync.WaitGroup) {
	wg.Add(1)
	openCustomerFile(path, cstChan, wg)

	folders := getFoldersInPath(path)
	for i := range folders {
		wg.Add(1)
		openProjectFile(folders[i], prjChan, wg)
	}

	wg.Done()
}

// openCustomerFile tries to open the `customer.yaml` file in the given folder path.
// If the file exists it will be parsed and the customer get added to the customer channel.
func openCustomerFile(path string, cstChan chan schema.Party, wg *sync.WaitGroup) {
	cstFile := filepath.Join(path, "customer.yaml")
	if _, err := os.Stat(cstFile); os.IsNotExist(err) {
		logrus.Error("the customer.yaml file does not exist in ", path)
	} else {
		var cst schema.Party
		util.OpenYaml(&cst, cstFile, "customer file")
		cstChan <- cst
	}
	wg.Done()
}

func openProjectFile(path string, prjChan chan schema.Project, wg *sync.WaitGroup) {
	prjFile := filepath.Join(path, "project.yaml")
	if _, err := os.Stat(prjFile); os.IsNotExist(err) {
		logrus.Error("the project.yaml file does not exist in ", path)
	} else {
		var prj schema.Project
		util.OpenYaml(&prj, prjFile, "project file")
		prjChan <- prj
	}
	wg.Done()
}
