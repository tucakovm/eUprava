package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"housing/domain"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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

	repo := &HousingRepo{DB: db}
	if err := repo.Migrate(); err != nil {
		return nil, err
	}
	if err := repo.InitData(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (dr *HousingRepo) Close() error { return dr.DB.Close() }

func (dr *HousingRepo) Migrate() error {
	stmts := []string{
		// Osnovne tabele
		`CREATE TABLE IF NOT EXISTS dom (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			naziv TEXT NOT NULL,
			adresa TEXT NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS soba (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			broj TEXT NOT NULL,
			slobodna BOOLEAN NOT NULL DEFAULT true,
			kapacitet INTEGER NOT NULL DEFAULT 1,
			dom_id UUID NOT NULL REFERENCES dom(id) ON DELETE CASCADE,
			CONSTRAINT soba_dom_broj_unq UNIQUE (dom_id, broj)
		);`,

		`CREATE TABLE IF NOT EXISTS student (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			ime TEXT NOT NULL,
			prezime TEXT NOT NULL,
			username TEXT NOT NULL UNIQUE,
			soba_id UUID NULL REFERENCES soba(id) ON DELETE SET NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS student_username_unq ON student(username);`,

		// Recenzije — veza po autor_username (TEXT) na student(username)
		`CREATE TABLE IF NOT EXISTS recenzija_sobe (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			ocena INTEGER NOT NULL CHECK (ocena BETWEEN 1 AND 5),
			komentar TEXT NULL,
			soba_id UUID NOT NULL REFERENCES soba(id) ON DELETE CASCADE,
			autor_username TEXT NOT NULL REFERENCES student(username) ON DELETE CASCADE
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS rec_soba_autor_unq
			ON recenzija_sobe (soba_id, autor_username);`,
		`CREATE INDEX IF NOT EXISTS recenzija_sobe_soba_idx
			ON recenzija_sobe (soba_id);`,
		`CREATE INDEX IF NOT EXISTS recenzija_sobe_autor_username_idx
			ON recenzija_sobe (autor_username);`,

		// Enum i kvarovi — prijavio_username (TEXT) na student(username)
		`CREATE TYPE IF NOT EXISTS status_kvara AS ENUM ('prijavljen','u_toku','resen');`,
		`CREATE TABLE IF NOT EXISTS kvar (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			opis TEXT NOT NULL,
			status status_kvara NOT NULL DEFAULT 'prijavljen',
			soba_id UUID NOT NULL REFERENCES soba(id) ON DELETE CASCADE,
			prijavio_username TEXT NOT NULL REFERENCES student(username) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS kvar_soba_idx ON kvar(soba_id);`,
		`CREATE INDEX IF NOT EXISTS kvar_prijavio_username_idx ON kvar(prijavio_username);`,

		// Studentska kartica — vezana na student(username)
		`CREATE TABLE IF NOT EXISTS studentska_kartica (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			stanje NUMERIC NOT NULL DEFAULT 0,
			student_username TEXT NOT NULL UNIQUE REFERENCES student(username) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS studentska_kartica_student_username_idx ON studentska_kartica(student_username);`,
	}

	// Retry-abilna transakcija (CockroachDB)
	return crdb.ExecuteTx(context.Background(), dr.DB, nil, func(tx *sql.Tx) error {
		for _, q := range stmts {
			if _, err := tx.Exec(q); err != nil {
				return err
			}
		}
		return nil
	})
}

// InitData — osnovni seed
func (dr *HousingRepo) InitData(ctx context.Context) error {
	var cnt int
	if err := dr.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM dom`).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}

	// Retry-abilna transakcija (CockroachDB)
	return crdb.ExecuteTx(ctx, dr.DB, nil, func(tx *sql.Tx) error {
		// IDs
		dom1ID := uuid.New()
		dom2ID := uuid.New()

		soba101ID := uuid.New()
		soba102ID := uuid.New()

		studentNikolaID := uuid.New()
		studentJovanaID := uuid.New()
		studentMarkoID := uuid.New()
		studentJelenaID := uuid.New()

		rec1ID := uuid.New()
		rec2ID := uuid.New()

		kvar1ID := uuid.New()
		kvar2ID := uuid.New()

		karticaNikolaID := uuid.New()
		karticaJovanaID := uuid.New()
		karticaMarkoID := uuid.New()
		karticaJelenaID := uuid.New()

		// Domovi
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO dom (id, naziv, adresa) VALUES 
			 ($1,'Dom Studenata 1','Bulevar Oslobođenja 12'),
			 ($2,'Dom Studenata 2','Cara Dušana 45')`,
			dom1ID, dom2ID,
		); err != nil {
			return err
		}

		//User iz user repoa
		userid := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		userCardId := uuid.New()

		// --- Dom (2) ---
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO dom (id, naziv, adresa) VALUES 
		 ($1,'Dom Studenata 1','Bulevar Oslobodjenja 12'),
		 ($2,'Dom Studenata 2','Cara Dusana 45')`,
			dom1ID, dom2ID,
		); err != nil {
			return err
		}

		// Sobe
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO soba (id, broj, slobodna, kapacitet, dom_id) VALUES
			 ($1,'101', true, 3, $3),
			 ($2,'102', false, 2, $3)`,
			soba101ID, soba102ID, dom1ID,
		); err != nil {
			return err
		}

		// Studenti (Nikola/Marko u 102, Jovana bez sobe)
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO student (id, ime, prezime, username, soba_id) VALUES
			 ($1,'Nikola','Nikolic','nikola123',$4),
			 ($2,'Jovana','Petrovic','jovana123',NULL),
			 ($3,'Marko','Ilic','marko123',$4)
			 ON CONFLICT (username) DO UPDATE
			   SET ime = EXCLUDED.ime,
			       prezime = EXCLUDED.prezime,
			       soba_id = EXCLUDED.soba_id`,
			studentNikolaID, studentJovanaID, studentMarkoID, soba102ID,
		); err != nil {
			return err
		}

		// --- Student (2) (Marko u sobi 102, Jelena bez sobe) ---
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO student (id, ime, prezime, soba_id) VALUES
	 ($1,'Marko','Markovic',$3),
	 ($2,'Jelena','Jovanovic',NULL),
	 ($4,'Asam','Arkom',NULL)`,
			studentMarkoID, studentJelenaID, soba102ID, userid,
		); err != nil {
			return err
		}

		// Recenzije (autor po username)
		kom1 := "Odlična soba, sve preporuke!"
		kom2 := "Moglo bi biti bolje, često je bučno."
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO recenzija_sobe (id, ocena, komentar, soba_id, autor_username) VALUES
			 ($1, 5, $3, $5, 'nikola123'),
			 ($2, 3, $4, $5, 'jovana123')`,
			rec1ID, rec2ID, &kom1, &kom2, soba102ID,
		); err != nil {
			return err
		}

		// Kvarovi (prijavio po username)
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO kvar (id, opis, status, soba_id, prijavio_username) VALUES
			 ($1, 'Pokvaren radijator', 'prijavljen', $3, 'marko123'),
			 ($2, 'Ne radi internet',  'u_toku',    $4, 'jovana123')`,
			kvar1ID, kvar2ID, soba102ID, soba101ID,
		); err != nil {
			return err
		}

		// Kartice po username
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO studentska_kartica (id, stanje, student_username) VALUES
			 ($1, 1500.00, 'nikola123'),
			 ($2,  800.00, 'jovana123')`,
			karticaNikolaID, karticaJovanaID,
		); err != nil {
			return err
		}

		// --- Studentska kartica (3) ---
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO studentska_kartica (id, stanje, student_id) VALUES
	 ($1, 1500.00, $4),
	 ($2,  800.00, $5),
	 ($3, 2000.00, $6)`,
			karticaMarkoID, karticaJelenaID, userCardId,
			studentMarkoID, studentJelenaID, userid,
		); err != nil {
			return err
		}

		return tx.Commit()

		return nil
	})

}

