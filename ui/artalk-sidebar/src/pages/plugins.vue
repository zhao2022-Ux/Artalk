<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { format } from 'timeago.js'
import { useNavStore } from '../stores/nav'
import { getArtalk } from '@/global'

const nav = useNavStore()
const { curtTab } = storeToRefs(nav)

interface PluginItem {
  id: string
  name: string
  description: string
  author_name: string
  author_link: string
  repo_name: string
  repo_link: string
  npm_name: string
  source: string
  integrity: string
  donate_link: string
  verified: boolean
  version: string
  local_version: string
  upgrade_available: boolean
  updated_at: string
  installed: boolean
  min_artalk_version: string
  compatible: boolean
  compatible_notice?: string
  enabled: boolean
}

const plugins = ref<PluginItem[]>([])
const search = ref('')
const onlyInstalled = ref(false)
const notification = ref('')

const data = ref<{ plugins: PluginItem[]; themes: PluginItem[] }>({
  plugins: [],
  themes: [],
})

onMounted(() => {
  nav.updateTabs(
    {
      plugins: '插件',
      themes: '主题',
    },
    'plugins',
  )

  // Plugin search
  nav.enableSearch(
    (value: string) => {
      search.value = value
      fetchList()
    },
    () => {
      if (search.value === '') return
      search.value = ''
      fetchList()
    },
  )

  fetchList()

  watch(curtTab, () => fetchList())
  watch(onlyInstalled, () => fetchList())
})

const listLoading = ref(false)
async function fetchList() {
  if (listLoading.value) return
  listLoading.value = true
  await getArtalk()!
    .ctx.getApi()
    .plugins.getPlugins({
      search: search.value,
      only_installed: onlyInstalled.value,
    })
    .then((res) => {
      data.value = {
        plugins: res.data.plugins,
        themes: res.data.themes,
      }

      if (curtTab.value === 'plugins') {
        plugins.value = data.value.plugins
      } else {
        plugins.value = data.value.themes
      }
    })
    .catch((err) => {
      showNotification(err.message)
    })
    .finally(() => {
      listLoading.value = false
    })
}

const registryUpdating = ref(false)
const registryUpdated = ref(false)
let registryUpdatedTimeout: any
function updateRegistry() {
  if (registryUpdating.value) return
  registryUpdating.value = true
  registryUpdated.value = false
  registryUpdatedTimeout && clearTimeout(registryUpdatedTimeout)
  getArtalk()!
    .ctx.getApi()
    .pluginRegistry.updatePluginRegistry({})
    .then(() => {
      fetchList().then(() => {
        registryUpdated.value = true
        registryUpdatedTimeout && clearTimeout(registryUpdatedTimeout)
        registryUpdatedTimeout = setTimeout(() => {
          registryUpdated.value = false
          registryUpdatedTimeout = null
        }, 3000)
      })
    })
    .catch((err) => {
      showNotification(err.message)
    })
    .finally(() => {
      registryUpdating.value = false
    })
}

function browsePlugin(plugin: PluginItem) {
  openLink(plugin.repo_link)
}

const pluginInstallingId = ref<string>('')

function installPlugin(plugin: PluginItem) {
  if (pluginInstallingId.value) return
  pluginInstallingId.value = plugin.id
  getArtalk()!
    .ctx.getApi()
    .plugins.installPlugin(plugin.id, {})
    .then(() => {
      plugin.installed = true
      plugin.enabled = true
    })
    .catch((err) => {
      showNotification(err.message)
    })
    .finally(() => {
      pluginInstallingId.value = ''
    })
}

function uninstallPlugin(plugin: PluginItem) {
  if (pluginInstallingId.value) return
  pluginInstallingId.value = plugin.id
  getArtalk()!
    .ctx.getApi()
    .plugins.uninstallPlugin(plugin.id, {})
    .then(() => {
      plugin.installed = false
      plugin.enabled = false
    })
    .catch((err) => {
      showNotification(err.message)
    })
    .finally(() => {
      pluginInstallingId.value = ''
    })
}

const pluginUpdatingId = ref<string>('')
async function updatePluginOptions(
  plugin: PluginItem,
  options: { enabled: boolean; client_options?: string },
) {
  if (pluginUpdatingId.value) return
  pluginUpdatingId.value = plugin.id
  return getArtalk()!
    .ctx.getApi()
    .plugins.updatePlugin(plugin.id, options)
    .then(({ data }) => {
      Object.entries(data.plugin).forEach(([key, value]) => {
        ;(plugin as any)[key] = value
      })
    })
    .catch((err) => {
      showNotification(err.message)
    })
    .finally(() => {
      pluginUpdatingId.value = ''
    })
}

