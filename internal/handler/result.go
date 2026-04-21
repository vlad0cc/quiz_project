package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"profil-math/internal/domain"
	"profil-math/internal/usecase"
)

type ResultHandler struct {
	useCase usecase.QuizUseCase
}

func NewResultHandler(useCase usecase.QuizUseCase) *ResultHandler {
	return &ResultHandler{useCase: useCase}
}

func (h *ResultHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionID, err := uuid.Parse(chi.URLParam(r, "sessionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	result, err := h.useCase.GetResult(r.Context(), sessionID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSessionNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrSessionNotFinished):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to load result")
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}
