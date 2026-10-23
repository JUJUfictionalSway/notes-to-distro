package main

import (
	"log"
	"net/http"
)

func main() {
	if err := initDB("notes.db"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /notes", handleList)
	mux.HandleFunc("GET /notes/new", handleNewForm)
	mux.HandleFunc("POST /notes", handleCreate)
	mux.HandleFunc("GET /notes/{id}", handleView)
	mux.HandleFunc("GET /notes/{id}/edit", handleEditForm)
	mux.HandleFunc("POST /notes/{id}", handleUpdate)
	mux.HandleFunc("POST /notes/{id}/delete", handleDelete)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// send bare "/" to the list so it's not a 404
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
	})

	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
