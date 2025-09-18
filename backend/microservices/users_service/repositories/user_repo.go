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

	// inicijalni korisnici sa hardkodovanim UUID-ovima
	users := []struct {
		ID        string
		FirstName string
		LastName  string
		Username  string
		Email     string
		Password  string
		Role      string
	}{
		{"550e8400-e29b-41d4-a716-446655440000", "AdminFN", "AdminLN", "admin", "admin@example.com", "Admin123!", "admin"},
		{"550e8400-e29b-41d4-a716-446655440001", "Nikola", "Nikolic", "nikola123", "nikola@example.com", "nikola123", "student"},
		{"550e8400-e29b-41d4-a716-446655440002", "Jovana", "Petrovic", "jovana123", "jovana@example.com", "jovana123", "student"}, // bez space
		{"550e8400-e29b-41d4-a716-446655440003", "Marko", "Ilic", "marko123", "marko@example.com", "marko123", "student"},
	}

	for _, u := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		userID, _ := uuid.Parse(u.ID)

		// INSERT u users (idempotentno po username)
		_, err := r.DB.ExecContext(ctx, `
			INSERT INTO users (id, firstname, lastname, username, email, password_hash, is_active, role)
			VALUES ($1,$2,$3,$4,$5,$6,true,$7)
			ON CONFLICT (username) DO NOTHING;
		`, userID, u.FirstName, u.LastName, u.Username, u.Email, hash, u.Role)
		if err != nil {
			return err
		}

		// Ako je student — kreiraj i red u housing.student (idempotentno)
		if u.Role == "student" {
			if _, err := r.DB.ExecContext(ctx, `
				INSERT INTO student (ime, prezime, username, soba_id)
				VALUES ($1, $2, $3, NULL)
				ON CONFLICT (username) DO NOTHING;
			`, u.FirstName, u.LastName, u.Username); err != nil {
				return err
			}
		}
	}

	return nil
}

// CreateUser: kreira user-a i, ako je rola 'student', kreira i red u housing.student.
// Sve u JEDNOJ transakciji radi konzistentnosti.
func (r *UserRepository) CreateUser(ctx context.Context, u *models.User, passwordHash []byte) error {
	const qInsertUser = `
INSERT INTO users (id, firstname, lastname, username, email, password_hash, is_active, role)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, firstname, lastname, username, email, is_active, role;
`
	const qInsertStudentIfRole = `
INSERT INTO student (ime, prezime, username, soba_id)
VALUES ($1, $2, $3, NULL)
ON CONFLICT (username) DO NOTHING;
`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1) upiši user-a
	err = tx.QueryRowContext(ctx, qInsertUser,
		u.Id, u.FirstName, u.LastName, u.Username, u.Email, passwordHash, u.IsActive, u.Role,
	).Scan(&u.Id, &u.FirstName, &u.LastName, &u.Username, &u.Email, &u.IsActive, &u.Role)
	if err != nil {
		return err
	}

	// 2) ako je student — kreiraj red i u 'student'
	if u.Role == "student" {
		if _, err = tx.ExecContext(ctx, qInsertStudentIfRole, u.FirstName, u.LastName, u.Username); err != nil {
			return err
		}
	}

	// 3) commit
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
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

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.UserDTO, error) {
	const q = `SELECT id, firstname, lastname, username, email, is_active, role 
               FROM users WHERE id = $1 LIMIT 1`

	var (
		u     models.UserDTO
		idStr string
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, q, id).Scan(
		&idStr,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.Email,
		&u.IsActive,
		&u.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		fmt.Printf("Database scan error: %v\n", err)
		return nil, fmt.Errorf("database error: %v", err)
	}

	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID in database: %v", err)
	}
	u.Id = parsedID

	return &u, nil
}
