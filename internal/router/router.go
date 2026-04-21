package router

import (
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"omniport-api/internal/modules/administration/access"
	"omniport-api/internal/modules/administration/auth"
	"omniport-api/internal/modules/administration/cargo"
	"omniport-api/internal/modules/administration/customer"
	"omniport-api/internal/modules/administration/dermaga"
	"omniport-api/internal/modules/administration/menu"
	"omniport-api/internal/modules/administration/pelabuhan"
	"omniport-api/internal/modules/administration/reference"
	"omniport-api/internal/modules/administration/role"
	"omniport-api/internal/modules/administration/branch"
	"omniport-api/internal/modules/administration/terminal"
	"omniport-api/internal/modules/administration/user"
	"omniport-api/internal/modules/administration/vessel"

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
	CustomerHandler  *customer.CustomerHandler
	PortHandler      *pelabuhan.PortHandler
	ReferenceHandler reference.ReferenceHandler
	VesselHandler     *vessel.VesselHandler
	CargoHandler      *cargo.CargoHandler
	BranchHandler     *branch.BranchHandler
	TerminalHandler   *terminal.TerminalHandler
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
				users.GET("/stats", cfg.UserHandler.GetStats)
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

			pelabuhan := master.Group("/pelabuhan")
			{
				pelabuhan.POST("/search", cfg.PortHandler.SearchPorts)
				pelabuhan.POST("", cfg.PortHandler.CreatePort)
				pelabuhan.PUT("/:id", cfg.PortHandler.UpdatePort)
				pelabuhan.DELETE("/:id", cfg.PortHandler.DeletePort)
			}

			customer := master.Group("/customer")
			{
				customer.POST("/search", cfg.CustomerHandler.SearchCustomers)
				customer.POST("", cfg.CustomerHandler.CreateCustomer)
				customer.PUT("/:id", cfg.CustomerHandler.UpdateCustomer)
				customer.DELETE("/:id", cfg.CustomerHandler.DeleteCustomer)
			}

			barang := master.Group("/barang")
			{
				barang.GET("/stats", cfg.CargoHandler.GetStats)
				barang.POST("/search", cfg.CargoHandler.Search)
				barang.POST("", cfg.CargoHandler.Create)
				barang.GET("/:id", cfg.CargoHandler.GetByID)
				barang.PUT("/:id", cfg.CargoHandler.Update)
				barang.DELETE("/:id", cfg.CargoHandler.Delete)
			}

			vessel := master.Group("/vessel")
			{
				vessel.GET("/stats", cfg.VesselHandler.GetStats)
				vessel.POST("/search", cfg.VesselHandler.Search)
				vessel.POST("", cfg.VesselHandler.Create)
				vessel.GET("/:id", cfg.VesselHandler.GetByID)
				vessel.PUT("/:id", cfg.VesselHandler.Update)
				vessel.DELETE("/:id", cfg.VesselHandler.Delete)
			}

			branches := master.Group("/branches")
			{
				branches.GET("/stats", cfg.BranchHandler.GetStats)
				branches.POST("/search", cfg.BranchHandler.Search)
				branches.POST("", cfg.BranchHandler.Create)
				branches.PUT("/:id", cfg.BranchHandler.Update)
				branches.DELETE("/:id", cfg.BranchHandler.Delete)
			}

			terminals := master.Group("/terminals")
			{
				terminals.GET("/stats", cfg.TerminalHandler.GetStats)
				terminals.POST("/search", cfg.TerminalHandler.Search)
				terminals.POST("", cfg.TerminalHandler.Create)
				terminals.PUT("/:id", cfg.TerminalHandler.Update)
				terminals.DELETE("/:id", cfg.TerminalHandler.Delete)
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
