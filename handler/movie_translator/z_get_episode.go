package movie_translator_handler

import (
	"net/http"
	"ro-backend/core"
	"strconv"

	"github.com/gorilla/mux"
)

func (m movieTranslatorHandler) GetEpisode(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	ss, err := strconv.ParseFloat(pathVars["ss"], 32)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	ep, err := strconv.ParseFloat(pathVars["ep"], 32)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	result, err := m.service.GetEpisode(float32(ss), float32(ep))
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, result)
}
