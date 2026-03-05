package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/afdalaziz52/to-do-list/config"
    "github.com/afdalaziz52/to-do-list/routes"
    "github.com/joho/godotenv"
)

func corsMiddleware(next http.Handler) http.Handler {
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
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("❌ Gagal load .env")
    }

    config.InitDB()

    r := routes.SetupRoutes()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
    fmt.Printf("   App: %s\n", os.Getenv("APP_NAME"))
    fmt.Printf("   Env: %s\n", os.Getenv("APP_ENV"))

    // wrap router dengan CORS
    if err := http.ListenAndServe(":"+port, corsMiddleware(r)); err != nil {
        log.Fatal("❌ Server gagal berjalan:", err)
    }
}