package handler

import (
	"dining/domain"
	"dining/service"
	"encoding/json"
	"net/http"
	"time"

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

func (dh *DiningHandler) renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
