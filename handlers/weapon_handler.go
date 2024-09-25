package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"io"
	"strings"
	"api/models"
)

func ListWeaponTypesHandler(weapons map[string][]models.Weapon) gin.HandlerFunc {
	return func(c *gin.Context) {
		var weaponTypes []string
		for weaponType := range weapons {
			weaponTypes = append(weaponTypes, weaponType)
		}
		c.JSON(http.StatusOK, gin.H{"types": weaponTypes})
	}
}

func ListWeaponsHandler(weapons map[string][]models.Weapon) gin.HandlerFunc {
	return func(c *gin.Context) {
		weaponType := strings.ToLower(c.Param("type"))

		weaponsOfType, ok := weapons[weaponType]
		if !ok {
			NotFoundHandler(c, "Weapon type not found")
			return
		}

		var weaponNames []string
		for _, weapon := range weaponsOfType {
			weaponNames = append(weaponNames, weapon.Name)
		}

		c.JSON(http.StatusOK, gin.H{"weapons": weaponNames})
	}
}

func GetWeaponHandler(weapons map[string][]models.Weapon) gin.HandlerFunc {
    return func(c *gin.Context) {
        weaponType := strings.ToLower(c.Param("type"))
        weaponName := strings.ToLower(strings.ReplaceAll(c.Param("name"), "_", " "))

        weaponsOfType, ok := weapons[weaponType]
        if !ok {
            NotFoundHandler(c, "Weapon type not found")
            return
        }

        for i := range weaponsOfType {
            if strings.ToLower(weaponsOfType[i].Name) == weaponName {
                c.JSON(http.StatusOK, weaponsOfType[i])
                return
            }
        }

        NotFoundHandler(c, "Weapon not found")
    }
}

func WeaponIconHandler(c *gin.Context) {
	weaponType := strings.ToLower(c.Param("type"))
	weaponName := strings.ToLower(strings.ReplaceAll(c.Param("name"), "_", " "))
	weaponName = strings.ReplaceAll(weaponName, " ", "_")

	remoteURL := fmt.Sprintf("%sweapons/%s/%s.png", cdnURL, weaponType, weaponName)
	resp, err := http.Get(remoteURL)
	if err != nil {
		NotFoundHandler(c, "Failed to fetch file")
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