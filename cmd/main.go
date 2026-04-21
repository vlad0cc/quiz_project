package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"profil-math/internal/bootstrap"
	"profil-math/internal/config"
	"profil-math/internal/handler"
	"profil-math/internal/repository"
	"profil-math/internal/service"
	"profil-math/internal/usecase"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx := context.Background()
	if err := bootstrap.RunMigrations(cfg.DatabaseURL(), cfg.MigrationsPath); err != nil {
		logger.Error("failed to run migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	sessionRepo := repository.NewSessionRepository(pool)
	questionRepo := repository.NewQuestionRepository(pool)
	answerRepo := repository.NewAnswerRepository(pool)
	resultRepo := repository.NewResultRepository(pool)

	selector := service.NewQuestionSelector(questionRepo, answerRepo, cfg.MaxSteps)
	scoring := service.NewScoringService()
	recommendation := service.NewRecommendationService()

	quizUseCase := usecase.NewQuizUseCase(
		sessionRepo,
		questionRepo,
		answerRepo,
		resultRepo,
		selector,
		scoring,
		recommendation,
		cfg.MaxSteps,
	)

	router := handler.NewRouter(cfg, logger, quizUseCase)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.AppPort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server started", slog.String("port", cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped unexpectedly", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
}
