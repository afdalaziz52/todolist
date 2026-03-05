package routes

import (
    "net/http"
    "github.com/afdalaziz52/to-do-list/handlers"
    "github.com/afdalaziz52/to-do-list/middleware"
    "github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
    r := mux.NewRouter()

    // ─── CORS Middleware ───
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            next.ServeHTTP(w, r)
        })
    })

    // ─── Static files ───
    fs := http.FileServer(http.Dir("./static"))
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    // login = halaman pertama (root)
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/index.html")
    }).Methods("GET")

    r.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/index.html")
    }).Methods("GET")

    // dashboard = halaman tasks setelah login
    r.HandleFunc("/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/dashboard.html")
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