package app

import (
	"URL_checker/internal/handler/url"
	"URL_checker/internal/repo/target"
	"URL_checker/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Run() error {
	dsn := "postgres://pavelpavlov@localhost:5432/postgres?sslmode=disable"

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

	router := gin.Default()
	url.RegisterRoutes(router, targetHandler)

	return router.Run(":8080")
}
