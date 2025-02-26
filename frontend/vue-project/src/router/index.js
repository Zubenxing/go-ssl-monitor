import { createRouter, createWebHistory } from 'vue-router'
import DomainManager from '../views/domain/DomainManager.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: DomainManager
    }
  ]
})

export default router 