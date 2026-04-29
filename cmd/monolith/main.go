package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "omniport-api/docs"
	"omniport-api/internal/config"
	"omniport-api/internal/database"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"omniport-api/internal/modules/administration/access"
	"omniport-api/internal/modules/administration/auth"
	"omniport-api/internal/modules/administration/customer"
	"omniport-api/internal/modules/administration/dermaga"
	"omniport-api/internal/modules/administration/dock"
	"omniport-api/internal/modules/administration/equipment"
	"omniport-api/internal/modules/administration/lookup"
	"omniport-api/internal/modules/administration/menu"
	"omniport-api/internal/modules/administration/pelabuhan"
	"omniport-api/internal/modules/administration/reference"
	"omniport-api/internal/modules/administration/role"
	"omniport-api/internal/modules/administration/tariff"
	"omniport-api/internal/modules/administration/user"
	"omniport-api/internal/modules/administration/warehouse"
	"omniport-api/internal/modules/plan/vesselschedule"
	"omniport-api/internal/router"

	"github.com/gin-gonic/gin"
)

// @title Omniport API
// @version 1.0
// @description API for Omniport administration backend.
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logLevel := parseLogLevel(cfg.App.LogLevel)
	var logHandler slog.Handler
	if cfg.App.Env == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	}

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	dbRegistry, err := database.NewRegistry(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	if dbRegistry.ADM == nil {
		slog.Error("Administration database connection is not configured")
		os.Exit(1)
	}
	if dbRegistry.PLAN == nil {
		slog.Error("Plan database connection is not configured")
		os.Exit(1)
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	jwtUtil := helper.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	userRepo := user.NewUserRepository(dbRegistry.ADM)
	roleRepo := role.NewRoleRepository(dbRegistry.ADM)
	accessRepo := access.NewAccessRepository(dbRegistry.ADM)
	dermagaRepo := dermaga.NewDermagaRepository(dbRegistry.ADM)
	referenceRepo := reference.NewReferenceRepository(dbRegistry.ADM)

	accessService := access.NewAccessService(accessRepo)
	authService := auth.NewAuthService(userRepo, jwtUtil)
	userService := user.NewUserService(userRepo)
	menuService := menu.NewMenuService(dbRegistry.ADM)
	roleService := role.NewRoleService(roleRepo)
	dermagaService := dermaga.NewDermagaService(dermagaRepo)
	referenceService := reference.NewReferenceService(referenceRepo)
	customerService := customer.NewCustomerService(dbRegistry.ADM)
	dockService := dock.NewDockService(dbRegistry.ADM)
	equipmentService := equipment.NewEquipmentService(dbRegistry.ADM)
	portService := pelabuhan.NewPortService(dbRegistry.ADM)
	warehouseService := warehouse.NewWarehouseService(dbRegistry.ADM)
	tariffService := tariff.NewTariffService(dbRegistry.ADM)
	lookupService := lookup.NewLookupService(dbRegistry.ADM, equipmentService)
	vesselScheduleService := vesselschedule.NewVesselScheduleService(dbRegistry.PLAN, dbRegistry.ADM)

	authHandler := auth.NewAuthHandler(authService)
	userHandler := user.NewUserHandler(userService)
	menuHandler := menu.NewMenuHandler(menuService)
	roleHandler := role.NewRoleHandler(roleService)
	accessHandler := access.NewAccessHandler(accessService)
	dermagaHandler := dermaga.NewDermagaHandler(dermagaService)
	referenceHandler := reference.NewReferenceHandler(referenceService)
	customerHandler := customer.NewCustomerHandler(customerService)
	dockHandler := dock.NewDockHandler(dockService)
	equipmentHandler := equipment.NewEquipmentHandler(equipmentService)
	portHandler := pelabuhan.NewPortHandler(portService)
	warehouseHandler := warehouse.NewWarehouseHandler(warehouseService)
	tariffHandler := tariff.NewTariffHandler(tariffService)
	lookupHandler := lookup.NewLookupHandler(lookupService)
	vesselScheduleHandler := vesselschedule.NewVesselScheduleHandler(vesselScheduleService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	router.SetupRouter(&router.RouterConfig{
		Engine:                r,
		JWTUtil:               jwtUtil,
		AuthHandler:           authHandler,
		UserHandler:           userHandler,
		MenuHandler:           menuHandler,
		RoleHandler:           roleHandler,
		AccessHandler:         accessHandler,
		DermagaHandler:        dermagaHandler,
		CustomerHandler:       customerHandler,
		DockHandler:           dockHandler,
		EquipmentHandler:      equipmentHandler,
		LookupHandler:         lookupHandler,
		PortHandler:           portHandler,
		ReferenceHandler:      referenceHandler,
		TariffHandler:         tariffHandler,
		VesselScheduleHandler: vesselScheduleHandler,
		WarehouseHandler:      warehouseHandler,
	})

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		slog.Info("Server running", "port", cfg.App.Port, "env", cfg.App.Env, "mode", cfg.App.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "VERBOSE", "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
