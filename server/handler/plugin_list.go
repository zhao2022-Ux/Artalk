package handler

import (
	"slices"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type ParamsPluginList struct {
	Search        string `query:"search" json:"search" validate:"optional"`                 // Search keywords
	OnlyInstalled bool   `query:"only_installed" json:"only_installed" validate:"optional"` // Only installed plugins
}

type ResponsePluginList struct {
	Plugins      []entity.CookedPlugin `json:"plugins"`
	Themes       []entity.CookedPlugin `json:"themes"`
	PluginsCount int64                 `json:"plugins_count"`
	ThemesCount  int64                 `json:"themes_count"`
}

// @Id           GetPlugins
// @Summary      Get Plugin List
// @Description  Get a list of plugins by some conditions
// @Tags         Plugin
// @Param        options  query  ParamsPluginList  true   "The options"
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponsePluginList
// @Failure      403  {object}  Map{msg=string}
// @Router       /plugins  [get]
func PluginList(app *core.App, router fiber.Router) {
	router.Get("/plugins", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsPluginList
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Find Installed Plugins
		var installed []entity.Plugin
		app.Dao().DB().Model(&entity.Plugin{}).Order("created_at DESC").Find(&installed)

		// Find Registry Plugins
		registry, err := getPluginRegistryCache(app.Dao())
		if err != nil {
			log.Error("[PluginList] Failed to get registry data: ", err)
			return common.RespError(c, 500, "Failed to get registry data.")
		}

		plugins := mergePluginsWithInstallInfo(registry.Plugins, installed)
		themes := mergePluginsWithInstallInfo(registry.Themes, installed)

		// Supply local plugins that are not in the registry
		extendPluginsOutRegistry(plugins, installed)
		extendPluginsOutRegistry(themes, installed)

		// Filter
		if p.Search != "" {
			filterFn := func(plugin entity.CookedPlugin, index int) bool {
				return strings.Contains(strings.ToLower(plugin.Name), strings.ToLower(p.Search)) || strings.EqualFold(plugin.ID, p.Search)
			}
			plugins = lo.Filter(plugins, filterFn)
			themes = lo.Filter(themes, filterFn)
		}

		if p.OnlyInstalled {
			filterFn := func(plugin entity.CookedPlugin, index int) bool {
				return plugin.Installed
			}
			plugins = lo.Filter(plugins, filterFn)
			themes = lo.Filter(themes, filterFn)
		}

		// Sort by name
		sortFn := func(a, b entity.CookedPlugin) int {
			return strings.Compare(a.Name, b.Name)
		}
		slices.SortFunc(plugins, sortFn)
		slices.SortFunc(themes, sortFn)

		return common.RespData(c, ResponsePluginList{
			Plugins:      plugins,
			PluginsCount: int64(len(plugins)),
			Themes:       themes,
			ThemesCount:  int64(len(themes)),
		})
	}))
}

func mergePluginsWithInstallInfo(
	registryPlugins []entity.PluginRegistryItem,
	installedPluginList []entity.Plugin,
) []entity.CookedPlugin {
	return lo.Map(registryPlugins, func(registryPlugin entity.PluginRegistryItem, _ int) entity.CookedPlugin {
		var installedPlugin *entity.Plugin
		if p, ok := checkPluginInstalled(installedPluginList, registryPlugin.ID); ok {
			installedPlugin = &p
		}

		return cookPlugin(&registryPlugin, installedPlugin)
	})
}

func extendPluginsOutRegistry(plugins []entity.CookedPlugin, installed []entity.Plugin) {
	for _, insPlugin := range installed {
		// Skip plugin already in the list
		if lo.ContainsBy(plugins, func(item entity.CookedPlugin) bool {
			return item.ID == insPlugin.PluginID
		}) {
			continue
		}

		// Extend plugin out of registry but installed
		plugins = append(plugins, entity.CookedPlugin{
			Enabled:    insPlugin.Enabled,
			Installed:  true,
			Compatible: true,
			PluginRegistryItem: entity.PluginRegistryItem{
				ID:        insPlugin.PluginID,
				Name:      insPlugin.Name,
				Type:      string(insPlugin.Type),
				Source:    insPlugin.Source,
				Integrity: insPlugin.Integrity,
				Version:   insPlugin.Version,
			},
		})
	}
}
