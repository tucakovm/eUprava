package service

import "dining/domain"

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
