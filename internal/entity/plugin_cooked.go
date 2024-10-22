package entity

type CookedPlugin struct {
	PluginRegistryItem

	Installed        bool   `json:"installed"`
	Enabled          bool   `json:"enabled"`
	LocalVersion     string `json:"local_version"`
	UpgradeAvailable bool   `json:"upgrade_available"`
	Compatible       bool   `json:"compatible"`
	CompatibleNotice string `json:"compatible_notice,omitempty" validate:"optional"`
}

type PluginRegistryItem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Description      string `json:"description"`
	AuthorName       string `json:"author_name"`
	AuthorLink       string `json:"author_link"`
	RepoName         string `json:"repo_name"`
	RepoLink         string `json:"repo_link"`
	NpmName          string `json:"npm_name"`
	Source           string `json:"source"`
	Integrity        string `json:"integrity"`
	OptionsSchema    string `json:"options_schema"`
	DonateLink       string `json:"donate_link"`
	Verified         bool   `json:"verified"`
	Version          string `json:"version"`
	UpdatedAt        string `json:"updated_at"`
	MinArtalkVersion string `json:"min_artalk_version"`
}

type PluginRegistryData struct {
	Plugins []PluginRegistryItem `json:"plugins"`
	Themes  []PluginRegistryItem `json:"themes"`
}
