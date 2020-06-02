package query

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

type ElementGroup []Element

func NewElements(ele interface{}) ElementGroup {
	var rsl []Element
	v := reflect.ValueOf(ele)
	if v.Kind() != reflect.Slice {
		logrus.Fatalf("\"%+v\" isn't a slice", ele)
	}
	for i := 0; i < v.Len(); i++ {
		rsl = append(rsl, NewElement(v.Index(i)))
	}
	return rsl
}

func (g ElementGroup) Match(terms SearchTerms) ElementGroup {
	var rsl ElementGroup
	for i := range g {
		if g[i].Match(terms) {
			rsl = append(rsl, g[i])
		}
	}
	return rsl
}

func (g ElementGroup) DateMatch(ranges DateTerms) ElementGroup {
	var rsl ElementGroup
	for i := range g {
		if g[i].DateMatch(ranges) {
			rsl = append(rsl, g[i])
		}
	}
	return rsl
}

type Element []KeyValue

func NewElement(v reflect.Value) Element {
	var rsl Element
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		rsl = append(rsl, KeyValue{
			Key:   t.Field(i).Name,
			Value: fmt.Sprint(v.Field(i)),
		})
	}
	return rsl
}

func (e Element) Match(terms SearchTerms) bool {
	for i := range terms {
		for j := range e {
			if terms[i].matchKey(e[j].Key) {
				if !terms[i].matchValue(e[j].Value) {
					return false
				}
			}
		}
	}
	return true
}

func (e Element) DateMatch(ranges DateTerms) bool {
	for i := range ranges {
		for j := range e {
			if ranges[i].matchKey(e[j].Key) {
				if !ranges[i].matchRange(e[j].Value) {
					return false
				}
			}
		}
	}
	return true
}

type KeyValue struct {
	Key   string
	Value string
}
