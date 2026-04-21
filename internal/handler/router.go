package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"profil-math/internal/config"
	"profil-math/internal/usecase"
)

func NewRouter(cfg config.Config, logger *slog.Logger, quizUseCase usecase.QuizUseCase) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(LoggingMiddleware(logger))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(cfg),
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	sessionHandler := NewSessionHandler(quizUseCase)
	questionHandler := NewQuestionHandler(quizUseCase)
	answerHandler := NewAnswerHandler(quizUseCase)
	resultHandler := NewResultHandler(quizUseCase)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/sessions", sessionHandler.Create)
		r.Get("/sessions/{sessionID}/question", questionHandler.GetNext)
		r.Post("/sessions/{sessionID}/answers", answerHandler.Submit)
		r.Get("/sessions/{sessionID}/result", resultHandler.Get)
	})

	if _, err := os.Stat(cfg.FrontendDist); err == nil {
		router.Handle("/*", NewSPAHandler(cfg.FrontendDist))
	}

	return router
}

func allowedOrigins(cfg config.Config) []string {
	if cfg.AppEnv == "development" {
		return []string{cfg.FrontendDevOrigin}
	}
	return []string{}
}
