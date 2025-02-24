package initializer

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/ai-flowx/drivex/pkg/config"
	"github.com/ai-flowx/drivex/pkg/mylog"
)

var (
	once sync.Once
)

func Setup(configName string) error {
	var err error

	once.Do(func() {
		err = config.InitConfig(configName)
		if err != nil {
			log.Println("Error initializing config:", err)
			return
		}
		log.Println("config.InitConfig ok")
		if !config.Debug {
			gin.SetMode(gin.ReleaseMode)
		}
		mylog.InitLog(config.LogLevel)
		log.Println("config.LogLevel ok")
	})

	return err
}

func Cleanup() {
	_ = mylog.Logger.Sync()
}
