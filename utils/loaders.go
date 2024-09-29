package utils

import (
	"encoding/json"
	"os"
	"api/models"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func loadJSONFile(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

func LoadCharacters(dirPath string) ([]models.Character, error) {
	var characters []models.Character

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		var character models.Character
		if err := loadJSONFile(filePath, &character); err != nil {
			return nil, fmt.Errorf("error loading character from %s: %v", file.Name(), err)
		}

		character.Name = strings.ReplaceAll(character.Name, " ", "%20")

		characters = append(characters, character)
	}

	return characters, nil
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