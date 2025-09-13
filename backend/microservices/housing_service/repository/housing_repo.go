package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq" // driver za sql.Open("postgres", ...)
	"github.com/google/uuid"
	"housing/domain"
)

type HousingRepo struct {
	DB *sql.DB
}

func NewHousingRepo() (*HousingRepo, error) {
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

	repo:= &HousingRepo{DB: db}	
	if err := repo.Migrate(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (dr *HousingRepo) Close() error { return dr.DB.Close() }

func (dr *HousingRepo) Migrate() error {
	stmts := []string{
		// === Housing model (CockroachDB, UUID svuda) ===
		`CREATE TABLE IF NOT EXISTS dom (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			naziv STRING NOT NULL,
			adresa STRING NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS soba (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			broj STRING NOT NULL,
			slobodna BOOL NOT NULL DEFAULT true,
			dom_id UUID NOT NULL REFERENCES dom(id) ON DELETE CASCADE,
			CONSTRAINT soba_dom_broj_unq UNIQUE (dom_id, broj)
		);`,

		`CREATE TABLE IF NOT EXISTS student (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			ime STRING NOT NULL,
			prezime STRING NOT NULL,
			soba_id UUID NULL REFERENCES soba(id) ON DELETE SET NULL
		);`,

		`CREATE TABLE IF NOT EXISTS recenzija_sobe (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			ocena INT NOT NULL CHECK (ocena BETWEEN 1 AND 5),
			komentar STRING NULL,
			soba_id UUID NOT NULL REFERENCES soba(id) ON DELETE CASCADE,
			autor_id UUID NOT NULL REFERENCES student(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS recenzija_sobe_soba_idx ON recenzija_sobe(soba_id);`,
		`CREATE INDEX IF NOT EXISTS recenzija_sobe_autor_idx ON recenzija_sobe(autor_id);`,

		`CREATE TYPE IF NOT EXISTS status_kvara AS ENUM ('prijavljen','u_toku','resen');`,
		`CREATE TABLE IF NOT EXISTS kvar (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			opis STRING NOT NULL,
			status status_kvara NOT NULL DEFAULT 'prijavljen',
			soba_id UUID NOT NULL REFERENCES soba(id) ON DELETE CASCADE,
			prijavio_id UUID NOT NULL REFERENCES student(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS kvar_soba_idx ON kvar(soba_id);`,
		`CREATE INDEX IF NOT EXISTS kvar_prijavio_idx ON kvar(prijavio_id);`,
	}

	tx, err := dr.DB.Begin()
	if err != nil {
		return err
	}
	for _, q := range stmts {
		if _, err := tx.Exec(q); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}


type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Dom
type DomRepository interface {
	Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Dom, error)
}

// Soba
type SobaRepository interface {
	Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Soba, error)
	GetByBroj(ctx context.Context, q DBTX, domID uuid.UUID, broj string, forUpdate bool) (domain.Soba, error)
	Create(ctx context.Context, q DBTX, s *domain.Soba) error
	SetSlobodna(ctx context.Context, q DBTX, sobaID uuid.UUID, slobodna bool) error
}

// Student
type StudentRepository interface {
	Create(ctx context.Context, q DBTX, st *domain.Student) error
	Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Student, error)
	AssignToSoba(ctx context.Context, q DBTX, studentID uuid.UUID, sobaID uuid.UUID) error
	UnassignSoba(ctx context.Context, q DBTX, studentID uuid.UUID) error
	ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.Student, error)
}

// Recenzija
type RecenzijaRepository interface {
	Create(ctx context.Context, q DBTX, r *domain.RecenzijaSobe) error
	ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.RecenzijaSobe, error)
}

// Kvar
type KvarRepository interface {
	Create(ctx context.Context, q DBTX, k *domain.Kvar) error
	UpdateStatus(ctx context.Context, q DBTX, kvarID uuid.UUID, status domain.StatusKvara) error
	ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.Kvar, error)
}

/* =========================================================
   Implementacije repo-a
   ========================================================= */

// ------- Dom -------
type domRepo struct{}

func NewDomRepo() DomRepository { return &domRepo{} }

func (r *domRepo) Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Dom, error) {
	var d domain.Dom
	err := q.QueryRowContext(ctx,
		`SELECT id, naziv, adresa FROM dom WHERE id = $1`, id,
	).Scan(&d.ID, &d.Naziv, &d.Adresa)
	return d, err
}

// ------- Soba -------
type sobaRepo struct{}

func NewSobaRepo() SobaRepository { return &sobaRepo{} }

func (r *sobaRepo) Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Soba, error) {
	var s domain.Soba
	err := q.QueryRowContext(ctx,
		`SELECT id, broj, slobodna, dom_id FROM soba WHERE id = $1`, id,
	).Scan(&s.ID, &s.Broj, &s.Slobodna, &s.DomID)
	return s, err
}

func (r *sobaRepo) GetByBroj(ctx context.Context, q DBTX, domID uuid.UUID, broj string, forUpdate bool) (domain.Soba, error) {
	sql := `SELECT id, broj, slobodna, dom_id FROM soba WHERE dom_id = $1 AND broj = $2`
	if forUpdate {
		sql += " FOR UPDATE"
	}
	var s domain.Soba
	err := q.QueryRowContext(ctx, sql, domID, broj).Scan(&s.ID, &s.Broj, &s.Slobodna, &s.DomID)
	return s, err
}

