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
	RoPresetService  service.RoPresetService
	PresetTagService service.PresetTagService
	UserService      service.UserService
}

func NewRoPresetHandler(p RoPresetHandlerParam) RoPresetHandler {
	return roPresetHandler{
		roPresetService:  p.RoPresetService,
		userService:      p.UserService,
		presetTagService: p.PresetTagService,
	}
}

type RoPresetHandler interface {
	GetMyPresetById(http.ResponseWriter, *http.Request)
	GetMyPresets(http.ResponseWriter, *http.Request)
	SearchPresetTags(http.ResponseWriter, *http.Request)
	CreatePreset(http.ResponseWriter, *http.Request)
	BulkCreatePresets(http.ResponseWriter, *http.Request)
	UpdateMyPreset(http.ResponseWriter, *http.Request)
	AddTags(http.ResponseWriter, *http.Request)
	RemoveTags(http.ResponseWriter, *http.Request)
	LikeTag(http.ResponseWriter, *http.Request)
	UnLikeTag(http.ResponseWriter, *http.Request)
	DeleteById(http.ResponseWriter, *http.Request)
}

type roPresetHandler struct {
	roPresetService  service.RoPresetService
	userService      service.UserService
	presetTagService service.PresetTagService
}

type PartialSearchRoPresetInput struct {
	Id      string   `json:"id,omitempty"`
	ClassId int      `json:"class_id,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Skip    int
	Take    int
}

type SearchPresetTagItem struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	Model     repository.PresetModel `json:"model"`
	Tags      map[string]int         `json:"tags"`
	Liked     bool                   `json:"Liked"`
	CreatedAt time.Time              `json:"createdAt"`
}

type SearchPresetTagsResponse struct {
	Items      []SearchPresetTagItem `json:"items"`
	TotalItems int                   `json:"totalItem"`
	Skip       int                   `json:"skip"`
	Take       int                   `json:"take"`
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

type CreateTagRequest struct {
	Tags []string `json:"tags"`
}

type LikeTagResponse struct {
	Id        string `json:"id"`
	Tag       string `json:"tag"`
	ClassId   int    `json:"classId"`
	PresetId  string `json:"presetId"`
	TotalLike int    `json:"totalLike"`
}

func (h roPresetHandler) AddTags(w http.ResponseWriter, r *http.Request) {
	var d CreateTagRequest
	json.NewDecoder(r.Body).Decode(&d)

	presetId := mux.Vars(r)["presetId"]
	userId := r.Header.Get("userId")

	res, err := h.presetTagService.CreateTags(repository.CreateTagInput{
		PublisherId: userId,
		PresetId:    presetId,
		Tags:        d.Tags,
	})
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
	pathVars := mux.Vars(r)
	presetId := pathVars["presetId"]
	tagId := pathVars["tagId"]
	userId := r.Header.Get("userId")

	res, err := h.presetTagService.DeleteTag(service.DeleteTagInput{
		TagId:    tagId,
		UserId:   userId,
		PresetId: presetId,
	})
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

func (h roPresetHandler) LikeTag(w http.ResponseWriter, r *http.Request) {
	tagId := mux.Vars(r)["tagId"]
	userId := r.Header.Get("userId")

	res, err := h.presetTagService.LikeTag(repository.LikeTagInput{
		Id:     tagId,
		UserId: userId,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := LikeTagResponse{
		Id:        res.Id,
		Tag:       res.Tag,
		ClassId:   res.ClassId,
		PresetId:  res.PresetId,
		TotalLike: res.TotalLike,
	}

	WriteOK(w, response)
}

func (h roPresetHandler) UnLikeTag(w http.ResponseWriter, r *http.Request) {
	tagId := mux.Vars(r)["tagId"]
	userId := r.Header.Get("userId")

	res, err := h.presetTagService.UnLikeTag(repository.LikeTagInput{
		Id:     tagId,
		UserId: userId,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := LikeTagResponse{
		Id:        res.Id,
		Tag:       res.Tag,
		ClassId:   res.ClassId,
		PresetId:  res.PresetId,
		TotalLike: res.TotalLike,
	}

	WriteOK(w, response)
}

func (h roPresetHandler) SearchPresetTags(w http.ResponseWriter, r *http.Request) {
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

	userId := r.Header.Get("userId")
	res, err := h.presetTagService.PartialSearchTags(repository.PartialSearchTagsInput{
		ClassId: classId,
		Tag:     tag,
	}, service.PartialSearchMetaInput{
		UserId: userId,
		Skip:   skip,
		Limit:  take,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	items := []SearchPresetTagItem{}
	for _, v := range res.Items {
		items = append(items, SearchPresetTagItem{
			Id:        v.Id,
			Name:      v.Name,
			Model:     v.Model,
			Tags:      v.Tags,
			Liked:     v.Liked,
			CreatedAt: v.CreatedAt,
		})
	}

	response := SearchPresetTagsResponse{
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

	for i := 0; i < len(res); i++ {
		response = append(response, BulkCreatePresetsResponse{
			Label: (res)[i].Label,
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
	for _, v := range res {
		response = append(response, GetMyPresetsResponse{
			Id:        v.Id,
			Label:     v.Label,
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
