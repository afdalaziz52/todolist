package routes

import (
    "net/http"
    "github.com/afdalaziz52/to-do-list/handlers"
    "github.com/afdalaziz52/to-do-list/middleware"
    "github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
    r := mux.NewRouter()

    // ─── Static file (frontend) ───
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "static/index.html")
    }).Methods("GET")

    // ─── Auth (public, tidak perlu token) ───
    auth := r.PathPrefix("/api/auth").Subrouter()
    auth.HandleFunc("/register", handlers.Register).Methods("POST")
    auth.HandleFunc("/login", handlers.Login).Methods("POST")

    // ─── Tasks (protected, butuh token) ───
    tasks := r.PathPrefix("/api/tasks").Subrouter()
    tasks.Use(middleware.AuthMiddleware)
    tasks.HandleFunc("", handlers.GetTasks).Methods("GET")
    tasks.HandleFunc("", handlers.CreateTask).Methods("POST")
    tasks.HandleFunc("/{id}", handlers.GetTaskByID).Methods("GET")
    tasks.HandleFunc("/{id}", handlers.UpdateTask).Methods("PATCH")
    tasks.HandleFunc("/{id}/status", handlers.UpdateStatus).Methods("PATCH")
    tasks.HandleFunc("/{id}", handlers.DeleteTask).Methods("DELETE")

    return r
}