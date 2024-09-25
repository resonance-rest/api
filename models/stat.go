package models

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