import { createRouter, createWebHashHistory } from 'vue-router'

import menus from './menus'

const router = createRouter({
  history: createWebHashHistory(),
  routes: menus,
})

export default router
