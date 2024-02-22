package handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/repository"
	"ro-backend/service"

	"github.com/gorilla/mux"
)

type RoPresetHandlerParam struct {
	RoPresetService service.RoPresetService
	UserService     service.UserService
}

func NewRoPresetHandler(p RoPresetHandlerParam) RoPresetHandler {
	return roPresetHandler{
		roPresetService: p.RoPresetService,
		userService:     p.UserService,
	}
}

type RoPresetHandler interface {
	GetPresetById(http.ResponseWriter, *http.Request)
	GetMyPresets(http.ResponseWriter, *http.Request)
	CreatePreset(http.ResponseWriter, *http.Request)
	UpdatePreset(http.ResponseWriter, *http.Request)
	BulkCreatePresets(http.ResponseWriter, *http.Request)
	DeleteById(http.ResponseWriter, *http.Request)
}

type roPresetHandler struct {
	roPresetService service.RoPresetService
	userService     service.UserService
}

// UpdatePreset implements RoPresetHandler.
func (h roPresetHandler) UpdatePreset(w http.ResponseWriter, r *http.Request) {
	var d repository.UpdatePresetInput
	json.NewDecoder(r.Body).Decode(&d)

	d.Id = mux.Vars(r)["presetId"]
	d.UserId = r.Header.Get("userId")

	res, err := h.roPresetService.UpdatePreset(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteById implements RoPresetHandler.
func (h roPresetHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	presetId := mux.Vars(r)["presetId"]

	_, err := h.roPresetService.DeletePresetById(service.CheckPresetOwnerInput{
		Id:     presetId,
		UserId: userId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type BulkCreatePresetsResponse struct {
	Label string `json:"label"`
}

// BulkCreatePresets implements RoPresetHandler.
func (h roPresetHandler) BulkCreatePresets(w http.ResponseWriter, r *http.Request) {
	var d repository.BulkCreatePresetInput
	json.NewDecoder(r.Body).Decode(&d)

	var response = []BulkCreatePresetsResponse{}

	if len(d.BulkData) == 0 {
		json.NewEncoder(w).Encode(response)
	}

	userId := r.Header.Get("userId")
	d.UserId = userId

	res, err := h.roPresetService.BulkCreatePresets(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for i := 0; i < len(*res); i++ {
		response = append(response, BulkCreatePresetsResponse{
			Label: (*res)[i].Label,
		})
	}

	json.NewEncoder(w).Encode(response)
}

// GetMyPresets implements RoPresetHandler.
func (h roPresetHandler) GetMyPresets(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	res, err := h.roPresetService.FindPresetsByUserId(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

// CreatePreset implements RoPresetHandler.
func (h roPresetHandler) CreatePreset(w http.ResponseWriter, r *http.Request) {
	var d repository.CreatePresetInput
	json.NewDecoder(r.Body).Decode(&d)

	d.UserId = r.Header.Get("userId")

	res, err := h.roPresetService.CreatePreset(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

// GetPresetById implements RoPresetHandler.
func (h roPresetHandler) GetPresetById(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	presetId := mux.Vars(r)["presetId"]

	res, err := h.roPresetService.FindPresetById(service.CheckPresetOwnerInput{
		Id:     presetId,
		UserId: userId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
