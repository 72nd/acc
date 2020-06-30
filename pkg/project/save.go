package project

import "gitlab.com/72th/acc/pkg/schema"

// Save saves a schema as a folder structure to the given folder (path).
// This function is only directly called when converting a project to folder mode.
func Save(s schema.Schema, path string) {

}

// SaveWithEnv does the same as Save() but uses the `ACC_FOLDER` env variable.
// This is the default use case.
func SaveWithEnv(s schema.Schema) {
	Save(s, repositoryPath())
}
