package database

import (
	"fmt"
	"time"

	"omniport-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Registry struct {
	ADM  *gorm.DB
	PLAN *gorm.DB
	RLS  *gorm.DB
	BILL *gorm.DB
	MREP *gorm.DB
}

func NewRegistry(cfg *config.Config) (*Registry, error) {
	gormConfig := &gorm.Config{}
	if cfg.App.Env == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	reg := &Registry{}
	var err error

	reg.ADM, err = openIfConfigured(cfg, cfg.Database.Adm, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("connect adm db: %w", err)
	}
	reg.PLAN, err = openIfConfigured(cfg, cfg.Database.Plan, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("connect plan db: %w", err)
	}
	reg.RLS, err = openIfConfigured(cfg, cfg.Database.Rls, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("connect rls db: %w", err)
	}
	reg.BILL, err = openIfConfigured(cfg, cfg.Database.Bill, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("connect bill db: %w", err)
	}
	reg.MREP, err = openIfConfigured(cfg, cfg.Database.Mrep, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("connect mrep db: %w", err)
	}

	return reg, nil
}

func openIfConfigured(cfg *config.Config, moduleCfg config.ModuleDBConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	if moduleCfg.User == "" {
		return nil, nil
	}
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN(moduleCfg)), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLife) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Database.MaxIdleTime) * time.Minute)

	return db, nil
}
