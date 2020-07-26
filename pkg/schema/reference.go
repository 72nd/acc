package schema

import "gopkg.in/yaml.v3"

// Reference is used to reference to other data types by using UUID's.
type Reference struct {
	Id          string
	Destination Identifiable
}

// NewReference returns a new Reference with the given id.
func NewReference(id string) Reference {
	return Reference{
		Id: id}
}

// Empty returns whether a referencing id is set
func (r Reference) Empty() bool {
	return r.Id == ""
}

// Match returns whether a given Referenceable matches the id of the reference.
func (r Reference) Match(val Identifiable) bool {
	return r.Id == val.GetId()
}

// SetDestination takes a slice of Referenceables and looks for the matching element
// in respective to the id. If found the Destination is set.
func (r *Reference) SetDestination(destinations []Identifiable) {
	for i := range destinations {
		if destinations[i].GetId() == r.Id {
			r.Destination = destinations[i]
		}
	}
}

// UnmarshalYAML implements the unmarshaling of a Reference for YAML files.
func (r *Reference) UnmarshalYAML(value *yaml.Node) error {
	r.Id = value.Value
	return nil
}

// MarshalYAML implements the marshalling of a Reference for YAML files.
func (r Reference) MarshalYAML() (interface{}, error) {
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
