package router

import (
	"gin-boilerplate/internal/handler"
	"gin-boilerplate/internal/middleware"
	"gin-boilerplate/pkg/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	Engine         *gin.Engine
	JWTUtil        *utils.JWTUtil
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	MenuHandler    *handler.MenuHandler
	RoleHandler    *handler.RoleHandler
	AccessHandler  *handler.AccessHandler
	DermagaHandler handler.DermagaHandler
}

func SetupRouter(cfg *RouterConfig) {
	r := cfg.Engine

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "gin-boilerplate",
		})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(cfg.JWTUtil))
	{
		// Auth routes (Public - handled by middleware skip logic)
		v1.POST("/auth/register", cfg.AuthHandler.Register)
		v1.POST("/auth/login", cfg.AuthHandler.Login)

		// User routes
		v1.GET("/users/profile", cfg.UserHandler.GetProfile)

		// Master routes
		master := v1.Group("/master")
		{
			// Menu routes
			menus := master.Group("/menus")
			{
				menus.GET("", cfg.MenuHandler.GetAllMenus)
				menus.POST("", cfg.MenuHandler.CreateMenu)
				menus.PUT("/:id", cfg.MenuHandler.UpdateMenu)
				menus.DELETE("/:id", cfg.MenuHandler.DeleteMenu)
			}

			// Role routes
			roles := master.Group("/roles")
			{
				roles.GET("", cfg.RoleHandler.GetAllRoles)
				roles.POST("", cfg.RoleHandler.CreateRole)
				roles.PUT("/:id", cfg.RoleHandler.UpdateRole)
				roles.DELETE("/:id", cfg.RoleHandler.DeleteRole)
				roles.GET("/:id/access", cfg.AccessHandler.GetRoleAccess)
				roles.POST("/:id/access", cfg.AccessHandler.UpdateRoleAccess)
			}
		}

		// Dermaga routes
		dermagas := v1.Group("/dermaga")
		{
			dermagas.POST("", cfg.DermagaHandler.Create)
			dermagas.GET("", cfg.DermagaHandler.FindAll)
			dermagas.GET("/:id", cfg.DermagaHandler.FindByID)
			dermagas.PUT("/:id", cfg.DermagaHandler.Update)
			dermagas.DELETE("/:id", cfg.DermagaHandler.Delete)
		}
	}
}
