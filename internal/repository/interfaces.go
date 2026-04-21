package repository

import (
	"context"

	"github.com/google/uuid"

	"profil-math/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
}

type QuestionFilter struct {
	ExcludeQuestionIDs []int
	Profile            *domain.Profile
	Difficulty         *int
	Limit              int
}

type QuestionRepository interface {
	GetByID(ctx context.Context, id int) (*domain.Question, error)
	List(ctx context.Context, filter QuestionFilter) ([]domain.Question, error)
}

type AnswerRepository interface {
	Create(ctx context.Context, answer *domain.Answer) error
	Exists(ctx context.Context, sessionID uuid.UUID, questionID int) (bool, error)
	ListAnsweredQuestionIDs(ctx context.Context, sessionID uuid.UUID) ([]int, error)
}

type ResultRepository interface {
	Create(ctx context.Context, result *domain.Result) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*domain.Result, error)
}