func (r *sobaRepo) Create(ctx context.Context, q DBTX, s *domain.Soba) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return q.QueryRowContext(ctx,
		`INSERT INTO soba (id, broj, slobodna, dom_id) 
		 VALUES ($1,$2,$3,$4) RETURNING id`,
		s.ID, s.Broj, s.Slobodna, s.DomID,
	).Scan(&s.ID)
}

func (r *sobaRepo) SetSlobodna(ctx context.Context, q DBTX, sobaID uuid.UUID, slobodna bool) error {
	_, err := q.ExecContext(ctx, `UPDATE soba SET slobodna = $1 WHERE id = $2`, slobodna, sobaID)
	return err
}

// ------- Student -------
type studentRepo struct{}

func NewStudentRepo() StudentRepository { return &studentRepo{} }

func (r *studentRepo) Create(ctx context.Context, q DBTX, st *domain.Student) error {
	if st.ID == uuid.Nil {
		st.ID = uuid.New()
	}
	return q.QueryRowContext(ctx,
		`INSERT INTO student (id, ime, prezime, soba_id) 
		 VALUES ($1,$2,$3,$4) RETURNING id`,
		st.ID, st.Ime, st.Prezime, st.SobaID,
	).Scan(&st.ID)
}

func (r *studentRepo) Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Student, error) {
	var s domain.Student
	err := q.QueryRowContext(ctx,
		`SELECT id, ime, prezime, soba_id FROM student WHERE id = $1`, id,
	).Scan(&s.ID, &s.Ime, &s.Prezime, &s.SobaID)
	return s, err
}

func (r *studentRepo) AssignToSoba(ctx context.Context, q DBTX, studentID uuid.UUID, sobaID uuid.UUID) error {
	_, err := q.ExecContext(ctx, `UPDATE student SET soba_id = $1 WHERE id = $2`, sobaID, studentID)
	return err
}

func (r *studentRepo) UnassignSoba(ctx context.Context, q DBTX, studentID uuid.UUID) error {
	_, err := q.ExecContext(ctx, `UPDATE student SET soba_id = NULL WHERE id = $1`, studentID)
	return err
}

func (r *studentRepo) ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.Student, error) {
	rows, err := q.QueryContext(ctx,
		`SELECT id, ime, prezime, soba_id FROM student WHERE soba_id = $1`, sobaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Student
	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(&s.ID, &s.Ime, &s.Prezime, &s.SobaID); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// ------- Recenzija -------
type recRepo struct{}

func NewRecRepo() RecenzijaRepository { return &recRepo{} }

func (r *recRepo) Create(ctx context.Context, q DBTX, rc *domain.RecenzijaSobe) error {
	if rc.ID == uuid.Nil {
		rc.ID = uuid.New()
	}
	if rc.Ocena < 1 || rc.Ocena > 5 {
		return errors.New("ocena mora biti između 1 i 5")
	}
	return q.QueryRowContext(ctx,
		`INSERT INTO recenzija_sobe (id, ocena, komentar, soba_id, autor_id)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		rc.ID, rc.Ocena, rc.Komentar, rc.SobaID, rc.AutorID,
	).Scan(&rc.ID)
}

func (r *recRepo) ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.RecenzijaSobe, error) {
	rows, err := q.QueryContext(ctx,
		`SELECT id, ocena, komentar, soba_id, autor_id
		   FROM recenzija_sobe
		  WHERE soba_id = $1
		  ORDER BY id DESC`, sobaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.RecenzijaSobe
	for rows.Next() {
		var x domain.RecenzijaSobe
		if err := rows.Scan(&x.ID, &x.Ocena, &x.Komentar, &x.SobaID, &x.AutorID); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

// ------- Kvar -------
type kvarRepo struct{}

func NewKvarRepo() KvarRepository { return &kvarRepo{} }

func (r *kvarRepo) Create(ctx context.Context, q DBTX, k *domain.Kvar) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return q.QueryRowContext(ctx,
		`INSERT INTO kvar (id, opis, status, soba_id, prijavio_id)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		k.ID, k.Opis, k.Status, k.SobaID, k.PrijavioID,
	).Scan(&k.ID)
}

func (r *kvarRepo) UpdateStatus(ctx context.Context, q DBTX, kvarID uuid.UUID, status domain.StatusKvara) error {
	_, err := q.ExecContext(ctx, `UPDATE kvar SET status = $1 WHERE id = $2`, status, kvarID)
	return err
}

func (r *kvarRepo) ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.Kvar, error) {
	rows, err := q.QueryContext(ctx,
		`SELECT id, opis, status, soba_id, prijavio_id
		   FROM kvar
		  WHERE soba_id = $1
		  ORDER BY id DESC`, sobaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Kvar
	for rows.Next() {
		var k domain.Kvar
		if err := rows.Scan(&k.ID, &k.Opis, &k.Status, &k.SobaID, &k.PrijavioID); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}