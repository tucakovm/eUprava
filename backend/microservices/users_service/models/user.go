package models

import (
	"fmt"

	"github.com/google/uuid"
)

// User struct represents a user in the database
type User struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsActive  bool      `json:"is_active"`
	Role      string    `json:"role"`
}

type UserDTO struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	Role      string    `json:"role"`
}

type RegisterRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ErrRespTmp struct {
	URL        string
	Method     string
	StatusCode int
}

func (e ErrRespTmp) Error() string {
	return fmt.Sprintf("temporary error [status code %d] for request: HTTP %s\t%s", e.StatusCode, e.Method, e.URL)
}

type ErrResp struct {
	URL        string
	Method     string
	StatusCode int
}

func (e ErrResp) Error() string {
	return fmt.Sprintf("error [status code %d] for request: HTTP %s\t%s", e.StatusCode, e.Method, e.URL)
}
