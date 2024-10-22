package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPluginRegistryUpdate struct {
}

// @Id           UpdatePluginRegistry
// @Summary      Update Plugin Registry
// @Description  Update the plugin registry data
// @Tags         Plugin
// @Security     ApiKeyAuth
// @Param        options  body  ParamsPluginRegistryUpdate  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /plugin_registry/update  [post]
func PluginRegistryUpdate(app *core.App, router fiber.Router) {
	router.Post("/plugin_registry/update", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsPluginRegistryUpdate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		registryBaseURL := strings.TrimSuffix(app.Conf().Plugin.RegistryURL, "/")
		registryURL := registryBaseURL + "/registry.json"

		// Http Get
		resp, err := http.Get(registryURL)
		if err != nil || resp.StatusCode != 200 {
			log.Error("[PluginRegistryUpdate] Failed to get registry data: ", resp.StatusCode, " ", err)
			return common.RespError(c, 500, fmt.Sprintf("Failed to get registry data. Got status code: %d", resp.StatusCode))
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("[PluginRegistryUpdate] Failed to read body: ", err)
			return common.RespError(c, 500, "Failed to read body.")
		}

		// Check JSON
		var data entity.PluginRegistryData
		if err := json.Unmarshal(body, &data); err != nil {
			log.Error("[PluginRegistryUpdate] Failed to parse registry data: ", err)
			return common.RespError(c, 400, "Failed to parse registry data.")
		}

		// Unmarshal JSON
		var jsonStr []byte
		if jsonStr, err = json.Marshal(data); err != nil {
			log.Error("[PluginRegistryUpdate] Failed to marshal registry data: ", err)
			return common.RespError(c, 500, "Failed to marshal registry data.")
		}

		// Save to DB
		var record entity.PluginOption
		if err := app.Dao().DB().Where(entity.PluginOption{
			PluginID: entity.PluginID_RegistryData,
		}).FirstOrCreate(&record).Error; err != nil {
			return common.RespError(c, 500, "Failed to find or create registry data.")
		}

		// Update value
		record.Value = string(jsonStr)
		if err := app.Dao().DB().Save(&record).Error; err != nil {
			return common.RespError(c, 500, "Failed to save registry data.")
		}

		return common.RespSuccess(c)
	}))
}
