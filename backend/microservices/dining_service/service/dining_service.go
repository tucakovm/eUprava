package service

import (
	"dining/domain"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DiningService struct {
	repo domain.DiningRepository
}

func NewDiningService(repo domain.DiningRepository) *DiningService {
	return &DiningService{
		repo: repo,
	}
}

func (ds *DiningService) GetAllCanteens() ([]domain.Canteen, error) {
	return ds.repo.GetAllCanteens()
}

func (ds *DiningService) GetCanteen(id string) (*domain.Canteen, error) {
	return ds.repo.GetCanteenByID(id)
}

func (ds *DiningService) DeleteCanteen(id string) error {
	return ds.repo.DeleteCanteenByID(id)
}

func (ds *DiningService) CreateCanteen(c *domain.Canteen) error {
	return ds.repo.CreateCanteen(c)
}

func (ds *DiningService) GetMenusByCanteenID(id string) ([]*domain.Menu, error) {
	return ds.repo.GetMenusByCanteenID(id)
}

func (ds *DiningService) CreateMenu(c *domain.MenuDTO) error {
	m := &domain.Menu{
		Name:      c.Name,
		CanteenId: c.CanteenId,
		Weekday:   c.Weekday,
		Breakfast: domain.Meal{
			Name:        c.Breakfast.Name,
			Description: c.Breakfast.Description,
			Price:       c.Breakfast.Price,
		},
		Lunch: domain.Meal{
			Name:        c.Lunch.Name,
			Description: c.Lunch.Description,
			Price:       c.Lunch.Price,
		},
		Dinner: domain.Meal{
			Name:        c.Dinner.Name,
			Description: c.Dinner.Description,
			Price:       c.Dinner.Price,
		},
	}
	return ds.repo.CreateMenu(m)
}

func (ds *DiningService) DeleteMenu(id string) error {
	return ds.repo.DeleteMenuAndMealsByID(id)
}

func (ds *DiningService) GetPopularMenus(id string) ([]domain.PopularMeal, error) {
	return ds.repo.GetPopularMealsByCanteen(id, 5)
}

func (ds *DiningService) GetMealHistory(id string) ([]domain.MealHistory, error) {
	cleanID := strings.Trim(strings.TrimSpace(id), "\"")

	if _, err := uuid.Parse(cleanID); err != nil {
		return nil, fmt.Errorf("invalid UUID format: %s", cleanID)
	}

	return ds.repo.GetMealHistoryByUser(cleanID)
}

func (ds *DiningService) GetMealHistoryWithReviewsByUser(id string) ([]domain.MealHistoryWithReview, error) {
	return ds.repo.GetMealHistoryWithReviewsByUser(id)
}

func (ds *DiningService) UpdateMenuReview(r *domain.MenuReview) error {
	return ds.repo.UpdateMenuReview(r)
}

func (ds *DiningService) CreateMenuReview(review *domain.MenuReview) error {
	return ds.repo.CreateMenuReview(review)
}

func (ds *DiningService) GetMenu(id string) (*domain.Menu, error) {
	return ds.repo.GetMenuWithMealsByID(id)
}

func (ds *DiningService) GetTopRatedMeals() ([]domain.MenuRating, error) {
	return ds.repo.GetTop3RatedMeals(3)
}

func (ds *DiningService) CreateMealHistory(mh *domain.MealHistory, userId string) error {
	return ds.repo.CreateMealHistory(mh, userId)
}

func (ds *DiningService) IncrementPopularMeal(menuId, canteenId uuid.UUID) error {
	return ds.repo.IncrementPopularMeal(menuId, canteenId)
}

func (ds *DiningService) GetMealHistoryForUsernames(usernames []string) ([]domain.MealRoomHistory, error) {
	return ds.repo.GetMealHistoryForUsernames(usernames)
}

// NEW: Svi meniji iz svih kantina za današnji dan (lokalno vreme servera)
func (ds *DiningService) GetMenusForToday() ([]*domain.Menu, error) {
	// mapiranje time.Weekday -> domain.Weekday
	switch time.Now().Local().Weekday() {
	case time.Monday:
		return ds.repo.GetMenusByWeekday(domain.Monday)
	case time.Tuesday:
		return ds.repo.GetMenusByWeekday(domain.Tuesday)
	case time.Wednesday:
		return ds.repo.GetMenusByWeekday(domain.Wednesday)
	case time.Thursday:
		return ds.repo.GetMenusByWeekday(domain.Thursday)
	case time.Friday:
		return ds.repo.GetMenusByWeekday(domain.Friday)
	case time.Saturday:
		return ds.repo.GetMenusByWeekday(domain.Saturday)
	case time.Sunday:
		return ds.repo.GetMenusByWeekday(domain.Sunday)
	default:
		// fallback, ne bi trebalo nikad da se desi
		return ds.repo.GetMenusByWeekday(domain.Weekday(time.Now().Weekday().String()))
	}

}
