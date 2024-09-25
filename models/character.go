package models

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