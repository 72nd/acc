package schema

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Identifiable describes types which are uniquely identifiable trough out the fonts structure.
type Identifiable interface {
	GetId() string
}

// Completable describes types which automatically can resolve some missing information atomically.
// Example is the setting of a unique Id.
type Completable interface {
	// SetId sets a unique id (UUID-4) for the object.
	SetId()
}

// SaveToYaml writes the element (fonts) as a json file to the given path.
// Indented states whether «prettify» the json output.
func SaveToYaml(data interface{}, path string) {
	var raw []byte
	var err error
	raw, err = yaml.Marshal(data)
	if err != nil {
		logrus.Fatalf("error converting data from file %s to YAML (marshalling): %s", path, err)
	}
	if err := ioutil.WriteFile(path, raw, 0644); err != nil {
		logrus.Fatalf("error writing file %s: %s", path, err)
	}
}
