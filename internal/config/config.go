package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Port     string
	Env      string
	LogLevel string
	Mode     string
	Ports    map[string]string
}

type ModuleDBConfig struct {
	User     string
	Password string
	Schema   string
}

type DatabaseConfig struct {
	Host    string
	Port    string
	Name    string
	SSLMode string

	Adm  ModuleDBConfig
	Plan ModuleDBConfig
	Rls  ModuleDBConfig
	Bill ModuleDBConfig
	Mrep ModuleDBConfig
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

func (d DatabaseConfig) DSN(module ModuleDBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta search_path=%s",
		d.Host, d.Port, module.User, module.Password, d.Name, d.SSLMode, module.Schema,
	)
}

func Load() (*Config, error) {
	cwd, _ := os.Getwd()
	fmt.Printf("INFO: Current working directory: %s\n", cwd)

	if err := godotenv.Load(); err != nil {
		fmt.Println("INFO: No .env file loaded from CWD, relying on system environment variables")
	}

	expiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		expiryHours = 24
	}

	if getEnv("DB_ADM_USER", "") == "" {
		fmt.Println("CRITICAL: DB_ADM_USER is missing from environment")
	}

	cfg := &Config{
		App: AppConfig{
			Port:     getEnv("APP_PORT", "8080"),
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("APP_LOG_LEVEL", "INFO"),
			Mode:     getEnv("APP_MODE", "monolith"),
			Ports: map[string]string{
				"ADM":  getEnv("APP_PORT_ADM", "8081"),
				"PLAN": getEnv("APP_PORT_PLAN", "8082"),
				"RLS":  getEnv("APP_PORT_RLS", "8083"),
				"BILL": getEnv("APP_PORT_BILL", "8084"),
				"MREP": getEnv("APP_PORT_MREP", "8085"),
			},
		},
		Database: DatabaseConfig{
			Host:    getEnv("DB_HOST", "localhost"),
			Port:    getEnv("DB_PORT", "5432"),
			Name:    getEnv("DB_NAME", "omniport"),
			SSLMode: getEnv("DB_SSLMODE", "disable"),
			Adm: ModuleDBConfig{
				User:     getEnv("DB_ADM_USER", ""),
				Password: getEnv("DB_ADM_PASSWORD", ""),
				Schema:   getEnv("DB_ADM_SCHEMA", "adm"),
			},
			Plan: ModuleDBConfig{
				User:     getEnv("DB_PLAN_USER", ""),
				Password: getEnv("DB_PLAN_PASSWORD", ""),
				Schema:   getEnv("DB_PLAN_SCHEMA", "plan"),
			},
			Rls: ModuleDBConfig{
				User:     getEnv("DB_RLS_USER", ""),
				Password: getEnv("DB_RLS_PASSWORD", ""),
				Schema:   getEnv("DB_RLS_SCHEMA", "rls"),
			},
			Bill: ModuleDBConfig{
				User:     getEnv("DB_BILL_USER", ""),
				Password: getEnv("DB_BILL_PASSWORD", ""),
				Schema:   getEnv("DB_BILL_SCHEMA", "bill"),
			},
			Mrep: ModuleDBConfig{
				User:     getEnv("DB_MREP_USER", ""),
				Password: getEnv("DB_MREP_PASSWORD", ""),
				Schema:   getEnv("DB_MREP_SCHEMA", "mrep"),
			},
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "secret"),
			ExpiryHours: expiryHours,
		},
	}

	return cfg, nil
}

func (a AppConfig) PortFor(service string) string {
	if a.Mode == "monolith" {
		return a.Port
	}
	if p, ok := a.Ports[service]; ok && p != "" {
		return p
	}
	return a.Port
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
