package handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/repository"
	"ro-backend/service"
	"strconv"
	"time"

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
	GetMyPresetById(http.ResponseWriter, *http.Request)
	GetMyPresets(http.ResponseWriter, *http.Request)
	GetByClassTag(http.ResponseWriter, *http.Request)
	CreatePreset(http.ResponseWriter, *http.Request)
	BulkCreatePresets(http.ResponseWriter, *http.Request)
	UpdateMyPreset(http.ResponseWriter, *http.Request)
	AddTags(http.ResponseWriter, *http.Request)
	RemoveTags(http.ResponseWriter, *http.Request)
	DeleteById(http.ResponseWriter, *http.Request)
}

type roPresetHandler struct {
	roPresetService service.RoPresetService
	userService     service.UserService
}

type PartialSearchRoPresetInput struct {
	Id      string   `json:"id,omitempty"`
	ClassId int      `json:"class_id,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Skip    int
	Take    int
}

type GetByClassTagItem struct {
	Id    string                 `json:"id"`
	Name  string                 `json:"name"`
	Model repository.PresetModel `json:"model"`
	Tags  []string               `json:"tags"`
}

type GetByClassTagResponse struct {
	Items      []GetByClassTagItem `json:"items"`
	TotalItems int                 `json:"totalItem"`
	Skip       int                 `json:"skip"`
	Take       int                 `json:"take"`
}

type BulkCreatePresetsResponse struct {
	Label string `json:"label"`
}

type GetMyPresetsResponse struct {
	Id        string    `json:"id"`
	Label     string    `json:"label"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (h roPresetHandler) AddTags(w http.ResponseWriter, r *http.Request) {
	var d service.AddTagsRequest
	json.NewDecoder(r.Body).Decode(&d)

	d.Id = mux.Vars(r)["presetId"]
	d.UserId = r.Header.Get("userId")

	res, err := h.roPresetService.AddTags(d)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := GetMyPresetsResponse{
		Id:        res.Id,
		Label:     res.Label,
		Tags:      res.Tags,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}

	WriteOK(w, response)
}

func (h roPresetHandler) RemoveTags(w http.ResponseWriter, r *http.Request) {
	var d service.AddTagsRequest
	json.NewDecoder(r.Body).Decode(&d)

	d.Id = mux.Vars(r)["presetId"]
	d.UserId = r.Header.Get("userId")

	res, err := h.roPresetService.RemoveTags(d)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := GetMyPresetsResponse{
		Id:        res.Id,
		Label:     res.Label,
		Tags:      res.Tags,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}

	WriteOK(w, response)
}

func (h roPresetHandler) GetByClassTag(w http.ResponseWriter, r *http.Request) {
	classId, err := strconv.Atoi(mux.Vars(r)["classId"])
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	tag := mux.Vars(r)["tag"]
	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	take, err := strconv.Atoi(r.URL.Query().Get("take"))
	if err != nil {
		WriteErr(w, err.Error())
		return
	}
	if take == 0 {
		take = 20
	}

	res, err := h.roPresetService.FindPresetsByTags(service.FindPresetsByTagsRequest{
		ClassId: classId,
		Tag:     tag,
		Skip:    skip,
		Take:    take,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	items := []GetByClassTagItem{}
	for _, v := range res.Items {
		items = append(items, GetByClassTagItem{
			Id:    v.Id,
			Name:  v.Name,
			Model: v.Model,
			Tags:  v.Tags,
		})
	}

	response := GetByClassTagResponse{
		Items:      items,
		TotalItems: int(res.Total),
		Skip:       skip,
		Take:       take,
	}

	WriteOK(w, response)
}

func (h roPresetHandler) UpdateMyPreset(w http.ResponseWriter, r *http.Request) {
	var d repository.UpdatePresetInput
	json.NewDecoder(r.Body).Decode(&d)

	d.Id = mux.Vars(r)["presetId"]
	d.UserId = r.Header.Get("userId")

	res, err := h.roPresetService.UpdatePreset(d)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := GetMyPresetsResponse{
		Id:        res.Id,
		Label:     res.Label,
		Tags:      res.Tags,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}

	WriteOK(w, response)
}

func (h roPresetHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	presetId := mux.Vars(r)["presetId"]

	_, err := h.roPresetService.DeletePresetById(service.CheckPresetOwnerRequest{
		Id:     presetId,
		UserId: userId,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	WriteNoContent(w, nil)
}

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
		WriteErr(w, err.Error())
		return
	}

	for i := 0; i < len(*res); i++ {
		response = append(response, BulkCreatePresetsResponse{
			Label: (*res)[i].Label,
		})
	}

	WriteCreated(w, response)
}

func (h roPresetHandler) GetMyPresets(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	res, err := h.roPresetService.FindPresetsByUserId(userId)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := []GetMyPresetsResponse{}
	for _, v := range *res {
		response = append(response, GetMyPresetsResponse{
			Id:        v.Id,
			Label:     v.Label,
			Tags:      v.Tags,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	WriteOK(w, response)
}

func (h roPresetHandler) CreatePreset(w http.ResponseWriter, r *http.Request) {
	var d repository.CreatePresetInput
	json.NewDecoder(r.Body).Decode(&d)

	d.UserId = r.Header.Get("userId")

	res, err := h.roPresetService.CreatePreset(d)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	WriteCreated(w, res)
}

func (h roPresetHandler) GetMyPresetById(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	presetId := mux.Vars(r)["presetId"]

	res, err := h.roPresetService.FindPresetById(service.CheckPresetOwnerRequest{
		Id:     presetId,
		UserId: userId,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	WriteOK(w, res)
}
