package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"api/models"
)

func ListStatsHandler(stats []models.Stat) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statNames []string
		for _, stat := range stats {
			statNames = append(statNames, stat.Name)
		}
		c.JSON(http.StatusOK, gin.H{"stats": statNames})
	}
}

func GetStatHandler(stats []models.Stat) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.ReplaceAll(strings.ToLower(c.Param("name")), "_", " ")
		for _, stat := range stats {
			if strings.ToLower(stat.Name) == name {
				c.JSON(http.StatusOK, stat)
				return
			}
		}
		NotFoundHandler(c, "Stat not found")
	}
}