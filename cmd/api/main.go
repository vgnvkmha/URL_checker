package main

import (
	"URL_checker/internal/handler/url"
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/queries"
	"URL_checker/internal/service"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "postgres://pavelpavlov@localhost:5432/postgres?sslmode=disable"
	repo, err := queries.NewPostgresRepo(dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	service := service.New(repo)
	handler := url.New(service)
	router := gin.Default()

	router.GET("/targets", func(c *gin.Context) {
		targets, err := handler.List(c.Request.Context())
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"current_targets": targets})
	})

	router.GET("/targets/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "id should be int"})
			return
		}
		target, err1 := handler.Get(c.Request.Context(), uint64(id))
		if err1 != nil {
			c.JSON(400, gin.H{"error": err1.Error()})
			return
		}
		c.JSON(200, gin.H{"target": target})
	})

	router.POST("/targets", func(c *gin.Context) {
		var params entities.Targets
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		target, err1 := handler.Create(c.Request.Context(), params)
		if err1 != nil {
			c.JSON(400, gin.H{"error": err1.Error()})
			return
		}
		c.JSON(200, gin.H{"created": target}) //TODO: посмотреть, на что ругается
	})

	router.PATCH("/targets/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "id should be int"})
			return
		}

		var req entities.PatchReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err1 := handler.Update(c.Request.Context(), uint64(id), req)
		if err1 != nil {
			c.JSON(400, gin.H{"error": err1.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "updated"})
	})

	router.DELETE("/targets/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "id should be int"})
			return
		}
		err1 := handler.Delete(c.Request.Context(), uint64(id))
		if err1 != nil {
			c.JSON(400, gin.H{"error": err1.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "deleted"})
	})

	router.Run() // listens on 0.0.0.0:8080 by default
}
