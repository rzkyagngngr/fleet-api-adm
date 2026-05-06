package helper

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// IdentityContext holds the full operational identity of a user request.
type IdentityContext struct {
	UserID       uint64
	UserEmail    string
	UserFullName string
	EmployeeID   string
	BranchCode   string
	BranchName   string
	TerminalCode string
	TerminalName string
	CompanyCode  string
	CompanyName  string
}

// GetBranchCodeInt returns the branch code as an integer (pointer for GORM nullable compatibility)
func (i IdentityContext) GetBranchCodeInt() *int {
	if i.BranchCode == "" {
		return nil
	}
	val, _ := strconv.Atoi(i.BranchCode)
	return &val
}

// GetTerminalCodeInt returns the terminal code as an integer (pointer for GORM nullable compatibility)
func (i IdentityContext) GetTerminalCodeInt() *int {
	if i.TerminalCode == "" {
		return nil
	}
	val, _ := strconv.Atoi(i.TerminalCode)
	return &val
}

// ExtractIdentity captures the identity context from a Gin request.
func ExtractIdentity(c *gin.Context) IdentityContext {
	// These keys are set by the AuthMiddleware
	return IdentityContext{
		UserID:       getUint64(c, "user_id"),
		UserEmail:    getString(c, "user_email"),
		UserFullName: getString(c, "full_name"),
		EmployeeID:   getString(c, "employee_id"),
		BranchCode:   getString(c, "branch_code"),
		BranchName:   getString(c, "branch_name"),
		TerminalCode: getString(c, "terminal_code"),
		TerminalName: getString(c, "terminal_name"),
		CompanyCode:  getString(c, "company_code"),
		CompanyName:  getString(c, "company_name"),
	}
}

func getString(c *gin.Context, key string) string {
	if val, ok := c.Get(key); ok {
		return val.(string)
	}
	return ""
}

func getUint64(c *gin.Context, key string) uint64 {
	if val, ok := c.Get(key); ok {
		return val.(uint64)
	}
	return 0
}

// IdentityAware defines an entity that needs its multi-tenant fields populated automatically.
type IdentityAware interface {
	SetIdentity(id IdentityContext)
}
