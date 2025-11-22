package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	health := h.healthService.CheckHealth(c.Request.Context())

	statusCode := http.StatusOK
	if health.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

func (h *Handler) ReadinessProbe(c *gin.Context) {
	dbHealth := h.healthService.CheckDatabase(c.Request.Context())

	if dbHealth.Status == "healthy" {
		helpers.SuccessResponse(c, gin.H{"ready": true}, "Application is ready")
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":  false,
			"reason": "Database is not ready",
		})
	}
}

func (h *Handler) LivenessProbe(c *gin.Context) {
	helpers.SuccessResponse(c, gin.H{"alive": true}, "Application is alive")
}
