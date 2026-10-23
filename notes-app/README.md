# notes-app

Simple Go + SQLite notes app: create, edit, delete notes with a subject and tags.

## Run it

```bash
go mod tidy      # downloads modernc.org/sqlite (pure Go, no C compiler needed)
go run .
```

Then open http://localhost:8080

## How it's laid out

- `main.go`      — starts the server, registers routes
- `db.go`        — opens SQLite, all CRUD queries live here
- `models.go`    — the `Note` struct
- `handlers.go`  — turns HTTP requests into db.go calls + template renders
- `schema.sql`   — table definitions, run automatically on startup
- `templates/`   — the three HTML pages (list, view, shared create/edit form)
- `static/`      — stylesheet

A `notes.db` SQLite file will be created next to the binary the first time you run it.

## Suggested build order

1. Get `go run .` serving the empty list page at `/notes`.
2. Wire up create (`/notes/new` form -> POST `/notes`).
3. Add view + edit + delete for a single note.
4. Confirm tags save/load correctly (they're the trickiest part — stored in their
   own table and joined through `note_tags`).
5. Polish: maybe a search box, or filter the list by clicking a tag.

## Notes for extending it later

- Right now `splitTags` lower-cases and trims tags from one comma-separated
  input field — good enough for v1, but a JS tag-input widget would feel nicer.
- There's no auth — every visitor sees the same notes. Fine for a personal
  project; would need sessions + a `user_id` column to support multiple users.
