package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"profil-math/internal/domain"
)

type resultRepository struct {
	pool *pgxpool.Pool
}

func NewResultRepository(pool *pgxpool.Pool) ResultRepository {
	return &resultRepository{pool: pool}
}

func (r *resultRepository) Create(ctx context.Context, result *domain.Result) error {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO results (session_id, score_ad, score_zi, recommendation, confidence, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (session_id) DO UPDATE
		 SET score_ad = EXCLUDED.score_ad,
		     score_zi = EXCLUDED.score_zi,
		     recommendation = EXCLUDED.recommendation,
		     confidence = EXCLUDED.confidence,
		     created_at = EXCLUDED.created_at`,
		result.SessionID,
		result.ScoreAD,
		result.ScoreZI,
		result.Recommendation,
		result.Confidence,
		result.CreatedAt,
	)
	return err
}

func (r *resultRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*domain.Result, error) {
	var result domain.Result
	err := r.pool.QueryRow(
		ctx,
		`SELECT id, session_id, score_ad, score_zi, recommendation, confidence, created_at
		 FROM results
		 WHERE session_id = $1`,
		sessionID,
	).Scan(
		&result.ID,
		&result.SessionID,
		&result.ScoreAD,
		&result.ScoreZI,
		&result.Recommendation,
		&result.Confidence,
		&result.CreatedAt,
	)
	if isNoRows(err) {
		return nil, domain.ErrResultNotFound
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}
