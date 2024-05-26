package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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

type Weapon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Rarity      int    `json:"rarity"`
	Stats       struct {
		Attack  int `json:"atk"`
		Substat struct {
			SubName  string `json:"name"`
			SubValue string `json:"value"`
		} `json:"substat"`
	} `json:"stats"`
	Skill struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Ranks       []struct {
			Zero  string `json:"0"`
			One   string `json:"1"`
			Three string `json:"3"`
			Four  string `json:"4"`
			Five  string `json:"5"`
		} `json:"ranks"`
	} `json:"skill,omitempty"`
}

var docsURL = "https://github.com/whosneksio/resonance.rest/blob/main/README.md"
var cdnURL = "http://cdn.resonance.rest/"

func main() {
	r := gin.Default()
	r.Use(Middleware())
	r.Use(rateLimitMiddleware())

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

	weapons, err := loadWeapons("data/weapons.json")
	if err != nil {
		fmt.Println("Error loading weapons:", err)
		return
	}

	totalWeapons := 0
	for _, weaponList := range weapons {
		totalWeapons += len(weaponList)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":    "1.0",
			"docs":       docsURL,
			"statistics": gin.H{"attributes": len(attributes), "characters": len(characters), "weapons": totalWeapons},
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found", "docs": docsURL})
	})

	r.GET("/codes", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"codes": []string{"WUTHERINGGIFT"}})
	})
	// EMOJIS

	r.GET("/characters/:name/emojis", func(c *gin.Context) {
		name := c.Param("name")
		emojiList, err := emojisLoad(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, emojiList)
	})

	r.GET("/characters/:name/emojis/:index", func(c *gin.Context) {
		name := c.Param("name")
		index := c.Param("index")

		emojiURL := fmt.Sprintf("%semojis/%s/%s.png", cdnURL, name, index)
		resp, err := http.Get(emojiURL)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch emoji"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch emoji"})
			return
		}

		c.Header("Content-Type", "image/png")

		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serve emoji"})
			return
		}
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
			remoteURL := fmt.Sprintf("%scharacters/%ss/%s.png", cdnURL, path, name)

			resp, err := http.Get(remoteURL)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%s not found", strings.Title(path)), "docs": docsURL})
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				c.JSON(resp.StatusCode, gin.H{"message": "Failed to fetch file from remote URL", "docs": docsURL})
				return
			}

			c.Header("Content-Type", "image/png")

			_, err = io.Copy(c.Writer, resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to serve file", "docs": docsURL})
				return
			}
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
			remoteURL := fmt.Sprintf("%sattributes/%s/%s.png", cdnURL, path, name)
			resp, err := http.Get(remoteURL)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%s not found", strings.Title(path)), "docs": docsURL})
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				c.JSON(resp.StatusCode, gin.H{"message": "Failed to fetch file from remote URL", "docs": docsURL})
				return
			}

			c.Header("Content-Type", "image/png")

			_, err = io.Copy(c.Writer, resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to serve file", "docs": docsURL})
				return
			}
		})
	}

	// WEAPONS

	r.GET("/weapons", func(c *gin.Context) {
		var weaponTypes []string
		for weaponType := range weapons {
			weaponTypes = append(weaponTypes, weaponType)
		}
		c.JSON(http.StatusOK, gin.H{"types": weaponTypes})
	})

	r.GET("/weapons/:type", func(c *gin.Context) {
		weaponType := strings.ToLower(c.Param("type"))

		weaponsOfType, ok := weapons[weaponType]
		if !ok {
			weaponType = strings.Title(weaponType)
			weaponsOfType, ok = weapons[weaponType]
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "Weapon type not found", "docs": docsURL})
				return
			}
		}

		var weaponNames []string
		for _, weapon := range weaponsOfType {
			weaponNames = append(weaponNames, weapon.Name)
		}

		c.JSON(http.StatusOK, gin.H{"weapons": weaponNames})
	})

	r.GET("/weapons/:type/:name", func(c *gin.Context) {
		weaponType := strings.ToLower(c.Param("type"))
		weaponName := strings.ToLower(strings.ReplaceAll(c.Param("name"), "_", " "))

		var foundWeapon *Weapon
		weaponsOfType, ok := weapons[weaponType]
		if !ok {
			weaponType = strings.Title(weaponType)
			weaponsOfType, ok = weapons[weaponType]
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "Weapon type not found", "docs": docsURL})
				return
			}
		}

		for _, w := range weaponsOfType {
			if strings.EqualFold(w.Name, weaponName) {
				foundWeapon = &w
				break
			}
		}

		if foundWeapon != nil {
			c.JSON(http.StatusOK, foundWeapon)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Weapon not found", "docs": docsURL})
		}
	})

	r.GET("/weapons/:type/:name/icon", func(c *gin.Context) {
		weaponType := strings.ToLower(c.Param("type"))
		weaponName := strings.ToLower(strings.ReplaceAll(c.Param("name"), "_", " "))
		weaponName = strings.ReplaceAll(weaponName, " ", "_")

		remoteURL := fmt.Sprintf("%sweapons/%s/%s.png", cdnURL, weaponType, weaponName)
		resp, err := http.Get(remoteURL)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Icon not found", "docs": docsURL})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch icon from remote URL", "docs": docsURL})
			return
		}

		c.Header("Content-Type", "image/png")

		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serve icon", "docs": docsURL})
			return
		}
	})

	r.Run(":8080")
}

func Middleware() gin.HandlerFunc {
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
		return "", fmt.Errorf("failed to fetch emoji")
	}

	return emojiURL, nil
}

func loadWeapons(filename string) (map[string][]Weapon, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var weapons map[string][]Weapon
	if err := json.NewDecoder(file).Decode(&weapons); err != nil {
		return nil, err
	}

	return weapons, nil
}

var (
	rateLimit    = 60
	requestCount = make(map[string]int)
	mutex        sync.Mutex
)

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mutex.Lock()
		defer mutex.Unlock()

		requestCount[ip]++

		if requestCount[ip] > rateLimit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		go func() {
			time.Sleep(time.Minute)
			mutex.Lock()
			defer mutex.Unlock()
			requestCount[ip] = 0
		}()

		c.Next()
	}
}
