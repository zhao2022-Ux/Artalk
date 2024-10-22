package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/blang/semver"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type ResponsePluginGet struct {
	Enabled       bool                `json:"enabled"`        // The plugin enabled status
	ClientOptions string              `json:"client_options"` // The plugin client options (JSON string)
	OptionsSchema string              `json:"options_schema"` // The plugin options schema (JSON string)
	Plugin        entity.CookedPlugin `json:"plugin"`         // The plugin info
}

// @Id           GetPlugin
// @Summary      Get Plugin Info
// @Description  Get a plugin info by ID
// @Tags         Plugin
// @Security     ApiKeyAuth
// @Param        plugin_id   path   string   true   "The plugin ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponsePluginGet
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /plugins/{plugin_id}  [get]
func PluginGet(app *core.App, router fiber.Router) {
	router.Get("/plugins/:plugin_id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		pluginID := c.Params("plugin_id")

		// Get Installed Data
		var installed entity.Plugin
		app.Dao().DB().Where(&entity.Plugin{PluginID: pluginID}).First(&installed)

		// If Plugin Installed
		clientOptionsJSON := ""
		optionsSchemaJSON := ""
		if !installed.IsEmpty() {
			// Get Options
			clientOptionsJSON = getPluginOptionRecord(app.Dao(), pluginID, entity.PluginOptionName_ClientOptions)

			// Get Options Schema
			optionsSchemaJSON = getPluginOptionRecord(app.Dao(), pluginID, entity.PluginOptionName_OptionsSchema)
		}

		// Get Registry Data
		registryPlugin := findRegistryPlugin(pluginID, app.Dao())

		cookedPlugin := cookPlugin(registryPlugin, &installed)

		return common.RespData(c, ResponsePluginGet{
			Enabled:       installed.Enabled,
			ClientOptions: clientOptionsJSON,
			OptionsSchema: optionsSchemaJSON,
			Plugin:        cookedPlugin,
		})
	}))
}

func getPluginRegistryCache(dao *dao.Dao) (entity.PluginRegistryData, error) {
	var data entity.PluginRegistryData

	var record entity.PluginOption
	dao.DB().Where(&entity.PluginOption{
		PluginID: entity.PluginID_RegistryData,
	}).First(&record)

	if !record.IsEmpty() && record.Value != "" {
		if err := json.Unmarshal([]byte(record.Value), &data); err != nil {
			return entity.PluginRegistryData{}, err
		}
	}
	return data, nil
}

func findRegistryPlugin(pluginID string, dao *dao.Dao) *entity.PluginRegistryItem {
	var registryPlugin *entity.PluginRegistryItem // if exists, the `registryPlugin` will be not nil
	if registryJSON, err := getPluginRegistryCache(dao); err == nil {
		registryPlugin = findPluginInRegistryList(pluginID, registryJSON.Plugins) // Find in Plugins
		if registryPlugin == nil {
			registryPlugin = findPluginInRegistryList(pluginID, registryJSON.Themes) // Find in Themes
		}
	}
	return registryPlugin
}

func findPluginInRegistryList(pluginID string, registryItems []entity.PluginRegistryItem) *entity.PluginRegistryItem {
	var found *entity.PluginRegistryItem
	if rp, ok := lo.Find(registryItems, func(item entity.PluginRegistryItem) bool {
		return item.ID == pluginID
	}); ok {
		found = &rp
	}
	return found
}

func checkPluginInstalled(installed []entity.Plugin, pluginID string) (entity.Plugin, bool) {
	return lo.Find(installed, func(plugin entity.Plugin) bool { return plugin.PluginID == pluginID })
}

func cookPlugin(registryPlugin *entity.PluginRegistryItem, installedPlugin *entity.Plugin) entity.CookedPlugin {
	if registryPlugin == nil && installedPlugin == nil {
		return entity.CookedPlugin{}
	}

	plugin := entity.CookedPlugin{
		Installed: false,
		Enabled:   false,
	}

	// Registry Data exists
	if registryPlugin != nil {
		plugin.PluginRegistryItem = *registryPlugin

		// Check compatibility
		plugin.Compatible = registryPlugin.MinArtalkVersion == "" || !semver.MustParse(registryPlugin.MinArtalkVersion).GT(semver.MustParse(strings.TrimPrefix(config.Version, "v")))
		plugin.CompatibleNotice = lo.If(!plugin.Compatible, fmt.Sprintf("The plugin required at least Artalk v%s.", registryPlugin.MinArtalkVersion)).Else("")
	}

	// Installed Data exists
	if installedPlugin != nil && !installedPlugin.IsEmpty() {
		plugin.Installed = true
		plugin.Enabled = installedPlugin.Enabled

		// Check upgrade
		plugin.UpgradeAvailable = semver.MustParse(registryPlugin.Version).GT(semver.MustParse(installedPlugin.Version))
		plugin.LocalVersion = installedPlugin.Version
	}

	return plugin
}

func getPluginOptionRecord(dao *dao.Dao, pluginID, optionName string) string {
	var record entity.PluginOption
	dao.DB().Where(&entity.PluginOption{PluginID: pluginID, Name: optionName}).First(&record)
	return record.Value
}
