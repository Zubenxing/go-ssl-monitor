import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建 axios 实例
const instance = axios.create({
  baseURL: '/api',
  timeout: 5000
})

// 请求拦截器
instance.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
instance.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // 清除 token 并跳转登录页
          localStorage.removeItem('token')
          localStorage.removeItem('username')
          router.push('/login')
          ElMessage.error('登录已过期，请重新登录')
          break
        default:
          ElMessage.error(error.response.data.error || '请求失败')
      }
    }
    return Promise.reject(error)
  }
)

export default instance 