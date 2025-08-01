package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"quest-manager/cmd"
)

func main() {
	_ = godotenv.Load(".env")
	configs := getConfigs()

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

	compositionRoot := cmd.NewCompositionRoot(
		configs,
		gormDb,
	)
	defer compositionRoot.CloseAll()

	router := cmd.NewRouter(compositionRoot)

	log.Printf("Server running on :%s", configs.HttpPort)
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
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing env var: %s", key)
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
