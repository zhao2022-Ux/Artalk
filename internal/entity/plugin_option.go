package entity

import (
	"gorm.io/gorm"
)

const PluginID_RegistryData = "__artalk_registry_data__"
const PluginOptionName_OptionsSchema = "options_schema"
const PluginOptionName_ClientOptions = "client_options"

type PluginOption struct {
	gorm.Model

	PluginID string `gorm:"index"`
	Name     string
	Value    string
}

func (n PluginOption) IsEmpty() bool {
	return n.ID == 0
}
