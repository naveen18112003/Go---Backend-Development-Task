package main

import (
	"context"
	"log"
	"os"

	"user-age-api/config"
	"user-age-api/db/sqlc"
	"user-age-api/internal/handler"
	"user-age-api/internal/logger"
	"user-age-api/internal/middleware"
	"user-age-api/internal/repository"
	"user-age-api/internal/routes"
	"user-age-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logr := logger.New(os.Getenv("APP_ENV"))
	defer logr.Sync() // nolint:errcheck

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{"error": e.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(logr))

	pool := setupDB(cfg, logr)
	defer pool.Close()

	queries := db.New(pool)
	repo := repository.NewUserRepository(queries)
	validate := validator.New(validator.WithRequiredStructEnabled())
	userService := service.NewUserService(repo, validate)
	userHandler := handler.NewUserHandler(userService, logr)

	routes.Register(app, userHandler)

	logr.Info("starting server", zap.String("port", cfg.AppPort))
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func setupDB(cfg config.Config, logr *zap.Logger) *pgxpool.Pool {
	ctx := context.Background()
	pgxCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to parse database url: %v", err)
	}
	pgxCfg.MaxConns = cfg.DBMaxConns
	pgxCfg.MaxConnIdleTime = cfg.DBMaxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		log.Fatalf("unable to create db pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}

	logr.Info("database connected")
	return pool
}



