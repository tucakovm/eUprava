package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"fmt"
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

/* ========================= Helpers ========================= */

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

/* ========================= Domovi (read) ========================= */

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

/* ========================= Students ========================= */

// POST /students
// Body: { "ime": "Marko", "prezime": "Markovic", "username": "marko123" }
func (h *HousingHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Ime      string `json:"ime"`
		Prezime  string `json:"prezime"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.Ime == "" || in.Prezime == "" || in.Username == "" {
		h.badRequest(w, "ime, prezime i username su obavezni")
		return
	}

	st, err := h.service.CreateStudent(r.Context(), in.Ime, in.Prezime, in.Username)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, st)
}

// POST /rooms/assign
// Body: { "domId": "...uuid...", "broj": "101", "username": "nikola123" }
func (h *HousingHandler) AssignStudentToRoom(w http.ResponseWriter, r *http.Request) {
	var in struct {
		DomID    string `json:"domId"`
		Broj     string `json:"broj"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.DomID == "" || in.Broj == "" || in.Username == "" {
		h.badRequest(w, "domId, broj i username su obavezni")
		return
	}
	domID, err := uuid.Parse(in.DomID)
	if err != nil {
		h.badRequest(w, "invalid domId")
		return
	}

	st, err := h.service.UpisiPostojecegStudentaUSobu(r.Context(), domID, in.Broj, in.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, st)
}

// POST /students/release
// Body: { "studentId": "...uuid..." }  // ovo ostaje po ID-u (interni admin-endpoint)
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

/* ========================= Studentska kartica (po username) ========================= */

// POST /students/cards
// Body: { "studentUsername": "nikola123" }
func (h *HousingHandler) CreateStudentCardIfMissing(w http.ResponseWriter, r *http.Request) {
	var in struct {
		StudentUsername string `json:"studentUsername"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.StudentUsername == "" {
		h.badRequest(w, "studentUsername je obavezan")
		return
	}

	card, err := h.service.KreirajStudentskuKarticuAkoNema(r.Context(), in.StudentUsername)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, card)
}

// GET /students/cards?studentUsername=<username>
func (h *HousingHandler) GetStudentCard(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("studentUsername")
	if u == "" {
		h.badRequest(w, "missing studentUsername")
		return
	}

	card, err := h.service.GetStudentskaKarticaByStudent(r.Context(), u)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	h.renderJSON(w, card)
}

// POST /students/cards/balance
// Body: { "studentUsername": "nikola123", "delta": 500.0 }
func (h *HousingHandler) UpdateStudentCardBalance(w http.ResponseWriter, r *http.Request) {
	var in struct {
		StudentUsername string  `json:"studentUsername"`
		Delta           float64 `json:"delta"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.badRequest(w, "bad json")
		return
	}
	if in.StudentUsername == "" {
		h.badRequest(w, "studentUsername je obavezan")
		return
	}

	card, err := h.service.AzurirajStanjeStudentskeKartice(r.Context(), in.StudentUsername, in.Delta)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}
	h.renderJSON(w, card)
}

/* ========================= Recenzije ========================= */

// POST /rooms/reviews
// Body: { "sobaId": "...uuid...", "autorUsername": "nikola123", "ocena": 5, "komentar": "..." }
func (h *HousingHandler) AddRoomReview(w http.ResponseWriter, r *http.Request) {
	var in struct {
		SobaID        string  `json:"sobaId"`
		AutorUsername string  `json:"autorUsername"`
		Ocena         int     `json:"ocena"`
		Komentar      *string `json:"komentar"`
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
	if in.AutorUsername == "" {
		h.badRequest(w, "autorUsername je obavezan")
		return
	}
	if in.Ocena < 1 || in.Ocena > 5 {
		h.badRequest(w, "ocena mora biti 1..5")
		return
	}

	rc, err := h.service.DodajRecenziju(r.Context(), sobaID, in.AutorUsername, in.Ocena, in.Komentar)
	if err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.renderJSON(w, rc)
}

/* ========================= Kvarovi ========================= */

// POST /rooms/faults
// Body: { "sobaId": "...uuid...", "prijavioUsername": "marko123", "opis": "..." }
func (h *HousingHandler) ReportFault(w http.ResponseWriter, r *http.Request) {
	var in struct {
		SobaID           string `json:"sobaId"`
		PrijavioUsername string `json:"prijavioUsername"`
		Opis             string `json:"opis"`
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
	if in.PrijavioUsername == "" {
		h.badRequest(w, "prijavioUsername je obavezan")
		return
	}

	k, err := h.service.PrijaviKvar(r.Context(), sobaID, in.PrijavioUsername, in.Opis)
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
	switch in.Status {
	case string(domain.StatusPrijavljen), string(domain.StatusUToku), string(domain.StatusResen):
	default:
		h.badRequest(w, "status mora biti: prijavljen | u_toku | resen")
		return
	}

	if err := h.service.PromeniStatusKvara(r.Context(), kid, domain.StatusKvara(in.Status)); err != nil {
		http.Error(w, "database exception", http.StatusInternalServerError)
		return
	}

	h.renderJSON(w, map[string]string{"status": "ok"})
}

/* ========================= Sobe (read) ========================= */

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

/* ========================= Dining proxy ========================= */

// GET /api/housing/notifications/menus
// Proxy ka Dining servisu: GET {TARGET}/api/dining/menus/today
func (h *HousingHandler) GetTodayDiningMenus(w http.ResponseWriter, r *http.Request) {
	// 1) kandidati za bazni URL
	candidates := []string{}
	if env := os.Getenv("DINING_BASE_URL"); env != "" {
		candidates = append(candidates, env)
	}
	// tipični default-i za dev / docker
	candidates = append(candidates,
		"http://localhost:8001",
		"localhost:8001",
		"http://dining-server:8001",
		"http://host.docker.internal:8001",
	)

	client := &http.Client{Timeout: 5 * time.Second}

	var resp *http.Response
	var lastErr error

	// 2) pokušaj redom dok jedan ne uspe
	for _, base := range candidates {
		req, err := http.NewRequest(http.MethodGet, base+"/api/dining/menus/today", nil)
		if err != nil {
			lastErr = err
			continue
		}

		// (opciono) propagiraj identitet/korisne headere
		if sid := r.Header.Get("X-Student-ID"); sid != "" {
			req.Header.Set("X-Student-ID", sid)
		}

		res, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		resp = res
		break
	}

	if resp == nil {
		http.Error(w, "failed to contact dining service: "+fmt.Sprintf("%v", lastErr), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 3) propagiraj relevantne headere i status
	for k, vals := range resp.Header {
		switch k {
		case "Content-Type", "Cache-Control":
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
	}
	w.WriteHeader(resp.StatusCode)

	// 4) prosledi telo kakvo jeste
	_, _ = io.Copy(w, resp.Body)
}
