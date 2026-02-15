package targetHandler

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	r *gin.Engine,
	h IURLHandler,
) {
	r.GET("/targets", h.List)
	r.GET("/target/:id", h.Get)
	r.POST("/targets", h.Create)
	r.PATCH("/targets/:id", h.Update)
	r.DELETE("/targets/:id", h.Delete)
	r.GET("targets/active", h.ListActive)

}