/* ================== DBTX ================== */

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

/* ================== Dom ================== */

type DomRepository interface {
	Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Dom, error)
	GetAll(ctx context.Context, q DBTX) ([]domain.Dom, error)
}

type domRepo struct{}

func NewDomRepo() DomRepository { return &domRepo{} }

func (r *domRepo) GetAll(ctx context.Context, q DBTX) ([]domain.Dom, error) {
	rows, err := q.QueryContext(ctx, `SELECT id, naziv, adresa FROM dom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domovi []domain.Dom
	for rows.Next() {
		var d domain.Dom
		if err := rows.Scan(&d.ID, &d.Naziv, &d.Adresa); err != nil {
			return nil, err
		}
		domovi = append(domovi, d)
	}
	return domovi, rows.Err()
}

func (r *domRepo) Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Dom, error) {
	var d domain.Dom
	err := q.QueryRowContext(ctx,
		`SELECT id, naziv, adresa FROM dom WHERE id = $1`, id,
	).Scan(&d.ID, &d.Naziv, &d.Adresa)
	return d, err
}

/* ================== Soba ================== */

type SobaRepository interface {
	Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Soba, error)
	GetByBroj(ctx context.Context, q DBTX, domID uuid.UUID, broj string, forUpdate bool) (domain.Soba, error)
	Create(ctx context.Context, q DBTX, s *domain.Soba) error
	SetSlobodna(ctx context.Context, q DBTX, sobaID uuid.UUID, slobodna bool) error
	ListSlobodne(ctx context.Context, q DBTX, domID uuid.UUID) ([]domain.Soba, error)
}

type sobaRepo struct{}

func NewSobaRepo() SobaRepository { return &sobaRepo{} }

func (r *sobaRepo) Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Soba, error) {
	var s domain.Soba
	err := q.QueryRowContext(ctx,
		`SELECT id, broj, slobodna, kapacitet, dom_id FROM soba WHERE id = $1`, id,
	).Scan(&s.ID, &s.Broj, &s.Slobodna, &s.Kapacitet, &s.DomID)
	return s, err
}

func (r *sobaRepo) GetByBroj(ctx context.Context, q DBTX, domID uuid.UUID, broj string, forUpdate bool) (domain.Soba, error) {
	sqlStr := `SELECT id, broj, slobodna, kapacitet, dom_id FROM soba WHERE dom_id = $1 AND broj = $2`
	if forUpdate {
		sqlStr += " FOR UPDATE"
	}
	var s domain.Soba
	err := q.QueryRowContext(ctx, sqlStr, domID, broj).Scan(&s.ID, &s.Broj, &s.Slobodna, &s.Kapacitet, &s.DomID)
	return s, err
}

func (r *sobaRepo) Create(ctx context.Context, q DBTX, s *domain.Soba) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return q.QueryRowContext(ctx,
		`INSERT INTO soba (id, broj, slobodna, kapacitet, dom_id) 
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		s.ID, s.Broj, s.Slobodna, s.Kapacitet, s.DomID,
	).Scan(&s.ID)
}

