package database

import (
	"fmt"

	"omniport-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}
	if cfg.App.Env == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	if cfg.Database.Adm.User == "" {
		return nil, fmt.Errorf("DB_ADM_USER is required")
	}

	return gorm.Open(postgres.Open(cfg.Database.DSN(cfg.Database.Adm)), gormConfig)
}
