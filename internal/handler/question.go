package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"profil-math/internal/domain"
	"profil-math/internal/usecase"
)

type QuestionHandler struct {
	useCase usecase.QuizUseCase
}

func NewQuestionHandler(useCase usecase.QuizUseCase) *QuestionHandler {
	return &QuestionHandler{useCase: useCase}
}

func (h *QuestionHandler) GetNext(w http.ResponseWriter, r *http.Request) {
	sessionID, err := uuid.Parse(chi.URLParam(r, "sessionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	question, err := h.useCase.GetNextQuestion(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load next question")
		return
	}

	if question == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeJSON(w, http.StatusOK, question)
}
