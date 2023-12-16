import './style/main.scss'

import type { ArtalkConfig, EventPayloadMap, ArtalkPlugin, ContextApi } from '~/types'
import type { EventHandler } from './lib/event-manager'
import ConcreteContext from './context'
import defaults from './defaults'
import { handelCustomConf, convertApiOptions } from './config'
import Services from './service'
import { DefaultPlugins } from './plugins'
import * as Stat from './plugins/stat'
import Api from './api'

/** Global Plugins for all instances */
const GlobalPlugins: ArtalkPlugin[] = [ ...DefaultPlugins ]

/**
 * Artalk
 *
 * @see https://artalk.js.org
 */
export default class Artalk {
  public conf!: ArtalkConfig
  public ctx!: ContextApi
  public $root!: HTMLElement

  /** Plugins */
  protected plugins: ArtalkPlugin[] = [ ...GlobalPlugins ]

  constructor(conf: Partial<ArtalkConfig>) {
    // Init Config
    this.conf = handelCustomConf(conf)
    if (this.conf.el instanceof HTMLElement) this.$root = this.conf.el

    // Init Context
    this.ctx = new ConcreteContext(this.conf, this.$root)

    // Init Services
    Object.entries(Services).forEach(([name, initService]) => {
      const obj = initService(this.ctx)
      if (obj) this.ctx.inject(name as any, obj) // auto inject deps to ctx
    })

    // Init Plugins
    this.plugins.forEach(plugin => {
      if (typeof plugin === 'function') plugin(this.ctx)
    })

    // Trigger inited event
    this.ctx.trigger('inited')
  }

  /** Use Plugin (plugin will be called in instance `use` func) */
  public use(plugin: ArtalkPlugin) {
    if (typeof plugin === 'function') plugin(this.ctx)
  }

  /** Update config of Artalk */
  public update(conf: Partial<ArtalkConfig>) {
    this.ctx.updateConf(conf)
    return this
  }

  /** Reload comment list of Artalk */
  public reload() {
    this.ctx.reload()
  }

  /** Destroy instance of Artalk */
  public destroy() {
    this.ctx.trigger('destroy')
    this.$root.remove()
  }

  /** Add an event listener */
  public on<K extends keyof EventPayloadMap>(name: K, handler: EventHandler<EventPayloadMap[K]>) {
    this.ctx.on(name, handler)
  }

  /** Remove an event listener */
  public off<K extends keyof EventPayloadMap>(name: K, handler: EventHandler<EventPayloadMap[K]>) {
    this.ctx.off(name, handler)
  }

  /** Trigger an event */
  public trigger<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K]) {
    this.ctx.trigger(name, payload)
  }

  /** Set dark mode */
  public setDarkMode(darkMode: boolean) {
    this.ctx.setDarkMode(darkMode)
  }

  // ===========================
  //       Static Members
  // ===========================

  /** Default Configs */
  public static readonly defaults: ArtalkConfig = defaults

  /** Init Artalk */
  public static init(conf: Partial<ArtalkConfig>): Artalk {
    return new Artalk(conf)
  }

  /** Use Plugin for all instances */
  public static use(plugin: ArtalkPlugin) {
    GlobalPlugins.push(plugin)
  }

  /** Load count widget */
  public static loadCountWidget(c: Partial<ArtalkConfig>) {
    const conf = handelCustomConf(c)

    Stat.initCountWidget({
      getApi: () => new Api(convertApiOptions(conf)),
      pageKey: conf.pageKey,
      countEl: conf.countEl,
      pvEl: conf.pvEl,
      pvAdd: false
    })
  }
}
