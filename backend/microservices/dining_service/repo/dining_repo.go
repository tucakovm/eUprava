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

	// TEST DATA----------->
	if err := repo.SeedCanteens(); err != nil {
		return nil, fmt.Errorf("seeding cantines failed: %w", err)
	}

	if err := repo.SeedMenusForCanteenA(); err != nil {
		return nil, fmt.Errorf("seeding menus for Canteen A failed: %w", err)
	}

	if err := repo.SeedMealHistory("550e8400-e29b-41d4-a716-446655440001"); err != nil {
		fmt.Println("Seeding meal history failed:", err)
	}

	if err := repo.SeedMenuReviews("550e8400-e29b-41d4-a716-446655440001"); err != nil {
		fmt.Println("Seeding meal reviews failed:", err)
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

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_canteens_name ON canteens(name);`,

		// Meals table
		`CREATE TABLE IF NOT EXISTS meals (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price NUMERIC NOT NULL DEFAULT 0
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
    		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			breakfast_review INT DEFAULT 0,
			lunch_review INT DEFAULT 0,
			dinner_review INT DEFAULT 0
		);`,

		`CREATE TABLE IF NOT EXISTS popular_meals (
			id UUID PRIMARY KEY,
			menu_id UUID NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
			canteen_id UUID NOT NULL REFERENCES canteens(id) ON DELETE CASCADE,
			times_selected INT NOT NULL DEFAULT 0,
			UNIQUE (menu_id, canteen_id)
		);`,

		// Meal history table
		`CREATE TABLE IF NOT EXISTS meal_history (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			menu_id UUID NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
			selected_at TIMESTAMP NOT NULL DEFAULT NOW()
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
	today := time.Now().Format("2006-01-02")

	open1, _ := time.Parse(layout, today+" 08:00")
	close1, _ := time.Parse(layout, today+" 16:00")

	open2, _ := time.Parse(layout, today+" 09:00")
	close2, _ := time.Parse(layout, today+" 17:00")

	open3, _ := time.Parse(layout, today+" 10:00")
	close3, _ := time.Parse(layout, today+" 18:00")

	testCanteens := []domain.Canteen{
		{Id: uuid.MustParse("3fd5f339-8d75-4eee-81c9-25e1fd967faa"), Name: "Canteen A", Address: "Street 1", OpenAt: open1, CloseAt: close1},
		{Id: uuid.MustParse("b2c4d6e8-1234-5678-9abc-def012345678"), Name: "Canteen B", Address: "Street 2", OpenAt: open2, CloseAt: close2},
		{Id: uuid.MustParse("f9e8d7c6-abcd-1234-5678-9abc12345678"), Name: "Canteen C", Address: "Street 3", OpenAt: open3, CloseAt: close3},
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

func (r *DiningRepo) SeedMenusForCanteenA() error {
	// ID kantine A
	canteenAID := uuid.MustParse("3fd5f339-8d75-4eee-81c9-25e1fd967faa")

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	for _, day := range days {
		for i := 1; i <= 2; i++ { // po 2 menija dnevno
			menu := domain.Menu{
				Name:      fmt.Sprintf("Menu %s #%d", day, i),
				CanteenId: canteenAID,
				Weekday:   domain.Weekday(day),
				Breakfast: domain.Meal{
					Name:        fmt.Sprintf("Breakfast %s #%d", day, i),
					Description: "Test breakfast",
					Price:       3.5,
				},
				Lunch: domain.Meal{
					Name:        fmt.Sprintf("Lunch %s #%d", day, i),
					Description: "Test lunch",
					Price:       5.0,
				},
				Dinner: domain.Meal{
					Name:        fmt.Sprintf("Dinner %s #%d", day, i),
					Description: "Test dinner",
					Price:       6.5,
				},
			}

			if err := r.CreateMenu(&menu); err != nil {
				return fmt.Errorf("failed to create menu for %s: %w", day, err)
			}
		}
	}

	return nil
}

func (r *DiningRepo) SeedMenuReviews(userId string) error {
	// Uzimamo nekoliko postojećih menija iz baze
	rows, err := r.DB.Query(`SELECT id FROM menus LIMIT 5`)
	if err != nil {
		return fmt.Errorf("failed to fetch menus for seeding reviews: %w", err)
	}
	defer rows.Close()

	var menuIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		menuIds = append(menuIds, id)
	}

	if len(menuIds) == 0 {
		return fmt.Errorf("no menus found for seeding reviews")
	}

	// Dodajemo review za svaki meni
	for i, menuId := range menuIds {
		reviewId := uuid.New()                // Možeš zakucati fiksni UUID npr: uuid.MustParse("11111111-1111-1111-1111-11111111111" + strconv.Itoa(i))
		breakfastReview := int64((i % 5) + 1) // ocene od 1 do 5
		lunchReview := int64(((i + 1) % 5) + 1)
		dinnerReview := int64(((i + 2) % 5) + 1)

		_, err := r.DB.Exec(
			`INSERT INTO menu_reviews (id, menu_id, user_id, breakfast_review, lunch_review, dinner_review)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (menu_id, user_id) DO NOTHING`,
			reviewId, menuId, userId, breakfastReview, lunchReview, dinnerReview,
		)
		if err != nil {
			return fmt.Errorf("failed to insert menu review: %w", err)
		}
	}

	return nil
}

func (r *DiningRepo) SeedMealHistory(userId string) error {
	// Uzimamo nekoliko postojećih menija iz baze (iz Canteen A)
	rows, err := r.DB.Query(`SELECT id FROM menus LIMIT 3`)
	if err != nil {
		return fmt.Errorf("failed to fetch menus for seeding history: %w", err)
	}
	defer rows.Close()

	var menuIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		menuIds = append(menuIds, id)
	}

	if len(menuIds) == 0 {
		return fmt.Errorf("no menus found for seeding history")
	}

	// Dodajemo par unosa u istoriju za datog korisnika
	for i, menuId := range menuIds {
		historyId := uuid.New()
		selectedAt := time.Now().Add(-time.Duration(i) * time.Hour) // svaki sat unazad

		_, err := r.DB.Exec(
			`INSERT INTO meal_history (id, user_id, menu_id, selected_at)
             VALUES ($1, $2, $3, $4)`,
			historyId, userId, menuId, selectedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert meal history: %w", err)
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
	_, err := r.DB.Exec(
		`INSERT INTO meals (id, name, description, price) VALUES ($1, $2, $3, $4)`,
		m.Id, m.Name, m.Description, m.Price,
	)
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
	// 1. Kreiraj breakfast meal
	if err := r.CreateMeal(&menu.Breakfast); err != nil {
		return fmt.Errorf("failed to create breakfast meal: %w", err)
	}

	// 2. Kreiraj lunch meal
	if err := r.CreateMeal(&menu.Lunch); err != nil {
		return fmt.Errorf("failed to create lunch meal: %w", err)
	}

	// 3. Kreiraj dinner meal
	if err := r.CreateMeal(&menu.Dinner); err != nil {
		return fmt.Errorf("failed to create dinner meal: %w", err)
	}

	// 4. Kreiraj sam meni
	menu.Id = uuid.New()
	_, err := r.DB.Exec(
		`INSERT INTO menus (id, name, canteen_id, weekday, breakfast_id, lunch_id, dinner_id) 
     VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		menu.Id, menu.Name, menu.CanteenId, menu.Weekday,
		menu.Breakfast.Id, menu.Lunch.Id, menu.Dinner.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to create menu: %w", err)
	}

	return nil
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
		`INSERT INTO menu_reviews (id, menu_id, user_id, breakfast_review, lunch_review, dinner_review)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		review.Id, review.MenuId, review.UserId, review.BreakfastReview, review.LunchReview, review.DinnerReview,
	)
	return err
}

func (r *DiningRepo) UpdateMenuReview(review *domain.MenuReview) error {
	result, err := r.DB.Exec(
		`UPDATE menu_reviews SET breakfast_review=$1, lunch_review=$2, dinner_review=$3
		 WHERE id=$4`,
		review.BreakfastReview, review.LunchReview, review.DinnerReview, review.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("menu review with id %s not found", review.Id)
	}

	return nil
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

func (r *DiningRepo) GetMenusByCanteenID(canteenID string) ([]*domain.Menu, error) {
	rows, err := r.DB.Query(`
		SELECT 
			m.id, m.name, m.canteen_id, m.weekday,
			b.id, b.name, b.description, b.price,
			l.id, l.name, l.description, l.price,
			d.id, d.name, d.description, d.price
		FROM menus m
		LEFT JOIN meals b ON m.breakfast_id = b.id
		LEFT JOIN meals l ON m.lunch_id = l.id
		LEFT JOIN meals d ON m.dinner_id = d.id
		WHERE m.canteen_id = $1
		ORDER BY m.weekday;
	`, canteenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []*domain.Menu

	for rows.Next() {
		var menu domain.Menu
		var breakfast, lunch, dinner domain.Meal

		err := rows.Scan(
			&menu.Id, &menu.Name, &menu.CanteenId, &menu.Weekday,
			&breakfast.Id, &breakfast.Name, &breakfast.Description, &breakfast.Price,
			&lunch.Id, &lunch.Name, &lunch.Description, &lunch.Price,
			&dinner.Id, &dinner.Name, &dinner.Description, &dinner.Price,
		)
		if err != nil {
			return nil, err
		}

		menu.Breakfast = breakfast
		menu.Lunch = lunch
		menu.Dinner = dinner

		menus = append(menus, &menu)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *DiningRepo) DeleteMenuAndMealsByID(id string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var breakfastID, lunchID, dinnerID *string
	err = tx.QueryRow(
		`SELECT breakfast_id, lunch_id, dinner_id FROM menus WHERE id=$1`, id,
	).Scan(&breakfastID, &lunchID, &dinnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("menu with id %s not found", id)
		}
		return err
	}

	result, err := tx.Exec(`DELETE FROM menus WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("menu with id %s not found", id)
	}

	mealIDs := []*string{breakfastID, lunchID, dinnerID}
	for _, mealID := range mealIDs {
		if mealID != nil {
			_, err = tx.Exec(`DELETE FROM meals WHERE id=$1`, *mealID)
			if err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *DiningRepo) AddMealSelection(userId, menuId string) error {
	historyId := uuid.New()
	_, err := r.DB.Exec(
		`INSERT INTO meal_history (id, user_id, menu_id, selected_at)
         VALUES ($1, $2, $3, $4)`,
		historyId, userId, menuId, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert meal history: %w", err)
	}

	// Dobavi canteen_id za dati menu
	var canteenId string
	err = r.DB.QueryRow(`SELECT canteen_id FROM menus WHERE id = $1`, menuId).Scan(&canteenId)
	if err != nil {
		return fmt.Errorf("failed to fetch canteen_id: %w", err)
	}

	_, err = r.DB.Exec(
		`INSERT INTO popular_meals (id, menu_id, canteen_id, times_selected)
         VALUES ($1, $2, $3, 1)
         ON CONFLICT (menu_id, canteen_id)
         DO UPDATE SET times_selected = popular_meals.times_selected + 1`,
		uuid.New(), menuId, canteenId,
	)
	if err != nil {
		return fmt.Errorf("failed to update popular meals: %w", err)
	}

	return nil
}

func (r *DiningRepo) GetPopularMealsByCanteen(canteenId string, limit int) ([]domain.PopularMeal, error) {
	rows, err := r.DB.Query(
		`SELECT pm.menu_id, m.name, pm.times_selected
		 FROM popular_meals pm
		 JOIN menus m ON pm.menu_id = m.id
		 WHERE pm.canteen_id = $1
		 ORDER BY pm.times_selected DESC
		 LIMIT $2`, canteenId, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var popular []domain.PopularMeal
	for rows.Next() {
		var p domain.PopularMeal
		if err := rows.Scan(&p.MenuId, &p.MenuName, &p.TimesSelected); err != nil {
			return nil, err
		}
		popular = append(popular, p)
	}
	return popular, nil
}

func (r *DiningRepo) GetMealHistoryByUser(userId string) ([]domain.MealHistory, error) {
	rows, err := r.DB.Query(
		`SELECT mh.id, mh.menu_id, m.name, mh.selected_at
		 FROM meal_history mh
		 JOIN menus m ON mh.menu_id = m.id
		 WHERE mh.user_id = $1
		 ORDER BY mh.selected_at DESC`, userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.MealHistory
	for rows.Next() {
		var h domain.MealHistory
		if err := rows.Scan(&h.Id, &h.MenuId, &h.MenuName, &h.SelectedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (r *DiningRepo) GetMenuReviewByMenuAndUser(menuId, userId string) (*domain.MenuReview, error) {
	var mr domain.MenuReview
	err := r.DB.QueryRow(
		`SELECT id, menu_id, user_id, breakfast_review, lunch_review, dinner_review 
		 FROM menu_reviews WHERE menu_id = $1 AND user_id = $2`,
		menuId, userId,
	).Scan(&mr.Id, &mr.MenuId, &mr.UserId, &mr.BreakfastReview, &mr.LunchReview, &mr.DinnerReview)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Vraći nil ako ocena ne postoji
		}
		return nil, err
	}
	return &mr, nil
}

func (r *DiningRepo) GetMealHistoryWithReviewsByUser(userId string) ([]domain.MealHistoryWithReview, error) {
	rows, err := r.DB.Query(`
		SELECT 
			mh.id, 
			mh.menu_id, 
			m.name, 
			mh.selected_at,
			mr.id,
			mr.breakfast_review,
			mr.lunch_review, 
			mr.dinner_review
		FROM meal_history mh
		JOIN menus m ON mh.menu_id = m.id
		LEFT JOIN menu_reviews mr ON mh.menu_id = mr.menu_id AND mr.user_id = $1
		WHERE mh.user_id = $1
		ORDER BY mh.selected_at DESC`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.MealHistoryWithReview
	for rows.Next() {
		var h domain.MealHistoryWithReview
		var reviewId sql.NullString
		var breakfastReview, lunchReview, dinnerReview sql.NullInt64

		err := rows.Scan(
			&h.Id, &h.MenuId, &h.MenuName, &h.SelectedAt,
			&reviewId, &breakfastReview, &lunchReview, &dinnerReview,
		)
		if err != nil {
			return nil, err
		}

		// Ako postoji ocena, dodeli je
		if reviewId.Valid {
			h.Review = &domain.MenuReview{
				Id:              uuid.MustParse(reviewId.String),
				MenuId:          uuid.MustParse(h.MenuId.String()),
				UserId:          uuid.MustParse(userId),
				BreakfastReview: breakfastReview.Int64,
				LunchReview:     lunchReview.Int64,
				DinnerReview:    dinnerReview.Int64,
			}
		}

		history = append(history, h)
	}
	return history, nil
}
