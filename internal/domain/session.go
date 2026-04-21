package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	FinishedAt *time.Time
	Step       int
	ScoreAD    float64
	ScoreZI    float64
}

func NewSession() *Session {
	return &Session{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
	}
}

func (s *Session) IsFinished(maxSteps int) bool {
	return s.FinishedAt != nil || s.Step >= maxSteps
}
