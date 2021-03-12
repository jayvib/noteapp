package copyutil

import "noteapp/notes"

// Shallow takes a note and then returns a copied note with
// a new address.
func Shallow(note *notes.Note) *notes.Note {
	cpyNote := new(notes.Note)
	*cpyNote = *note
	return cpyNote
}
