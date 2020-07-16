package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/logrusorgru/aurora"
)

type Queryables []Queryable

func (q Queryables) Types() []string {
	result := make([]string, len(q))
	for i := range q {
		result[i] = q[i].Name
	}
	return result
}

func (q Queryables) PPTypes() {
	typ := q.Types()
	fmt.Println(aurora.BrightMagenta("Types:"))
	for i := range typ {
		fmt.Println(aurora.Green(fmt.Sprintf("- %s", typ[i])))
	}
}

func (q Queryables) keys() Keys {
	var rsl Keys
	for i := range q {
		keys := q[i].keys()
		for j := range keys {
			rsl = rsl.Add(keys[j])
		}
	}
	return rsl
}

func (q Queryables) PPKeys(mode OutputMode) {
	kys := q.keys()
	fmt.Println(fmt.Sprintf("%s%s%s", aurora.BrightMagenta(aurora.Bold("Keys")), aurora.BrightMagenta(" for "), aurora.BrightMagenta(q.String())))
	out := Output{
		Header:  []string{"Name Key", "Element Type(s)"},
		Element: kys.Data(),
	}
	out.PPKeyValue(nil, mode, false)
}

func (q Queryables) QueryablesFromUserInput(input string) (Queryables, error) {
	if input == "" {
		return q, nil
	}

	inputs := strings.Split(input, ",")
	var rsl Queryables
	for i := range inputs {
		inp := strings.TrimPrefix(inputs[i], " ")
		inp = strings.TrimSuffix(inp, " ")
		for j := range q {
			if q[j].matchTypeFromUserInput(inp) {
				rsl = append(rsl, q[j])
			}
		}
	}
	if len(rsl) == 0 {
		return rsl, fmt.Errorf("no types for \"%s\" found", q.String())
	}
	return rsl, nil
}

func (q Queryables) String() string {
	if len(q) == 0 {
		return ""
	}
	rsl := q[0].Name
	for i := 1; i < len(q); i++ {
		rsl = fmt.Sprint(rsl, ", ", q[i].Name)
	}
	return rsl
}

type Keys []Key

func (k Keys) Add(key Key) Keys {
	for i := range k {
		if k[i].Name == key.Name {
			k[i].Queryables = append(k[i].Queryables, key.Queryables...)
			return k
		}
	}
	return append(k, key)
}

func (k Keys) Data() []KeyValue {
	rsl := make([]KeyValue, len(k))
	for i := range k {
		rsl[i] = KeyValue{
			Key:   k[i].Name,
			Value: k[i].Queryables.String()}
	}
	return rsl
}

type Key struct {
	Name       string
	Queryables Queryables
}

type Queryable struct {
	Name string
	Type interface{}
}

func (q Queryable) matchTypeFromUserInput(input string) bool {
	input = strings.ToLower(input)
	if input == q.Name {
		return true
	}
	if q.Name == strings.TrimSuffix(input, "s") {
		return true
	}
	return false
}

func (q Queryable) keys() Keys {
	var rsl Keys
	v := reflect.ValueOf(q.Type)
	for i := 0; i < v.NumField(); i++ {
		rsl = rsl.Add(Key{
			Name:       v.Type().Field(i).Name,
			Queryables: Queryables{q},
		})
	}
	return rsl
}
