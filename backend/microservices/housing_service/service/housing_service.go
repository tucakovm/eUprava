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
	Kartica repository.StudentskaKarticaRepository
}

func New(
	db *sql.DB,
	dom repository.DomRepository,
	soba repository.SobaRepository,
	student repository.StudentRepository,
	rec repository.RecenzijaRepository,
	kvar repository.KvarRepository,
	kartica repository.StudentskaKarticaRepository,
) *Services {
	return &Services{
		DB:      db,
		Dom:     dom,
		Soba:    soba,
		Student: student,
		Rec:     rec,
		Kvar:    kvar,
		Kartica: kartica,
	}
}

func ctxTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultTimeout)
}

/* ======================= Domovi ======================= */

func (s *Services) GetDom(ctx context.Context, domID uuid.UUID) (domain.Dom, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	return s.Dom.Get(ctx, s.DB, domID)
}

func (s *Services) GetAllDomovi(ctx context.Context) ([]domain.Dom, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	return s.Dom.GetAll(ctx, s.DB)
}

/* ======================= Studenti ======================= */

func (s *Services) CreateStudent(ctx context.Context, ime, prezime, username string) (domain.Student, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	st := domain.Student{
		ID:       uuid.New(),
		Ime:      ime,
		Prezime:  prezime,
		Username: username,
		SobaID:   nil,
	}
	if err := s.Student.Create(ctx, s.DB, &st); err != nil {
		return domain.Student{}, err
	}
	return st, nil
}

// Upis postojeceg studenta (po username) u sobu
func (s *Services) UpisiPostojecegStudentaUSobu(ctx context.Context, domID uuid.UUID, brojSobe, username string) (domain.Student, error) {
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

	// 1) Zaključaj sobu
	soba, err := s.Soba.GetByBroj(ctx, tx, domID, brojSobe, true)
	if err != nil {
		return domain.Student{}, err
	}

	// 2) Proveri popunjenost
	postojeci, err := s.Student.ListBySoba(ctx, tx, soba.ID)
	if err != nil {
		return domain.Student{}, err
	}
	if len(postojeci) >= soba.Kapacitet {
		return domain.Student{}, errors.New("soba je popunjena (nema slobodnih mesta)")
	}

	// 3) Nadji studenta po username
	st, err := s.Student.GetByUsername(ctx, tx, username)
	if err != nil {
		return domain.Student{}, errors.New("student ne postoji")
	}
	if st.SobaID != nil {
		return domain.Student{}, errors.New("student je već dodeljen nekoj sobi")
	}

	// 4) Povezi
	if err = s.Student.AssignToSoba(ctx, tx, st.ID, soba.ID); err != nil {
		return domain.Student{}, err
	}

	// 5) Ako je poslednje mesto, zatvori sobu
	if len(postojeci)+1 >= soba.Kapacitet {
		if err = s.Soba.SetSlobodna(ctx, tx, soba.ID, false); err != nil {
			return domain.Student{}, err
		}
	}

	// 6) Commit
	if err = tx.Commit(); err != nil {
		return domain.Student{}, err
	}

	st.SobaID = &soba.ID
	return st, nil
}

// Oslobodi sobu od studenta po ID-u (ovo ostaje po ID jer je to interni poziv)
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
		return tx.Commit()
	}

	if err = s.Student.UnassignSoba(ctx, tx, studentID); err != nil {
		return err
	}

	soba, err := s.Soba.Get(ctx, tx, *st.SobaID)
	if err != nil {
		return err
	}
	preostali, err := s.Student.ListBySoba(ctx, tx, soba.ID)
	if err != nil {
		return err
	}

	shouldBeFree := len(preostali) < soba.Kapacitet
	if err = s.Soba.SetSlobodna(ctx, tx, soba.ID, shouldBeFree); err != nil {
		return err
	}

	return tx.Commit()
}

/* ======================= Recenzije ======================= */

func (s *Services) DodajRecenziju(ctx context.Context, sobaID uuid.UUID, autorUsername string, ocena int, komentar *string) (domain.RecenzijaSobe, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	r := domain.RecenzijaSobe{
		ID:            uuid.New(),
		Ocena:         ocena,
		Komentar:      komentar,
		SobaID:        sobaID,
		AutorUsername: autorUsername,
	}
	if err := s.Rec.Create(ctx, s.DB, &r); err != nil {
		return domain.RecenzijaSobe{}, err
	}
	return r, nil
}

/* ======================= Kvarovi ======================= */

func (s *Services) PrijaviKvar(ctx context.Context, sobaID uuid.UUID, prijavioUsername, opis string) (domain.Kvar, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	k := domain.Kvar{
		ID:               uuid.New(),
		Opis:             opis,
		Status:           domain.StatusPrijavljen,
		SobaID:           sobaID,
		PrijavioUsername: prijavioUsername,
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

/* ======================= Sobe / DTO ======================= */

func (s *Services) GetSoba(ctx context.Context, sobaID uuid.UUID) (domain.Soba, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()
	return s.Soba.Get(ctx, s.DB, sobaID)
}

func (s *Services) GetSobaDetail(ctx context.Context, sobaID uuid.UUID) (domain.Soba, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	soba, err := s.Soba.Get(ctx, s.DB, sobaID)
	if err != nil {
		return domain.Soba{}, err
	}
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

	soba.Studenti = sts
	soba.Recenzije = recs
	soba.Kvarovi = kvars
	return soba, nil
}

/* ============ Studentske kartice po username ============ */

func (s *Services) KreirajStudentskuKarticuAkoNema(ctx context.Context, studentUsername string) (domain.StudentskaKartica, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()
	return s.Kartica.CreateIfNotExistsByUsername(ctx, s.DB, studentUsername)
}

func (s *Services) GetStudentskaKarticaByStudent(ctx context.Context, studentUsername string) (domain.StudentskaKartica, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()
	return s.Kartica.GetByStudentUsername(ctx, s.DB, studentUsername)
}

func (s *Services) AzurirajStanjeStudentskeKartice(ctx context.Context, studentUsername string, delta float64) (domain.StudentskaKartica, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()
	return s.Kartica.UpdateStanjeByUsername(ctx, s.DB, studentUsername, delta)
}

/* ======================= Slobodne sobe ======================= */

func (s *Services) ListSlobodneSobe(ctx context.Context, domID uuid.UUID) ([]domain.Soba, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()
	return s.Soba.ListSlobodne(ctx, s.DB, domID)
}

// service/services.go
func (s *Services) IsStudentAssignedToAnySoba(ctx context.Context, studentID string) (bool, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()
	return s.Student.IsAssignedToAnySoba(ctx, s.DB, studentID)
}
