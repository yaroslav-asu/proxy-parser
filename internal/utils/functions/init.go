package functions

import (
	"proxy-parser/internal/utils/logger"
	"proxy-parser/internal/utils/vars"
)

func Init() {
	vars.InitEnv()
	logger.InitLogger()
}
