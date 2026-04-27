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

	"omniport-api/internal/config"
	"omniport-api/internal/database"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"omniport-api/internal/modules/administration/access"
	"omniport-api/internal/modules/administration/auth"
	"omniport-api/internal/modules/administration/branch"
	"omniport-api/internal/modules/administration/cargo"
	"omniport-api/internal/modules/administration/company"
	"omniport-api/internal/modules/administration/dermaga"
	"omniport-api/internal/modules/administration/equipment"
	"omniport-api/internal/modules/administration/lookup"
	"omniport-api/internal/modules/administration/menu"
	"omniport-api/internal/modules/administration/reference"
	"omniport-api/internal/modules/administration/role"
	"omniport-api/internal/modules/administration/tariff"
	"omniport-api/internal/modules/administration/terminal"
	"omniport-api/internal/modules/administration/user"
	"omniport-api/internal/modules/administration/vessel"
	"omniport-api/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	setupLogger(cfg)

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	jwtUtil := helper.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	userRepo := user.NewUserRepository(db)
	roleRepo := role.NewRoleRepository(db)
	accessRepo := access.NewAccessRepository(db)
	dermagaRepo := dermaga.NewDermagaRepository(db)
	referenceRepo := reference.NewReferenceRepository(db)
	vesselRepo := vessel.NewVesselRepository(db)
	cargoRepo := cargo.NewCargoRepository(db)
	branchRepo := branch.NewBranchRepository(db)
	terminalRepo := terminal.NewTerminalRepository(db)
	companyRepo := company.NewCompanyRepository(db)

	accessService := access.NewAccessService(accessRepo)
	authService := auth.NewAuthService(userRepo, jwtUtil)
	userService := user.NewUserService(userRepo)
	menuService := menu.NewMenuService(db)
	roleService := role.NewRoleService(roleRepo)
	dermagaService := dermaga.NewDermagaService(dermagaRepo)
	referenceService := reference.NewReferenceService(referenceRepo)
	vesselService := vessel.NewVesselService(vesselRepo)
	cargoService := cargo.NewCargoService(cargoRepo)
	branchService := branch.NewBranchService(branchRepo)
	terminalService := terminal.NewTerminalService(terminalRepo, branchRepo)
	companyService := company.NewCompanyService(companyRepo)
	tariffService := tariff.NewTariffService(db)
	equipmentService := equipment.NewEquipmentService(db)
	lookupService := lookup.NewLookupService(db, equipmentService)

	authHandler := auth.NewAuthHandler(authService)
	userHandler := user.NewUserHandler(userService)
	menuHandler := menu.NewMenuHandler(menuService)
	roleHandler := role.NewRoleHandler(roleService)
	accessHandler := access.NewAccessHandler(accessService)
	dermagaHandler := dermaga.NewDermagaHandler(dermagaService)
	referenceHandler := reference.NewReferenceHandler(referenceService)
	vesselHandler := vessel.NewVesselHandler(vesselService)
	cargoHandler := cargo.NewCargoHandler(cargoService)

	// Break circular dependency using adapter
	userAdapter := &userProviderAdapter{s: userService}
	branchHandler := branch.NewBranchHandler(branchService, userAdapter)
	terminalHandler := terminal.NewTerminalHandler(terminalService, userAdapter)
	companyHandler := company.NewCompanyHandler(companyService)
	tariffHandler := tariff.NewTariffHandler(tariffService)
	lookupHandler := lookup.NewLookupHandler(lookupService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	router.SetupRouter(&router.RouterConfig{
		Engine:           r,
		JWTUtil:          jwtUtil,
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		MenuHandler:      menuHandler,
		RoleHandler:      roleHandler,
		AccessHandler:    accessHandler,
		DermagaHandler:   dermagaHandler,
		LookupHandler:    lookupHandler,
		ReferenceHandler: referenceHandler,
		TariffHandler:    tariffHandler,
		VesselHandler:    vesselHandler,
		CargoHandler:     cargoHandler,
		BranchHandler:    branchHandler,
		TerminalHandler:  terminalHandler,
		CompanyHandler:   companyHandler,
	})

	serve(cfg, "adm-service", cfg.App.PortFor("ADM"), r)
}

func serve(cfg *config.Config, service string, port string, handler http.Handler) {
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{Addr: addr, Handler: handler}

	go func() {
		slog.Info("Server running", "service", service, "port", port, "env", cfg.App.Env, "mode", cfg.App.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...", "service", service)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "service", service, "error", err)
	}

	slog.Info("Server exiting", "service", service)
}

func setupLogger(cfg *config.Config) {
	logLevel := parseLogLevel(cfg.App.LogLevel)
	var logHandler slog.Handler
	if cfg.App.Env == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	}
	slog.SetDefault(slog.New(logHandler))
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

// userProviderAdapter helps breaking circular dependency between user and branch/terminal
type userProviderAdapter struct {
	s user.UserService
}

func (a *userProviderAdapter) GetProfile(ctx context.Context, userID uint64) (any, error) {
	return a.s.GetProfile(ctx, userID)
}
