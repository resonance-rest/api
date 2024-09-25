package models

type Attribute struct {
	Name       string      `json:"name,omitempty"`
	Characters []Character `json:"characters,omitempty"`
}