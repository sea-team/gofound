import { Coin, DataLine, Document } from '@element-plus/icons-vue'

const menus = [
  {
    path: '/',
    name: 'dashboard',
    icon: Coin,
    label: '数据库',
    color: 'rgb(105, 192, 255)',
    component: () => import('./views/dashboard.vue'),
  }, {
    path: '/status',
    name: 'status',
    label: '服务器状态',
    color: 'rgb(149, 222, 100)',
    icon: DataLine,
    component: () => import('./views/status.vue'),
  }, {
    path: '/document',
    name: 'document',
    label: '帮助文档',
    icon: Document,
    color: 'rgb(255, 156, 110)',
    component: () => import('./views/document.vue'),
  },
]
export default menus
