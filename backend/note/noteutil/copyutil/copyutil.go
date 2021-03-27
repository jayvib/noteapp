package copyutil

import (
	"noteapp/note"
)

// Shallow takes a note and then returns a copied note with
// a new address.
func Shallow(n *note.Note) *note.Note {
	cpyNote := new(note.Note)
	*cpyNote = *n
	return cpyNote
}
