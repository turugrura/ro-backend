package movie_translator_handler

import (
	"net/http"
	"ro-backend/core"
)

func (m movieTranslatorHandler) GetAllEpisodes(w http.ResponseWriter, r *http.Request) {
	result, err := m.service.GetAllEpisodes()
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, result)
}
