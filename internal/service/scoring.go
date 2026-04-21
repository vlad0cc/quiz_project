package service

import "profil-math/internal/domain"

type ScoringService interface {
	ApplyAnswer(session *domain.Session, question *domain.Question, isCorrect bool)
}

type scoringService struct{}

func NewScoringService() ScoringService {
	return &scoringService{}
}

func (s *scoringService) ApplyAnswer(session *domain.Session, question *domain.Question, isCorrect bool) {
	if !isCorrect {
		return
	}

	scoreDelta := question.Weight * difficultyMultiplier(question.Difficulty)

	switch question.Profile {
	case domain.ProfileAD:
		session.ScoreAD += scoreDelta
	case domain.ProfileZI:
		session.ScoreZI += scoreDelta
	}
}

func difficultyMultiplier(difficulty int) float64 {
	switch difficulty {
	case 2:
		return 1.5
	case 3:
		return 2.0
	default:
		return 1.0
	}
}
