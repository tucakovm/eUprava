package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"housing/domain"
	"housing/service"
	"log"
)

type HousingHandler struct {
	service *service.Services
}

func NewHousingHandler(s *service.Services) *HousingHandler {
	return &HousingHandler{service: s}
}

/* =========================
   Helpers
   ========================= */

func (h *HousingHandler) renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func (h *HousingHandler) badRequest(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

/* =========================
   Domovi (read)
   ========================= */

// GET /dom?id=<uuid>
func (h *HousingHandler) GetDom(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.badRequest(w, "missing id")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.badRequest(w, "invalid id")
		return
	}

	dom, err := h.service.GetDom(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	h.renderJSON(w, dom)
}

// GET /doms
func (h *HousingHandler) ListDomovi(w http.ResponseWriter, r *http.Request) {
	domovi, err := h.service.GetAllDomovi(r.Context())
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}
	h.renderJSON(w, domovi)
}

/* =========================
   Students
   ========================= */

// POST /students
// Body: { "ime": "Marko", "prezime": "Markovic" }
func (h *HousingHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Ime     string `json:"ime"`
		Prezime string `json:"prezime"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.Ime == "" || in.Prezime == "" {
		h.badRequest(w, "ime i prezime su obavezni")
		return
	}

	st, err := h.service.CreateStudent(r.Context(), in.Ime, in.Prezime)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, st)
}

// POST /rooms/assign
// Body: { "domId": "...uuid...", "broj": "101", "ime": "Marko", "prezime": "Markovic" }
func (h *HousingHandler) AssignStudentToRoom(w http.ResponseWriter, r *http.Request) {
	var in struct {
		DomID   string `json:"domId"`
		Broj    string `json:"broj"`
		Ime     string `json:"ime"`
		Prezime string `json:"prezime"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.DomID == "" || in.Broj == "" || in.Ime == "" || in.Prezime == "" {
		h.badRequest(w, "domId, broj, ime i prezime su obavezni")
		return
	}
	domID, err := uuid.Parse(in.DomID)
	if err != nil {
		h.badRequest(w, "invalid domId")
		return
	}

	st, err := h.service.UpisiStudentaUSobu(r.Context(), domID, in.Broj, in.Ime, in.Prezime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict) // npr. "soba nije slobodna"
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, st)
}

// POST /students/release
// Body: { "studentId": "...uuid..." }
func (h *HousingHandler) ReleaseStudentRoom(w http.ResponseWriter, r *http.Request) {
	var in struct {
		StudentID string `json:"studentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	id, err := uuid.Parse(in.StudentID)
	if err != nil {
		h.badRequest(w, "invalid studentId")
		return
	}

	if err := h.service.OslobodiSobu(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.renderJSON(w, map[string]string{"status": "ok"})
}

/* =========================
   Studentska kartica (NOVO)
   ========================= */

// POST /students/cards
// Body: { "studentId": "...uuid..." }
// Kreira karticu ako ne postoji, inače vraća postojeću
func (h *HousingHandler) CreateStudentCardIfMissing(w http.ResponseWriter, r *http.Request) {
	var in struct {
		StudentID string `json:"studentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.StudentID == "" {
		h.badRequest(w, "studentId je obavezan")
		return
	}
	sid, err := uuid.Parse(in.StudentID)
	if err != nil {
		h.badRequest(w, "invalid studentId")
		return
	}

	card, err := h.service.KreirajStudentskuKarticuAkoNema(r.Context(), sid)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, card)
}

// GET /students/cards?studentId=<uuid>
// Vraća karticu za datog studenta
func (h *HousingHandler) GetStudentCard(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("studentId")
	if idStr == "" {
		h.badRequest(w, "missing studentId")
		return
	}
	sid, err := uuid.Parse(idStr)
	if err != nil {
		h.badRequest(w, "invalid studentId")
		return
	}

	card, err := h.service.GetStudentskaKarticaByStudent(r.Context(), sid)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	h.renderJSON(w, card)
}

