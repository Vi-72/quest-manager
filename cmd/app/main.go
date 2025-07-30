package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"quest-manager/cmd"
	"quest-manager/internal/adapters/out/postgres/questrepo"
	"quest-manager/internal/pkg/errs"
	"quest-manager/internal/web"

	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
)

func main() {
	configs := getConfigs()

	connectionString, err := makeConnectionString(
		configs.DbHost,
		configs.DbPort,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	createDbIfNotExists(configs.DbHost,
		configs.DbPort,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbSslMode)
	gormDb := mustGormOpen(connectionString)
	mustAutoMigrate(gormDb)

	compositionRoot := cmd.NewCompositionRoot(
		configs,
		gormDb,
	)
	defer compositionRoot.CloseAll()

	router := web.NewRouter(compositionRoot)

	log.Printf("Server running on :%s", configs.HttpPort)
	err = http.ListenAndServe(":"+configs.HttpPort, router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getConfigs() cmd.Config {
	return cmd.Config{
		HttpPort:   getEnv("HTTP_PORT"),
		DbHost:     getEnv("DB_HOST"),
		DbPort:     getEnv("DB_PORT"),
		DbUser:     getEnv("DB_USER"),
		DbPassword: getEnv("DB_PASSWORD"),
		DbName:     getEnv("DB_NAME"),
		DbSslMode:  getEnv("DB_SSLMODE"),
	}
}

func getEnv(key string) string {
	_ = godotenv.Load(".env")
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing env var: %s", key)
	}
	return val
}

func createDbIfNotExists(host string, port string, user string,
	password string, dbName string, sslMode string) {
	dsn, err := makeConnectionString(host, port, user, password, "postgres", sslMode)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Ошибка закрытия соединения с БД:", err)
		}
	}(db)

	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("БД уже существует, продолжаем.")
		} else {
			log.Fatalf("Ошибка создания БД: %v", err)
		}
	}
}

func makeConnectionString(host string, port string, user string,
	password string, dbName string, sslMode string) (string, error) {
	if host == "" {
		return "", errs.NewValueIsRequiredError("host")
	}
	if port == "" {
		return "", errs.NewValueIsRequiredError("port")
	}
	if user == "" {
		return "", errs.NewValueIsRequiredError("user")
	}
	if password == "" {
		return "", errs.NewValueIsRequiredError("password")
	}
	if dbName == "" {
		return "", errs.NewValueIsRequiredError("dbName")
	}
	if sslMode == "" {
		return "", errs.NewValueIsRequiredError("sslMode")
	}
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host,
		port,
		user,
		password,
		dbName,
		sslMode), nil
}

func mustGormOpen(connectionString string) *gorm.DB {
	pgGorm, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  connectionString,
			PreferSimpleProtocol: true,
		},
	), &gorm.Config{})
	if err != nil {
		log.Fatalf("connection to postgres through gorm\n: %s", err)
	}
	return pgGorm
}

func mustAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&questrepo.QuestDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции QuestDTO: %v", err)
	}
}
