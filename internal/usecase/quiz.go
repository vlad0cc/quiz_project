package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"profil-math/internal/domain"
	"profil-math/internal/dto"
	"profil-math/internal/repository"
	"profil-math/internal/service"
)

type QuizUseCase interface {
	StartSession(ctx context.Context) (uuid.UUID, error)
	GetNextQuestion(ctx context.Context, sessionID uuid.UUID) (*dto.QuestionResponse, error)
	SubmitAnswer(ctx context.Context, sessionID uuid.UUID, input dto.SubmitAnswerRequest) (*dto.SubmitAnswerResponse, error)
	GetResult(ctx context.Context, sessionID uuid.UUID) (*dto.ResultResponse, error)
}

type quizUseCase struct {
	sessions       repository.SessionRepository
	questions      repository.QuestionRepository
	answers        repository.AnswerRepository
	results        repository.ResultRepository
	selector       service.QuestionSelector
	scoring        service.ScoringService
	recommendation service.RecommendationService
	maxSteps       int
}

func NewQuizUseCase(
	sessions repository.SessionRepository,
	questions repository.QuestionRepository,
	answers repository.AnswerRepository,
	results repository.ResultRepository,
	selector service.QuestionSelector,
	scoring service.ScoringService,
	recommendation service.RecommendationService,
	maxSteps int,
) QuizUseCase {
	return &quizUseCase{
		sessions:       sessions,
		questions:      questions,
		answers:        answers,
		results:        results,
		selector:       selector,
		scoring:        scoring,
		recommendation: recommendation,
		maxSteps:       maxSteps,
	}
}

func (u *quizUseCase) StartSession(ctx context.Context) (uuid.UUID, error) {
	session := domain.NewSession()
	if err := u.sessions.Create(ctx, session); err != nil {
		return uuid.Nil, err
	}
	return session.ID, nil
}

func (u *quizUseCase) GetNextQuestion(ctx context.Context, sessionID uuid.UUID) (*dto.QuestionResponse, error) {
	session, err := u.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	question, err := u.selector.SelectNextQuestion(ctx, session)
	if err != nil {
		return nil, err
	}
	if question == nil {
		if !session.IsFinished(u.maxSteps) {
			now := time.Now().UTC()
			session.FinishedAt = &now
			if err := u.sessions.Update(ctx, session); err != nil {
				return nil, err
			}
			result := u.recommendation.BuildResult(session)
			if err := u.results.Create(ctx, result); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	return &dto.QuestionResponse{
		QuestionID: question.ID,
		Text:       question.Text,
		Options:    question.Options(),
		Step:       session.Step + 1,
		TotalSteps: u.maxSteps,
	}, nil
}

func (u *quizUseCase) SubmitAnswer(ctx context.Context, sessionID uuid.UUID, input dto.SubmitAnswerRequest) (*dto.SubmitAnswerResponse, error) {
	chosen := strings.ToUpper(strings.TrimSpace(input.Chosen))
	if chosen != "A" && chosen != "B" && chosen != "C" && chosen != "D" {
		return nil, domain.ErrInvalidAnswerOption
	}

	session, err := u.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.IsFinished(u.maxSteps) {
		return nil, domain.ErrSessionFinished
	}

	question, err := u.questions.GetByID(ctx, input.QuestionID)
	if err != nil {
		return nil, err
	}

	exists, err := u.answers.Exists(ctx, sessionID, input.QuestionID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrQuestionAlreadyAnswered
	}

	isCorrect := strings.EqualFold(question.CorrectOption, chosen)
	answer := &domain.Answer{
		SessionID:  sessionID,
		QuestionID: input.QuestionID,
		Chosen:     chosen,
		IsCorrect:  isCorrect,
		AnsweredAt: time.Now().UTC(),
	}
	if err := u.answers.Create(ctx, answer); err != nil {
		return nil, err
	}

	u.scoring.ApplyAnswer(session, question, isCorrect)
	session.Step++
	if session.Step >= u.maxSteps {
		now := time.Now().UTC()
		session.FinishedAt = &now
	}
	if err := u.sessions.Update(ctx, session); err != nil {
		return nil, err
	}

	if session.IsFinished(u.maxSteps) {
		result := u.recommendation.BuildResult(session)
		if err := u.results.Create(ctx, result); err != nil {
			return nil, err
		}
	}

	nextStep := session.Step + 1
	if nextStep > u.maxSteps {
		nextStep = u.maxSteps
	}

	return &dto.SubmitAnswerResponse{
		IsCorrect:  isCorrect,
		Step:       nextStep,
		TotalSteps: u.maxSteps,
	}, nil
}

func (u *quizUseCase) GetResult(ctx context.Context, sessionID uuid.UUID) (*dto.ResultResponse, error) {
	session, err := u.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if !session.IsFinished(u.maxSteps) {
		return nil, domain.ErrSessionNotFinished
	}

	result, err := u.results.GetBySessionID(ctx, sessionID)
	if err != nil {
		if err != domain.ErrResultNotFound {
			return nil, err
		}
		result = u.recommendation.BuildResult(session)
		if err := u.results.Create(ctx, result); err != nil {
			return nil, err
		}
	}

	adPercent, ziPercent := u.recommendation.BuildBreakdown(result)

	return &dto.ResultResponse{
		Recommendation:      string(result.Recommendation),
		RecommendationLabel: result.Recommendation.Label(),
		Confidence:          result.Confidence,
		ScoreAD:             result.ScoreAD,
		ScoreZI:             result.ScoreZI,
		Breakdown: dto.ResultBreakdown{
			ADPercent: adPercent,
			ZIPercent: ziPercent,
		},
	}, nil
}
