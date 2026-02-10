package checksHandler

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	r *gin.Engine,
	h ICheckHandler,
) {
	r.POST("/targets/checks", h.Insert)
	r.GET("/targets/:id/status", h.LatestByTarget)
	r.GET("/targets/:id/checks", h.ListByTarget)
}
