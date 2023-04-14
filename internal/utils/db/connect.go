package db

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"proxy-parser/internal/utils/vars"
	"proxy-parser/models"
)

func Connect() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", vars.DbUser, vars.DbPassword, vars.DbName)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&models.Proxy{}, &models.Site{})
	if err != nil {
		zap.L().Fatal("failed to auto migrate database")
	}
	return db
}
