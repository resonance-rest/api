package handlers

import (
	"strings"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"api/models" 
)

func ListAttributesHandler(attributes []models.Attribute) gin.HandlerFunc {
	return func(c *gin.Context) {
		var attributeNames []string
		for _, attribute := range attributes {
			attributeNames = append(attributeNames, attribute.Name)
		}
		c.JSON(http.StatusOK, gin.H{"attributes": attributeNames})
	}
}

func GetAttributeHandler(attributes []models.Attribute) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.ToLower(c.Param("name"))

		for _, attribute := range attributes {
			if strings.ToLower(attribute.Name) == name {
				c.JSON(http.StatusOK, attribute)
				return
			}
		}

		NotFoundHandler(c, "Attribute not found")
	}
}

func AttributeIconHandler(c *gin.Context) {
	name := strings.ToLower(c.Param("name"))
	remoteURL := fmt.Sprintf("%sattributes/icon/%s.webp", cdnURL, name)
	
	resp, err := http.Get(remoteURL)
	if err != nil {
		NotFoundHandler(c, "Failed to fetch icon")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		NotFoundHandler(c, "Failed to fetch icon")
		return
	}

	c.Header("Content-Type", "image/webp")
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		NotFoundHandler(c, "Failed to serve icon")
		return
	}
}