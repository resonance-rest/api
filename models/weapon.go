package models

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