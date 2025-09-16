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
	// Novo: repo za studentske kartice
	Kartica repository.StudentskaKarticaRepository
}

func New(
	db *sql.DB,
	dom repository.DomRepository,
	soba repository.SobaRepository,
	student repository.StudentRepository,
	rec repository.RecenzijaRepository,
	kvar repository.KvarRepository,
	kartica repository.StudentskaKarticaRepository, // novo
) *Services {
	return &Services{
		DB:      db,
		Dom:     dom,
		Soba:    soba,
		Student: student,
		Rec:     rec,
		Kvar:    kvar,
		Kartica: kartica, // novo
	}
}

// Helper za kontekst sa timeout-om
func ctxTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultTimeout)
}

// GetDom vraća jedan dom po ID-u.
func (s *Services) GetDom(ctx context.Context, domID uuid.UUID) (domain.Dom, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	return s.Dom.Get(ctx, s.DB, domID)
}

// GetAllDomovi vraća sve domove.
func (s *Services) GetAllDomovi(ctx context.Context) ([]domain.Dom, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	return s.Dom.GetAll(ctx, s.DB)
}

/* ======================= Studentske operacije ======================= */

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

// UpisiStudentaUSobu kreira studenta i upisuje ga u sobu ako ima mesta.
// Zaključava sobu (FOR UPDATE) i postavlja slobodna=false tek kada popunjenost stigne do limita.
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
	soba, err := s.Soba.GetByBroj(ctx, tx, domID, brojSobe, true) // SELECT ... FOR UPDATE
	if err != nil {
		return domain.Student{}, err
	}

	// 2) Proveri trenutnu popunjenost
	postojeci, err := s.Student.ListBySoba(ctx, tx, soba.ID)
	if err != nil {
		return domain.Student{}, err
	}
	if len(postojeci) >= soba.Kapacitet {
		return domain.Student{}, errors.New("soba je popunjena (nema slobodnih mesta)")
	}

	// 3) Kreiraj studenta
	st := domain.Student{
		ID:      uuid.New(),
		Ime:     ime,
		Prezime: prezime,
	}
	if err = s.Student.Create(ctx, tx, &st); err != nil {
		return domain.Student{}, err
	}

	// 4) Poveži studenta sa sobom
	if err = s.Student.AssignToSoba(ctx, tx, st.ID, soba.ID); err != nil {
		return domain.Student{}, err
	}

	// 5) Ako je ovo poslednje mesto (len+1 == limit) označi sobu kao neslobodnu
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

// OslobodiSobu uklanja dodelu i podešava slobodna=true ako nakon odjave ima mesta.
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

	// 1) Nađi studenta i njegovu sobu
	st, err := s.Student.Get(ctx, tx, studentID)
	if err != nil {
		return err
	}
	if st.SobaID == nil {
		// nema sobu — no-op
		return tx.Commit()
	}

	// 2) Ukloni dodelu
	if err = s.Student.UnassignSoba(ctx, tx, studentID); err != nil {
		return err
	}

	// 3) Pročitaj sobu (sa limitom) i prebroj preostale studente
	soba, err := s.Soba.Get(ctx, tx, *st.SobaID)
	if err != nil {
		return err
	}
	preostali, err := s.Student.ListBySoba(ctx, tx, soba.ID)
	if err != nil {
		return err
	}

	// 4) Ako posle odjave ima mesta (preostali < limit) označi sobu kao slobodnu
	//    (Ako želiš finiju logiku: možeš ostaviti slobodna=false dok ne padne ispod limita,
	//     ali ovde je dovoljno da čim ima makar jedno mesto, soba bude slobodna.)
	shouldBeFree := len(preostali) < soba.Kapacitet
	if err = s.Soba.SetSlobodna(ctx, tx, soba.ID, shouldBeFree); err != nil {
		return err
	}

	return tx.Commit()
}

/* ======================= Recenzije ======================= */

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

/* ======================= Kvarovi ======================= */

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

/* ======================= Sobe / DTO-i ======================= */

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

// KreirajStudentskuKarticuAkoNema: kreira studentsku karticu za datog studenta ako ne postoji,
// u suprotnom vrati postojeću.
func (s *Services) KreirajStudentskuKarticuAkoNema(ctx context.Context, studentID uuid.UUID) (domain.StudentskaKartica, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	return s.Kartica.CreateIfNotExists(ctx, s.DB, studentID)
}

// (opciono) Može zatrebati i direktno čitanje kartice po studentu
func (s *Services) GetStudentskaKarticaByStudent(ctx context.Context, studentID uuid.UUID) (domain.StudentskaKartica, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	return s.Kartica.GetByStudentID(ctx, s.DB, studentID)
}

// (opciono) Ažuriranje stanja (pozitivno ili negativno)
func (s *Services) AzurirajStanjeStudentskeKartice(ctx context.Context, studentID uuid.UUID, delta float64) (domain.StudentskaKartica, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	return s.Kartica.UpdateStanje(ctx, s.DB, studentID, delta)
}

// ListSlobodneSobe vraća sve slobodne sobe u okviru zadatog doma.
func (s *Services) ListSlobodneSobe(ctx context.Context, domID uuid.UUID) ([]domain.Soba, error) {
	ctx, cancel := ctxTimeout(ctx)
	defer cancel()

	return s.Soba.ListSlobodne(ctx, s.DB, domID)
}
