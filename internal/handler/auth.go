package handler

import (
	"github.com/gin-gonic/gin"
	"my_app/internal/model"
	"net/http"
	"strings"
)

func (h *Handler) signUp(c *gin.Context) {
	var userInput model.CreateUserInput
	if err := c.BindJSON(&userInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, invalidRequest)
		return
	}

	newUserId, err := h.service.Authorization.CreateUser(userInput)
	if err != nil {
		// TODO: дописать логику обработки ошибки "already exists"
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newCreatedResponse(c, newUserId)
	return
}

func (h *Handler) signIn(c *gin.Context) {
	var signInUserInput model.SignInUserInput

	if err := c.BindJSON(&signInUserInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, invalidRequest)
		return
	}

	userId, err := h.service.Authorization.UserAuthorize(signInUserInput)
	if err != nil {
		// TODO: дописать логику обработки ошибки "doesn't exists or wrong password"
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	jwtPair, err := h.service.Authorization.GenerateJwtTokenPair(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, jwtPair)
	return
}

func (h *Handler) authenticationMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	userId, err := h.service.Authorization.GetUserIdFromAuthHeader(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	c.Set("userID", userId)
	c.Next()
}
