package schema

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"regexp"
	"strconv"
)

// Identifiable describes types which are uniquely identifiable trough out the utils structure.
type Identifiable interface {
	GetId() string
	GetIdentifier() string
	String() string
}

func SuggestNextIdentifier(idt []Identifiable, prefix string) string {
	r := regexp.MustCompile(`(\d+)$`)
	max := 0
	for i := range idt {
		rsl := r.FindAllString(idt[i].GetIdentifier(), -1)
		if len(rsl) != 1 {
			continue
		}
		val, err := strconv.Atoi(rsl[0])
		if err != nil {
			logrus.Debugf("regex to find last number in identifier returned something else than a int (%+v), take a look: %s", rsl, err)
		}
		if max < val {
			max = val
		}
	}
	return fmt.Sprintf("%s%d", prefix, max+1)
}

// Completable describes types which automatically can resolve some missing information atomically.
// Example is the setting of a unique Id.
type Completable interface {
	// SetId sets a unique id (UUID-4) for the object.
	SetId()
}

// SaveToYaml writes the element (utils) as a json file to the given path.
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

func GetUuid() string {
	return uuid.Must(uuid.NewRandom()).String()
}
