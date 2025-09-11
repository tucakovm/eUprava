package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"users_module/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository() (*UserRepository, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbHost == "" || dbPort == "" || dbName == "" || dbUser == "" {
		return nil, fmt.Errorf("env variables DB_HOST, DB_PORT, DB_NAME, DB_USER must be set")
	}

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &UserRepository{DB: db}

	// pokreni migraciju i seed
	if err := repo.migrateAndSeed(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *UserRepository) migrateAndSeed() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// napravi tabelu ako ne postoji
	_, err := r.DB.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  firstname  STRING NOT NULL,
  lastname   STRING NOT NULL,
  username   STRING UNIQUE NOT NULL,
  email      STRING UNIQUE NOT NULL,
  password_hash BYTES NOT NULL,
  is_active  BOOL NOT NULL DEFAULT true,
  role       STRING NOT NULL DEFAULT 'user',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`)
	if err != nil {
		return err
	}

	// inicijalni korisnici
	users := []struct {
		FirstName string
		LastName  string
		Username  string
		Email     string
		Password  string
		Role      string
	}{
		{"Admin", "User", "admin", "admin@example.com", "Admin123!", "admin"},
		{"Test", "User", "testuser", "test@example.com", "Test123!", "user"},
	}

	for _, u := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

		_, err := r.DB.ExecContext(ctx, `
INSERT INTO users (id, firstname, lastname, username, email, password_hash, is_active, role)
VALUES ($1,$2,$3,$4,$5,$6,true,$7)
ON CONFLICT (username) DO NOTHING;
		`,
			uuid.New(), u.FirstName, u.LastName, u.Username, u.Email, hash, u.Role,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) CreateUser(ctx context.Context, u *models.User, passwordHash []byte) error {
	const q = `
INSERT INTO users (id, firstname, lastname, username, email, password_hash, is_active, role)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, firstname, lastname, username, email, is_active, role;
`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, q,
		u.Id, u.FirstName, u.LastName, u.Username, u.Email, passwordHash, u.IsActive, u.Role,
	).Scan(&u.Id, &u.FirstName, &u.LastName, &u.Username, &u.Email, &u.IsActive, &u.Role)
}

func (r *UserRepository) GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, []byte, error) {
	const q = `
SELECT id, firstname, lastname, username, email, password_hash, is_active, role
FROM users
WHERE email = $1 OR username = $2
LIMIT 1;
`
	var (
		u   models.User
		pwd []byte
	)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, q, identifier, identifier).Scan(
		&u.Id, &u.FirstName, &u.LastName, &u.Username, &u.Email, &pwd, &u.IsActive, &u.Role,
	)
	if err != nil {
		return nil, nil, err
	}
	return &u, pwd, nil
}
