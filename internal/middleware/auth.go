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
	// Public paths that do not require JWT authentication
	publicPaths := map[string]bool{
		"/api/v1/auth/login":    true,
		"/api/v1/auth/register": true,
		"/health":               true,
	}

	return func(c *gin.Context) {
		// Skip authentication for public paths
		if publicPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		// Skip authentication for swagger or other public top-level paths
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			c.Next()
			return
		}

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

// GetUserID retrieves the user ID from the Gin context
func GetUserID(c *gin.Context) uint64 {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}
	return val.(uint64)
}

// GetUserEmail retrieves the user email from the Gin context
func GetUserEmail(c *gin.Context) string {
	val, exists := c.Get(UserEmailKey)
	if !exists {
		return ""
	}
	return val.(string)
}
