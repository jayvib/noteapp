package note

import "bytes"

// SortByIdSorter implements sort.Interface which
// sort the note by its ID.
type SortByIdSorter []*Note

// Len returns the length of notes.
func (n SortByIdSorter) Len() int { return len(n) }

// Less compare the adjacent IDs of the note.
func (n SortByIdSorter) Less(i, j int) bool {
	return bytes.Compare(n[i].ID[:], n[j].ID[:]) < 0
}

// Swap swaps the note i, and note j.
func (n SortByIdSorter) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// SortByTitleSorter implements sort.Interface which
// sort the note by title.
type SortByTitleSorter []*Note

// Len returns the length of notes.
func (n SortByTitleSorter) Len() int { return len(n) }

// Less compare the adjacent IDs of the note.
func (n SortByTitleSorter) Less(i, j int) bool {
	return n[i].GetTitle() < n[j].GetTitle()
}

// Swap swaps the note i, and note j.
func (n SortByTitleSorter) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
