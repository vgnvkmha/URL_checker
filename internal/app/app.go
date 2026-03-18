package app

import (
	"URL_checker/internal/checker"
	"URL_checker/internal/configs"
	checksHandler "URL_checker/internal/handler/checks"
	url "URL_checker/internal/handler/target"
	"URL_checker/internal/logger"
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/target"
	"URL_checker/internal/scheduler"
	"URL_checker/internal/service/cache"
	serviceChecker "URL_checker/internal/service/check"
	serviceTarget "URL_checker/internal/service/target"
	"URL_checker/internal/workerpool"
	"URL_checker/internal/writer"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

func Run() error {
	postgresCfg, err := configs.LoadPostgres()
	if err != nil {
		return err
	}
	pgDb, err := configs.NewPostgres(*postgresCfg)
	if err != nil {
		return err
	}

	targetRepo, err := target.New(pgDb)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisCfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := configs.NewRedisClient(redisCfg)
	if err != nil {
		log.Fatal(err)
	}
	cache := cache.NewCache(redisClient)
	if cache == nil {
		log.Fatal("Nil cache")
	}

	var group *singleflight.Group

	logger, loggerErr := logger.New()
	if loggerErr != nil {
		return loggerErr
	}
	targetService := serviceTarget.New(targetRepo, cache, group, logger)
	targetHandler := url.New(targetService)

	checkRepo, errCheck := checks.New(pgDb)
	if errCheck != nil {
		return errCheck
	}

	// ===== checker pipeline =====
	queue := make(chan dto.Targets, 100)
	results := make(chan dto.Checks, 100)

	httpChecker := checker.NewHTTPChecker(logger)
	writer := writer.NewWriter(checkRepo, results, logger)
	workers := workerpool.New(httpChecker, results, 50)
	scheduler := scheduler.New(queue, targetRepo)

	go writer.Run(ctx)
	workers.Start(ctx, queue)
	go scheduler.Run(ctx)

	checkService := serviceChecker.NewCheckService(checkRepo, cache, logger)
	checkHandler := checksHandler.NewCheckHandler(checkService)

	router := gin.Default()
	url.RegisterRoutes(router, targetHandler)
	checksHandler.RegisterRoutes(router, checkHandler)

	return router.RunTLS(":8080", "cert.pem", "key.pem")
}
