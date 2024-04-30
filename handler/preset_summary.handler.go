package handler

import (
	"net/http"
	"ro-backend/appError"
	"ro-backend/core"
	"ro-backend/service"
)

type PresetSummaryHandler interface {
	GenerateSummary(http.ResponseWriter, *http.Request)
}

func NewPresetSummaryHandler(s service.PresetSummaryService) PresetSummaryHandler {
	return presetSummaryHandler{s: s}
}

type presetSummaryHandler struct {
	s service.PresetSummaryService
}

type GenerateSummaryResponse struct {
	Id      string                `json:"id"`
	Summary service.PresetSummary `json:"summary"`
}

func (h presetSummaryHandler) GenerateSummary(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	_, err := h.s.GenerateSummary()
	if err != nil {
		core.WriteErr(w, appError.ErrUnAuthentication)
		return
	}

	var response = GenerateSummaryResponse{
		Id: userId,
	}

	core.WriteOK(w, response)
}
