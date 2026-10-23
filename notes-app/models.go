package main

import "time"

// Note is one saved note.
type Note struct {
	ID        int
	Subject   string
	Body      string
	Tags      []string // e.g. []string{"go", "school"}
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TagsString joins tags for display, e.g. "go, school"
func (n Note) TagsString() string {
	out := ""
	for i, t := range n.Tags {
		if i > 0 {
			out += ", "
		}
		out += t
	}
	return out
}
