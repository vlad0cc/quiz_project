package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"profil-math/internal/domain"
)

type answerRepository struct {
	pool *pgxpool.Pool
}

func NewAnswerRepository(pool *pgxpool.Pool) AnswerRepository {
	return &answerRepository{pool: pool}
}

func (r *answerRepository) Create(ctx context.Context, answer *domain.Answer) error {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO answers (session_id, question_id, chosen, is_correct, answered_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		answer.SessionID,
		answer.QuestionID,
		answer.Chosen,
		answer.IsCorrect,
		answer.AnsweredAt,
	)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrQuestionAlreadyAnswered
	}

	return err
}

func (r *answerRepository) Exists(ctx context.Context, sessionID uuid.UUID, questionID int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM answers WHERE session_id = $1 AND question_id = $2)`,
		sessionID,
		questionID,
	).Scan(&exists)
	return exists, err
}

func (r *answerRepository) ListAnsweredQuestionIDs(ctx context.Context, sessionID uuid.UUID) ([]int, error) {
	rows, err := r.pool.Query(ctx, `SELECT question_id FROM answers WHERE session_id = $1 ORDER BY answered_at`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return ids, nil
}
