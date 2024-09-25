package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"api/models"
)

func ListCharactersHandler(characters []models.Character) gin.HandlerFunc {
	return func(c *gin.Context) {
		var characterNames []string
		for _, character := range characters {
			characterNames = append(characterNames, character.Name)
		}
		c.JSON(http.StatusOK, gin.H{"characters": characterNames})
	}
}

func GetCharacterHandler(characters []models.Character) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.ToLower(c.Param("name"))

		for _, character := range characters {
			if strings.ToLower(character.Name) == name {
				c.JSON(http.StatusOK, character)
				return
			}
		}

		NotFoundHandler(c, "Character not found") 
	}
}

func CharacterEmojisHandler(c *gin.Context) {
	name := c.Param("name")
	emojiList, err := loadEmojis(name)
	if err != nil {
		
		return
	}
	c.JSON(http.StatusOK, emojiList)
}

func CharacterEmojiHandler(c *gin.Context) {
	name := c.Param("name")
	index := c.Param("index")

	emojiURL := fmt.Sprintf("%semojis/%s/%s.png", cdnURL, name, index)
	resp, err := http.Get(emojiURL)
	if err != nil {
		NotFoundHandler(c, "Failed to fetch emoji")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		NotFoundHandler(c, "Failed to fetch emoji")
		return
	}

	c.Header("Content-Type", "image/png")

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		NotFoundHandler(c, "Failed to serve emoji")
		return
	}
}

func CharacterImageHandler(c *gin.Context) {
	name := strings.ToLower(c.Param("name"))
	imageType := c.Param("imagetype")

	validTypes := map[string]bool{
		"portrait": true,
		"icon":     true,
		"circle":   true,
		"card":     true,
	}

	if !validTypes[imageType] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid image type", "docs": docsURL})
		return
	}

	remoteURL := fmt.Sprintf("%scharacters/%ss/%s.png", cdnURL, imageType, name)

	resp, err := http.Get(remoteURL)
	if err != nil {
		NotFoundHandler(c, "Character image not found")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		NotFoundHandler(c, "Failed to fetch file")
		return
	}

	c.Header("Content-Type", "image/png")

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		NotFoundHandler(c, "Failed to serve file")
		return
	}
}

func loadEmojis(charName string) (*models.Emojis, error) {
	emojiList := &models.Emojis{}
	for i := 0; ; i++ {
		emojiPath, err := emojisPerCharacter(charName, fmt.Sprintf("%d", i))
		if err != nil {
			break
		}
		emojiList.Emojis = append(emojiList.Emojis, emojiPath)
	}
	return emojiList, nil
}

func emojisPerCharacter(name string, index string) (string, error) {
	emojiURL := fmt.Sprintf("%semojis/%s/%s.png", cdnURL, name, index)
	resp, err := http.Get(emojiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch emoji")
	}

	return emojiURL, nil
}