import { createRouter, createWebHistory } from 'vue-router'
import DomainManager from '../views/domain/DomainManager.vue'
import Login from '../views/Login.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { 
        requiresAuth: false,
        layout: 'blank'
      }
    },
    {
      path: '/',
      name: 'home',
      component: DomainManager,
      meta: { 
        requiresAuth: true,
        layout: 'default'
      }
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !token) {
    // 需要认证但没有token，重定向到登录页
    next({ name: 'login' })
  } else if (to.name === 'login' && token) {
    // 已登录用户访问登录页，重定向到首页
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router 