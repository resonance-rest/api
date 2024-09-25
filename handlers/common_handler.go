package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"api/models"
)

var (
	docsURL = "https://github.com/whosneksio/resonance.rest/blob/main/README.md"
	cdnURL  = "http://cdn.resonance.rest/"
)

func HomeHandler(characters []models.Character, attributes []models.Attribute, weapons map[string][]models.Weapon, echoes []models.Echo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "2.0",
			"statistics": gin.H{
				"attributes": len(attributes),
				"characters": len(characters),
				"weapons":    len(weapons),
				"echoes":     len(echoes),
			},
		})
	}
}

func NotFoundHandler(c *gin.Context, message ...string) {
    msg := "Not found"
    if len(message) > 0 {
        msg = message[0]
    }
    c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": msg})
}





