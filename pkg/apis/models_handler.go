package apis

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ai-flowx/drivex/pkg/config"
)

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

func ModelsHandler(c *gin.Context) {
	var models []Model

	keys := make([]string, 0, len(config.ModelToService))

	for k := range config.SupportModels {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	t := time.Now()

	for _, k := range keys {
		models = append(models, Model{
			ID:      k,
			Object:  "model",
			Created: t.Unix(),
			OwnedBy: "openai",
		})
	}

	if len(models) > 0 {
		models = append(models, Model{
			ID:      "random",
			Object:  "model",
			Created: t.Unix(),
			OwnedBy: "openai",
		})
	}

	if len(models) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No models found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"object": "list",
		"data":   models,
	})
}
