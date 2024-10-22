package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/blang/semver"
	"github.com/gofiber/fiber/v2"
)

type ParamsPluginInstall struct {
}

// @Id           InstallPlugin
// @Summary      Install Plugin
// @Description  Install a plugin by ID
// @Tags         Plugin
// @Security     ApiKeyAuth
// @Param        plugin_id  path  string               true  "The plugin ID"
// @Param        options    body  ParamsPluginInstall  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /plugins/{plugin_id}/install  [post]
func PluginInstall(app *core.App, router fiber.Router) {
	router.Post("/plugins/:plugin_id/install", installPlugin(app, InstallPluginTypeInstall))
}

type InstallPluginType string

const (
	InstallPluginTypeInstall InstallPluginType = "install"
	InstallPluginTypeUpgrade InstallPluginType = "upgrade"
)

func installPlugin(app *core.App, mode InstallPluginType) func(*fiber.Ctx) error {
	return common.AdminGuard(app, func(c *fiber.Ctx) error {
		isUpgradeMode := mode == InstallPluginTypeUpgrade
		pluginID := c.Params("plugin_id")

		// Check had installed
		var installed entity.Plugin
		app.Dao().DB().Where(&entity.Plugin{PluginID: pluginID}).First(&installed)

		// Find Registry Plugins
		plugin := findRegistryPlugin(pluginID, app.Dao())
		if plugin == nil {
			return common.RespError(c, 400, "Plugin not found.")
		}

		if !isUpgradeMode {
			if !installed.IsEmpty() {
				return common.RespError(c, 400, "Plugin already installed.")
			}

			// Install Plugin
			app.Dao().DB().Save(&entity.Plugin{
				PluginID:  plugin.ID,
				Name:      plugin.Name,
				Type:      entity.PluginType(plugin.Type),
				Source:    plugin.Source,
				Integrity: plugin.Integrity,
				Version:   plugin.Version,
				Enabled:   true,
			})

			downloadRemoteOptionsSchema(app, plugin.ID, plugin.OptionsSchema)
		} else {
			if installed.IsEmpty() {
				return common.RespError(c, 400, "Plugin not installed.")
			}

			// Upgrade Plugin
			if semver.MustParse(plugin.Version).GT(semver.MustParse(installed.Version)) {
				installed.Version = plugin.Version
				installed.Source = plugin.Source
				installed.Integrity = plugin.Integrity
				app.Dao().DB().Save(&installed)

				downloadRemoteOptionsSchema(app, plugin.ID, plugin.OptionsSchema)
			} else {
				return common.RespError(c, 400, "No new version available.")
			}
		}

		return common.RespSuccess(c)
	})
}

// Download remote options schema (artalk-plugin-options.schema.json)
// Save to database table plugin_options, (plugin_id, name=options_schema, value=JSON)
func downloadRemoteOptionsSchema(app *core.App, pluginID string, schemaURL string) {
	if schemaURL == "" {
		return
	}

	// Download from url
	resp, err := http.Get(schemaURL)
	if err != nil || resp.StatusCode != 200 {
		log.Error("[PluginOptionsSchema] Failed to get options schema from '", schemaURL, "': ", resp.StatusCode, " ", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("[PluginOptionsSchema] Failed to read body from '", schemaURL, "': ", err)
		return
	}

	// Check is json valid
	if err := json.Unmarshal(body, &map[string]any{}); err != nil {
		log.Error("[PluginOptionsSchema] Options schema '", schemaURL, "' is not valid JSON: ", err)
		return
	}

	// Find Plugin Options Schema
	var record entity.PluginOption
	app.Dao().DB().Where(&entity.PluginOption{
		PluginID: pluginID,
		Name:     entity.PluginOptionName_OptionsSchema,
	}).First(&record)

	if record.IsEmpty() {
		// Create Plugin Options Schema
		record = entity.PluginOption{
			PluginID: pluginID,
			Name:     entity.PluginOptionName_OptionsSchema,
			Value:    string(body),
		}
	} else {
		// Update Plugin Options Schema
		record.Value = string(body)
	}

	app.Dao().DB().Save(&record)
}
