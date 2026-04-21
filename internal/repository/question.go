package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"profil-math/internal/domain"
)

type questionRepository struct {
	pool *pgxpool.Pool
}

func NewQuestionRepository(pool *pgxpool.Pool) QuestionRepository {
	return &questionRepository{pool: pool}
}

func (r *questionRepository) GetByID(ctx context.Context, id int) (*domain.Question, error) {
	var question domain.Question
	err := r.pool.QueryRow(
		ctx,
		`SELECT id, text, option_a, option_b, option_c, option_d, correct_option, profile, difficulty, weight
		 FROM questions
		 WHERE id = $1`,
		id,
	).Scan(
		&question.ID,
		&question.Text,
		&question.OptionA,
		&question.OptionB,
		&question.OptionC,
		&question.OptionD,
		&question.CorrectOption,
		&question.Profile,
		&question.Difficulty,
		&question.Weight,
	)
	if isNoRows(err) {
		return nil, domain.ErrQuestionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) List(ctx context.Context, filter QuestionFilter) ([]domain.Question, error) {
	query := `SELECT id, text, option_a, option_b, option_c, option_d, correct_option, profile, difficulty, weight
			  FROM questions`
	conditions := make([]string, 0, 3)
	args := make([]any, 0, 4)

	if len(filter.ExcludeQuestionIDs) > 0 {
		ids := make([]int32, 0, len(filter.ExcludeQuestionIDs))
		for _, id := range filter.ExcludeQuestionIDs {
			ids = append(ids, int32(id))
		}
		args = append(args, ids)
		conditions = append(conditions, fmt.Sprintf("id <> ALL($%d)", len(args)))
	}

	if filter.Profile != nil {
		args = append(args, *filter.Profile)
		conditions = append(conditions, fmt.Sprintf("profile = $%d", len(args)))
	}

	if filter.Difficulty != nil {
		args = append(args, *filter.Difficulty)
		conditions = append(conditions, fmt.Sprintf("difficulty = $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY random()"
	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := make([]domain.Question, 0)
	for rows.Next() {
		var question domain.Question
		if err := rows.Scan(
			&question.ID,
			&question.Text,
			&question.OptionA,
			&question.OptionB,
			&question.OptionC,
			&question.OptionD,
			&question.CorrectOption,
			&question.Profile,
			&question.Difficulty,
			&question.Weight,
		); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}