const pluginUpgradingId = ref<string>('')
function upgradePlugin(plugin: PluginItem) {
  if (pluginUpgradingId.value) return
  if (!plugin.compatible) {
    showNotification(plugin.compatible_notice || 'Plugin is not compatible')
    return
  }
  pluginUpgradingId.value = plugin.id
  getArtalk()!
    .ctx.getApi()
    .plugins.upgradePlugin(plugin.id, {})
    .then(() => {
      fetchList()
    })
    .catch((err) => {
      showNotification(err.message)
    })
    .finally(() => {
      pluginUpgradingId.value = ''
    })
}

function enablePlugin(plugin: PluginItem) {
  updatePluginOptions(plugin, { enabled: true })
}

function disablePlugin(plugin: PluginItem) {
  updatePluginOptions(plugin, { enabled: false })
}

function openLink(link: string) {
  window.open(link, '_blank', 'noopener noreferrer')
}

let notificationTimeout: any

function showNotification(msg: string) {
  notificationTimeout && clearTimeout(notificationTimeout)
  notification.value = msg
  notificationTimeout = setTimeout(() => {
    notification.value = ''
    notificationTimeout = null
  }, 3000)
}

const isSemverGreater = (a: string, b: string) => {
  return a.localeCompare(b, undefined, { numeric: true }) === 1
}

const optionsEditor = reactive({
  visible: false,
  schema: {},
  options: {},
  plugin: null as PluginItem | null,
})

const closeOptionsEditor = () => {
  optionsEditor.visible = false
  optionsEditor.schema = {}
  optionsEditor.options = {}
  optionsEditor.plugin = null
}

const openOptionsEditor = (plugin: PluginItem) => {
  optionsEditor.visible = true
  optionsEditor.plugin = plugin

  getArtalk()!
    .ctx.getApi()
    .plugins.getPlugin(plugin.id)
    .then((res) => {
      const data = res.data
      console.log(data.options_schema)
      try {
        if (typeof data.options_schema === 'string')
          optionsEditor.schema = JSON.parse(data.options_schema)
        if (typeof data.client_options === 'string')
          optionsEditor.options = JSON.parse(data.client_options)
      } catch {
        void 0
      }
    })
}

const handleOptionsEditorSubmit = (value: Record<string, any>) => {
  if (optionsEditor.plugin === null) throw new Error('optionsEditor.plugin not found')

  updatePluginOptions(optionsEditor.plugin, {
    enabled: optionsEditor.plugin.enabled,
    client_options: JSON.stringify(value),
  }).then(() => {
    showNotification('Options saved')
  })
}
</script>

