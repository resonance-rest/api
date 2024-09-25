package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"api/models"
)

func ListEchoesHandler(echoes []models.Echo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var echoNames []string
		for _, echo := range echoes {
			echoNames = append(echoNames, echo.Name)
		}
		c.JSON(http.StatusOK, gin.H{"echoes": echoNames})
	}
}

func GetEchoHandler(echoes []models.Echo) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.ReplaceAll(strings.ToLower(c.Param("name")), "_", " ")
		for _, echo := range echoes {
			if strings.ToLower(echo.Name) == name {
				c.JSON(http.StatusOK, echo)
				return
			}
		}
		NotFoundHandler(c, "Echo not found")
	}
}