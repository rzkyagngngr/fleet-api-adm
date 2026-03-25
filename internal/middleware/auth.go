package middleware

import (
	"net/http"
	"strings"

	"gin-boilerplate/pkg/utils"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"
const UserEmailKey = "user_email"
const EmployeeIDKey = "employee_id"
const FullNameKey = "full_name"
const BranchCodeKey = "branch_code"

func AuthMiddleware(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid authorization format, use: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := jwtUtil.ValidateToken(parts[1])
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(EmployeeIDKey, claims.EmployeeID)
		c.Set(FullNameKey, claims.FullName)
		c.Set(BranchCodeKey, claims.BranchCode)
		c.Next()
	}
}
