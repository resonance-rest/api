package models

type Echo struct {
	Name          string        `json:"name,omitempty"`
	Cost          int           `json:"cost,omitempty"`
	SonataEffects []string      `json:"sonataEffects,omitempty"`
	Outline       string        `json:"outline,omitempty"`
	Description   string        `json:"description,omitempty"`
	Ranks         []interface{} `json:"ranks,omitempty"`
	Cooldown      string        `json:"cooldown,omitempty"`
}