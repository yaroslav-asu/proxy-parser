package db

import (
	"fmt"
	models2 "github.com/yaroslav-asu/proxy-parser/internal/models"
	"github.com/yaroslav-asu/proxy-parser/internal/utils/vars"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", vars.DbUser, vars.DbPassword, vars.DbName)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&models2.Proxy{}, &models2.Site{})
	if err != nil {
		zap.L().Fatal("failed to auto migrate database")
	}
	return db
}

func Close(db *gorm.DB) {
	postgresDB, err := db.DB()
	if err != nil {
		zap.L().Error("Failed to get db instance: " + err.Error())
		zap.L().Info("DB connection wasn't close")
		return
	}
	err = postgresDB.Close()
	if err != nil {
		zap.L().Info("DB connection wasn't close: " + err.Error())
	}
}
