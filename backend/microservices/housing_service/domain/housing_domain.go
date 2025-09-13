package domain

import (
	"github.com/google/uuid"
)

type Dom struct {
	ID     uuid.UUID `json:"id"`
	Naziv  string    `json:"naziv"`
	Adresa string    `json:"adresa"`
	Sobe   []Soba    `json:"sobe,omitempty"`
}

type Soba struct {
	ID        uuid.UUID       `json:"id"`
	Broj      string          `json:"broj"`
	Slobodna  bool            `json:"slobodna"`
	DomID     uuid.UUID       `json:"domId"` 
	Studenti  []Student       `json:"studenti,omitempty"`   
	Recenzije []RecenzijaSobe `json:"recenzije,omitempty"`
	Kvarovi   []Kvar          `json:"kvarovi,omitempty"`
}

type Student struct {
	ID      uuid.UUID  `json:"id"`
	Ime     string     `json:"ime"`
	Prezime string     `json:"prezime"`
	SobaID  *uuid.UUID `json:"sobaId,omitempty"` 
}

type RecenzijaSobe struct {
	ID       uuid.UUID `json:"id"`
	Ocena    int       `json:"ocena"`
	Komentar *string   `json:"komentar,omitempty"`
	SobaID   uuid.UUID `json:"sobaId"`  
	AutorID  uuid.UUID `json:"autorId"` 
}

type StatusKvara string

const (
	StatusPrijavljen StatusKvara = "prijavljen"
	StatusUToku      StatusKvara = "u_toku"
	StatusResen      StatusKvara = "resen"
)

type Kvar struct {
	ID         uuid.UUID   `json:"id"`
	Opis       string      `json:"opis"`
	Status     StatusKvara `json:"status"`
	SobaID     uuid.UUID   `json:"sobaId"`
	PrijavioID uuid.UUID   `json:"prijavioId"`  
}
