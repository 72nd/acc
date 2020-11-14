// REST API Server
package server

import (
	"github.com/72nd/acc/pkg/schema"
)

// Defines the endpoint for the REST interface.
type Endpoint struct {
	// The Schema to operate on.
	schema *schema.Schema
}

// NewEndpoint returns a new endpoint. Takes a Schema as parameter. The request
// received by the server will be applied on this given
// data.
func NewEndpoint(s *schema.Schema) Endpoint {
	return Endpoint{
		schema: s,
	}
}

// Serve Runs the REST endpoint on the given port.
func (e *Endpoint) Serve(port int) {
	//router := e.buildRouter()
	//http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

/*
// GenDocs generates the API documentation and writes the result into the given folder.
// As this function generates two different kind of documentation (JSON and Markdown)
// two files are written. The two file extensions are added to the given path.
func (e Endpoint) GenDocs(path string, do-overwrite bool) {
	router := e.buildRouter()
	docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{
	})
}

// buildRouter builds the chi router.
func (e *Endpoint) buildRouter() chi.Router {
	r := chi.NewRouter();
	r.Use(render.SetContentType(render.ContentTypeJSON))
	return r
}
*/
