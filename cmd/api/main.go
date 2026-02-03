package main

import (
	"URL_checker/internal/handler/url"
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	service := service.New()
	handler := url.New(service)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/targets", func(c *gin.Context) {
		c.JSON(200, handler.GetTargets())
	})

	router.POST("/targets", func(c *gin.Context) {
		var params entities.URL
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		if handler.PostTarget(
			params.URL,
			params.IntervalSec,
			params.TimeoutMS,
		) {
			c.JSON(200, gin.H{
				"status": "created",
			})
		} else {
			c.JSON(200, gin.H{
				"status": "not valid params",
			})
		}

	})

	router.PATCH("/targets/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}

		var req entities.PatchReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		handler.PatchTarget(req, id)
		c.JSON(200, gin.H{"status": "updated"})
	})

	router.Run() // listens on 0.0.0.0:8080 by default
}
