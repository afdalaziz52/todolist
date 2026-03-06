package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/go-playground/validator/v10"
)

var DB *sql.DB
var Validate *validator.Validate

func InitDB() {
    var err error
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal(err)
    }

    Validate = validator.New()

    fmt.Println("✅ Database connected!")
    fmt.Printf("   Database: %s\n", os.Getenv("DB_NAME"))
}