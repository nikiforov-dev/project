package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	invalidRequest = "invalid request body"
)

type StatusResponse struct {
	Status string `json:"status"`
}

func newStatusOkResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, StatusResponse{Status: "ok"})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{Error: message})
}

type CreatedResponse struct {
	ID int `json:"id"`
}

func newCreatedResponse(c *gin.Context, id int) {
	c.AbortWithStatusJSON(http.StatusCreated, CreatedResponse{ID: id})
}
