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

func main() {
    // load .env
    if err := godotenv.Load(); err != nil {
        log.Fatal("❌ Gagal load .env")
    }

    // koneksi database
    config.InitDB()

    // setup router
    r := routes.SetupRoutes()

    // jalankan server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
    fmt.Printf("   App: %s\n", os.Getenv("APP_NAME"))
    fmt.Printf("   Env: %s\n", os.Getenv("APP_ENV"))

    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatal("❌ Server gagal berjalan:", err)
    }
}