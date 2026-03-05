package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Task represents a todo item
type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Category  string `json:"category"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

// Store holds all tasks in memory
type Store struct {
	mu     sync.RWMutex
	tasks  []Task
	nextID int
}

var store = &Store{
	tasks:  []Task{},
	nextID: 1,
}

// Response helpers
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /api/tasks - list all tasks
func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	writeJSON(w, http.StatusOK, store.tasks)
}

// POST /api/tasks - create new task
func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text     string `json:"text"`
		Category string `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Request tidak valid")
		return
	}

	text := strings.TrimSpace(input.Text)
	if text == "" {
		writeError(w, http.StatusBadRequest, "Teks tugas tidak boleh kosong")
		return
	}
	if len(text) > 120 {
		writeError(w, http.StatusBadRequest, "Teks tugas terlalu panjang (maks 120 karakter)")
		return
	}

	category := input.Category
	if category == "" {
		category = "lainnya"
	}

	store.mu.Lock()
	task := Task{
		ID:        store.nextID,
		Text:      text,
		Category:  category,
		Done:      false,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	store.tasks = append(store.tasks, task)
	store.nextID++
	store.mu.Unlock()

	writeJSON(w, http.StatusCreated, task)
}

// PUT /api/tasks/{id} - update task (toggle done or edit text)
func handleUpdateTask(w http.ResponseWriter, r *http.Request, id int) {
	var input struct {
		Done     *bool   `json:"done"`
		Text     *string `json:"text"`
		Category *string `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Request tidak valid")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for i, t := range store.tasks {
		if t.ID == id {
			if input.Done != nil {
				store.tasks[i].Done = *input.Done
			}
			if input.Text != nil {
				text := strings.TrimSpace(*input.Text)
				if text == "" {
					writeError(w, http.StatusBadRequest, "Teks tugas tidak boleh kosong")
					return
				}
				store.tasks[i].Text = text
			}
			if input.Category != nil {
				store.tasks[i].Category = *input.Category
			}
			writeJSON(w, http.StatusOK, store.tasks[i])
			return
		}
	}
	writeError(w, http.StatusNotFound, "Tugas tidak ditemukan")
}

// DELETE /api/tasks/{id} - delete a single task
func handleDeleteTask(w http.ResponseWriter, r *http.Request, id int) {
	store.mu.Lock()
	defer store.mu.Unlock()

	for i, t := range store.tasks {
		if t.ID == id {
			store.tasks = append(store.tasks[:i], store.tasks[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": "Tugas dihapus"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "Tugas tidak ditemukan")
}

// DELETE /api/tasks/done - delete all completed tasks
func handleDeleteDone(w http.ResponseWriter, r *http.Request) {
	store.mu.Lock()
	defer store.mu.Unlock()

	remaining := []Task{}
	deleted := 0
	for _, t := range store.tasks {
		if !t.Done {
			remaining = append(remaining, t)
		} else {
			deleted++
		}
	}
	store.tasks = remaining
	writeJSON(w, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("%d tugas selesai dihapus", deleted),
		"deleted": deleted,
	})
}

// Router for /api/tasks and /api/tasks/{id}
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/tasks")
	path = strings.TrimPrefix(path, "/")

	// DELETE /api/tasks/done
	if path == "done" && r.Method == http.MethodDelete {
		handleDeleteDone(w, r)
		return
	}

	// /api/tasks (no id)
	if path == "" {
		switch r.Method {
		case http.MethodGet:
			handleGetTasks(w, r)
		case http.MethodPost:
			handleCreateTask(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		}
		return
	}

	// /api/tasks/{id}
	id, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}
	switch r.Method {
	case http.MethodPut:
		handleUpdateTask(w, r, id)
	case http.MethodDelete:
		handleDeleteTask(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
	}
}

// Serve static files from current directory
func staticHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "index.html")
		return
	}
	http.ServeFile(w, r, "."+r.URL.Path)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks/", tasksHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.HandleFunc("/", staticHandler)

	port := ":8080"
	fmt.Println("🚀 Server berjalan di http://localhost" + port)
	fmt.Println("📋 Endpoints:")
	fmt.Println("   GET    /api/tasks        - Ambil semua tugas")
	fmt.Println("   POST   /api/tasks        - Buat tugas baru")
	fmt.Println("   PUT    /api/tasks/{id}   - Update tugas")
	fmt.Println("   DELETE /api/tasks/{id}   - Hapus satu tugas")
	fmt.Println("   DELETE /api/tasks/done   - Hapus semua tugas selesai")

	if err := http.ListenAndServe(port, mux); err != nil {
		fmt.Println("Error:", err)
	}
}
