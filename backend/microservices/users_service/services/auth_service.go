package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"users_module/models"
	"users_module/repositories"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("email or username already exists")
	ErrUserDisabled       = errors.New("user disabled")
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (models.User, error)
	Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error)
	GetUser(ctx context.Context, id string) (*models.UserDTO, error)
}

type authService struct {
	repo      repositories.UserRepository
	jwtSecret []byte
}

func (s *authService) GetUser(ctx context.Context, id string) (*models.UserDTO, error) {
	cleanID := strings.Trim(strings.TrimSpace(id), "\"")

	if _, err := uuid.Parse(cleanID); err != nil {
		return nil, fmt.Errorf("invalid UUID format: %s", cleanID)
	}

	return s.repo.GetUserByID(ctx, cleanID)
}

func NewAuthService(repo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (models.User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	username := strings.TrimSpace(req.Username)
	if req.FirstName == "" || req.LastName == "" || username == "" || email == "" || len(req.Password) < 8 {
		return models.User{}, errors.New("missing fields or weak password (min 8 chars)")
	}

	u := models.User{
		Id:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  username,
		Email:     email,
		IsActive:  true,
		Role:      "user",
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	if err := s.repo.CreateUser(ctx, &u, hash); err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate") || strings.Contains(msg, "23505") {
			return models.User{}, ErrUserExists
		}
		return models.User{}, err
	}

	return u, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	ident := strings.TrimSpace(strings.ToLower(req.Identifier))
	if ident == "" || req.Password == "" {
		return models.LoginResponse{}, errors.New("missing fields")
	}

	u, hash, err := s.repo.GetByEmailOrUsername(ctx, ident)
	if err != nil {
		return models.LoginResponse{}, ErrInvalidCredentials
	}
	if !u.IsActive {
		return models.LoginResponse{}, ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte(req.Password)); err != nil {
		return models.LoginResponse{}, ErrInvalidCredentials
	}

	token, err := s.signJWT(*u)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{Token: token, User: *u}, nil
}

func (s *authService) signJWT(u models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   u.Id.String(),
		"usr":   u.Username,
		"role":  u.Role,
		"email": u.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.jwtSecret)
}
