package handler

import (
	"net/http"

	"profil-math/internal/dto"
	"profil-math/internal/usecase"
)

type SessionHandler struct {
	useCase usecase.QuizUseCase
}

func NewSessionHandler(useCase usecase.QuizUseCase) *SessionHandler {
	return &SessionHandler{useCase: useCase}
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	sessionID, err := h.useCase.StartSession(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateSessionResponse{SessionID: sessionID.String()})
}
