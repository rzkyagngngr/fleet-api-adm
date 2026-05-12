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
	Storage  StorageConfig
}

type AppConfig struct {
	Port              string
	Env               string
	LogLevel          string
	Mode              string
	Ports             map[string]string
	ReadHeaderTimeout int
	ReadTimeout       int
	WriteTimeout      int
	IdleTimeout       int
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
	MaxOpen int
	MaxIdle int
	MaxLife int
	MaxIdleTime int

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

type StorageConfig struct {
	S3Bucket   string
	S3Region   string
	S3BaseURL  string
	S3Endpoint string
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
			ReadHeaderTimeout: getEnvAsInt("APP_READ_HEADER_TIMEOUT_SEC", 5),
			ReadTimeout:       getEnvAsInt("APP_READ_TIMEOUT_SEC", 15),
			WriteTimeout:      getEnvAsInt("APP_WRITE_TIMEOUT_SEC", 30),
			IdleTimeout:       getEnvAsInt("APP_IDLE_TIMEOUT_SEC", 60),
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
			MaxOpen: getEnvAsInt("DB_MAX_OPEN_CONNS", 40),
			MaxIdle: getEnvAsInt("DB_MAX_IDLE_CONNS", 20),
			MaxLife: getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 30),
			MaxIdleTime: getEnvAsInt("DB_CONN_MAX_IDLE_TIME_MIN", 10),
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
		Storage: StorageConfig{
			S3Bucket:   getEnv("AWS_BUCKET_NAME", "omniport-assets"),
			S3Region:   getEnv("AWS_REGION", "ap-southeast-3"),
			S3BaseURL:  getEnv("AWS_S3_BASE_URL", ""),
			S3Endpoint: getEnv("AWS_S3_ENDPOINT", ""),
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

func getEnvAsInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
