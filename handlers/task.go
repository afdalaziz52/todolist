package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/afdalaziz52/to-do-list/config"
    "github.com/afdalaziz52/to-do-list/helpers"
    "github.com/afdalaziz52/to-do-list/models"
    "github.com/gorilla/mux"
)

// GET /tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)

    rows, err := config.DB.Query(
        "SELECT id, user_id, title, category, custom_category, status, created_at, updated_at FROM tasks WHERE user_id = ?", userID)
    if err != nil {
        helpers.WriteJSON(w, 500, map[string]any{"status": "error", "message": "Gagal ambil data"})
        return
    }
    defer rows.Close()

    tasks := []models.Task{}
    for rows.Next() {
        var t models.Task
        rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Category, &t.CustomCategory, &t.Status, &t.CreatedAt, &t.UpdatedAt)
        tasks = append(tasks, t)
    }

    helpers.WriteJSON(w, 200, map[string]any{"status": "success", "data": tasks})
}

// GET /tasks/{id}
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)
    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    var t models.Task
    err := config.DB.QueryRow(
        "SELECT id, user_id, title, category, custom_category, status, created_at, updated_at FROM tasks WHERE id = ? AND user_id = ?", id, userID).
        Scan(&t.ID, &t.UserID, &t.Title, &t.Category, &t.CustomCategory, &t.Status, &t.CreatedAt, &t.UpdatedAt)

    if err != nil {
        helpers.WriteJSON(w, 404, map[string]any{"status": "error", "message": "Task tidak ditemukan"})
        return
    }

    helpers.WriteJSON(w, 200, map[string]any{"status": "success", "data": t})
}

// POST /tasks
func CreateTask(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)

    var req models.TaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        helpers.WriteJSON(w, 400, map[string]any{"status": "error", "message": "Format tidak valid"})
        return
    }

    if err := config.Validate.Struct(req); err != nil {
        helpers.WriteJSON(w, 400, map[string]any{"status": "error", "message": err.Error()})
        return
    }

    // validasi custom_category
    if req.Category == "other" && req.CustomCategory == nil {
        helpers.WriteJSON(w, 400, map[string]any{"status": "error", "message": "Custom category wajib diisi"})
        return
    }
    if req.Category != "other" {
        req.CustomCategory = nil
    }

    result, err := config.DB.Exec(
        "INSERT INTO tasks (user_id, title, category, custom_category, status) VALUES (?, ?, ?, ?, 'pending')",
        userID, req.Title, req.Category, req.CustomCategory)
    if err != nil {
        helpers.WriteJSON(w, 500, map[string]any{"status": "error", "message": "Gagal membuat task"})
        return
    }

    id, _ := result.LastInsertId()
    helpers.WriteJSON(w, 201, map[string]any{"status": "success", "message": "Task berhasil dibuat", "id": id})
}

// PATCH /tasks/{id}
func UpdateTask(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)
    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    var req models.TaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        helpers.WriteJSON(w, 400, map[string]any{"status": "error", "message": "Format tidak valid"})
        return
    }

    if err := config.Validate.Struct(req); err != nil {
        helpers.WriteJSON(w, 400, map[string]any{"status": "error", "message": err.Error()})
        return
    }

    if req.Category == "other" && req.CustomCategory == nil {
        helpers.WriteJSON(w, 400, map[string]any{"status": "error", "message": "Custom category wajib diisi"})
        return
    }
    if req.Category != "other" {
        req.CustomCategory = nil
    }

    result, err := config.DB.Exec(
        "UPDATE tasks SET title=?, category=?, custom_category=? WHERE id=? AND user_id=?",
        req.Title, req.Category, req.CustomCategory, id, userID)
    if err != nil {
        helpers.WriteJSON(w, 500, map[string]any{"status": "error", "message": "Gagal update task"})
        return
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        helpers.WriteJSON(w, 404, map[string]any{"status": "error", "message": "Task tidak ditemukan"})
        return
    }

    helpers.WriteJSON(w, 200, map[string]any{"status": "success", "message": "Task berhasil diupdate"})
}

// PATCH /tasks/{id}/status
func UpdateStatus(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)
    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    var body struct {
        Status string `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        helpers.WriteJSON(w, 400, map[string]any{"status": "error", "message": "Format tidak valid"})
        return
    }

    if body.Status != "pending" && body.Status != "done" {
        helpers.WriteJSON(w, 400, map[string]any{"status": "error", "message": "Status harus pending atau done"})
        return
    }

    result, err := config.DB.Exec(
        "UPDATE tasks SET status=? WHERE id=? AND user_id=?", body.Status, id, userID)
    if err != nil {
        helpers.WriteJSON(w, 500, map[string]any{"status": "error", "message": "Gagal update status"})
        return
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        helpers.WriteJSON(w, 404, map[string]any{"status": "error", "message": "Task tidak ditemukan"})
        return
    }

    helpers.WriteJSON(w, 200, map[string]any{"status": "success", "message": "Status berhasil diupdate"})
}

// DELETE /tasks/{id}
func DeleteTask(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)
    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    result, err := config.DB.Exec(
        "DELETE FROM tasks WHERE id=? AND user_id=?", id, userID)
    if err != nil {
        helpers.WriteJSON(w, 500, map[string]any{"status": "error", "message": "Gagal hapus task"})
        return
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        helpers.WriteJSON(w, 404, map[string]any{"status": "error", "message": "Task tidak ditemukan"})
        return
    }

    helpers.WriteJSON(w, 200, map[string]any{"status": "success", "message": "Task berhasil dihapus"})
}