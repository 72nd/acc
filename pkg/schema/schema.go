package schema

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
