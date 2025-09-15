package repo

import (
	"database/sql"
	"dining/domain"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
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

	repo := &DiningRepo{DB: db}

	if err := repo.Migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	// Dodavanje test kantina
	if err := repo.SeedCanteens(); err != nil {
		return nil, fmt.Errorf("seeding cantines failed: %w", err)
	}

	return repo, nil
}

func (r *DiningRepo) Migrate() error {
	queries := []string{
		// Canteens table
		`CREATE TABLE IF NOT EXISTS canteens (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			open_at TIMESTAMP NOT NULL,
			close_at TIMESTAMP NOT NULL
		);`,

		// Dodaj UNIQUE index na name da ON CONFLICT radi
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_canteens_name ON canteens(name);`,

		// Meals table
		`CREATE TABLE IF NOT EXISTS meals (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT
		);`,

		// Menus table
		`CREATE TABLE IF NOT EXISTS menus (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			canteen_id UUID REFERENCES canteens(id) ON DELETE CASCADE,
			weekday TEXT NOT NULL,
			breakfast_id UUID REFERENCES meals(id) ON DELETE SET NULL,
			lunch_id UUID REFERENCES meals(id) ON DELETE SET NULL,
			dinner_id UUID REFERENCES meals(id) ON DELETE SET NULL
		);`,

		// Menu reviews table
		`CREATE TABLE IF NOT EXISTS menu_reviews (
			id UUID PRIMARY KEY,
			menu_id UUID REFERENCES menus(id) ON DELETE CASCADE,
			breakfast_review INT DEFAULT 0,
			lunch_review INT DEFAULT 0,
			dinner_review INT DEFAULT 0
		);`,
	}

	for _, q := range queries {
		if _, err := r.DB.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (r *DiningRepo) SeedCanteens() error {
	layout := "2006-01-02 15:04"
	today := time.Now().Format("2006-01-02") // danasnji datum

	open1, _ := time.Parse(layout, today+" 08:00")
	close1, _ := time.Parse(layout, today+" 16:00")

	open2, _ := time.Parse(layout, today+" 09:00")
	close2, _ := time.Parse(layout, today+" 17:00")

	open3, _ := time.Parse(layout, today+" 10:00")
	close3, _ := time.Parse(layout, today+" 18:00")

	testCanteens := []domain.Canteen{
		{Id: uuid.New(), Name: "Canteen A", Address: "Street 1", OpenAt: open1, CloseAt: close1},
		{Id: uuid.New(), Name: "Canteen B", Address: "Street 2", OpenAt: open2, CloseAt: close2},
		{Id: uuid.New(), Name: "Canteen C", Address: "Street 3", OpenAt: open3, CloseAt: close3},
	}

	for _, c := range testCanteens {
		_, err := r.DB.Exec(
			`INSERT INTO canteens (id, name, address, open_at, close_at)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (name) DO NOTHING`,
			c.Id, c.Name, c.Address, c.OpenAt, c.CloseAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
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

func (dr *DiningRepo) CreateCanteen(c *domain.Canteen) error {
	c.Id = uuid.New()

	_, err := dr.DB.Exec(
		`INSERT INTO canteens (id, name, address, open_at, close_at) VALUES ($1, $2, $3, $4, $5)`,
		c.Id, c.Name, c.Address, c.OpenAt, c.CloseAt,
	)
	return err
}

func (dr *DiningRepo) DeleteCanteenByID(id string) error {
	result, err := dr.DB.Exec(`DELETE FROM canteens WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("canteen with id %s not found", id)
	}

	return nil
}

func (r *DiningRepo) CreateMeal(m *domain.Meal) error {
	m.Id = uuid.New()
	_, err := r.DB.Exec(`INSERT INTO meals (id, name, description) VALUES ($1, $2, $3)`,
		m.Id, m.Name, m.Description)
	return err
}

func (r *DiningRepo) UpdateMeal(m *domain.Meal) error {
	_, err := r.DB.Exec(`UPDATE meals SET name=$1, description=$2 WHERE id=$3`,
		m.Name, m.Description, m.Id)
	return err
}

