package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Note struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type NoteStore struct {
	db *sql.DB
}

func NewNoteStore() *NoteStore {
	mysqlPassword := os.Getenv("MYSQL_ROOT_PASSWORD")
	dataSourceName := fmt.Sprintf("root:%s@tcp(mysql:3306)/notesapp?parseTime=true", mysqlPassword)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v", err)
	}
	defer db.Close()

	return &NoteStore{db: db}
}

func (store *NoteStore) createNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := store.db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", note.Title, note.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	note.ID = int(id)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	log.Printf("created note: %d", note.ID)
}

func (store *NoteStore) deleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert idStr to int
	ID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	_, err = store.db.Exec("DELETE FROM notes WHERE id = ?", ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("deleted note: %d", ID)
}

func (store *NoteStore) getNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := store.db.Query("SELECT * FROM notes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notes = append(notes, note)
	}

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (store *NoteStore) getNoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert idStr to int
	ID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	row := store.db.QueryRow("SELECT * FROM notes WHERE id = ?", ID)

	var note Note

	err = row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (store *NoteStore) updateNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateTitle := note.Title != ""
	updateContent := note.Content != ""

	if updateTitle {
		_, err := store.db.Query("UPDATE notes SET title = ? WHERE id = ?", note.Title, note.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	if updateContent {
		_, err := store.db.Query("UPDATE notes SET content = ? WHERE id = ?", note.Content, note.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	err := json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	log.Printf("updated note: %d", note.ID)
}

func Router() *mux.Router {
	store := NewNoteStore()
	r := mux.NewRouter()
	r.HandleFunc("/create", store.createNote).Methods("POST")
	r.HandleFunc("/update", store.updateNote).Methods("POST")
	r.HandleFunc("/delete/{id}", store.deleteNote).Methods("POST")
	r.HandleFunc("/get/{id}", store.getNoteByID).Methods("GET")
	r.HandleFunc("/get-all", store.getNotes).Methods("GET")

	return r
}

func main() {
	r := Router()

	log.Println("Starting server on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
