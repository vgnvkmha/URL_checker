package url

import (
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IURLHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	ListActive(c *gin.Context)
}

type URLHandler struct {
	service service.IURLService
}

func New(service service.IURLService) IURLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) Create(c *gin.Context) {
	var params entities.Targets

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	target, err := h.service.Create(c.Request.Context(), params)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"created": target,
	})
}

func (h *URLHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "id should be int"})
		return
	}
	target, err1 := h.service.Get(c.Request.Context(), uint64(id))
	if err1 != nil {
		c.JSON(400, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(200, gin.H{"target": target})
}

func (h *URLHandler) List(c *gin.Context) {
	targets, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"current_targets": targets})
}

func (h *URLHandler) Update(c *gin.Context) {
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

	err1 := h.service.Update(c.Request.Context(), uint64(id), req)
	if err1 != nil {
		c.JSON(400, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "updated"})

}

// TODO: Сделать реализацию
func (h *URLHandler) ListActive(c *gin.Context) {
	return
}

func (h *URLHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "id should be int"})
		return
	}
	err1 := h.service.Delete(c.Request.Context(), uint64(id))
	if err1 != nil {
		c.JSON(400, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "deleted"})
}
