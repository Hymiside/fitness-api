package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserIdNotFound = errors.New("userId not found")
	ErrParseJSON      = errors.New("error to parse json")
	ErrInvalidRequest = errors.New("error invalid request")
)

const (
	userCtx             = "userID"
	authorizationHeader = "Authorization"
)

// schoolIdentity инденцифицирует пользователя при запросах в пути /api/...
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}
	if len(headerParts[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
		return
	}
	userID, role, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Set(userCtx, map[string]interface{}{
		"userID": userID,
		"role":   role,
	})
}

func getData(c *gin.Context) (map[string]interface{}, error) {
	data, ok := c.Get(userCtx)
	if !ok {
		return nil, ErrUserIdNotFound
	}
	return data.(map[string]interface{}), nil
}