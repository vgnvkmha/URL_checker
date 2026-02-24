package app

import (
	"URL_checker/internal/checker"
	"URL_checker/internal/configs"
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
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

	targetService := serviceTarget.New(targetRepo)
	targetHandler := url.New(targetService)

	checkRepo, errCheck := checks.New(pgDb)
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
	redisCfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := NewRedis(redisCfg)
	if err != nil {
		log.Fatal(err)
	}

	checkService := serviceChecker.NewCheckService(checkRepo, redisClient)
	checkHandler := checksHandler.NewCheckHandler(checkService)

	router := gin.Default()
	url.RegisterRoutes(router, targetHandler)
	checksHandler.RegisterRoutes(router, checkHandler)

	return router.Run(":" + postgresCfg.Port)
}

func NewRedis(cfg *configs.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
