package main

import (
	"html/template"
	"net/http"
	"strconv"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// GET /notes  -- list every note
func handleList(w http.ResponseWriter, r *http.Request) {
	notes, err := listNotes()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpl.ExecuteTemplate(w, "list.html", notes)
}

// GET /notes/new -- blank create form
func handleNewForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "form.html", Note{})
}

// POST /notes -- create a note, then redirect to the list
func handleCreate(w http.ResponseWriter, r *http.Request) {
	subject := r.FormValue("subject")
	body := r.FormValue("body")
	tags := splitTags(r.FormValue("tags"))

	if _, err := createNote(subject, body, tags); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/notes", http.StatusSeeOther)
}

// GET /notes/{id} -- view one note
func handleView(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	n, err := getNote(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, "view.html", n)
}

// GET /notes/{id}/edit -- pre-filled edit form
func handleEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	n, err := getNote(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, "form.html", n)
}

// POST /notes/{id} -- save edits, then back to the note
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	subject := r.FormValue("subject")
	body := r.FormValue("body")
	tags := splitTags(r.FormValue("tags"))

	if err := updateNote(id, subject, body, tags); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/notes/"+strconv.Itoa(id), http.StatusSeeOther)
}

// POST /notes/{id}/delete -- remove, then back to the list
func handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := deleteNote(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/notes", http.StatusSeeOther)
}

// idFromPath pulls the numeric id out of paths like /notes/3 or /notes/3/edit.
// Uses Go 1.22's PathValue since routes are registered with "{id}".
func idFromPath(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}
