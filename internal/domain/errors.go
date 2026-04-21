package domain

import "errors"

var (
	ErrSessionNotFound         = errors.New("session not found")
	ErrQuestionNotFound        = errors.New("question not found")
	ErrResultNotFound          = errors.New("result not found")
	ErrSessionNotFinished      = errors.New("session is not finished")
	ErrQuestionAlreadyAnswered = errors.New("question already answered in this session")
	ErrInvalidAnswerOption     = errors.New("invalid answer option")
	ErrSessionFinished         = errors.New("session already finished")
)
