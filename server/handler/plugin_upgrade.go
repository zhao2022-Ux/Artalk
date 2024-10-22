package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/gofiber/fiber/v2"
)

type ParamsPluginUpgrade struct {
}

// @Id           UpgradePlugin
// @Summary      Upgrade Plugin
// @Description  Upgrade a plugin by ID
// @Tags         Plugin
// @Security     ApiKeyAuth
// @Param        plugin_id path  string               true  "The plugin ID"
// @Param        options   body  ParamsPluginUpgrade  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /plugins/{plugin_id}/upgrade  [post]
func PluginUpgrade(app *core.App, router fiber.Router) {
	router.Post("/plugins/:plugin_id/upgrade", installPlugin(app, InstallPluginTypeUpgrade))
}
