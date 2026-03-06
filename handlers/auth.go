package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "github.com/afdalaziz52/to-do-list/config"
    "github.com/afdalaziz52/to-do-list/models"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "os"
)

// POST /auth/register
func Register(w http.ResponseWriter, r *http.Request) {
    var req models.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, 400, map[string]any{"status": "error", "message": "Format tidak valid"})
        return
    }

    if err := config.Validate.Struct(req); err != nil {
        writeJSON(w, 400, map[string]any{"status": "error", "message": err.Error()})
        return
    }

    // cek email sudah terdaftar
    var exists int
    config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&exists)
    if exists > 0 {
        writeJSON(w, 409, map[string]any{"status": "error", "message": "Email sudah terdaftar"})
        return
    }

    // hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        writeJSON(w, 500, map[string]any{"status": "error", "message": "Gagal memproses password"})
        return
    }

    _, err = config.DB.Exec(
        "INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
        req.Name, req.Email, string(hashed))
    if err != nil {
        writeJSON(w, 500, map[string]any{"status": "error", "message": "Gagal mendaftarkan user"})
        return
    }

    writeJSON(w, 201, map[string]any{"status": "success", "message": "Registrasi berhasil"})
}

// POST /auth/login
func Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, 400, map[string]any{"status": "error", "message": "Format tidak valid"})
        return
    }

    if err := config.Validate.Struct(req); err != nil {
        writeJSON(w, 400, map[string]any{"status": "error", "message": err.Error()})
        return
    }

    // cari user by email
    var user models.User
    err := config.DB.QueryRow(
        "SELECT id, name, email, password FROM users WHERE email = ?", req.Email).
        Scan(&user.ID, &user.Name, &user.Email, &user.Password)

    if err == sql.ErrNoRows {
        writeJSON(w, 401, map[string]any{"status": "error", "message": "Email atau password salah"})
        return
    }
    if err != nil {
        writeJSON(w, 500, map[string]any{"status": "error", "message": "Gagal login"})
        return
    }

    // cek password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        writeJSON(w, 401, map[string]any{"status": "error", "message": "Email atau password salah"})
        return
    }

    // buat JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": user.ID,
        "exp":    time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        writeJSON(w, 500, map[string]any{"status": "error", "message": "Gagal membuat token"})
        return
    }

    writeJSON(w, 200, map[string]any{
        "status":  "success",
        "message": "Login berhasil",
        "token":   tokenString,
        "user": map[string]any{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
        },
    })
}
