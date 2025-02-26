<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-box">
        <div class="login-header">
          <img src="../assets/logo.svg" alt="Logo" class="login-logo">
          <h2 class="login-title">SSL证书监控系统</h2>
          <p class="login-subtitle">安全监控您的网站证书</p>
        </div>
        
        <el-form :model="loginForm" :rules="rules" ref="loginFormRef" class="login-form">
          <el-form-item prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="请输入用户名"
              size="large">
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input 
              v-model="loginForm.password" 
              type="password" 
              placeholder="请输入密码"
              show-password
              size="large"
              @keyup.enter="handleLogin">
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <div class="login-options">
            <el-checkbox v-model="rememberMe">记住用户名</el-checkbox>
            <a href="javascript:void(0)" @click="forgotPassword" class="forgot-password">忘记密码?</a>
          </div>
          
          <el-button 
            type="primary" 
            @click="handleLogin" 
            class="login-button"
            :loading="loading"
            size="large">
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form>
        
        <div class="login-footer">
          <p>© {{ new Date().getFullYear() }} SSL证书监控系统</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const loginForm = reactive({
  username: '',
  password: ''
})

const loading = ref(false)
const rememberMe = ref(false)
const loginFormRef = ref()

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const response = await axios.post('/api/auth/login', loginForm)
        const { token } = response.data
        
        localStorage.setItem('token', token)
        if (rememberMe.value) {
          localStorage.setItem('username', loginForm.username)
        } else {
          localStorage.removeItem('username')
        }
        
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const forgotPassword = () => {
  ElMessage.info('密码重置功能开发中，敬请期待...')
}

onMounted(() => {
  const savedUsername = localStorage.getItem('username')
  if (savedUsername) {
    loginForm.username = savedUsername
    rememberMe.value = true
  }
})
</script>

<style>
.login-page {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(-45deg, #6366f1, #8b5cf6, #4cadc5, #06b6d4);
  background-size: 400% 400%;
  animation: gradient 15s ease infinite;
  z-index: 9999;
}

@keyframes gradient {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.login-container {
  width: 100%;
  max-width: 480px;
  margin: 0 auto;
  padding: 0 20px;
}

.login-box {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 40px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-logo {
  width: 80px;
  height: 80px;
  margin-bottom: 20px;
  filter: brightness(0) invert(1);
}

.login-title {
  color: white;
  font-size: 28px;
  margin: 0 0 10px;
  font-weight: 600;
}

.login-subtitle {
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-size: 16px;
}

.login-form {
  margin-bottom: 30px;
}

:deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.9) !important;
  border: none !important;
  box-shadow: none !important;
  border-radius: 8px !important;
  height: 48px;
  transition: all 0.3s ease;
}

:deep(.el-input__wrapper:hover) {
  background: rgba(255, 255, 255, 1) !important;
}

:deep(.el-input__inner) {
  height: 48px !important;
  font-size: 15px !important;
  color: #1f2937 !important;
}

:deep(.el-input__prefix-inner) {
  font-size: 18px;
  color: #6366f1;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 20px 0;
}

:deep(.el-checkbox__label) {
  color: rgba(255, 255, 255, 0.9) !important;
  font-size: 14px;
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: white !important;
  border-color: white !important;
}

:deep(.el-checkbox__inner::after) {
  border-color: #6366f1 !important;
}

:deep(.el-checkbox__inner) {
  background-color: rgba(255, 255, 255, 0.2) !important;
  border-color: rgba(255, 255, 255, 0.4) !important;
}

.forgot-password {
  color: rgba(255, 255, 255, 0.9);
  text-decoration: none;
  font-size: 14px;
  transition: all 0.2s ease;
}

.forgot-password:hover {
  color: white;
  text-decoration: underline;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
  background: white !important;
  color: #6366f1 !important;
  border: none;
  border-radius: 8px;
  margin-top: 20px;
  transition: all 0.3s ease;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
}

.login-button:active {
  transform: translateY(0);
}

.login-footer {
  text-align: center;
  margin-top: 30px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
}

@media (max-width: 480px) {
  .login-container {
    padding: 0 15px;
  }
  
  .login-box {
    padding: 30px 20px;
  }

  .login-logo {
    width: 60px;
    height: 60px;
  }

  .login-title {
    font-size: 24px;
  }
}
</style> 