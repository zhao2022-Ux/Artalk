import { handleConfFormServer } from './config'
import { DefaultPlugins } from './plugins'
import { mergeDeep } from './lib/merge-deep'
import { MountError } from './plugins/mount-error'
import type { CommonPluginItem } from './api/v2'
import type { ConfigPartial, ArtalkPlugin, Context } from '@/types'

/**
 * Global Plugins for all Artalk instances
 */
export const GlobalPlugins: Set<ArtalkPlugin> = new Set([...DefaultPlugins])

/**
 * Plugin options for plugin initialization
 */
export const PluginOptions: WeakMap<ArtalkPlugin, any> = new WeakMap()

export async function mount(localConf: ConfigPartial, ctx: Context) {
  const loaded = new Set<ArtalkPlugin>()
  const loadPlugins = (plugins: Set<ArtalkPlugin>) => {
    plugins.forEach((plugin) => {
      if (typeof plugin !== 'function') return
      if (loaded.has(plugin)) return
      plugin(ctx, PluginOptions.get(plugin))
      loaded.add(plugin)
    })
  }

  // Load local plugins
  loadPlugins(GlobalPlugins)

  // Get conf from server
  const { data } = await ctx
    .getApi()
    .conf.conf()
    .catch((err) => {
      MountError(ctx, { err, onRetry: () => mount(localConf, ctx) })
      throw err
    })

  // Merge remote and local config
  let conf: ConfigPartial = {
    ...localConf,
    apiVersion: data.version?.version, // server version info
  }

  const remoteConf = handleConfFormServer(data.frontend_conf || {})
  conf = conf.preferRemoteConf ? mergeDeep(conf, remoteConf) : mergeDeep(remoteConf, conf)

  // Apply local + remote conf
  ctx.updateConf(conf)

  // Load remote plugins
  const remotePlugins = data.plugins || []
  remotePlugins &&
    (await loadNetworkPlugins(remotePlugins, ctx.getConf().server)
      .then((plugins) => {
        loadPlugins(plugins)
      })
      .catch((err) => {
        console.error('Failed to load plugin', err)
      }))
}

/**
 * Dynamically load plugins from Network
 */
async function loadNetworkPlugins(
  plugins: CommonPluginItem[],
  apiBase: string,
): Promise<Set<ArtalkPlugin>> {
  const networkPlugins = new Set<ArtalkPlugin>()
  if (!plugins || !Array.isArray(plugins)) return networkPlugins

  const tasks: Promise<void>[] = []

  const addPlugin = (targetPlugin: CommonPluginItem) => {
    Object.entries(window.ArtalkPlugins || {}).forEach(([pluginName, plugin]) => {
      if (typeof plugin !== 'function' || networkPlugins.has(plugin)) return

      // Add plugin to list
      networkPlugins.add(plugin)

      // Parse plugin options
      let options: any = {}
      try {
        options = JSON.parse(targetPlugin.options || '{}')
      } catch (err) {
        console.error(
          `[artalk] Failed to parse plugin '${pluginName}' options: '${targetPlugin.options}'.`,
          err,
        )
      }
      PluginOptions.set(plugin, targetPlugin.options)
    })
  }

  plugins.forEach((plugin) => {
    if (!plugin.source) return

    // check url valid
    if (!/^(http|https):\/\//.test(plugin.source))
      plugin.source = `${apiBase.replace(/\/$/, '')}/${plugin.source.replace(/^\//, '')}`

    tasks.push(
      new Promise<void>((resolve) => {
        // check if loaded
        if (document.querySelector(`script[src="${plugin.source}"]`)) {
          resolve()
          return
        }

        // load script
        const script = document.createElement('script')
        script.src = plugin.source
        script.crossOrigin = 'anonymous'
        if (plugin.integrity) script.integrity = plugin.integrity

        document.head.appendChild(script)
        script.onload = () => {
          addPlugin(plugin)
          resolve()
        }
        script.onerror = () => {
          console.error(`[artalk] Failed to load plugin script from '${plugin.source}'.`)
          resolve()
        }
      }),
    )
  })

  await Promise.all(tasks)

  return networkPlugins
}
