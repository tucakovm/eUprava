package domain

import "time"

type Canteen struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	OpenAt  time.Time `json:"open_at"`
	CloseAt time.Time `json:"close_at"`
}

type DiningRepository interface {
	GetAllCanteens() ([]Canteen, error)
}
