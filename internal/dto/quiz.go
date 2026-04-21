package dto

type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
}

type QuestionResponse struct {
	QuestionID int               `json:"question_id"`
	Text       string            `json:"text"`
	Options    map[string]string `json:"options"`
	Step       int               `json:"step"`
	TotalSteps int               `json:"total_steps"`
}

type SubmitAnswerRequest struct {
	QuestionID int    `json:"question_id"`
	Chosen     string `json:"chosen"`
}

type SubmitAnswerResponse struct {
	IsCorrect  bool `json:"is_correct"`
	Step       int  `json:"step"`
	TotalSteps int  `json:"total_steps"`
}

type ResultBreakdown struct {
	ADPercent int `json:"ad_percent"`
	ZIPercent int `json:"zi_percent"`
}

type ResultResponse struct {
	Recommendation      string          `json:"recommendation"`
	RecommendationLabel string          `json:"recommendation_label"`
	Confidence          float64         `json:"confidence"`
	ScoreAD             float64         `json:"score_ad"`
	ScoreZI             float64         `json:"score_zi"`
	Breakdown           ResultBreakdown `json:"breakdown"`
}
