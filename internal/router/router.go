package router

import (
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"omniport-api/internal/modules/administration/access"
	"omniport-api/internal/modules/administration/auth"
	"omniport-api/internal/modules/administration/dermaga"
	"omniport-api/internal/modules/administration/menu"
	"omniport-api/internal/modules/administration/reference"
	"omniport-api/internal/modules/administration/role"
	"omniport-api/internal/modules/administration/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	Engine           *gin.Engine
	JWTUtil          *helper.JWTUtil
	AuthHandler      *auth.AuthHandler
	UserHandler      *user.UserHandler
	MenuHandler      *menu.MenuHandler
	RoleHandler      *role.RoleHandler
	AccessHandler    *access.AccessHandler
	DermagaHandler   dermaga.DermagaHandler
	ReferenceHandler reference.ReferenceHandler
}

func SetupRouter(cfg *RouterConfig) {
	r := cfg.Engine

	r.Use(middleware.TraceMiddleware())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "omniport-api",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(cfg.JWTUtil))
	{
		v1.POST("/auth/register", cfg.AuthHandler.Register)
		v1.POST("/auth/login", cfg.AuthHandler.Login)
		v1.POST("/auth/change-terminal", cfg.AuthHandler.ChangeTerminal)

		v1.GET("/users/profile", cfg.UserHandler.GetProfile)

		master := v1.Group("/master")
		{
			menus := master.Group("/menus")
			{
				menus.GET("", cfg.MenuHandler.GetAllMenus)
				menus.POST("/search", cfg.MenuHandler.SearchMenus)
				menus.POST("", cfg.MenuHandler.CreateMenu)
				menus.PUT("/:id", cfg.MenuHandler.UpdateMenu)
				menus.DELETE("/:id", cfg.MenuHandler.DeleteMenu)
			}

			users := master.Group("/users")
			{
				users.GET("", cfg.UserHandler.FindAll)
				users.POST("/search", cfg.UserHandler.Search)
				users.GET("/:id", cfg.UserHandler.FindByID)
				users.POST("", cfg.UserHandler.Create)
				users.PUT("/:id", cfg.UserHandler.Update)
				users.DELETE("/:id", cfg.UserHandler.Delete)
			}

			roles := master.Group("/roles")
			{
				roles.GET("", cfg.RoleHandler.GetAllRoles)
				roles.POST("", cfg.RoleHandler.CreateRole)
				roles.PUT("/:id", cfg.RoleHandler.UpdateRole)
				roles.DELETE("/:id", cfg.RoleHandler.DeleteRole)
				roles.GET("/:id/access", cfg.AccessHandler.GetRoleAccess)
				roles.GET("/:id/all-menu-access", cfg.AccessHandler.GetAllMenuByRole)
				roles.POST("/:id/access", cfg.AccessHandler.UpdateRoleAccess)
			}

			references := master.Group("/references")
			{
				references.GET("", cfg.ReferenceHandler.GetAllReferences)
				references.GET("/:id", cfg.ReferenceHandler.GetReferenceDetail)
				references.POST("", cfg.ReferenceHandler.SaveReference)
				references.DELETE("/:id", cfg.ReferenceHandler.DeleteReference)
			}
		}

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
