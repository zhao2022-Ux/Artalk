import type { MessageSchema } from '../i18n'

export const zhCN: MessageSchema = {
  ctrlCenter: '控制中心',
  msgCenter: '通知中心',
  noContent: '无内容',
  searchHint: '搜索关键字...',
  allSites: '所有站点',
  siteManage: '站点管理',
  comment: '评论',
  page: '页面',
  user: '用户',
  site: '站点',
  transfer: '迁移',
  settings: '设置',
  all: '全部',
  pending: '待审',
  personal: '个人',
  mentions: '提及',
  mine: '我的',
  admin: '管理员',
  create: '新增',
  import: '导入',
  export: '导出',
  settingSaved: '配置已保存',
  settingSaveFailed: '配置保存失败',
  settingNotice: '注：某些配置项可能需手动重启才能生效',
  apply: '应用',
  updateComplete: '更新完毕',
  updateReady: '开始更新...',
  opFailed: '操作失败',
  updateTitle: '更新标题',
  uploading: '上传中',
  cancel: '取消',
  back: '返回',
  cacheClear: '缓存清除',
  cacheWarm: '缓存预热',
  editTitle: '标题修改',
  switchKey: 'KEY 变更',
  commentAllowAll: '所有人可评',
  commentOnlyAdmin: '仅管理员可评',
  config: '配置文件',
  envVarControlHint: '由环境变量 {key} 控制',
  userAdminHint: '该用户具有管理员权限',
  userInConfHint: '该用户存在于配置文件中',
  edit: '编辑',
  delete: '删除',
  siteCount: '共 {count} 个站点',
  createSite: '新增站点',
  siteName: '站点名称',
  siteUrls: '站点 URLs',
  multiSepHint: '多个用逗号隔开',
  add: '新增',
  rename: '重命名',
  inputHint: '输入内容...',
  userCreate: '用户创建',
  userEdit: '用户编辑',
  userInConfCannotEditHint: '暂不支持在线编辑配置文件中的用户，请手动修改配置文件',
  userDeleteConfirm:
    '该操作将删除 用户："{name}" 邮箱："{email}" 所有评论，包括其评论下面他人的回复评论，是否继续？',
  userDeleteManuallyHint: '用户已从数据库删除，请手动编辑配置文件并删除用户',
  pageDeleteConfirm: '确认删除页面 "{title}"？将会删除所有相关数据',
  siteDeleteConfirm: '确认删除站点 "{name}"？将会删除所有相关数据',
  siteNameInputHint: '请输入站点名称',
  comments: '评论',
  last: '近期',
  show: '展开',
  username: '用户名',
  email: '邮箱',
  link: '链接',
  badgeText: '徽章文字',
  badgeColor: '徽章颜色',
  role: '身份角色',
  normal: '普通',
  password: '密码',
  passwordEmptyHint: '留空不修改密码',
  emailNotify: '邮件通知',
  enabled: '开启',
  disabled: '关闭',
  save: '保存',
  dataFile: '数据文件',
  artransfer: '转换工具',
  targetSiteName: '目标站点名',
  targetSiteURL: '目标站点 URL',
  payload: '启动参数',
  optional: '可选',
  uploadReadyToImport: '文件已成功上传，可以开始导入',
  artransferToolHint: '使用 {link} 将评论数据转为 Artrans 格式',
  moreDetails: '查看详情',
  loginFailure: '登录发生异常',
  login: '登录',
  logout: '退出登录',
  logoutConfirm: '确定要退出登录吗？',
  loginSelectHint: '请选择您要登录的账户：',
  plugins: '插件',
}

export default zhCN
