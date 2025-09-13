package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"housing/domain"
	"housing/repository"
)

const defaultTimeout = 5 * time.Second

type Services struct {
	DB      *sql.DB
	Dom     repository.DomRepository
	Soba    repository.SobaRepository
	Student repository.StudentRepository
	Rec     repository.RecenzijaRepository
	Kvar    repository.KvarRepository
}

func New(db *sql.DB, dom repository.DomRepository, soba repository.SobaRepository,
	student repository.StudentRepository, rec repository.RecenzijaRepository, kvar repository.KvarRepository) *Services {
	return &Services{
		DB:      db,
		Dom:     dom,
		Soba:    soba,
		Student: student,
		Rec:     rec,
		Kvar:    kvar,
	}
}

// Helper za kontekst sa timeout-om
func ctxTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultTimeout)
}

/* =======================
   Studentske operacije
   ======================= */

// CreateStudent kreira studenta bez dodele sobe (SobaID = nil).
func (s *Services) CreateStudent(ctx context.Context, ime, prezime string) (domain.Student, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	st := domain.Student{
		ID:      uuid.New(),
		Ime:     ime,
		Prezime: prezime,
		SobaID:  nil,
	}
	if err := s.Student.Create(ctx, s.DB, &st); err != nil {
		return domain.Student{}, err
	}
	return st, nil
}

// UpisiStudentaUSobu kreira studenta i dodeljuje ga zadatoj sobi (broj u okviru doma).
// Radi u transakciji i zaključava red sobe (SELECT ... FOR UPDATE).
func (s *Services) UpisiStudentaUSobu(ctx context.Context, domID uuid.UUID, brojSobe, ime, prezime string) (domain.Student, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.Student{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1) Zaključaj sobu po broju u domu
	soba, err := s.Soba.GetByBroj(ctx, tx, domID, brojSobe, true)
	if err != nil {
		return domain.Student{}, err
	}
	if !soba.Slobodna {
		return domain.Student{}, errors.New("soba nije slobodna")
	}

	// 2) Kreiraj studenta
	st := domain.Student{
		ID:      uuid.New(),
		Ime:     ime,
		Prezime: prezime,
	}
	if err = s.Student.Create(ctx, tx, &st); err != nil {
		return domain.Student{}, err
	}

	// 3) Poveži i označi sobu zauzetom
	if err = s.Student.AssignToSoba(ctx, tx, st.ID, soba.ID); err != nil {
		return domain.Student{}, err
	}
	if err = s.Soba.SetSlobodna(ctx, tx, soba.ID, false); err != nil {
		return domain.Student{}, err
	}

	// 4) Commit
	if err = tx.Commit(); err != nil {
		return domain.Student{}, err
	}
	st.SobaID = &soba.ID
	return st, nil
}

// OslobodiSobu uklanja dodelu sobe studentu i označava sobu kao slobodnu.
func (s *Services) OslobodiSobu(ctx context.Context, studentID uuid.UUID) (err error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	st, err := s.Student.Get(ctx, tx, studentID)
	if err != nil {
		return err
	}
	if st.SobaID == nil {
		// već nema sobu — no-op
		return tx.Commit()
	}

	if err = s.Student.UnassignSoba(ctx, tx, studentID); err != nil {
		return err
	}
	if err = s.Soba.SetSlobodna(ctx, tx, *st.SobaID, true); err != nil {
		return err
	}

	return tx.Commit()
}

/* =======================
   Recenzije
   ======================= */

func (s *Services) DodajRecenziju(ctx context.Context, sobaID, autorID uuid.UUID, ocena int, komentar *string) (domain.RecenzijaSobe, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	r := domain.RecenzijaSobe{
		ID:       uuid.New(),
		Ocena:    ocena,
		Komentar: komentar,
		SobaID:   sobaID,
		AutorID:  autorID,
	}
	if err := s.Rec.Create(ctx, s.DB, &r); err != nil {
		return domain.RecenzijaSobe{}, err
	}
	return r, nil
}

/* =======================
   Kvarovi
   ======================= */

func (s *Services) PrijaviKvar(ctx context.Context, sobaID, prijavioID uuid.UUID, opis string) (domain.Kvar, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	k := domain.Kvar{
		ID:         uuid.New(),
		Opis:       opis,
		Status:     domain.StatusPrijavljen,
		SobaID:     sobaID,
		PrijavioID: prijavioID,
	}
	if err := s.Kvar.Create(ctx, s.DB, &k); err != nil {
		return domain.Kvar{}, err
	}
	return k, nil
}

func (s *Services) PromeniStatusKvara(ctx context.Context, kvarID uuid.UUID, status domain.StatusKvara) error {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()
	return s.Kvar.UpdateStatus(ctx, s.DB, kvarID, status)
}

/* =======================
   Sobe / DTO-i
   ======================= */

// GetSoba vraća osnovne podatke sobe.
func (s *Services) GetSoba(ctx context.Context, sobaID uuid.UUID) (domain.Soba, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()
	return s.Soba.Get(ctx, s.DB, sobaID)
}

// GetSobaDetail vraća sobu sa studentima, recenzijama i kvarovima (DTO za API).
func (s *Services) GetSobaDetail(ctx context.Context, sobaID uuid.UUID) (domain.Soba, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	// 1) osnovni entitet
	soba, err := s.Soba.Get(ctx, s.DB, sobaID)
	if err != nil {
		return domain.Soba{}, err
	}

	// 2) rela liste
	sts, err := s.Student.ListBySoba(ctx, s.DB, sobaID)
	if err != nil {
		return domain.Soba{}, err
	}
	recs, err := s.Rec.ListBySoba(ctx, s.DB, sobaID)
	if err != nil {
		return domain.Soba{}, err
	}
	kvars, err := s.Kvar.ListBySoba(ctx, s.DB, sobaID)
	if err != nil {
		return domain.Soba{}, err
	}

	// 3) popuni DTO polja
	soba.Studenti = sts
	soba.Recenzije = recs
	soba.Kvarovi = kvars
	return soba, nil
}
