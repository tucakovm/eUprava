package handler

import (
	"dining/service"
	"encoding/json"
	"net/http"
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

func (dh *DiningHandler) renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
