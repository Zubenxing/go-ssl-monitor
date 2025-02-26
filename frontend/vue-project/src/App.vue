<script setup>
</script>

<template>
  <div class="app-container">
    <el-container>
      <el-header>
        <h1 class="app-title">SSL证书监控系统</h1>
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
      <el-footer>
        <p class="footer-text">© 2024 SSL证书监控系统</p>
      </el-footer>
    </el-container>
  </div>
</template>

<style>
body {
  margin: 0;
  padding: 0;
  background-color: #f5f7fa;
}

.app-container {
  min-height: 100vh;
}

.el-header {
  background-color: #409EFF;
  color: white;
  line-height: 60px;
  text-align: center;
  padding: 0 20px;
}

.app-title {
  margin: 0;
  font-size: 24px;
}

.el-main {
  padding: 20px;
}

.el-footer {
  text-align: center;
  color: #909399;
  padding: 20px 0;
}

.footer-text {
  margin: 0;
  font-size: 14px;
}
</style>
