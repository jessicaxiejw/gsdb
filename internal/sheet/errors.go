package sheet

import "fmt"

type FileNotFoundError struct {
	kind string
	name string
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("cannnot find a %s named %s", e.kind, e.name)
}

type DuplicatedFilesError struct {
	kind string
	name string
}

func (e *DuplicatedFilesError) Error() string {
	return fmt.Sprintf("there can only be one %s named %s shared with this service account", e.kind, e.name)
}

type DuplicatedTableError struct {
	kind string
	name string
}

func (e *DuplicatedTableError) Error() string {
	return fmt.Sprintf("there can only be one %s named %s shared with this service account. Are you creating a table that already exists?", e.kind, e.name)
}
