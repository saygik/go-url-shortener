package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBServer   string `json:"CLCK_DB_SERVER"`
	DBName     string `json:"CLCK_DB_NAME"`
	DBUser     string `json:"CLCK_DB_USER"`
	DBPassword string `json:"CLCK_DB_PASS"`
}
type Config struct {
	Env  string
	Port int
	DB   DBConfig
}

func Load() (cfg Config, err error) {
	err = godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.Env = getEnv("CLCK_ENV", "local")
	cfg.Port = getEnvAsInt("CLCK_PORT", 8080)
	cfg.DB.DBServer = getEnv("CLCK_DB_SERVER", "localhost")
	cfg.DB.DBName = getEnv("CLCK_DB_NAME", "local")
	cfg.DB.DBUser = getEnv("CLCK_DB_USER", "")
	cfg.DB.DBPassword = getEnv("CLCK_DB_PASS", "")
	return cfg, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// func NewConfig(path string) (config Config, err error) {
// 	viper.AutomaticEnv()
// 	viper.AddConfigPath(path)
// 	viper.SetConfigName("app")
// 	viper.SetConfigType("env")
// 	config.Port = "9090"
// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		log.Fatalf("error reading config file: %s", err)
// 	}

// 	err = viper.Unmarshal(&config)
// 	return
// }
