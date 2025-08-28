package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	GoogleID  string    `json:"google_id" db:"google_id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Picture   string    `json:"picture" db:"picture"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type News struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Content     string    `json:"content" db:"content"`
	Link        *string   `json:"link,omitempty" db:"link"`
	PhotoURL    *string   `json:"photo_url,omitempty" db:"photo_url"`
	Status      string    `json:"status" db:"status"`
	Explanation *string   `json:"explanation,omitempty" db:"explanation"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type NewsSubmission struct {
	Content  string  `json:"content" binding:"required"`
	Link     *string `json:"link,omitempty"`
	PhotoURL *string `json:"photo_url,omitempty"`
}

type NewsVerification struct {
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	Explanation string    `json:"explanation"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type LoginResponse struct {
	AuthURL string `json:"auth_url"`
}

type AuthCallbackResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