<template>
  <div class="atk-page-plugins">
    <AppDialog v-if="optionsEditor.visible" @close="closeOptionsEditor()">
      <PluginOptionsEditor :schema="optionsEditor.schema" :default-value="optionsEditor.options" @save="handleOptionsEditorSubmit" />
    </AppDialog>

    <div v-if="notification" class="notification" @click="notification = ''">
      {{ notification }}
    </div>

    <div class="action-bar">
      <div class="action-item">
        <input id="show-installed" v-model="onlyInstalled" type="checkbox" />
        <label for="show-installed">仅显示已安装</label>
      </div>
      <div class="action-item">
        <button class="action-btn" @click="updateRegistry()">
          <template v-if="!registryUpdated">
            <div v-if="registryUpdating" class="atk-loading-icon" />
            检查更新
          </template>
          <template v-else> 列表更新完毕 </template>
        </button>
        <button class="action-btn">设置</button>
      </div>
    </div>

    <div class="plugins">
      <div
        v-for="plugin in plugins"
        :key="plugin.id"
        class="plugin-item"
        @click="browsePlugin(plugin)"
      >
        <div class="header">
          <div class="icon">
            <span>{{ plugin.name.substring(0, 1) }}</span>
          </div>
          <div class="name">{{ plugin.name }}</div>
        </div>
        <div class="meta">
          <div class="meta-item author" @click.stop="openLink(plugin.author_link)">
            {{ plugin.author_name }}
          </div>
          <div class="meta-item version">v{{ plugin.version }}</div>
          <div
            class="meta-item updated-at"
            :datetime="plugin.updated_at"
            :title="plugin.updated_at"
          >
            {{ format(plugin.updated_at) }}
          </div>
        </div>
        <div class="description">{{ plugin.description }}</div>
        <div class="actions">
          <div
            class="plugin-id"
            :title="`${plugin.id}${plugin.upgrade_available ? ' [Upgrade Available]' : ''}`"
          >
            <template v-if="plugin.upgrade_available"
              >{{ plugin.local_version }} -&gt; {{ plugin.version }}</template
            >
            <template v-else>{{ plugin.id }}</template>
          </div>

          <span class="buttons">
            <button class="action-btn" @click.stop="openOptionsEditor(plugin)">配置</button>

            <template v-if="!plugin.installed">
              <button
                v-if="plugin.compatible"
                class="action-btn install-btn"
                @click.stop="installPlugin(plugin)"
              >
                <div v-if="plugin.id === pluginInstallingId" class="atk-loading-icon" />
                安装
              </button>
              <div v-else class="not-compatible" :title="plugin.compatible_notice">不兼容</div>
            </template>

            <template v-if="plugin.installed">
              <button
                v-if="plugin.enabled"
                class="action-btn disable-btn"
                @click.stop="disablePlugin(plugin)"
              >
                <div v-if="plugin.id === pluginUpdatingId" class="atk-loading-icon" />
                禁用
              </button>
              <button v-else class="action-btn enable-btn" @click.stop="enablePlugin(plugin)">
                <div v-if="plugin.id === pluginUpdatingId" class="atk-loading-icon" />
                启用
              </button>

              <button class="action-btn uninstall-btn" @click.stop="uninstallPlugin(plugin)">
                <div v-if="plugin.id === pluginInstallingId" class="atk-loading-icon" />
                卸载
              </button>

              <button
                v-if="plugin.upgrade_available"
                class="action-btn upgrade-btn"
                @click.stop="upgradePlugin(plugin)"
              >
                <div v-if="plugin.id === pluginUpgradingId" class="atk-loading-icon" />
                更新
              </button>
            </template>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.atk-page-plugins {
  position: relative;
  background: var(--at-color-bg-grey);
  min-height: calc(100vh - 112px);
  margin-bottom: -50px;
  padding: 20px;

  .notification {
    z-index: 100000;
    position: fixed;
    top: 130px;
    left: 50%;
    transform: translateX(-50%);
    background: var(--at-color-deep);
    color: var(--at-color-bg);
    padding: 10px 20px;
    border-radius: 4px;
  }

  .atk-loading-icon {
    width: 16px;
    height: 16px;
    margin-right: 8px;
    border-top-color: #fff;
    border-left-color: #fff;
  }

  .atk-button {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 4px 14px;
    border-radius: 3px;
    background: #36abcf;
    color: #fff;
    cursor: pointer;
    user-select: none;
    transition: background 0.2s;
    font-size: 14px;
    border: 0;

    &:not(:last-child) {
      margin-right: 10px;
    }

    &:hover {
      opacity: 0.8;
    }
  }

  .action-bar {
    padding: 10px 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: -10px;
    margin-bottom: 5px;
    position: sticky;
    top: 0;
    background: var(--at-color-bg-grey);

    .action-item {
      display: flex;
      align-items: center;

      input[type='checkbox'] {
        margin-right: 5px;
      }

      label {
        font-size: 14px;
      }

      button {
        @extend .atk-button;
      }
    }
  }

  .plugins {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    justify-content: center;
    gap: 15px;

    .plugin-item {
      height: 100%;
      display: flex;
      flex-direction: column;
      padding: 14px 18px;
      min-height: 175px;
      background: var(--at-color-bg);
      border: 1px solid var(--at-color-border);
      border-radius: 6px;
      cursor: pointer;

      .header {
        display: flex;
        align-items: center;
        flex-direction: row;
        margin-bottom: 15px;

        .icon {
          height: 30px;
          width: 30px;
          margin-right: 15px;

          span,
          img {
            border-radius: 4px;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100%;
            width: 100%;
            color: #fff;
            background: #5b6f7e;
          }
        }

        .name {
          font-size: 18px;
        }
      }

      .meta {
        display: flex;
        flex-wrap: wrap;
        justify-content: flex-start;
        font-size: 13px;

        .meta-item {
          color: var(--at-color-meta);
          background: var(--at-color-bg-grey);
          padding: 2px 10px;
          border-radius: 50px;
          margin-bottom: 5px;

          &:not(:last-child) {
            margin-right: 5px;
          }

          a {
            color: var(--at-color-meta);
          }
        }
      }

      .description {
        margin-top: 8px;
        font-size: 14px;
        color: var(--at-color-sub);
      }

      .actions {
        display: flex;
        justify-content: space-between;
        align-items: flex-end;
        margin-top: auto;
        padding-top: 10px;

        .plugin-id {
          font-size: 12px;
          color: var(--at-color-meta);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          padding-right: 10px;
        }

        .buttons {
          flex: none;
          display: flex;
          align-items: center;
          justify-content: flex-end;
          margin-left: auto;
          margin-bottom: -5px;
          margin-right: -8px;
        }

        .not-compatible {
          font-size: 12px;
          color: #fff;
          background: #fcb024;
          padding: 2px 8px;
          border-radius: 50px;
        }

        .action-btn {
          font-size: 13px;
          padding: 2px 8px;

          .atk-loading-icon {
            height: 12px;
            width: 12px;
            margin-right: 4px;
          }

          &:not(:last-child) {
            margin-right: 4px;
          }

          &.enable-btn {
            background: #36abcf;
          }

          &.disable-btn {
            background: #fcb024;
          }

          &.uninstall-btn {
            background: #fb4646;
          }

          &.upgrade-btn {
            background: #29c332;
          }

          @extend .atk-button;
        }
      }
    }
  }
}
</style>
