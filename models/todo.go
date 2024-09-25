package models

import "time"

type TodoData struct {
	ID          string    `json:"id" db:"id"`                   // UUID型
	UserID      string    `json:"user_id" db:"user_id"`         // ユーザーID(UUID型)
	Description string    `json:"description" db:"description"` // 説明
	Completed   bool      `json:"completed" db:"completed"`     // 完了フラグ
	CreatedAt   time.Time `json:"created_at" db:"created_at"`   // タイムスタンプ
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`   // タイムスタンプ
}
