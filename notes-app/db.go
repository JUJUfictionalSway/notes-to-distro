package main

import (
	"database/sql"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// initDB opens the SQLite file and makes sure the schema exists.
func initDB(path string) error {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	return err
}

// listNotes returns every note, newest first, tags included.
func listNotes() ([]Note, error) {
	rows, err := db.Query(`SELECT id, subject, body, created_at, updated_at
	                        FROM notes ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Subject, &n.Body, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		n.Tags, err = tagsForNote(n.ID)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

// getNote fetches one note by id, or sql.ErrNoRows if it doesn't exist.
func getNote(id int) (Note, error) {
	var n Note
	row := db.QueryRow(`SELECT id, subject, body, created_at, updated_at
	                     FROM notes WHERE id = ?`, id)
	err := row.Scan(&n.ID, &n.Subject, &n.Body, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return n, err
	}
	n.Tags, err = tagsForNote(id)
	return n, err
}

// createNote inserts a note and its tags, returns the new id.
func createNote(subject, body string, tags []string) (int, error) {
	res, err := db.Exec(`INSERT INTO notes (subject, body) VALUES (?, ?)`, subject, body)
	if err != nil {
		return 0, err
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	id := int(id64)
	return id, setTags(id, tags)
}

// updateNote overwrites subject/body/tags and bumps updated_at.
func updateNote(id int, subject, body string, tags []string) error {
	_, err := db.Exec(`UPDATE notes SET subject = ?, body = ?, updated_at = CURRENT_TIMESTAMP
	                    WHERE id = ?`, subject, body, id)
	if err != nil {
		return err
	}
	return setTags(id, tags)
}

// deleteNote removes a note; note_tags rows go too via ON DELETE CASCADE.
func deleteNote(id int) error {
	_, err := db.Exec(`DELETE FROM notes WHERE id = ?`, id)
	return err
}

// tagsForNote returns tag names attached to a note.
func tagsForNote(noteID int) ([]string, error) {
	rows, err := db.Query(`SELECT t.name FROM tags t
	                        JOIN note_tags nt ON nt.tag_id = t.id
	                        WHERE nt.note_id = ? ORDER BY t.name`, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tags = append(tags, name)
	}
	return tags, nil
}

// setTags replaces all tags on a note with the given list, creating
// any tag rows that don't exist yet.
func setTags(noteID int, tags []string) error {
	if _, err := db.Exec(`DELETE FROM note_tags WHERE note_id = ?`, noteID); err != nil {
		return err
	}
	for _, raw := range tags {
		name := strings.ToLower(strings.TrimSpace(raw))
		if name == "" {
			continue
		}
		if _, err := db.Exec(`INSERT OR IGNORE INTO tags (name) VALUES (?)`, name); err != nil {
			return err
		}
		var tagID int
		if err := db.QueryRow(`SELECT id FROM tags WHERE name = ?`, name).Scan(&tagID); err != nil {
			return err
		}
		if _, err := db.Exec(`INSERT OR IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)`,
			noteID, tagID); err != nil {
			return err
		}
	}
	return nil
}

// splitTags turns "go, school, project" into []string{"go","school","project"}
func splitTags(raw string) []string {
	parts := strings.Split(raw, ",")
	var tags []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}
	return tags
}
