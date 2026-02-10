package checksHandler

import (
	"URL_checker/internal/repo/dto"
	serviceChecker "URL_checker/internal/service/check"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ICheckHandler interface {
	Insert(c *gin.Context)
	LatestByTarget(c *gin.Context)
	ListByTarget(c *gin.Context)
}

type CheckHandler struct {
	service serviceChecker.ICheckService
}

func NewCheckHandler(service serviceChecker.ICheckService) ICheckHandler {
	return &CheckHandler{
		service: service,
	}
}

func (h *CheckHandler) Insert(c *gin.Context) {
	var params dto.Checks

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	target, err := h.service.Insert(c.Request.Context(), params)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"created": target,
	})
}

func (h *CheckHandler) LatestByTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	target, err1 := h.service.LatestByTarget(c.Request.Context(), id)
	if err1 != nil {
		c.JSON(400, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(200, gin.H{"target": target})
}

// TODO: вывод немного другой
func (h *CheckHandler) ListByTarget(c *gin.Context) {

	var req dto.ChecksList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	targets, err1 := h.service.ListByTarget(
		c.Request.Context(),
		req.TargetId,
		req.Limit,
		req.From,
		req.To,
	)
	if err1 != nil {
		c.JSON(400, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(200, gin.H{"targets": targets})
}
