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
	// Load .env (tidak fatal jika tidak ada, bisa pakai env system)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env tidak ditemukan, menggunakan environment system")
	}

	config.InitDB()

	r := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
	fmt.Printf("   App: %s\n", os.Getenv("APP_NAME"))
	fmt.Printf("   Env: %s\n", os.Getenv("APP_ENV"))

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("❌ Server gagal berjalan:", err)
	}
}
