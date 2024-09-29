package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"api/handlers"
	"api/models"
	"api/utils"
)

var docsURL = "https://github.com/whosneksio/resonance.rest/blob/main/README.md"

func main() {
	r := gin.Default()

	r.Use(utils.LowercaseMiddleware())
	r.Use(utils.RateLimitMiddleware())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// Load data
	characters, err := utils.LoadCharacters("data/characters")
	if err != nil {
		log.Fatalf("Error loading characters: %v", err)
	}

	attributes, err := utils.LoadAttributes("data/attributes.json")
	if err != nil {
		log.Fatalf("Error loading attributes: %v", err)
	}

	weapons, err := utils.LoadWeapons("data/weapons.json")
	if err != nil {
		log.Fatalf("Error loading weapons: %v", err)
	}

	echoes, err := utils.LoadEchoes("data/echoes.json")
	if err != nil {
		log.Fatalf("Error loading echoes: %v", err)
	}

	sonatas, err := utils.LoadSonatas("data/sonatas.json")
	if err != nil {
		log.Fatalf("Error loading sonatas: %v", err)
	}

	stats, err := utils.LoadStats("data/echoes/stats.json")
	if err != nil {
		log.Fatalf("Error loading stats: %v", err)
	}

	substats, err := utils.LoadSubstats("data/echoes/substats.json")
	if err != nil {
		log.Fatalf("Error loading substats: %v", err)
	}

	codes, err := utils.LoadCodes("data/codes.json")
	if err != nil {
		log.Fatalf("Error loading codes: %v", err)
	}

	setupRoutes(r, characters, attributes, weapons, echoes, sonatas, stats, substats, codes)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRoutes(r *gin.Engine, characters []models.Character, attributes []models.Attribute, weapons map[string][]models.Weapon, echoes []models.Echo, sonatas []models.Sonata, stats []models.Stat, substats []models.Substat, codes []models.Code) {
	r.GET("/", handlers.HomeHandler(characters, attributes, weapons, echoes))
	
	r.NoRoute(func(c *gin.Context) {handlers.NotFoundHandler(c, "Route not found")})

	r.GET("/codes", handlers.CodesHandler(codes))


	// Character routes
	r.GET("/characters", handlers.ListCharactersHandler(characters))
	r.GET("/characters/:name", handlers.GetCharacterHandler(characters))
	r.GET("/characters/:name/emojis", handlers.CharacterEmojisHandler)
	r.GET("/characters/:name/emojis/:index", handlers.CharacterEmojiHandler)
	r.GET("/characters/:name/:imagetype", handlers.CharacterImageHandler)

	// Attribute routes
	r.GET("/attributes", handlers.ListAttributesHandler(attributes))
	r.GET("/attributes/:name", handlers.GetAttributeHandler(attributes))
	r.GET("/attributes/:name/icon", handlers.AttributeIconHandler)

	// Weapon routes
	r.GET("/weapons", handlers.ListWeaponTypesHandler(weapons))
	r.GET("/weapons/:type", handlers.ListWeaponsHandler(weapons))
	r.GET("/weapons/:type/:name", handlers.GetWeaponHandler(weapons))
	r.GET("/weapons/:type/:name/icon", handlers.WeaponIconHandler)

	// Echo routes
	r.GET("/echoes", handlers.ListEchoesHandler(echoes))
	r.GET("/echoes/:name", handlers.GetEchoHandler(echoes))

	// Sonata routes
	r.GET("/echoes/sonatas", handlers.ListSonatasHandler(sonatas))
	r.GET("/echoes/sonatas/:name", handlers.GetSonataHandler(sonatas))

	// Stat routes
	r.GET("/echoes/stats", handlers.ListStatsHandler(stats))
	r.GET("/echoes/stats/:name", handlers.GetStatHandler(stats))

	// Substat routes
	r.GET("/echoes/substats", handlers.ListSubstatsHandler(substats))
	r.GET("/echoes/substats/:name", handlers.GetSubstatHandler(substats))


}