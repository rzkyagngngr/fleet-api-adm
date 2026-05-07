package router

import (
	_ "omniport-api/docs"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"omniport-api/internal/modules/administration/access"
	"omniport-api/internal/modules/administration/auth"
	"omniport-api/internal/modules/administration/branch"
	"omniport-api/internal/modules/administration/cargo"
	"omniport-api/internal/modules/administration/company"
	"omniport-api/internal/modules/administration/customer"
	"omniport-api/internal/modules/administration/dermaga"
	"omniport-api/internal/modules/administration/dock"
	"omniport-api/internal/modules/administration/equipment"
	"omniport-api/internal/modules/administration/file"
	"omniport-api/internal/modules/administration/lookup"
	"omniport-api/internal/modules/administration/menu"
	"omniport-api/internal/modules/administration/pelabuhan"
	"omniport-api/internal/modules/administration/reference"
	"omniport-api/internal/modules/administration/role"
	"omniport-api/internal/modules/administration/tariff"
	"omniport-api/internal/modules/administration/terminal"
	"omniport-api/internal/modules/administration/user"
	"omniport-api/internal/modules/administration/vessel"
	"omniport-api/internal/modules/administration/warehouse"
	"omniport-api/internal/modules/plan/op"
	"omniport-api/internal/modules/plan/postrequest"
	"omniport-api/internal/modules/plan/vesselschedule"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	Engine                *gin.Engine
	JWTUtil               *helper.JWTUtil
	AuthHandler           *auth.AuthHandler
	UserHandler           *user.UserHandler
	MenuHandler           *menu.MenuHandler
	RoleHandler           *role.RoleHandler
	AccessHandler         *access.AccessHandler
	DermagaHandler        dermaga.DermagaHandler
	CustomerHandler       *customer.CustomerHandler
	DockHandler           *dock.DockHandler
	EquipmentHandler      *equipment.EquipmentHandler
	LookupHandler         *lookup.LookupHandler
	PortHandler           *pelabuhan.PortHandler
	ReferenceHandler      reference.ReferenceHandler
	TariffHandler         *tariff.TariffHandler
	VesselHandler         *vessel.VesselHandler
	VesselScheduleHandler *vesselschedule.VesselScheduleHandler
	CargoHandler          *cargo.CargoHandler
	BranchHandler         *branch.BranchHandler
	TerminalHandler       *terminal.TerminalHandler
	WarehouseHandler      *warehouse.WarehouseHandler
	CompanyHandler        *company.CompanyHandler
	PostRequestHandler    *postrequest.PostRequestHandler
	FileHandler           *file.FileHandler
	OpsPlanHandler        *op.OpsPlanHandler
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
		v1.POST("/auth/refresh-token", cfg.AuthHandler.RefreshToken)
		v1.GET("/auth/me", cfg.AuthHandler.Me)

		v1.GET("/users/profile", cfg.UserHandler.GetProfile)
		v1.GET("/users/me/locations", cfg.UserHandler.GetMyLocations)

		master := v1.Group("/master")
		{
			menus := master.Group("/menus")
			{
				menus.GET("", func(c *gin.Context) {
					if c.Query("id") != "" {
						cfg.MenuHandler.GetMenuDetail(c)
						return
					}
					cfg.MenuHandler.GetAllMenus(c)
				})
				menus.POST("/search", cfg.MenuHandler.SearchMenus)
				menus.POST("", cfg.MenuHandler.CreateMenu)
				menus.PUT("", cfg.MenuHandler.UpdateMenu)
				menus.DELETE("", cfg.MenuHandler.DeleteMenu)
			}

			users := master.Group("/users")
			{
				users.GET("/stats", cfg.UserHandler.GetStats)
				users.GET("", func(c *gin.Context) {
					if c.Query("id") != "" {
						cfg.UserHandler.FindByID(c)
						return
					}
					cfg.UserHandler.FindAll(c)
				})
				users.POST("/search", cfg.UserHandler.Search)
				users.POST("", cfg.UserHandler.Create)
				users.PUT("", cfg.UserHandler.Update)
				users.DELETE("", cfg.UserHandler.Delete)
			}

			roles := master.Group("/roles")
			{
				roles.GET("", func(c *gin.Context) {
					if c.Query("id") != "" {
						cfg.RoleHandler.GetRoleDetail(c)
						return
					}
					cfg.RoleHandler.GetAllRoles(c)
				})
				roles.POST("", cfg.RoleHandler.CreateRole)
				roles.PUT("", cfg.RoleHandler.UpdateRole)
				roles.DELETE("", cfg.RoleHandler.DeleteRole)
				roles.GET("/:id/access", cfg.AccessHandler.GetRoleAccess)
				roles.GET("/:id/all-menu-access", cfg.AccessHandler.GetAllMenuByRole)
				roles.POST("/:id/access", cfg.AccessHandler.UpdateRoleAccess)
			}

			references := master.Group("/references")
			{
				references.GET("", func(c *gin.Context) {
					if c.Query("id") != "" {
						cfg.ReferenceHandler.GetReferenceDetail(c)
						return
					}
					cfg.ReferenceHandler.GetAllReferences(c)
				})
				references.POST("", cfg.ReferenceHandler.SaveReference)
				references.DELETE("", cfg.ReferenceHandler.DeleteReference)
			}

			tariff := master.Group("/tariff")
			{
				tariff.POST("/search", cfg.TariffHandler.Search)
				tariff.POST("/status-zero/search", cfg.TariffHandler.SearchStatusZero)
				tariff.GET("", cfg.TariffHandler.GetByID)
				tariff.POST("", cfg.TariffHandler.Create)
				tariff.PUT("", cfg.TariffHandler.Update)
				tariff.PUT("/:id/status", cfg.TariffHandler.UpdateStatus)
				tariff.DELETE("", cfg.TariffHandler.Delete)
			}

			pelabuhan := master.Group("/pelabuhan")
			{
				pelabuhan.POST("/search", cfg.PortHandler.SearchPorts)
				pelabuhan.GET("", cfg.PortHandler.GetPortDetail)
				pelabuhan.POST("", cfg.PortHandler.CreatePort)
				pelabuhan.PUT("", cfg.PortHandler.UpdatePort)
				pelabuhan.DELETE("", cfg.PortHandler.DeletePort)
			}

			equipment := master.Group("/equipment")
			{
				equipment.POST("/group-options/search", cfg.EquipmentHandler.ListEquipmentGroupOptions)
				equipment.POST("/customer-options/search", cfg.EquipmentHandler.ListCustomerOptions)
				equipment.POST("/search", cfg.EquipmentHandler.SearchEquipments)
				equipment.GET("", cfg.EquipmentHandler.GetEquipmentDetail)
				equipment.POST("", cfg.EquipmentHandler.CreateEquipment)
				equipment.PUT("", cfg.EquipmentHandler.UpdateEquipment)
				equipment.DELETE("", cfg.EquipmentHandler.DeleteEquipment)
			}

			lookup := master.Group("/lookup")
			{
				lookup.POST("/equipment-groups/search", cfg.LookupHandler.ListEquipmentGroupOptions)
				lookup.POST("/customers/search", cfg.LookupHandler.ListCustomerOptions)
				lookup.POST("/equipments/search", cfg.LookupHandler.ListEquipmentOptions)
				lookup.POST("/cargos/search", cfg.LookupHandler.ListCargoOptions)
				lookup.POST("/cargo-packages/search", cfg.LookupHandler.ListCargoPackageOptions)
				lookup.POST("/cargo-units/search", cfg.LookupHandler.ListCargoUnitOptions)
				lookup.POST("/billing-services/search", cfg.LookupHandler.ListBillingServiceOptions)
				lookup.POST("/docks/search", cfg.LookupHandler.ListDockOptions)
				lookup.POST("/vessels/search", cfg.LookupHandler.ListVesselOptions)
				lookup.POST("/ports/search", cfg.LookupHandler.ListPortOptions)
			}

			dock := master.Group("/dock")
			{
				dock.POST("/search", cfg.DockHandler.SearchDock)
				dock.GET("", cfg.DockHandler.GetDockDetail)
				dock.POST("", cfg.DockHandler.CreateDock)
				dock.PUT("", cfg.DockHandler.UpdateDock)
				dock.DELETE("", cfg.DockHandler.DeleteDock)
			}

			customer := master.Group("/customer")
			{
				customer.POST("/search", cfg.CustomerHandler.SearchCustomers)
				customer.POST("", cfg.CustomerHandler.CreateCustomer)
				customer.GET("", cfg.CustomerHandler.GetCustomerDetail)
				customer.PUT("", cfg.CustomerHandler.UpdateCustomer)
				customer.DELETE("", cfg.CustomerHandler.DeleteCustomer)
			}

			barang := master.Group("/barang")
			{
				barang.GET("/stats", cfg.CargoHandler.GetStats)
				barang.POST("/search", cfg.CargoHandler.Search)
				barang.POST("", cfg.CargoHandler.Create)
				barang.GET("", cfg.CargoHandler.GetByID)
				barang.PUT("", cfg.CargoHandler.Update)
				barang.DELETE("", cfg.CargoHandler.Delete)
			}

			warehouse := master.Group("/warehouse")
			{
				warehouse.POST("/search", cfg.WarehouseHandler.SearchWarehouse)
				warehouse.GET("", cfg.WarehouseHandler.GetWarehouseDetail)
				warehouse.POST("", cfg.WarehouseHandler.CreateWarehouse)
				warehouse.PUT("", cfg.WarehouseHandler.UpdateWarehouse)
				warehouse.DELETE("", cfg.WarehouseHandler.DeleteWarehouse)
			}

			vessel := master.Group("/vessel")
			{
				vessel.GET("/stats", cfg.VesselHandler.GetStats)
				vessel.POST("/search", cfg.VesselHandler.Search)
				vessel.POST("", cfg.VesselHandler.Create)
				vessel.GET("", cfg.VesselHandler.GetByID)
				vessel.GET("/:id", cfg.VesselHandler.GetByID)
				vessel.PUT("", cfg.VesselHandler.Update)
				vessel.PUT("/:id", cfg.VesselHandler.Update)
				vessel.DELETE("", cfg.VesselHandler.Delete)
				vessel.DELETE("/:id", cfg.VesselHandler.Delete)
			}

			branches := master.Group("/branches")
			{
				branches.GET("/stats", cfg.BranchHandler.GetStats)
				branches.POST("/search", cfg.BranchHandler.Search)
				branches.POST("", cfg.BranchHandler.Create)
				branches.GET("", cfg.BranchHandler.GetByID)
				branches.PUT("", cfg.BranchHandler.Update)
				branches.DELETE("", cfg.BranchHandler.Delete)
			}

			companies := master.Group("/companies")
			{
				companies.POST("/search", cfg.CompanyHandler.Search)
				companies.POST("", cfg.CompanyHandler.Create)
				companies.GET("", cfg.CompanyHandler.GetByID)
				companies.PUT("", cfg.CompanyHandler.Update)
				companies.DELETE("", cfg.CompanyHandler.Delete)
			}

			terminals := master.Group("/terminals")
			{
				terminals.GET("/stats", cfg.TerminalHandler.GetStats)
				terminals.POST("/search", cfg.TerminalHandler.Search)
				terminals.POST("", cfg.TerminalHandler.Create)
				terminals.GET("", cfg.TerminalHandler.GetByID)
				terminals.PUT("", cfg.TerminalHandler.Update)
				terminals.DELETE("", cfg.TerminalHandler.Delete)
			}
		}

		plan := v1.Group("/plan")
		{
			// Vessel Schedule (Shared)
			vesselSchedule := plan.Group("/vessel-schedule")
			{
				vesselSchedule.POST("/search", cfg.VesselScheduleHandler.Search)
				vesselSchedule.POST("", cfg.VesselScheduleHandler.Create)
				vesselSchedule.GET("", cfg.VesselScheduleHandler.GetByScheduleCode)
				vesselSchedule.PUT("", cfg.VesselScheduleHandler.Update)
				vesselSchedule.GET("/detail", cfg.VesselScheduleHandler.GetByScheduleCode)
				vesselSchedule.PUT("/detail", cfg.VesselScheduleHandler.Update)
				vesselSchedule.GET("/:schedule_code", cfg.VesselScheduleHandler.GetByScheduleCode)
				vesselSchedule.PUT("/:schedule_code", cfg.VesselScheduleHandler.Update)
				vesselSchedule.DELETE("/:id", cfg.VesselScheduleHandler.Delete)
			}

			// Permohonan Jasa Barang (PJB)
			reqBarang := plan.Group("/request/barang")
			{
				reqBarang.GET("/stats", cfg.PostRequestHandler.GetStats)
				reqBarang.POST("/search", cfg.PostRequestHandler.Search)
				reqBarang.POST("", cfg.PostRequestHandler.Create)
				reqBarang.GET("", cfg.PostRequestHandler.GetByID)
				reqBarang.GET("/:id", cfg.PostRequestHandler.GetByID)
				reqBarang.PUT("", cfg.PostRequestHandler.Update)
				reqBarang.PUT("/:id", cfg.PostRequestHandler.Update)
				reqBarang.PUT("/status", cfg.PostRequestHandler.UpdateStatus)
				reqBarang.PUT("/:id/status", cfg.PostRequestHandler.UpdateStatus)
				reqBarang.DELETE("", cfg.PostRequestHandler.Delete)
				reqBarang.DELETE("/:id", cfg.PostRequestHandler.Delete)
			}

			opsPlan := plan.Group("/op")
			{
				opsPlan.POST("", cfg.OpsPlanHandler.Create)
				opsPlan.POST("/update", cfg.OpsPlanHandler.Update)
				opsPlan.POST("/readyOp", cfg.OpsPlanHandler.ReadyOpsPlan)
				opsPlan.POST("/getDataRequest", cfg.OpsPlanHandler.GetDataRequest)
				opsPlan.POST("/getDataOp", cfg.OpsPlanHandler.GetDataOp)
				opsPlan.POST("/getDetailOp", cfg.OpsPlanHandler.GetDetailOp)
				opsPlan.POST("/getDataVesselSchedule", cfg.OpsPlanHandler.GetDataVesselSchedule)
				opsPlan.POST("/getDataVesel", cfg.OpsPlanHandler.GetDataVesel)
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

		storage := v1.Group("/storage")
		{
			storage.POST("/upload-signature", cfg.FileHandler.GetUploadSignature)
			storage.POST("/commit/:id", cfg.FileHandler.CommitUpload)
			storage.GET("/file/:id", cfg.FileHandler.GetFileDetail)
		}
	}
}
