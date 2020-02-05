package schema

type Party struct {
	Name       string `json:"name" default:"Max Mustermann"`
	Street     string `json:"street" default:"Main Street"`
	StreetNr   int    `json:"streetNr" default:"1"`
	Place      string `json:"place" default:"Zurich"`
	PostalCode int    `json:"postalCode" default:"8000"`
}
