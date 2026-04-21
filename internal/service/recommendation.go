package service

import (
	"math"
	"time"

	"profil-math/internal/domain"
)

type RecommendationService interface {
	BuildResult(session *domain.Session) *domain.Result
	BuildBreakdown(result *domain.Result) (int, int)
}

type recommendationService struct{}

func NewRecommendationService() RecommendationService {
	return &recommendationService{}
}

func (s *recommendationService) BuildResult(session *domain.Session) *domain.Result {
	recommendation := domain.ProfileAD
	if session.ScoreZI > session.ScoreAD {
		recommendation = domain.ProfileZI
	}

	total := session.ScoreAD + session.ScoreZI
	confidence := 0.5
	if total > 0 {
		maxScore := math.Max(session.ScoreAD, session.ScoreZI)
		confidence = maxScore / total
	}

	return &domain.Result{
		SessionID:      session.ID,
		ScoreAD:        session.ScoreAD,
		ScoreZI:        session.ScoreZI,
		Recommendation: recommendation,
		Confidence:     confidence,
		CreatedAt:      time.Now().UTC(),
	}
}

func (s *recommendationService) BuildBreakdown(result *domain.Result) (int, int) {
	total := result.ScoreAD + result.ScoreZI
	if total == 0 {
		return 50, 50
	}

	adPercent := int(math.Round((result.ScoreAD / total) * 100))
	if adPercent < 0 {
		adPercent = 0
	}
	if adPercent > 100 {
		adPercent = 100
	}

	return adPercent, 100 - adPercent
}
