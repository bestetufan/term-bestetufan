package controller

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Message string `json:"message" example:"message"`
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}

func successResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, response{msg})
}
