package domain

import (
	"time"

	"github.com/google/uuid"
)

type Canteen struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	OpenAt  time.Time `json:"open_at"`
	CloseAt time.Time `json:"close_at"`
}

type StudentCard struct {
	ID        uuid.UUID `json:"id"`
	Stanje    float64   `json:"stanje"`
	StudentID uuid.UUID `json:"studentID"`
}

type CanteenDTO struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	OpenAt  string `json:"open_at"` // dolazi iz frontenda kao "HH:mm"
	CloseAt string `json:"close_at"`
}

type Canteens []Canteen

type Meal struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}

type Weekday string

const (
	Monday    Weekday = "Monday"
	Tuesday   Weekday = "Tuesday"
	Wednesday Weekday = "Wednesday"
	Thursday  Weekday = "Thursday"
	Friday    Weekday = "Friday"
	Saturday  Weekday = "Saturday"
	Sunday    Weekday = "Sunday"
)

type Menu struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CanteenId uuid.UUID `json:"canteen_id"`
	Weekday   Weekday   `json:"weekday"`
	Breakfast Meal      `json:"breakfast"`
	Lunch     Meal      `json:"lunch"`
	Dinner    Meal      `json:"dinner"`
}

type MenuDTO struct {
	Name      string    `json:"name"`
	CanteenId uuid.UUID `json:"canteen_id"`
	Weekday   Weekday   `json:"weekday"`
	Breakfast Meal      `json:"breakfast"`
	Lunch     Meal      `json:"lunch"`
	Dinner    Meal      `json:"dinner"`
}

type MenuReview struct {
	Id              uuid.UUID `json:"id"`
	MenuId          uuid.UUID `json:"menu_id"`
	UserId          uuid.UUID `json:"user_id"`
	BreakfastReview int64     `json:"breakfast_review"`
	LunchReview     int64     `json:"lunch_review"`
	DinnerReview    int64     `json:"dinner_review"`
}

type MenuReviewDTO struct {
	MenuId          string `json:"menu_id"`
	UserId          string `json:"user_id"`
	BreakfastReview int64  `json:"breakfast_review"`
	LunchReview     int64  `json:"lunch_review"`
	DinnerReview    int64  `json:"dinner_review"`
}

type PopularMeal struct {
	MenuId        string `json:"menu_id"`
	MenuName      string `json:"menu_name"`
	TimesSelected int    `json:"times_selected"`
}

type MealHistory struct {
	Id         string    `json:"id"`
	MenuId     string    `json:"menu_id"`
	MenuName   string    `json:"menu_name"`
	SelectedAt time.Time `json:"selected_at"`
}

type MealRoomHistory struct {
	UserName   string    `json:"user_name"`
	MenuId     string    `json:"menu_id"`
	MenuName   string    `json:"menu_name"`
	SelectedAt time.Time `json:"selected_at"`
}

type MealHistoryWithReview struct {
	Id         uuid.UUID   `json:"id"`
	MenuId     uuid.UUID   `json:"menu_id"`
	MenuName   string      `json:"menu_name"`
	SelectedAt time.Time   `json:"selected_at"`
	Review     *MenuReview `json:"review,omitempty"`
}

type MenuRating struct {
	MenuId   uuid.UUID `json:"menu_id"`
	MenuName string    `json:"menu_name"`
	Score    float64   `json:"score"`
}

type DiningRepository interface {
	GetAllCanteens() ([]Canteen, error)
	CreateCanteen(c *Canteen) error
	DeleteCanteenByID(id string) error
	CreateMeal(m *Meal) error
	UpdateMeal(m *Meal) error
	DeleteMealByID(id string) error
	CreateMenu(menu *Menu) error
	UpdateMenu(menu *Menu) error
	DeleteMenuByID(id string) error
	CreateMenuReview(review *MenuReview) error
	DeleteMenuReviewByID(id string) error
	GetMenuReviewByID(id string) (*MenuReview, error)
	GetMenuByID(id string) (*Menu, error)
	GetMealByID(id string) (*Meal, error)
	GetCanteenByID(id string) (*Canteen, error)
	GetMenusByCanteenID(canteenID string) ([]*Menu, error)
	DeleteMenuAndMealsByID(id string) error
	GetPopularMealsByCanteen(canteenId string, limit int) ([]PopularMeal, error)
	GetMealHistoryByUser(userId string) ([]MealHistory, error)
	GetMealHistoryWithReviewsByUser(userId string) ([]MealHistoryWithReview, error)
	UpdateMenuReview(review *MenuReview) error
	GetMenuWithMealsByID(menuId string) (*Menu, error)
	GetTop3RatedMeals(limit int) ([]MenuRating, error)
	CreateMealHistory(mh *MealHistory, userId string) error
	IncrementPopularMeal(menuId, canteenId uuid.UUID) error
	GetMealHistoryForUsernames(usernames []string) ([]MealRoomHistory, error)
}
