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

	// ─── Static assets (CSS, JS, dll) ───
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// ─── Halaman HTML ───
	// Root → login page (halaman pertama)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	}).Methods("GET")

	r.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	}).Methods("GET")

	r.HandleFunc("/register.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/register.html")
	}).Methods("GET")

	// Dashboard = halaman utama setelah login (protected di sisi frontend dengan JWT check)
	r.HandleFunc("/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	}).Methods("GET")

	// ─── Auth API (public) ───
	auth := r.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/register", handlers.Register).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")

	// ─── Auth API (protected) ───
	authProtected := r.PathPrefix("/api/auth").Subrouter()
	authProtected.Use(middleware.AuthMiddleware)
	authProtected.HandleFunc("/change-password", handlers.ChangePassword).Methods("POST")

	// ─── Tasks API (protected, butuh JWT token) ───
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
