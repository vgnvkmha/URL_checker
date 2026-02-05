package url

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	r *gin.Engine,
	h IURLHandler,
) {
	r.GET("/targets", h.List)
	r.POST("/targets", h.Create)
	r.PATCH("/targets/:id", h.Update)
	r.DELETE("/targets/:id", h.Delete)
}
