<template>
  <div class="user-profile">
    <el-dropdown trigger="click" @command="handleCommand">
      <div class="user-info">
        <el-avatar :size="32" :src="userInfo.avatar || defaultAvatar" />
        <span class="username">{{ userInfo.username }}</span>
        <el-icon class="el-icon--right"><arrow-down /></el-icon>
      </div>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="profile">
            <el-icon><User /></el-icon>个人信息
          </el-dropdown-item>
          <el-dropdown-item command="changePassword">
            <el-icon><Lock /></el-icon>修改密码
          </el-dropdown-item>
          <el-dropdown-item divided command="logout">
            <el-icon><SwitchButton /></el-icon>退出登录
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <!-- 个人信息对话框 -->
    <el-dialog
      v-model="profileDialogVisible"
      title="个人信息"
      width="400px">
      <div class="profile-content">
        <div class="avatar-container">
          <el-avatar :size="80" :src="userInfo.avatar || defaultAvatar" />
          <el-upload
            class="avatar-uploader"
            action="/api/user/avatar"
            :headers="uploadHeaders"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload">
            <el-button size="small">更换头像</el-button>
          </el-upload>
        </div>
        <el-form :model="userInfo" label-width="80px">
          <el-form-item label="用户名">
            <el-input v-model="userInfo.username" disabled />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="userInfo.email" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="profileDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="updateProfile">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="passwordDialogVisible"
      title="修改密码"
      width="400px">
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px">
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input
            v-model="passwordForm.currentPassword"
            type="password"
            show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            show-password />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="passwordDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="updatePassword">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User, Lock, SwitchButton, ArrowDown } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import axios from '@/utils/axios'

const router = useRouter()
const defaultAvatar = '/default-avatar.png'

const userInfo = reactive({
  username: localStorage.getItem('username') || '',
  email: '',
  avatar: ''
})

const profileDialogVisible = ref(false)
const passwordDialogVisible = ref(false)
const passwordFormRef = ref(null)

// 密码表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码验证规则
const passwordRules = {
  currentPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 上传头像的请求头
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 处理下拉菜单命令
const handleCommand = (command) => {
  console.log('Command:', command) // 添加日志
  switch (command) {
    case 'profile':
      profileDialogVisible.value = true
      break
    case 'changePassword':
      passwordDialogVisible.value = true
      break
    case 'logout':
      handleLogout()
      break
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await axios.post('/auth/logout')  // 可选：调用后端登出接口
  } catch (error) {
    console.error('Logout error:', error)
  } finally {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    router.push('/login')
  }
}

// 更新个人信息
const updateProfile = async () => {
  try {
    const response = await axios.put('/api/user/profile', {
      email: userInfo.email
    }, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    
    if (response.data.success) {
      ElMessage.success('个人信息更新成功')
      profileDialogVisible.value = false
    }
  } catch (error) {
    ElMessage.error('更新失败：' + error.message)
  }
}

// 更新密码
const updatePassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const response = await axios.put('/api/user/password', {
          currentPassword: passwordForm.currentPassword,
          newPassword: passwordForm.newPassword
        }, {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
        
        if (response.data.success) {
          ElMessage.success('密码修改成功，请重新登录')
          handleLogout()
        }
      } catch (error) {
        ElMessage.error('密码修改失败：' + error.message)
      }
    }
  })
}

// 头像上传相关方法
const handleAvatarSuccess = (response) => {
  if (response.success) {
    userInfo.avatar = response.data.url
    ElMessage.success('头像上传成功')
  }
}

const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('头像只能是 JPG 或 PNG 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('头像大小不能超过 2MB!')
  }
  return isJPG && isLt2M
}

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const response = await axios.get('/api/user/profile', {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    
    if (response.data.success) {
      Object.assign(userInfo, response.data.data)
    }
  } catch (error) {
    console.error('获取用户信息失败：', error)
  }
}

// 组件挂载时获取用户信息
onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.user-profile {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 12px;
}

.username {
  margin: 0 8px;
  font-size: 14px;
  color: #606266;
}

.profile-content {
  padding: 20px;
}

.avatar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 20px;
  gap: 10px;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dialog-footer {
  padding-top: 20px;
  text-align: right;
}
</style> 