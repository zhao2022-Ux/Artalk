package handler

import (
	"encoding/json"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPluginUpdate struct {
	Enabled       bool   `json:"enabled" validate:"required"`        // The plugin enabled status
	ClientOptions string `json:"client_options" validate:"optional"` // The plugin client options (JSON string)
}

type ResponsePluginUpdate struct {
	Plugin entity.CookedPlugin `json:"plugin"`
}

// @Id           UpdatePlugin
// @Summary      Update Plugin
// @Description  Update a plugin status by ID
// @Tags         Plugin
// @Security     ApiKeyAuth
// @Param        plugin_id  path  string              true  "The plugin ID"
// @Param        plugin     body  ParamsPluginUpdate  true  "The plugin status"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponsePluginUpdate
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /plugins/{plugin_id}  [put]
func PluginUpdate(app *core.App, router fiber.Router) {
	router.Put("/plugins/:plugin_id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		pluginID := c.Params("plugin_id")
		var p ParamsPluginUpdate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Get Registry Data
		registryPlugin := findRegistryPlugin(pluginID, app.Dao())

		// Check had installed
		var installed entity.Plugin
		app.Dao().DB().Where(&entity.Plugin{PluginID: pluginID}).First(&installed)

		if installed.IsEmpty() {
			return common.RespError(c, 400, "Plugin not installed.")
		}

		// Update Plugin
		installed.Enabled = p.Enabled
		app.Dao().DB().Save(&installed)

		// Update Options
		if p.ClientOptions != "" {
			// Validate JSON
			if err := json.Unmarshal([]byte(p.ClientOptions), &map[string]any{}); err != nil {
				return common.RespError(c, 400, "Param 'client_options' is not a valid JSON string.")
			}

			// (PluginID, Name, Value) = (pluginID, "client_options", p.ClientOptions)
			var clientOptionsRecord entity.PluginOption
			app.Dao().DB().Where(&entity.PluginOption{PluginID: pluginID, Name: entity.PluginOptionName_ClientOptions}).First(&clientOptionsRecord)
			if clientOptionsRecord.IsEmpty() {
				clientOptionsRecord = entity.PluginOption{
					PluginID: pluginID,
					Name:     entity.PluginOptionName_ClientOptions,
					Value:    p.ClientOptions,
				}
			} else {
				clientOptionsRecord.Value = p.ClientOptions
			}
			app.Dao().DB().Save(&clientOptionsRecord)
		}

		cookedPlugin := cookPlugin(registryPlugin, &installed)

		return common.RespData(c, ResponsePluginUpdate{
			Plugin: cookedPlugin,
		})
	}))
}
