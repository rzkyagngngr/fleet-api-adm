package middleware

import (
	"net/http"
	"strings"

	"omniport-api/internal/helper"

	"github.com/gin-gonic/gin"
)

const InternalServiceTokenHeader = "X-Internal-Service-Token"

func InternalServiceAuth(token string) gin.HandlerFunc {
	expected := strings.TrimSpace(token)

	return func(c *gin.Context) {
		if expected == "" {
			c.Next()
			return
		}

		if c.GetHeader(InternalServiceTokenHeader) != expected {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid internal service token")
			c.Abort()
			return
		}

		c.Next()
	}
}
