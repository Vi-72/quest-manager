package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"quest-manager/cmd"
	authclient "quest-manager/internal/adapters/out/client/auth"
)

func main() {
	_ = godotenv.Load(".env")
	configs := getConfigs()

	// Database setup
	connectionString, err := cmd.MakeConnectionString(
		configs.DbHost,
		configs.DbPort,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd.CreateDbIfNotExists(configs.DbHost,
		configs.DbPort,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbSslMode)
	gormDb := cmd.MustGormOpen(connectionString)
	cmd.MustAutoMigrate(gormDb)

	// Create container
	container, err := cmd.NewContainer(configs, gormDb)
	if err != nil {
		log.Fatalf("failed to create container: %v", err)
	}

	// Build and validate container
	ctx := context.Background()
	if err := container.Build(ctx); err != nil {
		log.Fatalf("failed to build container: %v", err)
	}

	defer container.CloseAll()

	// Create router
	router := cmd.NewRouter(container)

	// Start server
	log.Printf("🚀 Server starting on :%s", configs.HttpPort)

	err = http.ListenAndServe(":"+configs.HttpPort, router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getConfigs() cmd.Config {
	return cmd.Config{
		HttpPort:            getEnv("HTTP_PORT"),
		DbHost:              getEnv("DB_HOST"),
		DbPort:              getEnv("DB_PORT"),
		DbUser:              getEnv("DB_USER"),
		DbPassword:          getEnv("DB_PASSWORD"),
		DbName:              getEnv("DB_NAME"),
		DbSslMode:           getEnv("DB_SSLMODE"),
		EventGoroutineLimit: getEnvInt("EVENT_GOROUTINE_LIMIT"),

		// Auth factory
		AuthFactory: &authclient.Factory{
			Addr: getEnv("AUTH_GRPC"),
		},

		// Middleware configuration
		Middleware: cmd.MiddlewareConfig{
			EnableAuth: getEnvBool("ENABLE_AUTH_MIDDLEWARE", true),
			DevAuth: cmd.DevAuthConfig{
				HeaderName:   getEnvWithDefault("DEV_AUTH_HEADER_NAME", cmd.DefaultDevAuthHeaderName),
				StaticUserID: getEnvWithDefault("DEV_AUTH_STATIC_USER_ID", cmd.DefaultDevAuthStaticUserID),
			},
		},
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing env var: %s", key)
	}
	return val
}

func getEnvWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func getEnvInt(key string) int {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing env var: %s", key)
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Invalid integer value for env var %s: %s", key, val)
	}
	return intVal
}

func getEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf("Invalid boolean value for env var %s: %s, using default: %v", key, val, defaultValue)
		return defaultValue
	}
	return boolVal
}
