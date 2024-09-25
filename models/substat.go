package models

type Substat struct {
	Name string  `json:"name,omitempty"`
	Min  float64 `json:"min,omitempty"`
	Max  float64 `json:"max,omitempty"`
}