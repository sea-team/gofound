
const menus = [
  {
    path: '/',
    name: 'dashboard',
    label:"数据库",
    component: () => import('./views/dashboard.vue'),
  },
  {
    path: '/document',
    name: 'document',
    label:"帮助文档",
    component: () => import('./views/document.vue'),
  },
  {
    path: '/status',
    name: 'status',
    label:"服务器状态",
    component: () => import('./views/status.vue'),
  },
]
export default menus
