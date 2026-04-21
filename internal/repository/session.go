package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"profil-math/internal/domain"
)

type sessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) SessionRepository {
	return &sessionRepository{pool: pool}
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) error {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO sessions (id, created_at, step, score_ad, score_zi) VALUES ($1, $2, $3, $4, $5)`,
		session.ID,
		session.CreatedAt,
		session.Step,
		session.ScoreAD,
		session.ScoreZI,
	)
	return err
}

func (r *sessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	var session domain.Session
	err := r.pool.QueryRow(
		ctx,
		`SELECT id, created_at, finished_at, step, score_ad, score_zi FROM sessions WHERE id = $1`,
		id,
	).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.FinishedAt,
		&session.Step,
		&session.ScoreAD,
		&session.ScoreZI,
	)
	if isNoRows(err) {
		return nil, domain.ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *domain.Session) error {
	_, err := r.pool.Exec(
		ctx,
		`UPDATE sessions SET finished_at = $2, step = $3, score_ad = $4, score_zi = $5 WHERE id = $1`,
		session.ID,
		session.FinishedAt,
		session.Step,
		session.ScoreAD,
		session.ScoreZI,
	)
	return err
}
