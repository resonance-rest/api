package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"api/models"
)

func ListSonatasHandler(sonatas []models.Sonata) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sonataNames []string
		for _, sonata := range sonatas {
			sonataNames = append(sonataNames, sonata.Name)
		}
		c.JSON(http.StatusOK, gin.H{"sonatas": sonataNames})
	}
}

func GetSonataHandler(sonatas []models.Sonata) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.ReplaceAll(strings.ToLower(c.Param("name")), "_", " ")
		for _, sonata := range sonatas {
			if strings.ToLower(sonata.Name) == name {
				c.JSON(http.StatusOK, sonata)
				return
			}
		}
		NotFoundHandler(c, "Sonata not found")
	}
}