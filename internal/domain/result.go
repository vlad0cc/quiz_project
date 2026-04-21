package domain

import (
	"time"

	"github.com/google/uuid"
)

type Result struct {
	ID             int
	SessionID      uuid.UUID
	ScoreAD        float64
	ScoreZI        float64
	Recommendation Profile
	Confidence     float64
	CreatedAt      time.Time
}
