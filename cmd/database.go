package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/adapters/out/postgres/locationrepo"
	"quest-manager/internal/adapters/out/postgres/questrepo"
	"quest-manager/internal/pkg/errs"

	_ "github.com/lib/pq"
)

// MustConnectDB creates database connection for tests
func MustConnectDB(databaseURL string) (*gorm.DB, *sql.DB, error) {
	db := MustGormOpen(databaseURL)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	return db, sqlDB, nil
}

func CreateDbIfNotExists(host string, port string, user string,
	password string, dbName string, sslMode string) {
	dsn, err := MakeConnectionString(host, port, user, password, "postgres", sslMode)
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

func MakeConnectionString(host string, port string, user string,
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

func MustGormOpen(connectionString string) *gorm.DB {
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

func MustAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&questrepo.QuestDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции QuestDTO: %v", err)
	}
	err = db.AutoMigrate(&locationrepo.LocationDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции LocationDTO: %v", err)
	}
	err = db.AutoMigrate(&eventrepo.EventDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции EventDTO: %v", err)
	}
}
