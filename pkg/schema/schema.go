package schema

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// Identifiable describes types which are uniquely identifiable trough out the data structure.
type Identifiable interface {
	GetId() string
}

// Completable describes types which automatically can resolve some missing information atomically.
// Example is the setting of a unique Id.
type Completable interface {
	// SetId sets a unique id (UUID-4) for the object.
	SetId()
}

// SaveToJson writes the element (data) as a json file to the given path.
// Indented states whether «prettify» the json output.
func SaveToJson(data interface{}, path string, indented bool) {
	var raw []byte
	var err error
	if indented {
		raw, err = json.MarshalIndent(data, "", "    ")
	} else {
		raw, err = json.Marshal(data)
	}
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(path)
	if err := ioutil.WriteFile(path, raw, 0644); err != nil {
		logrus.Fatal(err)
	}
}
