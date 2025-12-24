import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '../utils/auth'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import OAuthCallback from '../views/OAuthCallback.vue'
import Admin from '../views/Admin.vue'
import AdminTiers from '../views/AdminTiers.vue'
import AdminSettings from '../views/AdminSettings.vue'
import NotFound from '../views/NotFound.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/callback',
    name: 'callback',
    component: OAuthCallback
  },
  {
    path: '/admin',
    name: 'admin',
    component: Admin
  },
  {
    path: '/admin/tiers',
    name: 'admin-tiers',
    component: AdminTiers
  },
  {
    path: '/admin/settings',
    name: 'admin-settings',
    component: AdminSettings
  },
  {
    // 捕获所有未定义的路径
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫（可选，根据需要启用）
// router.beforeEach((to, from, next) => {
//   const token = getToken()
//   if (!token && to.path !== '/login' && to.path !== '/callback') {
//     next('/login')
//   } else {
//     next()
//   }
// })

export default router
