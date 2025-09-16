package service

import (
	"dining/domain"
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
