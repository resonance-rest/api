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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Stat struct {
	Cost    int    `json:"cost,omitempty"`
	Name    string `json:"name,omitempty"`
	Primary []struct {
		Name  string    `json:"name,omitempty"`
		Ranks []float64 `json:"ranks,omitempty"`
	} `json:"primary,omitempty"`
	Secondary []struct {
		Name  string    `json:"name,omitempty"`
		Ranks []float64 `json:"ranks,omitempty"`
	} `json:"secondary,omitempty"`
}

type Substat struct {
	Name string  `json:"name,omitempty"`
	Min  float64 `json:"min,omitempty"`
	Max  float64 `json:"max,omitempty"`
}

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

type Sonata struct {
	Name      string `json:"name,omitempty"`
	TwoPiece  string `json:"twoPiece,omitempty"`
	FivePiece string `json:"fivePiece,omitempty"`
}

type Echo struct {
	Name          string        `json:"name,omitempty"`
	Cost          int           `json:"cost,omitempty"`
	SonataEffects []string      `json:"sonataEffects,omitempty"`
	Outline       string        `json:"outline,omitempty"`
	Description   string        `json:"description,omitempty"`
	Ranks         []interface{} `json:"ranks,omitempty"`
	Cooldown      string        `json:"cooldown,omitempty"`
}

type Attribute struct {
	Name       string      `json:"name,omitempty"`
	Characters []Character `json:"characters,omitempty"`
}

type Emojis struct {
	Emojis []string `json:"emojis,omitempty"`
}

type Weapon struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Rarity      int    `json:"rarity,omitempty"`
	Stats       struct {
		Attack  int `json:"atk,omitempty"`
		Substat struct {
			SubName  string `json:"name,omitempty"`
			SubValue string `json:"value,omitempty"`
		} `json:"substat,omitempty"`
	} `json:"stats,omitempty"`
	Skill struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Ranks       []struct {
			Zero  string `json:"0,omitempty"`
			One   string `json:"1,omitempty"`
			Three string `json:"3,omitempty"`
			Four  string `json:"4,omitempty"`
			Five  string `json:"5,omitempty"`
		} `json:"ranks,omitempty"`
	} `json:"skill,omitempty"`
}

var docsURL = "https://github.com/whosneksio/resonance.rest/blob/main/README.md"
var cdnURL = "http://cdn.resonance.rest/"

func main() {
	r := gin.Default()

	r.Use(Middleware())
	r.Use(rateLimitMiddleware())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

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

	echoes, err := echoesLoad("data/echoes.json")
	if err != nil {
		fmt.Println("Error loading echoes:", err)
		return
	}

	sonatas, err := sonatasLoad("data/echoes/sonatas.json")
	if err != nil {
		fmt.Println("Error loading sonatas:", err)
		return
	}

	stats, err := statsLoad("data/echoes/stats.json")
	if err != nil {
		fmt.Println("Error loading stats:", err)
		return
	}

	substats, err := substatsLoad("data/echoes/substats.json")
	if err != nil {
		fmt.Println("Error loading substats:", err)
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
			"statistics": gin.H{"attributes": len(attributes), "characters": len(characters), "weapons": totalWeapons, "echoes": len(echoes)},
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
			remoteURL := fmt.Sprintf("%sattributes/%s/%s.webp", cdnURL, path, name)
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

	r.GET("/echoes", func(c *gin.Context) {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var echoNames []string
		for _, echo := range echoes {
			echoNames = append(echoNames, echo.Name)
		}
		c.JSON(http.StatusOK, gin.H{"echoes": echoNames})
	})

	r.GET("/echoes/:name", func(c *gin.Context) {
		name := strings.ReplaceAll(strings.ToLower(c.Param("name")), "_", " ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, echo := range echoes {
			if strings.ToLower(echo.Name) == name {
				c.JSON(http.StatusOK, echo)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Echo not found", "docs": docsURL})
	})

	r.GET("/echoes/sonatas", func(c *gin.Context) {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var sonataNames []string
		for _, sonata := range sonatas {
			sonataNames = append(sonataNames, sonata.Name)
		}
		c.JSON(http.StatusOK, gin.H{"sonatas": sonataNames})
	})

	r.GET("/echoes/sonatas/:name", func(c *gin.Context) {
		name := strings.ReplaceAll(strings.ToLower(c.Param("name")), "_", " ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, sonata := range sonatas {
			if strings.ToLower(sonata.Name) == name {
				c.JSON(http.StatusOK, sonata)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Sonata not found", "docs": docsURL})
	})

	r.GET("/echoes/stats", func(c *gin.Context) {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var statNames []string
		for _, stat := range stats {
			statNames = append(statNames, stat.Name)
		}
		c.JSON(http.StatusOK, gin.H{"stats": statNames})
	})

	r.GET("/echoes/stats/:name", func(c *gin.Context) {
		name := strings.ReplaceAll(strings.ToLower(c.Param("name")), "_", " ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, stat := range stats {
			if strings.ToLower(stat.Name) == name {
				c.JSON(http.StatusOK, stat)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Stat not found", "docs": docsURL})
	})

	r.GET("/echoes/substats", func(c *gin.Context) {
		substats, err := substatsLoad("data/substats.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var substatNames []string
		for _, substat := range substats {
			substatNames = append(substatNames, substat.Name)
		}
		c.JSON(http.StatusOK, gin.H{"substats": substatNames})
	})

	r.GET("/echoes/substats/:name", func(c *gin.Context) {
		name := strings.ReplaceAll(strings.ToLower(c.Param("name")), "_", " ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, substat := range substats {
			if strings.ToLower(substat.Name) == name {
				c.JSON(http.StatusOK, substat)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Substat not found", "docs": docsURL})
	})

	r.Run(":8080")
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

func echoesLoad(filename string) ([]Echo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var echoes []Echo
	if err := json.NewDecoder(file).Decode(&echoes); err != nil {
		return nil, err
	}

	return echoes, nil
}

func sonatasLoad(filename string) ([]Sonata, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sonatas []Sonata
	if err := json.NewDecoder(file).Decode(&sonatas); err != nil {
		return nil, err
	}

	return sonatas, nil
}

func statsLoad(filename string) ([]Stat, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stats []Stat
	if err := json.NewDecoder(file).Decode(&stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func substatsLoad(filename string) ([]Substat, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var substats []Substat
	if err := json.NewDecoder(file).Decode(&substats); err != nil {
		return nil, err
	}

	return substats, nil
}

var (
	rateLimit    = 200
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

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.ToLower(c.Request.URL.Path)
		c.Next()
	}
}
