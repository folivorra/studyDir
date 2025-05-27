package controller

import (
	"encoding/json"
	"github.com/folivorra/studyDir/tree/develop/goRedis/internal/model"
	"github.com/folivorra/studyDir/tree/develop/goRedis/internal/storage"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ItemController struct {
	Store storage.Storager
}

func NewItemController(store storage.Storager) *ItemController {
	return &ItemController{Store: store}
}

func (i *ItemController) AddItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.Name == "" || item.ID <= 0 || item.Price < 0 {
		http.Error(w, "Invalid item data", http.StatusBadRequest)
		return
	}

	if err := i.Store.AddItem(item); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (i *ItemController) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err = i.Store.DeleteItem(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (i *ItemController) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := i.Store.GetAllItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (i *ItemController) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, err := i.Store.GetItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (i *ItemController) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.Name == "" || item.Price < 0 {
		http.Error(w, "Invalid item data", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID in URL", http.StatusBadRequest)
		return
	}

	item.ID = id

	// если id в json и в url не совпадают - приоритет отдается url-значению

	if err = i.Store.UpdateItem(item); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (i *ItemController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/items", i.AddItem).Methods("POST")
	r.HandleFunc("/items", i.GetAllItems).Methods("GET")
	r.HandleFunc("/items/{id}", i.GetItem).Methods("GET")
	r.HandleFunc("/items/{id}", i.UpdateItem).Methods("PUT")
	r.HandleFunc("/items/{id}", i.DeleteItem).Methods("DELETE")
}
