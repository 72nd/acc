package schema

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"github.com/72nd/acc/pkg/util"
)

const DefaultProjectsFile = "projects.yaml"
const DefaultProjectPrefix = "p-"

// Projects represents a collection of multiple Project elements.
type Projects []Project

// NewProjects returns an empty Projects collection.
func NewProjects() Projects {
	return Projects{}
}

// OpenProjects opens the Projects saved in the YAML file given by the path.
func OpenProjects(path string) Projects {
	var prj Projects
	util.OpenYaml(&prj, path, "projects")
	return prj
}

// Save writes the element as YAML file to the given path.
func (p Projects) Save(path string) {
	util.SaveToYaml(p, path, "projects")
}

// ProjectById returns the Project with the given id. If no record could be found
// an error will be returned.
func (p Projects) ProjectById(id string) (*Project, error) {
	for i := range p {
		if p[i].Id == id {
			return &p[i], nil
		}
	}
	return nil, fmt.Errorf("no project for id \"%s\" found", id)
}

// ProjectByIdent returns the Project with the given identifier.
// If no record could be found an error will be returned.
func (p Projects) ProjectByIdent(ident string) (*Project, error) {
	for i := range p {
		if p[i].Identifier == ident {
			return &p[i], nil
		}
	}
	return nil, fmt.Errorf("no project for ident \"%s\" found", ident)
}

// GetIdentifiables returns the a sliche of all identifiers. This is used for the
// identifier suggestion while interactively adding a new Project.
func (p Projects) GetIdentifiables() []Identifiable {
	rsl := make([]Identifiable, len(p))
	for i := range p {
		rsl[i] = p[i]
	}
	return rsl
}

// SearchItems returns a searchable structure of the Projects. So the user
// can search for Projects in the interactive mode.
func (p Projects) SearchItems() util.SearchItems {
	rsl := make(util.SearchItems, len(p))
	for i := range p {
		rsl[i] = p[i].SearchItem()
	}
	return rsl
}

// Validate all Projects.
func (p Projects) Validate() util.ValidateResults {
	var rsl util.ValidateResults
	for i := range p {
		rsl = append(rsl, util.Check(p[i]))
	}
	return rsl
}

// Project represents a project for a customer.
type Project struct {
	Id         string `yaml:"id" default:"1"`
	Identifier string `yaml:"identifier" default:"p-1"`
	Name       string `yaml:"name" default:"Building a space rocket"`
	CustomerId string `yaml:"customerId" default:""`
}

// NewProject returns a new Project element with the default values.
func NewProject() Project {
	prj := Project{}
	if err := defaults.Set(&prj); err != nil {
		logrus.Fatal("error setting defaults for project: ", err)
	}
	prj.Id = GetUuid()
	return prj
}

// InteractiveNewProject returns a new Project based on the user input.
func InteractiveNewProject(s Schema) Project {
	var prj Project
	prj.Id = GetUuid()
	prj.Identifier = util.AskString(
		"Value",
		"Unique human readable identifier",
		SuggestNextIdentifier(s.Projects.GetIdentifiables(), DefaultProjectPrefix))
	prj.Name = util.AskString(
		"Name",
		"Name of the project",
		"Aktion neue soziale Marktwirtschaft")
	prj.CustomerId = util.AskStringFromSearch(
		"Associated customer",
		"Customer which the project is associated",
		s.Parties.CustomersSearchItems())
	return prj
}

// SetId generates a unique id for the element if there isn't already one defined.
func (p *Project) SetId() {
	if p.Id != "" {
		return
	}
	p.Id = GetUuid()
}

// GetId returns the id of the Project.
func (p Project) GetId() string {
	return p.Id
}

// GetIdentifier return the identifier of the Project.
func (p Project) GetIdentifier() string {
	return p.Identifier
}

// SearchItem returns a searchable representation of the Project.
func (p Project) SearchItem() util.SearchItem {
	return util.SearchItem{
		Name:        fmt.Sprintf("%s (%s)", p.Name, p.Identifier),
		Type:        p.Type(),
		Value:       p.Id,
		SearchValue: fmt.Sprintf("%s %s", p.Identifier, p.Name)}
}

// String returns a human readable representation of the element.
func (p Project) String() string {
	return fmt.Sprintf("project %s (%s)", p.Name, p.Identifier)
}

// Short returns a short represenation of the element.
func (p Project) Short() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Identifier)
}

// Type returns a string with the type name of the element.
func (p Project) Type() string {
	return "Project"
}

// Conditions returns the validation conditions.
func (p Project) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: p.Id == "",
			Message:   "unique identifier not set (Id is empty)",
		},
		{
			Condition: p.Identifier == "",
			Message:   "human readable identifier not set (Identifier is empty)",
		},
		{
			Condition: p.Name == "",
			Message:   "name not set (Name is empty)",
		},
		{
			Condition: p.CustomerId == "",
			Message:   "customer id not set (CustomerId is empty)",
		},
	}
}
