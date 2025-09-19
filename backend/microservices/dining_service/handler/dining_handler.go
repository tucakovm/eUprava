package handler

import (
	"bytes"
	"dining/domain"
	"dining/service"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DiningHandler struct {
	service service.DiningService
}

func NewDiningHandler(service service.DiningService) *DiningHandler {
	return &DiningHandler{
		service: service,
	}
}

func (dh *DiningHandler) GetAllCanteens(rw http.ResponseWriter, r *http.Request) {
	allProducts, err := dh.service.GetAllCanteens()

	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	dh.renderJSON(rw, allProducts)

}

func (dh *DiningHandler) GetCanteen(rw http.ResponseWriter, r *http.Request) {
	canteenId := mux.Vars(r)["id"]

	canteen, err := dh.service.GetCanteen(canteenId)
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(rw).Encode(canteen); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (dh *DiningHandler) DeleteCanteen(rw http.ResponseWriter, r *http.Request) {
	canteenId := mux.Vars(r)["id"]

	err := dh.service.DeleteCanteen(canteenId)
	if err != nil {
		http.Error(rw, "Failed to delete canteen", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Canteen deleted successfully",
		"id":      canteenId,
	})
}

func (dh *DiningHandler) CreateCanteen(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	var dto domain.CanteenDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	today := time.Now().Format("2006-01-02")
	openAt, err := time.Parse("2006-01-02 15:04", today+" "+dto.OpenAt)
	if err != nil {
		http.Error(rw, "Invalid open_at format", http.StatusBadRequest)
		return
	}
	closeAt, err := time.Parse("2006-01-02 15:04", today+" "+dto.CloseAt)
	if err != nil {
		http.Error(rw, "Invalid close_at format", http.StatusBadRequest)
		return
	}

	canteen := &domain.Canteen{
		Name:    dto.Name,
		Address: dto.Address,
		OpenAt:  openAt,
		CloseAt: closeAt,
	}

	if err := dh.service.CreateCanteen(canteen); err != nil {
		http.Error(rw, "Failed to create canteen", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(rw).Encode(canteen)
}

func (dh *DiningHandler) GetMenusByCanteenID(rw http.ResponseWriter, r *http.Request) {
	canteenId := mux.Vars(r)["id"]

	menus, err := dh.service.GetMenusByCanteenID(canteenId)
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	dh.renderJSON(rw, menus)
}

func (dh *DiningHandler) CreateMenu(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	var menu domain.MenuDTO
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := dh.service.CreateMenu(&menu); err != nil {
		http.Error(rw, "Failed to create menu", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(rw).Encode(menu)
}

func (dh *DiningHandler) DeleteMenu(rw http.ResponseWriter, r *http.Request) {
	manuId := mux.Vars(r)["id"]

	err := dh.service.DeleteMenu(manuId)
	if err != nil {
		http.Error(rw, "Failed to delete manu", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Menu deleted successfully",
		"id":      manuId,
	})
}

func (dh *DiningHandler) GetPopularMeals(rw http.ResponseWriter, r *http.Request) {
	canteenId := mux.Vars(r)["id"]

	popMenus, err := dh.service.GetPopularMenus(canteenId)
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	dh.renderJSON(rw, popMenus)
}

func (dh *DiningHandler) GetMealHistory(rw http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	history, err := dh.service.GetMealHistory(userId)
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	dh.renderJSON(rw, history)
}

func (h *DiningHandler) GetMealHistoryWithReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := strings.Trim(vars["userId"], `"`)
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	history, err := h.service.GetMealHistoryWithReviewsByUser(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (h *DiningHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	var review domain.MenuReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateMenuReview(&review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *DiningHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var input domain.MenuReviewDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	review := domain.MenuReview{
		Id:              uuid.New(),
		MenuId:          uuid.MustParse(input.MenuId),
		UserId:          uuid.MustParse(input.UserId),
		BreakfastReview: input.BreakfastReview,
		LunchReview:     input.LunchReview,
		DinnerReview:    input.DinnerReview,
	}

	if err := h.service.CreateMenuReview(&review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (dh *DiningHandler) fetchStudentCard(studentID string) (*domain.StudentCard, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	fmt.Println("Fetching student card: ", studentID)

	// URL sa query parametrom studentId
	url := fmt.Sprintf("http://housing-server:8003/api/housing/students/cards?studentUsername=%s", studentID)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var card domain.StudentCard
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, fmt.Errorf("failed to decode student card: %w", err)
	}

	return &card, nil
}

func (dh *DiningHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuId := strings.Trim(vars["id"], `"`)
	if menuId == "" {
		http.Error(w, "Menu ID is required", http.StatusBadRequest)
		return
	}

	menu, err := dh.service.GetMenu(menuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	studentIDStr := r.Header.Get("X-Student-ID")
	if studentIDStr == "" {
		http.Error(w, "missing student ID", http.StatusBadRequest)
		return
	}

	card, err := dh.fetchStudentCard(studentIDStr)
	if err != nil {
		fmt.Println("Warning: failed to fetch student card:", err)
		// možemo i ovde samo logovati, a ne prekidati
	}

	// Odgovor sa menijem i karticom
	response := struct {
		Menu *domain.Menu        `json:"menu"`
		Card *domain.StudentCard `json:"card,omitempty"`
	}{
		Menu: menu,
		Card: card,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (dh *DiningHandler) GetTopRatedMeals(w http.ResponseWriter, r *http.Request) {
	topMeals, err := dh.service.GetTopRatedMeals()
	if err != nil {
		log.Println("GetTopRatedMeals error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(topMeals); err != nil {
		log.Println("JSON encode error:", err)
	}
}

func (dh *DiningHandler) TakeMeal(w http.ResponseWriter, r *http.Request) {
	var in struct {
		StudentUsername string  `json:"studentUsername"`
		Delta           float64 `json:"delta"`
		MenuId          string  `json:"menuId"`
		StudentID       string  `json:"studentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	client := &http.Client{Timeout: 3 * time.Second}

	url := "http://housing-server:8003/api/housing/students/cards/balance"

	reqBody, _ := json.Marshal(in)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "failed to contact housing service", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "failed to update balance", resp.StatusCode)
		return
	}

	menu, _ := dh.service.GetMenu(in.MenuId)

	mh := &domain.MealHistory{
		MenuId:     in.MenuId,
		MenuName:   menu.Name,
		SelectedAt: time.Time{},
	}
	err = dh.service.CreateMealHistory(mh, in.StudentID)
	if err != nil {
		fmt.Println("Warning: failed to create meal history:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = dh.service.IncrementPopularMeal(menu.Id, menu.CanteenId)
	if err != nil {
		fmt.Println("Warning: failed to increment popular meal:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func (dh *DiningHandler) CheckDoesStudentInRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := strings.Trim(vars["userId"], `"`)
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("http://housing-server:8003/api/housing/rooms/checkStudent/%s", userId)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to contact housing service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Housing service error: %s", resp.Status), resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read housing service response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (dh *DiningHandler) GetMealRoomHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		Usernames []string `json:"usernames"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Invalid JSON format: %v", err)
		// vraćamo praznu listu
		json.NewEncoder(w).Encode([]domain.MealRoomHistory{})
		return
	}

	if len(request.Usernames) == 0 {
		log.Printf("Usernames array empty")
		// vraćamo praznu listu
		json.NewEncoder(w).Encode([]domain.MealRoomHistory{})
		return
	}

	history, err := dh.service.GetMealHistoryForUsernames(request.Usernames)
	if err != nil {
		log.Printf("Failed to get meal history for usernames %v: %v", request.Usernames, err)
		// vraćamo praznu listu
		json.NewEncoder(w).Encode([]domain.MealRoomHistory{})
		return
	}

	if history == nil {
		history = []domain.MealRoomHistory{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

func (dh *DiningHandler) renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
