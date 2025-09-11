package repo

import (
	"database/sql"
	"dining/domain"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DiningRepo struct {
	DB *sql.DB
}

func NewDiningRepo() (*DiningRepo, error) {

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

	return &DiningRepo{DB: db}, nil
}

func (dr *DiningRepo) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS canteens (
		id STRING PRIMARY KEY,
		name STRING NOT NULL,
		address STRING NOT NULL,
		open_at TIMESTAMP NOT NULL,
		close_at TIMESTAMP NOT NULL
	);`
	_, err := dr.DB.Exec(query)
	return err
}

func (r *DiningRepo) GetAllCanteens() ([]domain.Canteen, error) {
	rows, err := r.DB.Query(`SELECT id, name, address, open_at, close_at FROM canteens`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var canteens []domain.Canteen
	for rows.Next() {
		var c domain.Canteen
		if err := rows.Scan(&c.Id, &c.Name, &c.Address, &c.OpenAt, &c.CloseAt); err != nil {
			return nil, err
		}
		canteens = append(canteens, c)
	}

	return canteens, nil
}
