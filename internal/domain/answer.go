package domain

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         int
	SessionID  uuid.UUID
	QuestionID int
	Chosen     string
	IsCorrect  bool
	AnsweredAt time.Time
}
