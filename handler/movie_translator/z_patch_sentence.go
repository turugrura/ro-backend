package movie_translator_handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/core"
	"ro-backend/repository"
	"strconv"

	"github.com/gorilla/mux"
)

func (m movieTranslatorHandler) PatchSentence(w http.ResponseWriter, r *http.Request) {
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

	var requestInput repository.PatchSentenceInput
	json.NewDecoder(r.Body).Decode(&requestInput)

	requestInput.UpdatedBy = r.Header.Get("userId")

	result, err := m.service.PatchSentence(float32(ss), float32(ep), requestInput)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, &result)
}
