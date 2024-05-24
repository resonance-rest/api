package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Character struct {
	Name       string `json:"name"`
	Quote      string `json:"quote,omitempty"`
	Attribute  string `json:"attribute,omitempty"`
	Weapon     string `json:"weapon,omitempty"`
	Rarity     int    `json:"rarity,omitempty"`
	Class      string `json:"class,omitempty"`
	Birthplace string `json:"birthplace,omitempty"`
	Birthday   string `json:"birthday,omitempty"`
}

type Attribute struct {
	Name       string      `json:"name"`
	Characters []Character `json:"characters,omitempty"`
}

type Emojis struct {
	Emojis []string `json:"emojis"`
}

func main() {
	docsURL := "https://github.com/whosneksio/wuwa.api/blob/main/README.md"

	r := gin.Default()
	r.Use(middleware())

	characters, err := charactersLoad("data/characters.json")
	if err != nil {
		fmt.Println("Error loading characters:", err)
		return
	}

	attributes, err := attributesLoad("data/attributes.json")
	if err != nil {
		fmt.Println("Error loading attributes:", err)
		return
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":    "1.0",
			"docs":       docsURL,
			"statistics": gin.H{"attributes": len(attributes), "characters": len(characters)},
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found", "docs": docsURL})
	})

	// EMOJIS

	r.GET("/characters/:name/emojis", func(c *gin.Context) {
		name := strings.ToLower(c.Param("name"))
		emojiList, err := emojisLoad(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, emojiList)
	})

	r.GET("/characters/:name/emojis/:index", func(c *gin.Context) {
		name := strings.ToLower(c.Param("name"))
		indexStr := c.Param("index")
		number, err := strconv.Atoi(indexStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
			return
		}

		emojiPath, err := emojisPerCharacter(name, number)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Emoji not found"})
			return
		}

		c.File(emojiPath)
	})

	//r.StaticFS("/cdn/emojis/", http.Dir("./cdn/emojis/"))

	// CHARACTERS

	r.GET("/characters", func(c *gin.Context) {
		var characterNames []string
		for _, character := range characters {
			characterNames = append(characterNames, character.Name)
		}
		c.JSON(http.StatusOK, gin.H{"characters": characterNames})
	})

	r.GET("/characters/:name", func(c *gin.Context) {
		name := strings.ToLower(c.Param("name"))

		for _, character := range characters {
			if strings.ToLower(character.Name) == name {
				c.JSON(http.StatusOK, character)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Character not found", "docs": docsURL})
	})

	for _, path := range []string{"portrait", "icon", "circle"} {
		r.GET(fmt.Sprintf("/characters/:name/%s", path), func(c *gin.Context) {
			name := strings.ToLower(c.Param("name"))
			filePath := fmt.Sprintf("./cdn/characters/%ss/%s.png", path, name)
			if _, err := os.Stat(filePath); err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%s not found", strings.Title(path)), "docs": docsURL})
				return
			}
			c.File(filePath)
		})
	}

	// ATTRIBUTES

	r.GET("/attributes", func(c *gin.Context) {
		var attributeNames []string
		for _, attribute := range attributes {
			attributeNames = append(attributeNames, attribute.Name)
		}
		c.JSON(http.StatusOK, gin.H{"attributes": attributeNames})
	})

	r.GET("/attributes/:name", func(c *gin.Context) {
		name := strings.ToLower(c.Param("name"))

		for _, attribute := range attributes {
			if strings.ToLower(attribute.Name) == name {
				c.JSON(http.StatusOK, attribute)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Attribute not found", "docs": docsURL})
	})

	for _, path := range []string{"icon"} {
		r.GET(fmt.Sprintf("/attributes/:name/%s", path), func(c *gin.Context) {
			name := strings.ToLower(c.Param("name"))
			filePath := fmt.Sprintf("./cdn/attributes/%s/%s.png", path, name)
			if _, err := os.Stat(filePath); err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%s not found", strings.Title(path)), "docs": docsURL})
				return
			}
			c.File(filePath)
		})
	}

	r.Run(":8080")
}

func middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.ToLower(c.Request.URL.Path)
		c.Next()
	}
}

func charactersLoad(filename string) ([]Character, error) {
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

func attributesLoad(filename string) ([]Attribute, error) {
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

func emojisLoad(charName string) (*Emojis, error) {
	emojiList := &Emojis{}
	files, err := ioutil.ReadDir(fmt.Sprintf("./cdn/emojis/%s", charName))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".png") {
			emojiList.Emojis = append(emojiList.Emojis, strings.TrimSuffix(file.Name(), ".png"))
		}
	}

	return emojiList, nil
}

func emojisPerCharacter(name string, index int) (string, error) {
	emojisFolder := fmt.Sprintf("./cdn/emojis/%s", name)
	files, err := ioutil.ReadDir(emojisFolder)
	if err != nil {
		return "", err
	}

	var pngFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".png" {
			pngFiles = append(pngFiles, file.Name())
		}
	}

	if index < 0 || index >= len(pngFiles) {
		return "", fmt.Errorf("index out of range")
	}

	emojiPath := filepath.Join(emojisFolder, pngFiles[index])
	return emojiPath, nil
}
