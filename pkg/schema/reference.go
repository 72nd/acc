package schema

import (
	"gopkg.in/yaml.v3"
)

// Ref is used to reference to other data types by using UUID's.
type Ref struct {
	Id          string
	Destination Identifiable
}

// NewRef returns a new Reference with the given id.
func NewRef(id string) Ref {
	return Ref{
		Id: id,
	}
}

// Empty returns whether a referencing id is set
func (r Ref) Empty() bool {
	return r.Id == ""
}

// Match returns whether a given Identifiable matches the id of the reference.
func (r Ref) Match(val Identifiable) bool {
	return r.Id == val.GetId()
}

// SetDestination takes a slice of Identifiables and looks for the matching element
// in respective to the id. If found the Destination is set.
func (r *Ref) SetDestination(destinations []Identifiable) {
	for i := range destinations {
		if destinations[i].GetId() == r.Id {
			r.Destination = destinations[i]
		}
	}
}

// UnmarshalYAML implements the unmarshaling of a Reference for YAML files.
func (r *Ref) UnmarshalYAML(value *yaml.Node) error {
	r.Id = value.Value
	return nil
}

// MarshalYAML implements the marshalling of a Reference for YAML files.
func (r Ref) MarshalYAML() (interface{}, error) {
	n := yaml.Node{}
	n.Kind = yaml.ScalarNode
	n.Value = r.Id
	if r.Destination != nil {
		n.LineComment = r.Destination.String()
	} else if r.Id != "" {
		// This omits pointless comments when the referencing id is empty.
		n.LineComment = "no element found for this id"
	}
	return &n, nil
}
