package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func writeError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": msg,
	})
}

func writeOK(c *gin.Context, data interface{}) {
	ret, ok := data.(gin.H)
	if !ok {
		ret = gin.H{}
		ret["data"] = data
	}
	ret["code"] = 0
	ret["message"] = "success"
	c.JSON(http.StatusOK, ret)
}
