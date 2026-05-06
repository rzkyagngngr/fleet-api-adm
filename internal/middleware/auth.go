package middleware

import (
	"net/http"
	"strings"

	"omniport-api/internal/helper"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"
const UserEmailKey = "user_email"
const EmployeeIDKey = "employee_id"
const FullNameKey = "full_name"
const BranchCodeKey = "branch_code"
const BranchNameKey = "branch_name"
const TerminalCodeKey = "terminal_code"
const TerminalNameKey = "terminal_name"
const CompanyCodeKey = "company_code"
const CompanyNameKey = "company_name"

func AuthMiddleware(jwtUtil *helper.JWTUtil) gin.HandlerFunc {
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
			helper.ErrorResponse(c, http.StatusUnauthorized, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid authorization format, use: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := jwtUtil.ValidateToken(parts[1])
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(EmployeeIDKey, claims.EmployeeID)
		c.Set(FullNameKey, claims.FullName)
		c.Set(BranchCodeKey, claims.BranchCode)
		c.Set(BranchNameKey, claims.BranchName)
		c.Set(TerminalCodeKey, claims.TerminalCode)
		c.Set(TerminalNameKey, claims.TerminalName)
		c.Set(CompanyCodeKey, claims.CompanyCode)
		c.Set(CompanyNameKey, claims.CompanyName)
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

func GetFullName(c *gin.Context) string {
	val, exists := c.Get(FullNameKey)
	if !exists {
		return ""
	}
	return val.(string)
}

func GetBranchCode(c *gin.Context) string {
	val, exists := c.Get(BranchCodeKey)
	if !exists {
		return ""
	}
	return val.(string)
}

func GetBranchName(c *gin.Context) string {
	val, exists := c.Get(BranchNameKey)
	if !exists {
		return ""
	}
	return val.(string)
}

func GetTerminalCode(c *gin.Context) string {
	val, exists := c.Get(TerminalCodeKey)
	if !exists {
		return ""
	}
	return val.(string)
}

func GetTerminalName(c *gin.Context) string {
	val, exists := c.Get(TerminalNameKey)
	if !exists {
		return ""
	}
	return val.(string)
}

func GetCompanyCode(c *gin.Context) string {
	val, exists := c.Get(CompanyCodeKey)
	if !exists {
		return ""
	}
	return val.(string)
}

func GetCompanyName(c *gin.Context) string {
	val, exists := c.Get(CompanyNameKey)
	if !exists {
		return ""
	}
	return val.(string)
}
