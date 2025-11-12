package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email" binding:"required,email"`
	Password  string    `json:"password,omitempty" bson:"password" binding:"required,min=6"`
	Name      *string   `json:"name,omitempty" bson:"name,omitempty"`
	Bio       *string   `json:"bio,omitempty" bson:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	Email     string    `json:"email"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ProfileResponse struct {
	Email string  `json:"email"`
	Name  *string `json:"name"`
	Bio   *string `json:"bio"`
}

type UpdateProfileRequest struct {
	Name *string `json:"name"`
	Bio  *string `json:"bio"`
}

// HashPassword hashes the user password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
