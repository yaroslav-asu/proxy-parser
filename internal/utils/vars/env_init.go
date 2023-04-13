package vars

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

var (
	RunningMode string
	DbUser      string
	DbPassword  string
	DbName      string
)

func initDefaultEnv() {
	zap.L().Info("Started to initialize environmental vars")
	err := godotenv.Load(".env")
	if err != nil {
		zap.L().Warn("Failed to load .env file")
	}
	RunningMode = os.Getenv("RUNNING_MODE")
	zap.L().Info("Environmental vars successfully initialized")

}

func initDbEnv() {
	zap.L().Info("Started to initialize environmental vars for db")
	err := godotenv.Load(".env.db")
	if err != nil {
		zap.L().Warn("Failed to load .env.db file")
	}
	DbUser = os.Getenv("POSTGRES_USER")
	DbPassword = os.Getenv("POSTGRES_PASSWORD")
	DbName = os.Getenv("POSTGRES_DB")
}

func InitEnv() {
	initDefaultEnv()
	initDbEnv()
}
