package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ai-flowx/drivex/pkg/apis"
	"github.com/ai-flowx/drivex/pkg/config"
	"github.com/ai-flowx/drivex/pkg/handler"
	"github.com/ai-flowx/drivex/pkg/initializer"
	"github.com/ai-flowx/drivex/pkg/mylog"
)

const (
	configName = "config.json"
	valueTrue  = "true"
)

func main() {
	var cfg string

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) > 1 {
		cfg = os.Args[1]
	} else {
		cfg = configName
	}

	if err := initializer.Setup(cfg); err != nil {
		return
	}
	defer initializer.Cleanup()

	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Request-Private-Network"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Private-Network"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		if c.GetHeader("Access-Control-Request-Private-Network") == valueTrue {
			c.Header("Access-Control-Allow-Private-Network", valueTrue)
		}
		if c.GetHeader("Access-Control-Request-Credentials") == valueTrue {
			c.Header("Access-Control-Request-Credentials", valueTrue)
		}
		c.Status(http.StatusNoContent)
	})

	mylog.Logger.Info("check EnableWeb config", zap.Bool("config.GSOAConf.EnableWeb", config.GSOAConf.EnableWeb))

	if config.GSOAConf.EnableWeb {
		mylog.Logger.Info("web enabled")
		r.Static("/static", "./static")
		r.StaticFile("/", "./static/index.html")
		r.GET("/:filename", func(c *gin.Context) {
			filename := c.Param("filename")
			if strings.HasSuffix(filename, ".html") {
				c.File("./static/" + filename)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			}
		})
	}

	r.GET("/v1/models", apis.ModelsHandler)

	v1 := r.Group("/v1")

	v1.POST("/*path", func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "/v1/chat/completions") ||
			strings.HasSuffix(c.Request.URL.Path, "/chat/completions") ||
			strings.HasSuffix(c.Request.URL.Path, "/v1") {
			handler.OpenAIHandler(c)
			return
		} else if strings.HasSuffix(c.Request.URL.Path, "/v1/embeddings") {
			// TBD: FIXME
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Path not found"})
	})

	if err := r.Run(config.ServerPort); err != nil {
		mylog.Logger.Error(err.Error())
		return
	}
}
