package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespSuccess(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    message,
	})
}

func RespFail(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": 500,
		"msg":    message,
	})
}
