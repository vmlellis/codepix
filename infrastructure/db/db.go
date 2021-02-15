package db

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/vmlellis/imersao/codepix-go/domain/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../../.env")

	if err != nil {
		log.Fatalf("Error loading .env files")
	}
}

func ConnectDB(env string) *gorm.DB {
	var dsn, dbType string
	var db *gorm.DB
	var err error

	config := &gorm.Config{}

	if os.Getenv("debug") == "true" {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	if env != "test" {
		dsn = os.Getenv("dsn")
		dbType = os.Getenv("dbType")
	} else {
		dsn = os.Getenv("dsnTest")
		dbType = os.Getenv("dbType")
	}

	if dbType == "postgres" {
		db, err = gorm.Open(postgres.Open(dsn), config)
	} else if dbType == "postgres" {
		db, err = gorm.Open(sqlite.Open(dsn), config)
	} else {
		err = errors.New("db type not suported")
	}

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		panic(err)
	}

	if os.Getenv("AutoMigrateDb") == "true" {
		db.AutoMigrate(&model.Bank{}, &model.Account{}, &model.PixKey{}, &model.Transaction{})
	}

	return db
}