func (r *sobaRepo) SetSlobodna(ctx context.Context, q DBTX, sobaID uuid.UUID, slobodna bool) error {
	_, err := q.ExecContext(ctx, `UPDATE soba SET slobodna = $1 WHERE id = $2`, slobodna, sobaID)
	return err
}

func (r *sobaRepo) ListSlobodne(ctx context.Context, q DBTX, domID uuid.UUID) ([]domain.Soba, error) {
	rows, err := q.QueryContext(ctx,
		`SELECT id, broj, slobodna, kapacitet, dom_id
		   FROM soba
		  WHERE dom_id = $1 AND slobodna = true
		  ORDER BY broj`, domID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Soba
	for rows.Next() {
		var s domain.Soba
		if err := rows.Scan(&s.ID, &s.Broj, &s.Slobodna, &s.Kapacitet, &s.DomID); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

/* ================== Student ================== */

type StudentRepository interface {
	Create(ctx context.Context, q DBTX, st *domain.Student) error
	Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Student, error)
	GetByUsername(ctx context.Context, q DBTX, username string) (domain.Student, error)
	AssignToSoba(ctx context.Context, q DBTX, studentID uuid.UUID, sobaID uuid.UUID) error
	UnassignSoba(ctx context.Context, q DBTX, studentID uuid.UUID) error
	ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.Student, error)
}

type studentRepo struct{}

func NewStudentRepo() StudentRepository { return &studentRepo{} }

func (r *studentRepo) Create(ctx context.Context, q DBTX, st *domain.Student) error {
	if st.ID == uuid.Nil {
		st.ID = uuid.New()
	}
	return q.QueryRowContext(ctx,
		`INSERT INTO student (id, ime, prezime, username, soba_id) 
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		st.ID, st.Ime, st.Prezime, st.Username, st.SobaID,
	).Scan(&st.ID)
}

func (r *studentRepo) Get(ctx context.Context, q DBTX, id uuid.UUID) (domain.Student, error) {
	var s domain.Student
	err := q.QueryRowContext(ctx,
		`SELECT id, ime, prezime, username, soba_id FROM student WHERE id = $1`, id,
	).Scan(&s.ID, &s.Ime, &s.Prezime, &s.Username, &s.SobaID)
	return s, err
}

func (r *studentRepo) GetByUsername(ctx context.Context, q DBTX, username string) (domain.Student, error) {
	var s domain.Student
	err := q.QueryRowContext(ctx,
		`SELECT id, ime, prezime, username, soba_id FROM student WHERE username = $1`, username,
	).Scan(&s.ID, &s.Ime, &s.Prezime, &s.Username, &s.SobaID)
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
		`SELECT id, ime, prezime, username, soba_id FROM student WHERE soba_id = $1`, sobaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Student
	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(&s.ID, &s.Ime, &s.Prezime, &s.Username, &s.SobaID); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

/* ================== Recenzija ================== */

type RecenzijaRepository interface {
	Create(ctx context.Context, q DBTX, r *domain.RecenzijaSobe) error
	ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.RecenzijaSobe, error)
}

type recRepo struct{}

func NewRecRepo() RecenzijaRepository { return &recRepo{} }

func (r *recRepo) Create(ctx context.Context, q DBTX, rc *domain.RecenzijaSobe) error {
	if rc.ID == uuid.Nil {
		rc.ID = uuid.New()
	}
	if rc.Ocena < 1 || rc.Ocena > 5 {
		return errors.New("ocena mora biti između 1 i 5")
	}
	// Očekujemo rc.AutorUsername kao string u domain-u!
	return q.QueryRowContext(ctx,
		`INSERT INTO recenzija_sobe (id, ocena, komentar, soba_id, autor_username)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		rc.ID, rc.Ocena, rc.Komentar, rc.SobaID, rc.AutorUsername,
	).Scan(&rc.ID)
}

func (r *recRepo) ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.RecenzijaSobe, error) {
	rows, err := q.QueryContext(ctx,
		`SELECT id, ocena, komentar, soba_id, autor_username
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
		// Očekujemo AutorUsername kao string u domain-u!
		if err := rows.Scan(&x.ID, &x.Ocena, &x.Komentar, &x.SobaID, &x.AutorUsername); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

/* ================== Kvar ================== */

type KvarRepository interface {
	Create(ctx context.Context, q DBTX, k *domain.Kvar) error
	UpdateStatus(ctx context.Context, q DBTX, kvarID uuid.UUID, status domain.StatusKvara) error
	ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.Kvar, error)
}

type kvarRepo struct{}

func NewKvarRepo() KvarRepository { return &kvarRepo{} }

func (r *kvarRepo) Create(ctx context.Context, q DBTX, k *domain.Kvar) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	// Očekujemo PrijavioUsername kao string u domain-u!
	return q.QueryRowContext(ctx,
		`INSERT INTO kvar (id, opis, status, soba_id, prijavio_username)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		k.ID, k.Opis, k.Status, k.SobaID, k.PrijavioUsername,
	).Scan(&k.ID)
}

