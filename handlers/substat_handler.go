package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"api/models"
)

func ListSubstatsHandler(substats []models.Substat) gin.HandlerFunc {
	return func(c *gin.Context) {
		var substatNames []string
		for _, substat := range substats {
			substatNames = append(substatNames, substat.Name)
		}
		c.JSON(http.StatusOK, gin.H{"substats": substatNames})
	}
}

func GetSubstatHandler(substats []models.Substat) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.ReplaceAll(strings.ToLower(c.Param("name")), "_", " ")
		for _, substat := range substats {
			if strings.ToLower(substat.Name) == name {
				c.JSON(http.StatusOK, substat)
				return
			}
		}
		NotFoundHandler(c, "Substat not found")
	}
}