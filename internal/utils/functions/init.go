package functions

import (
	"github.com/yaroslav-asu/proxy-parser/internal/utils/logger"
	"github.com/yaroslav-asu/proxy-parser/internal/utils/vars"
)

func Init() {
	vars.InitEnv()
	logger.InitLogger()
}
