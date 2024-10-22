package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPluginUninstall struct {
}

// @Id           UninstallPlugin
// @Summary      Uninstall Plugin
// @Description  Uninstall a plugin by ID
// @Tags         Plugin
// @Security     ApiKeyAuth
// @Param        plugin_id  path  string                 true  "The plugin ID"
// @Param        options    body  ParamsPluginUninstall  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /plugins/{plugin_id}/uninstall  [post]
func PluginUninstall(app *core.App, router fiber.Router) {
	router.Post("/plugins/:plugin_id/uninstall", common.AdminGuard(app, func(c *fiber.Ctx) error {
		pluginID := c.Params("plugin_id")

		var installed entity.Plugin
		app.Dao().DB().Where(&entity.Plugin{PluginID: pluginID}).First(&installed)

		if installed.IsEmpty() {
			return common.RespError(c, 400, "Plugin not installed.")
		}

		if err := app.Dao().DB().Unscoped().Delete(&installed).Error; err != nil {
			return common.RespError(c, 500, "Failed to uninstall plugin.")
		}

		return common.RespSuccess(c)
	}))
}
