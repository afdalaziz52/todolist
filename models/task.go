package models

import "time"

type Task struct {
    ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    Title          string    `json:"title"`
    Category       string    `json:"category"`
    CustomCategory *string   `json:"custom_category"` 
    Status         string    `json:"status"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

// untuk request tambah & edit task
type TaskRequest struct {
    Title          string  `json:"title"           validate:"required"`
    Category       string  `json:"category"        validate:"required,oneof=work study personal finance social other"`
    CustomCategory *string `json:"custom_category"`
}