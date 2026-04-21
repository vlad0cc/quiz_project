package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"profil-math/internal/domain"
	"profil-math/internal/dto"
	"profil-math/internal/usecase"
)

type AnswerHandler struct {
	useCase usecase.QuizUseCase
}

func NewAnswerHandler(useCase usecase.QuizUseCase) *AnswerHandler {
	return &AnswerHandler{useCase: useCase}
}

func (h *AnswerHandler) Submit(w http.ResponseWriter, r *http.Request) {
	sessionID, err := uuid.Parse(chi.URLParam(r, "sessionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	var input dto.SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.useCase.SubmitAnswer(r.Context(), sessionID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidAnswerOption):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrSessionNotFound), errors.Is(err, domain.ErrQuestionNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrQuestionAlreadyAnswered), errors.Is(err, domain.ErrSessionFinished):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to submit answer")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}
