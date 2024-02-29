package handler

import (
	"encoding/json"
	"fmt"
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
	GetMyEntirePresets(http.ResponseWriter, *http.Request)
	SearchPresetTags(http.ResponseWriter, *http.Request)
	CreatePreset(http.ResponseWriter, *http.Request)
	BulkCreatePresets(http.ResponseWriter, *http.Request)
	UpdateMyPreset(http.ResponseWriter, *http.Request)
	PublishMyPreset(http.ResponseWriter, *http.Request)
	UnPublishMyPreset(http.ResponseWriter, *http.Request)
	AddTags(http.ResponseWriter, *http.Request)
	BulkOperationTags(http.ResponseWriter, *http.Request)
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
	Id          string                 `json:"id"`
	PublishName string                 `json:"publishName"`
	Model       repository.PresetModel `json:"model"`
	Tags        map[string]int         `json:"tags"`
	Liked       bool                   `json:"liked"`
	CreatedAt   time.Time              `json:"createdAt"`
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

type TagWithLiked struct {
	Id          string    `json:"id"`
	PublisherId string    `json:"publisherId"`
	Tag         string    `json:"tag"`
	ClassId     int       `json:"classId"`
	PresetId    string    `json:"presetId"`
	TotalLike   int       `json:"totalLike"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Liked       bool      `json:"liked"`
}

type GetMyPresetsResponse struct {
	Id          string         `json:"id"`
	Label       string         `json:"label"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	PublishName string         `json:"publishName"`
	IsPublished bool           `json:"isPublished"`
	PublishedAt time.Time      `json:"publishedAt"`
	Tags        []TagWithLiked `json:"tags"`
}

func (r *GetMyPresetsResponse) From(p service.PresetWithTags) {
	r.Id = p.Id
	r.Label = p.Label
	r.CreatedAt = p.CreatedAt
	r.UpdatedAt = p.UpdatedAt
	r.PublishName = p.PublishName
	r.IsPublished = p.IsPublished
	r.PublishedAt = p.PublishedAt

	tags := []TagWithLiked{}
	for _, v := range p.Tags {
		tags = append(tags, TagWithLiked{
			Id:          v.Id,
			PublisherId: v.PublisherId,
			Tag:         v.Tag,
			ClassId:     v.ClassId,
			PresetId:    v.PresetId,
			TotalLike:   v.TotalLike,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			Liked:       v.Liked,
		})
	}
	r.Tags = tags
}

type GetMyEntirePresetsResponse struct {
	Id          string                 `json:"id"`
	Label       string                 `json:"label"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	PublishName string                 `json:"publishName"`
	IsPublished bool                   `json:"isPublished"`
	PublishedAt time.Time              `json:"publishedAt"`
	Tags        []TagWithLiked         `json:"tags"`
	Model       repository.PresetModel `json:"model"`
}

func (r *GetMyEntirePresetsResponse) From(p service.PresetWithTags) {
	r.Id = p.Id
	r.Label = p.Label
	r.CreatedAt = p.CreatedAt
	r.UpdatedAt = p.UpdatedAt
	r.PublishName = p.PublishName
	r.IsPublished = p.IsPublished
	r.PublishedAt = p.PublishedAt
	r.Model = p.Model

	tags := []TagWithLiked{}
	for _, v := range p.Tags {
		tags = append(tags, TagWithLiked{
			Id:          v.Id,
			PublisherId: v.PublisherId,
			Tag:         v.Tag,
			ClassId:     v.ClassId,
			PresetId:    v.PresetId,
			TotalLike:   v.TotalLike,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			Liked:       v.Liked,
		})
	}
	r.Tags = tags
}

type UpsertTagResponse struct {
	Id        string    `json:"id"`
	Label     string    `json:"label"`
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

type PublishPresetRequest struct {
	PublishName string `json:"publishName"`
}

type BulkErrResponse struct {
	ErrMsg string `json:"errorMessage"`
}

type BulkOperationRequest struct {
	PublisherId string   `json:"publisherId"`
	ClassId     int      `json:"classId"`
	PresetId    string   `json:"presetId"`
	CreateTags  []string `json:"createTags"`
	DeleteTags  []string `json:"deleteTags"`
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

	response := GetMyPresetsResponse{}
	response.From(*res)

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

	response := GetMyPresetsResponse{}
	response.From(*res)

	WriteOK(w, response)
}

func (h roPresetHandler) BulkOperationTags(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	presetId := pathVars["presetId"]
	userId := r.Header.Get("userId")

	var req BulkOperationRequest
	json.NewDecoder(r.Body).Decode(&req)

	req.PublisherId = userId
	req.PresetId = presetId

	res, err := h.presetTagService.BulkOperationTags(service.BulkOperationInput(req))
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := GetMyPresetsResponse{}
	response.From(*res)

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
			Id:          v.Id,
			PublishName: v.PublishName,
			Model:       v.Model,
			Tags:        v.Tags,
			Liked:       v.Liked,
			CreatedAt:   v.CreatedAt,
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

	presetId := mux.Vars(r)["presetId"]
	d.UserId = r.Header.Get("userId")

	res, err := h.roPresetService.UpdatePreset(presetId, d)
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

func (h roPresetHandler) PublishMyPreset(w http.ResponseWriter, r *http.Request) {
	var d PublishPresetRequest
	json.NewDecoder(r.Body).Decode(&d)

	presetId := mux.Vars(r)["presetId"]
	userId := r.Header.Get("userId")

	res, err := h.roPresetService.PublishPreset(presetId, repository.UpdatePresetInput{
		PublishName: d.PublishName,
		UserId:      userId,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	var response GetMyPresetsResponse
	response.From(service.PresetWithTags{
		RoPreset: *res,
	})

	WriteOK(w, response)
}

func (h roPresetHandler) UnPublishMyPreset(w http.ResponseWriter, r *http.Request) {
	presetId := mux.Vars(r)["presetId"]
	userId := r.Header.Get("userId")

	res, err := h.roPresetService.UnPublishPreset(presetId, repository.UpdatePresetInput{
		UserId: userId,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	var response GetMyPresetsResponse
	response.From(service.PresetWithTags{
		RoPreset: *res,
	})

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

	if len(d.BulkData) == 0 {
		WriteOK(w, []repository.RoPreset{})
		return
	}

	errResponse := []BulkErrResponse{}
	for i, v := range d.BulkData {
		if v.Label == "" {
			errResponse = append(errResponse, BulkErrResponse{
				ErrMsg: fmt.Sprintf("preset number '%v' has empty label", i),
			})
		} else if err := v.Model.Validate(); err != nil {
			errResponse = append(errResponse, BulkErrResponse{
				ErrMsg: fmt.Sprintf("preset '%v' is invalid", v.Label),
			})
		}
	}
	if len(errResponse) > 0 {
		WriteErrObj(w, http.StatusBadRequest, errResponse)
		return
	}

	userId := r.Header.Get("userId")
	d.UserId = userId

	res, err := h.roPresetService.BulkCreatePresets(d)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	WriteCreated(w, res)
}

func (h roPresetHandler) GetMyPresets(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	res, err := h.roPresetService.FindPresetsByUserId(userId, false)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	presetWithTags, err := h.presetTagService.AttachTags(userId, res)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := []GetMyPresetsResponse{}
	for _, v := range presetWithTags {
		var r GetMyPresetsResponse
		r.From(v)
		response = append(response, r)
	}

	WriteOK(w, response)
}

func (h roPresetHandler) GetMyEntirePresets(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	res, err := h.roPresetService.FindPresetsByUserId(userId, true)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	presetWithTags, err := h.presetTagService.AttachTags(userId, res)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	response := []GetMyEntirePresetsResponse{}
	for _, v := range presetWithTags {
		var r GetMyEntirePresetsResponse
		r.From(v)
		response = append(response, r)
	}

	WriteOK(w, response)
}

func (h roPresetHandler) CreatePreset(w http.ResponseWriter, r *http.Request) {
	var d repository.CreatePresetInput
	json.NewDecoder(r.Body).Decode(&d)
	if err := d.Validate(); err != nil {
		WriteErr(w, err.Error())
		return
	}

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
