package checksHandler

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	r *gin.Engine,
	h ICheckHandler,
) {
	r.POST("/targets/checks", h.Insert) //Это будет вызывается внутренне, когда сервис проверки будет делать запрос и сохранять данные
	r.GET("/targets/:id/status", h.LatestByTarget)
	r.GET("/targets/:id/checks", h.ListByTarget)
}