// POST /students/cards/balance
// Body: { "studentId": "...uuid...", "delta": 500.0 }  // pozitivan = dopuna, negativan = zaduženje
func (h *HousingHandler) UpdateStudentCardBalance(w http.ResponseWriter, r *http.Request) {
	var in struct {
		StudentID string  `json:"studentId"`
		Delta     float64 `json:"delta"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.StudentID == "" {
		h.badRequest(w, "studentId je obavezan")
		return
	}
	sid, err := uuid.Parse(in.StudentID)
	if err != nil {
		h.badRequest(w, "invalid studentId")
		return
	}

	card, err := h.service.AzurirajStanjeStudentskeKartice(r.Context(), sid, in.Delta)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}
	h.renderJSON(w, card)
}

/* =========================
   Recenzije
   ========================= */

// POST /rooms/reviews
// Body: { "sobaId": "...uuid...", "autorId": "...uuid...", "ocena": 5, "komentar": "..." }
func (h *HousingHandler) AddRoomReview(w http.ResponseWriter, r *http.Request) {
	var in struct {
		SobaID   string  `json:"sobaId"`
		AutorID  string  `json:"autorId"`
		Ocena    int     `json:"ocena"`
		Komentar *string `json:"komentar"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	sobaID, err := uuid.Parse(in.SobaID)
	if err != nil {
		h.badRequest(w, "invalid sobaId")
		return
	}
	autorID, err := uuid.Parse(in.AutorID)
	if err != nil {
		h.badRequest(w, "invalid autorId")
		return
	}
	if in.Ocena < 1 || in.Ocena > 5 {
		h.badRequest(w, "ocena mora biti 1..5")
		return
	}

	rc, err := h.service.DodajRecenziju(r.Context(), sobaID, autorID, in.Ocena, in.Komentar)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, rc)
}

/* =========================
   Kvarovi
   ========================= */

// POST /rooms/faults
// Body: { "sobaId": "...uuid...", "prijavioId": "...uuid...", "opis": "..." }
func (h *HousingHandler) ReportFault(w http.ResponseWriter, r *http.Request) {
	var in struct {
		SobaID     string `json:"sobaId"`
		PrijavioID string `json:"prijavioId"`
		Opis       string `json:"opis"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.Opis == "" {
		h.badRequest(w, "opis je obavezan")
		return
	}
	sobaID, err := uuid.Parse(in.SobaID)
	if err != nil {
		h.badRequest(w, "invalid sobaId")
		return
	}
	prijavioID, err := uuid.Parse(in.PrijavioID)
	if err != nil {
		h.badRequest(w, "invalid prijavioId")
		return
	}

	k, err := h.service.PrijaviKvar(r.Context(), sobaID, prijavioID, in.Opis)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, k)
}

// POST /faults/status
// Body: { "kvarId": "...uuid...", "status": "u_toku|resen|prijavljen" }
func (h *HousingHandler) ChangeFaultStatus(w http.ResponseWriter, r *http.Request) {
	var in struct {
		KvarID string `json:"kvarId"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	kid, err := uuid.Parse(in.KvarID)
	if err != nil {
		h.badRequest(w, "invalid kvarId")
		return
	}
	status := domain.StatusKvara(in.Status)
	switch status {
	case domain.StatusPrijavljen, domain.StatusUToku, domain.StatusResen:
	default:
		h.badRequest(w, "status mora biti: prijavljen | u_toku | resen")
		return
	}

	if err := h.service.PromeniStatusKvara(r.Context(), kid, status); err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}

	h.renderJSON(w, map[string]string{"status": "ok"})
}

/* =========================
   Sobe (read)
   ========================= */

// GET /rooms?id=<uuid>
func (h *HousingHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.badRequest(w, "missing id")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.badRequest(w, "invalid id")
		return
	}

	room, err := h.service.GetSoba(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	h.renderJSON(w, room)
}

// GET /rooms/detail?id=<uuid>
func (h *HousingHandler) GetRoomDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.badRequest(w, "missing id")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.badRequest(w, "invalid id")
		return
	}

	room, err := h.service.GetSobaDetail(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	h.renderJSON(w, room)
}

// GET /rooms/free?domId=<uuid>
// Vraća sve slobodne sobe u okviru zadatog doma
func (h *HousingHandler) ListFreeRooms(w http.ResponseWriter, r *http.Request) {
	domIDStr := r.URL.Query().Get("domId")
	if domIDStr == "" {
		h.badRequest(w, "missing domId")
		return
	}
	domID, err := uuid.Parse(domIDStr)
	if err != nil {
		h.badRequest(w, "invalid domId")
		return
	}

	rooms, err := h.service.ListSlobodneSobe(r.Context(), domID)
	if err != nil {
		log.Printf("ListSlobodneSobe failed: %v", err)
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}

	h.renderJSON(w, rooms)
}
