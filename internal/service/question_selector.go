package service

import (
	"context"
	"math"

	"profil-math/internal/domain"
	"profil-math/internal/repository"
)

type QuestionSelector interface {
	SelectNextQuestion(ctx context.Context, session *domain.Session) (*domain.Question, error)
}

type questionSelector struct {
	questions repository.QuestionRepository
	answers   repository.AnswerRepository
	maxSteps  int
}

func NewQuestionSelector(
	questions repository.QuestionRepository,
	answers repository.AnswerRepository,
	maxSteps int,
) QuestionSelector {
	return &questionSelector{
		questions: questions,
		answers:   answers,
		maxSteps:  maxSteps,
	}
}

func (s *questionSelector) SelectNextQuestion(ctx context.Context, session *domain.Session) (*domain.Question, error) {
	if session.IsFinished(s.maxSteps) {
		return nil, nil
	}

	answeredIDs, err := s.answers.ListAnsweredQuestionIDs(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	if len(answeredIDs) >= s.maxSteps {
		return nil, nil
	}

	if len(answeredIDs) < 4 {
		return s.selectByFilter(ctx, repository.QuestionFilter{
			ExcludeQuestionIDs: answeredIDs,
			Difficulty:         intPtr(1),
			Limit:              1,
		})
	}

	targetDifficulty := int(math.Ceil(float64(session.Step) / 3.0))
	if targetDifficulty < 1 {
		targetDifficulty = 1
	}
	if targetDifficulty > 3 {
		targetDifficulty = 3
	}

	laggingProfile := domain.ProfileAD
	if session.ScoreZI < session.ScoreAD {
		laggingProfile = domain.ProfileZI
	} else if session.ScoreAD == session.ScoreZI {
		laggingProfile = ""
	}

	if laggingProfile != "" {
		question, err := s.selectByFilter(ctx, repository.QuestionFilter{
			ExcludeQuestionIDs: answeredIDs,
			Profile:            &laggingProfile,
			Difficulty:         &targetDifficulty,
			Limit:              1,
		})
		if err != nil || question != nil {
			return question, err
		}
	}

	question, err := s.selectByFilter(ctx, repository.QuestionFilter{
		ExcludeQuestionIDs: answeredIDs,
		Difficulty:         &targetDifficulty,
		Limit:              1,
	})
	if err != nil || question != nil {
		return question, err
	}

	return s.selectByFilter(ctx, repository.QuestionFilter{
		ExcludeQuestionIDs: answeredIDs,
		Limit:              1,
	})
}

func (s *questionSelector) selectByFilter(ctx context.Context, filter repository.QuestionFilter) (*domain.Question, error) {
	questions, err := s.questions.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, nil
	}
	return &questions[0], nil
}

func intPtr(value int) *int {
	return &value
}
