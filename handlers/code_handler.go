package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"api/models"
)

func CodesHandler(codes []models.Code) gin.HandlerFunc {
	return func(c *gin.Context) {
		var codesWithRewards []gin.H 

		for _, code := range codes {
			codesWithRewards = append(codesWithRewards, gin.H{
				"name":   code.Name,
				"reward": code.Reward, 
			})
		}

		c.JSON(http.StatusOK, gin.H{"codes": codesWithRewards})
	}
}