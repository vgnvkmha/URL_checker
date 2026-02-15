package app

import (
	"URL_checker/internal/checker"
	checksHandler "URL_checker/internal/handler/checks"
	url "URL_checker/internal/handler/target"
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/target"
	"URL_checker/internal/scheduler"
	serviceChecker "URL_checker/internal/service/check"
	serviceTarget "URL_checker/internal/service/target"
	"URL_checker/internal/workerpool"
	"URL_checker/internal/writer"
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func Run() error {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	targetRepo, err := target.New(db)
	if err != nil {
		return err
	}

	targetService := serviceTarget.New(targetRepo)
	targetHandler := url.New(targetService)

	checkRepo, errCheck := checks.New(db)
	if errCheck != nil {
		return errCheck
	}

	// ===== checker pipeline =====
	queue := make(chan dto.Targets, 100)
	results := make(chan dto.Checks, 100)

	httpChecker := checker.NewHTTPChecker()
	writer := writer.NewWriter(checkRepo, results)
	workers := workerpool.New(httpChecker, results, 50)
	scheduler := scheduler.New(queue, targetRepo)

	go writer.Run(ctx)
	workers.Start(ctx, queue)
	go scheduler.Run(ctx)

	checkService := serviceChecker.NewCheckService(checkRepo)
	checkHandler := checksHandler.NewCheckHandler(checkService)

	router := gin.Default()
	url.RegisterRoutes(router, targetHandler)
	checksHandler.RegisterRoutes(router, checkHandler)

	return router.Run(":" + port)
}
