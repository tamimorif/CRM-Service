package helpers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InternalServerError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": "internal server error",
	})
}

func NotFound(c *gin.Context, model string) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"message": model + " not found",
	})
}

func UnprocessableEntity(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": "validation error",
		"errors":  err.Error(),
	})
}