func (r *kvarRepo) UpdateStatus(ctx context.Context, q DBTX, kvarID uuid.UUID, status domain.StatusKvara) error {
	_, err := q.ExecContext(ctx, `UPDATE kvar SET status = $1 WHERE id = $2`, status, kvarID)
	return err
}

func (r *kvarRepo) ListBySoba(ctx context.Context, q DBTX, sobaID uuid.UUID) ([]domain.Kvar, error) {
	rows, err := q.QueryContext(ctx,
		`SELECT id, opis, status, soba_id, prijavio_username
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
		// Očekujemo PrijavioUsername kao string u domain-u!
		if err := rows.Scan(&k.ID, &k.Opis, &k.Status, &k.SobaID, &k.PrijavioUsername); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

/* ============ Studentska kartica (po username) ============ */

type StudentskaKarticaRepository interface {
	CreateIfNotExistsByUsername(ctx context.Context, q DBTX, studentUsername string) (domain.StudentskaKartica, error)
	GetByStudentUsername(ctx context.Context, q DBTX, studentUsername string) (domain.StudentskaKartica, error)
	UpdateStanjeByUsername(ctx context.Context, q DBTX, studentUsername string, delta float64) (domain.StudentskaKartica, error)
}

type karticaRepo struct{}

func NewStudentskaKarticaRepo() StudentskaKarticaRepository { return &karticaRepo{} }

func (r *karticaRepo) GetByStudentUsername(ctx context.Context, q DBTX, studentUsername string) (domain.StudentskaKartica, error) {
	var k domain.StudentskaKartica
	err := q.QueryRowContext(ctx,
		`SELECT id, stanje, student_username
		   FROM studentska_kartica
		  WHERE student_username = $1`, studentUsername).
		Scan(&k.ID, &k.Stanje, &k.StudentUsername)
	return k, err
}

func (r *karticaRepo) CreateIfNotExistsByUsername(ctx context.Context, q DBTX, studentUsername string) (domain.StudentskaKartica, error) {
	var k domain.StudentskaKartica
	err := q.QueryRowContext(ctx,
		`INSERT INTO studentska_kartica (student_username)
		 VALUES ($1)
		 ON CONFLICT (student_username) DO NOTHING
		 RETURNING id, stanje, student_username`, studentUsername).
		Scan(&k.ID, &k.Stanje, &k.StudentUsername)

	if err == sql.ErrNoRows {
		return r.GetByStudentUsername(ctx, q, studentUsername)
	}
	return k, err
}

func (r *karticaRepo) UpdateStanjeByUsername(ctx context.Context, q DBTX, studentUsername string, delta float64) (domain.StudentskaKartica, error) {
	var k domain.StudentskaKartica
	err := q.QueryRowContext(ctx,
		`UPDATE studentska_kartica
		    SET stanje = stanje + $1
		  WHERE student_username = $2
		  RETURNING id, stanje, student_username`,
		delta, studentUsername).
		Scan(&k.ID, &k.Stanje, &k.StudentUsername)
	return k, err
}
