package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Character struct {
	Name      string `json:"name"`
	Quote     string `json:"quote,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	Weapon    string `json:"weapon,omitempty"`
	Rarity    int    `json:"rarity,omitempty"`
	Class     string `json:"class,omitempty"`
	Birthplace string `json:"birthplace,omitempty"`
	Birthday  string `json:"birthday,omitempty"`
}

type Attribute struct {
	Name       string      `json:"name"`
	Characters []Character `json:"characters"`
}

func main() {
	r := gin.Default()

	characters, err := loadCharactersFromFile("api/characters.json")
	if err != nil {
		fmt.Println("Error loading characters:", err)
		return
	}

	attributes, err := loadAttributesFromFile("api/attributes.json")
	if err != nil {
		fmt.Println("Error loading attributes:", err)
		return
	}

	r.Use(ToLowerMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "1.0",
			"statistics": gin.H{
				"attributes": len(attributes),
				"characters": len(characters),
			},
			"endpoints": gin.H{
				"characters": "/characters",
				"attributes": "/attributes",
			},
		})
	})

	// CHARACTERS

	r.GET("/characters", func(c *gin.Context) {
		var characterNames []string
		for _, character := range characters {
			characterNames = append(characterNames, character.Name)
		}
		c.JSON(http.StatusOK, gin.H{"characters": characterNames})
	})

	r.GET("/characters/:name", func(c *gin.Context) {
		name := c.Param("name")

		for _, character := range characters {
			if strings.EqualFold(character.Name, name) {
				c.JSON(http.StatusOK, character)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Character not found"})
	})

	r.GET("/characters/:name/portrait", func(c *gin.Context) {
		name := c.Param("name")
		filePath := fmt.Sprintf("./cdn/characters/portraits/%s.webp", name)
		_, err := os.Stat(filePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Portrait not found"})
			return
		}

		c.File(filePath)
	})

	// r.StaticFS("/cdn/characters/portraits/", http.Dir("./cdn/characters/portraits"))

	r.GET("/characters/:name/icon", func(c *gin.Context) {
		name := c.Param("name")
		filePath := fmt.Sprintf("./cdn/characters/icons/%s.webp", name)
		_, err := os.Stat(filePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Icon not found"})
			return
		}

		c.File(filePath)
	})

	// r.StaticFS("/cdn/characters/icons/", http.Dir("./cdn/characters/icons"))

	// ATTRIBUTES

	r.GET("/attributes", func(c *gin.Context) {
		var attributeNames []string
		for _, attribute := range attributes {
			attributeNames = append(attributeNames, attribute.Name)
		}
		c.JSON(http.StatusOK, gin.H{"attributes": attributeNames})
	})

	r.GET("/attributes/:name", func(c *gin.Context) {
		name := c.Param("name")

		for _, attribute := range attributes {
			if strings.EqualFold(attribute.Name, name) {
				c.JSON(http.StatusOK, attribute)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Attribute not found"})
	})

	r.GET("/attributes/:name/icon", func(c *gin.Context) {
		name := c.Param("name")
		filePath := fmt.Sprintf("./cdn/attributes/icons/%s.webp", name)
		_, err := os.Stat(filePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Attribute not found"})
			return
		}

		c.File(filePath)
	})

	// r.StaticFS("/cdn/attributes/icons/", http.Dir("./cdn/attributes/icons"))


	r.Run(":8080")
}

func ToLowerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.ToLower(c.Request.URL.Path)
		c.Next()
	}
}

func loadCharactersFromFile(filename string) ([]Character, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var characters []Character
	if err := json.NewDecoder(file).Decode(&characters); err != nil {
		return nil, err
	}

	return characters, nil
}

func loadAttributesFromFile(filename string) ([]Attribute, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var attributes []Attribute
	if err := json.NewDecoder(file).Decode(&attributes); err != nil {
		return nil, err
	}

	return attributes, nil
}
