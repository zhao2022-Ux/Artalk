package entity

import (
	"gorm.io/gorm"
)

type PluginType string

const (
	PluginTypePlugin PluginType = "plugin"
	PluginTypeTheme  PluginType = "theme"
)

type Plugin struct {
	gorm.Model

	PluginID  string `gorm:"index"`
	Name      string
	Type      PluginType
	Source    string
	Integrity string
	Version   string
	Enabled   bool
}

func (n Plugin) IsEmpty() bool {
	return n.ID == 0
}
