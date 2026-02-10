package app

import (
	"URL_checker/internal/checker"
	checksHandler "URL_checker/internal/handler/checks"
	"URL_checker/internal/handler/url"
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/target"
	"URL_checker/internal/scheduler"
	"URL_checker/internal/service"
	"URL_checker/internal/workerpool"
	"URL_checker/internal/writer"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Run() error {
	dsn := "postgres://pavelpavlov@localhost:5432/postgres?sslmode=disable"
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

	targetService := service.New(targetRepo)
	targetHandler := url.New(targetService)

	checkRepo, errCheck := checks.New(db)
	if errCheck != nil {
		return err
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

	checkService := service.NewCheckService(checkRepo)
	checkHandler := checksHandler.NewCheckHandler(checkService)

	router := gin.Default()
	url.RegisterRoutes(router, targetHandler)
	checksHandler.RegisterRoutes(router, checkHandler)

	return router.Run(":8080")
}
