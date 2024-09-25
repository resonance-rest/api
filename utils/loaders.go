package utils

import (
	"encoding/json"
	"os"
	"api/models"
)

func loadJSONFile(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

func LoadCharacters(filename string) ([]models.Character, error) {
	var characters []models.Character
	err := loadJSONFile(filename, &characters)
	return characters, err
}

func LoadAttributes(filename string) ([]models.Attribute, error) {
	var attributes []models.Attribute
	err := loadJSONFile(filename, &attributes)
	return attributes, err
}

func LoadWeapons(filename string) (map[string][]models.Weapon, error) {
	var weapons map[string][]models.Weapon
	err := loadJSONFile(filename, &weapons)
	return weapons, err
}

func LoadEchoes(filename string) ([]models.Echo, error) {
	var echoes []models.Echo
	err := loadJSONFile(filename, &echoes)
	return echoes, err
}

func LoadSonatas(filename string) ([]models.Sonata, error) {
	var sonatas []models.Sonata
	err := loadJSONFile(filename, &sonatas)
	return sonatas, err
}

func LoadStats(filename string) ([]models.Stat, error) {
	var stats []models.Stat
	err := loadJSONFile(filename, &stats)
	return stats, err
}

func LoadSubstats(filename string) ([]models.Substat, error) {
	var substats []models.Substat
	err := loadJSONFile(filename, &substats)
	return substats, err
}

func LoadCodes(filename string) ([]models.Code, error) {
	var codes []models.Code
	err := loadJSONFile(filename, &codes)
	return codes, err
}