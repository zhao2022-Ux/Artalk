package common

import (
	"fmt"
	"slices"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type ApiVersionData struct {
	App        string `json:"app"`
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
}

func GetApiVersionDataMap() ApiVersionData {
	return ApiVersionData{
		App:        "artalk",
		Version:    strings.TrimPrefix(config.Version, "v"),
		CommitHash: config.CommitHash(),
	}
}

type ConfData struct {
	FrontendConf Map            `json:"frontend_conf"`
	Plugins      []PluginItem   `json:"plugins"`
	Version      ApiVersionData `json:"version"`
}

func GetApiPublicConfDataMap(app *core.App, c *fiber.Ctx) ConfData {
	isAdmin := CheckIsAdminReq(app, c)
	imgUpload := app.Conf().ImgUpload.Enabled
	if isAdmin {
		imgUpload = true // 管理员始终允许上传图片
	}

	frontendConfSrc := app.Conf().Frontend
	if frontendConfSrc == nil {
		frontendConfSrc = make(map[string]interface{})
	}

	frontendConf := make(map[string]interface{})
	utils.CopyStruct(&frontendConfSrc, &frontendConf)

	frontendConf["imgUpload"] = &imgUpload
	if app.Conf().Locale != "" {
		frontendConf["locale"] = app.Conf().Locale
	}

	customPluginURLsRaw := frontendConf["pluginURLs"].([]any)
	customPluginURLs := lo.Map(customPluginURLsRaw, func(url any, _ int) string {
		u, ok := url.(string)
		if !ok {
			return ""
		}
		return u
	})
	frontendConf["pluginURLs"] = []any{} // Cleared, only for forward compatibility

	// Plugin system in Artalk v2.10.0+
	// TODO: this is a db query, should be cached

	return ConfData{
		FrontendConf: frontendConf,
		Plugins:      getEnabledPlugins(app, customPluginURLs),
		Version:      GetApiVersionDataMap(),
	}
}

func handleCustomPluginURLs(app *core.App, urls []string) []string {
	return utils.RemoveDuplicates(lo.Filter(urls, func(u string, _ int) bool {
		if strings.TrimSpace(u) == "" {
			return false
		}
		if !utils.ValidateURL(u) {
			return true
		}
		if trusted, _, _ := middleware.CheckURLTrusted(app, u); trusted {
			return true
		}
		return false
	}))
}

type PluginItem struct {
	Source    string            `json:"source"`
	Type      entity.PluginType `json:"type"`
	Integrity string            `json:"integrity" validate:"optional"`
	Options   string            `json:"options" validate:"optional"`
}

func getEnabledPlugins(app *core.App, customURLs []string) []PluginItem {
	var plugins []PluginItem

	// User plugins
	if app.Conf().Plugin.Enabled {
		var dbPlugins []entity.Plugin
		app.Dao().DB().Where(&entity.Plugin{Enabled: true}).Find(&dbPlugins)

		plugins = append(plugins, lo.Map(dbPlugins, func(plugin entity.Plugin, _ int) PluginItem {
			return PluginItem{
				Source:    plugin.Source,
				Integrity: plugin.Integrity,
				Type:      plugin.Type,
			}
		})...)

		// Custom plugins
		for _, url := range handleCustomPluginURLs(app, customURLs) {
			plugins = append(plugins, PluginItem{
				Source: url,
				Type:   entity.PluginTypePlugin,
			})
		}
	}

	// Import internal plugins
	if app.Conf().Auth.Enabled {
		plugins = append(plugins, PluginItem{
			Source: "dist/plugins/artalk-plugin-auth.js",
			Type:   entity.PluginTypePlugin,
		})
	}

	if !slices.Contains([]string{"en", "zh-CN", ""}, app.Conf().Locale) {
		plugins = append(plugins, PluginItem{
			Source: fmt.Sprintf("dist/i18n/%s.js", app.Conf().Locale),
			Type:   entity.PluginTypePlugin,
		})
	}

	return plugins
}
