<script setup>
import { ref, computed } from 'vue'
import { HomeFilled, Setting, Document, QuestionFilled } from '@element-plus/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import UserProfile from '@/components/UserProfile.vue'

const route = useRoute()
const router = useRouter()
const isCollapse = ref(false)
const isLoginPage = computed(() => route.path === '/login')
const isLoggedIn = computed(() => {
  return !!localStorage.getItem('token')
})

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  ElMessage.success('退出登录成功')
  router.push('/login')
}
</script>

<template>
  <div class="app-container" :class="{ 'login-layout': isLoginPage }">
    <router-view v-if="isLoginPage"></router-view>
    <el-container v-else>
      <!-- 侧边栏 -->
      <el-aside :width="isCollapse ? '80px' : '240px'" class="aside">
        <div class="logo-container">
          <img src="./assets/logo.svg" alt="Logo" class="logo" v-if="!isCollapse" />
          <span v-if="!isCollapse">SSL证书监控</span>
        </div>
        <el-menu
          default-active="1"
          class="el-menu-vertical"
          :collapse="isCollapse"
          background-color="#6366f1"
          text-color="#fff"
          active-text-color="#fff"
          router
        >
          <el-menu-item index="/">
            <el-icon><HomeFilled /></el-icon>
            <template #title>仪表盘</template>
          </el-menu-item>
          <el-menu-item index="/domains">
            <el-icon><Document /></el-icon>
            <template #title>域名管理</template>
          </el-menu-item>
          <el-menu-item index="/lincoln">
            <el-icon><Document /></el-icon>
            <template #title>同步备份</template>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <template #title>系统设置</template>
          </el-menu-item>
        </el-menu>
        <div class="help-section" :class="{ 'collapsed': isCollapse }">
          <el-icon><QuestionFilled /></el-icon>
          <div class="help-content" v-if="!isCollapse">
            <h4>帮助中心</h4>
            <p>遇到问题？请联系我们获取更多帮助</p>
            <el-button type="primary" class="contact-btn">联系支持</el-button>
          </div>
        </div>
      </el-aside>

      <el-container>
        <!-- 顶部导航 -->
        <el-header class="header">
          <div class="header-left">
            <el-button
              type="text"
              @click="isCollapse = !isCollapse"
              class="collapse-btn"
            >
              <el-icon>
                <component :is="isCollapse ? 'Expand' : 'Fold'" />
              </el-icon>
            </el-button>
            <span class="welcome-text">欢迎使用SSL证书监控系统</span>
          </div>
          <div class="header-right">
            <UserProfile v-if="isLoggedIn" />
          </div>
        </el-header>

        <!-- 主要内容区 -->
        <el-main class="main-content">
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style>
:root {
  --primary-color: #6366f1;
  --text-color: #374151;
  --bg-color: #f3f4f6;
  --header-height: 60px;
  --aside-width: 240px;
  --gradient-start: #4e44bd;
  --gradient-middle: #8F7FF7;
  --gradient-end: #4cadc5;
}

body {
  margin: 0;
  padding: 0;
  background-color: var(--bg-color);
  color: var(--text-color);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.app-container {
  min-height: 100vh;
}

/* 登录页面布局 */
.login-layout {
  height: 100vh;
  overflow: hidden;
}

.login-layout .aside,
.login-layout .header,
.login-layout .main-content {
  display: none !important;
}

/* 侧边栏样式 */
.aside {
  background: linear-gradient(180deg, 
    var(--gradient-start) 0%,
    var(--gradient-middle) 45%,
    var(--gradient-end) 100%
  );
  transition: width 0.3s;
  position: fixed;
  height: 100vh;
  z-index: 1000;
  box-shadow: 4px 0 10px rgba(0, 0, 0, 0.05);
}

.logo-container {
  height: var(--header-height);
  display: flex;
  align-items: center;
  padding: 0 24px;
  color: white;
  font-size: 18px;
  font-weight: bold;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(5px);
}

.logo {
  width: 32px;
  height: 32px;
  margin-right: 10px;
  filter: brightness(0) invert(1);
}

.el-menu-vertical {
  border-right: none !important;
  background: transparent !important;
}

.el-menu-vertical:not(.el-menu--collapse) {
  width: var(--aside-width);
  padding: 12px;
}

:deep(.el-menu-item) {
  height: 50px;
  line-height: 50px;
  background: transparent !important;
  margin: 8px 0;
  border-radius: 8px;
}

:deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.15) !important;
}

:deep(.el-menu-item.is-active) {
  background: rgba(255, 255, 255, 0.2) !important;
  border-radius: 8px;
  margin-right: 0;
}

:deep(.el-menu-item .el-icon) {
  font-size: 20px;
  margin-right: 16px;
}

:deep(.el-menu--collapse .el-menu-item .el-icon) {
  margin-right: 0;
  font-size: 22px;
}

/* 帮助区域样式 */
.help-section {
  position: absolute;
  bottom: 20px;
  left: 0;
  right: 0;
  padding: 20px;
  color: white;
  transition: all 0.3s;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(5px);
  margin: 0 12px;
  border-radius: 12px;
}

.help-section.collapsed {
  padding: 10px;
  text-align: center;
}

.help-content h4 {
  margin: 0 0 10px 0;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
}

.help-content p {
  margin: 0 0 15px 0;
  font-size: 12px;
  opacity: 0.8;
}

.contact-btn {
  width: 100%;
  background-color: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white !important;
}

.contact-btn:hover {
  background-color: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.4);
}

/* 顶部导航样式 */
.header {
  background-color: white;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: var(--header-height);
  position: fixed;
  top: 0;
  right: 0;
  left: var(--aside-width);
  z-index: 999;
  transition: left 0.3s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  padding: 0;
  margin-right: 20px;
  font-size: 20px;
  color: var(--text-color);
}

.welcome-text {
  font-size: 16px;
  color: var(--text-color);
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: var(--bg-color);
}

.user-avatar {
  background-color: var(--primary-color);
  margin-right: 8px;
}

.username {
  font-size: 14px;
  color: var(--text-color);
}

/* 主要内容区样式 */
.main-content {
  margin-left: var(--aside-width);
  margin-top: var(--header-height);
  padding: 20px;
  min-height: calc(100vh - var(--header-height));
  transition: margin-left 0.3s;
}

/* 响应侧边栏折叠 */
.el-container:has(.aside:not([style*="width: 240px"])) {
  .header {
    left: 80px;
  }
  
  .main-content {
    margin-left: 80px;
  }
}
</style>