func (r *DiningRepo) DeleteMealByID(id string) error {
	result, err := r.DB.Exec(`DELETE FROM meals WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("meal with id %s not found", id)
	}
	return nil
}

func (r *DiningRepo) CreateMenu(menu *domain.Menu) error {
	menu.Id = uuid.New()
	_, err := r.DB.Exec(
		`INSERT INTO menus (id, name, canteen_id, weekday, breakfast_id, lunch_id, dinner_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		menu.Id, menu.Name, menu.CanteenId, menu.Weekday, menu.Breakfast.Id, menu.Lunch.Id, menu.Dinner.Id,
	)
	return err
}

func (r *DiningRepo) UpdateMenu(menu *domain.Menu) error {
	_, err := r.DB.Exec(
		`UPDATE menus SET name=$1, canteen_id=$2, weekday=$3, breakfast_id=$4, lunch_id=$5, dinner_id=$6
		 WHERE id=$7`,
		menu.Name, menu.CanteenId, menu.Weekday, menu.Breakfast.Id, menu.Lunch.Id, menu.Dinner.Id, menu.Id,
	)
	return err
}

func (r *DiningRepo) DeleteMenuByID(id string) error {
	result, err := r.DB.Exec(`DELETE FROM menus WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("menu with id %s not found", id)
	}
	return nil
}

func (r *DiningRepo) CreateMenuReview(review *domain.MenuReview) error {
	review.Id = uuid.New()
	_, err := r.DB.Exec(
		`INSERT INTO menu_reviews (id, menu_id, breakfast_review, lunch_review, dinner_review)
		 VALUES ($1, $2, $3, $4, $5)`,
		review.Id, review.MenuId, review.BreakfastReview, review.LunchReview, review.DinnerReview,
	)
	return err
}

func (r *DiningRepo) DeleteMenuReviewByID(id string) error {
	result, err := r.DB.Exec(`DELETE FROM menu_reviews WHERE id=$1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("menu review with id %s not found", id)
	}

	return nil
}

func (r *DiningRepo) GetCanteenByID(id string) (*domain.Canteen, error) {
	var c domain.Canteen
	err := r.DB.QueryRow(
		`SELECT id, name, address, open_at, close_at FROM canteens WHERE id = $1`, id,
	).Scan(&c.Id, &c.Name, &c.Address, &c.OpenAt, &c.CloseAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("canteen with id %s not found", id)
		}
		return nil, err
	}
	return &c, nil
}

func (r *DiningRepo) GetMealByID(id string) (*domain.Meal, error) {
	var m domain.Meal
	err := r.DB.QueryRow(
		`SELECT id, name, description FROM meals WHERE id = $1`, id,
	).Scan(&m.Id, &m.Name, &m.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("meal with id %s not found", id)
		}
		return nil, err
	}
	return &m, nil
}

func (r *DiningRepo) GetMenuByID(id string) (*domain.Menu, error) {
	var m domain.Menu
	err := r.DB.QueryRow(
		`SELECT id, name, canteen_id, weekday, breakfast_id, lunch_id, dinner_id 
		 FROM menus WHERE id = $1`, id,
	).Scan(&m.Id, &m.Name, &m.CanteenId, &m.Weekday, &m.Breakfast.Id, &m.Lunch.Id, &m.Dinner.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("menu with id %s not found", id)
		}
		return nil, err
	}
	return &m, nil
}

func (r *DiningRepo) GetMenuReviewByID(id string) (*domain.MenuReview, error) {
	var mr domain.MenuReview
	err := r.DB.QueryRow(
		`SELECT id, menu_id, breakfast_review, lunch_review, dinner_review 
		 FROM menu_reviews WHERE id = $1`, id,
	).Scan(&mr.Id, &mr.MenuId, &mr.BreakfastReview, &mr.LunchReview, &mr.DinnerReview)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("menu review with id %s not found", id)
		}
		return nil, err
	}
	return &mr, nil
}
