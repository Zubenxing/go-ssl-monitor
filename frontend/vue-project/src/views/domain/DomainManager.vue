<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import moment from 'moment'

const API_BASE_URL = 'http://localhost:8080/api'
const MAX_RETRIES = 3
const RETRY_DELAY = 1000 // 1 second

// 数据状态
const domains = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const filterStatus = ref('')
const addingDomain = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const isUpdating = ref(false)
const newDomain = ref({
  domainName: '',
  notificationEmail: '',
  autoRenewal: true
})

// 工具函数
const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))

const loadDomains = async (retryCount = 0) => {
  loading.value = true
  try {
    const response = await axios.get(`${API_BASE_URL}/domains`)
    domains.value = response.data
  } catch (error) {
    console.error('Failed to load domains:', error)
    if (error.message.includes('Network Error') && retryCount < MAX_RETRIES) {
      ElMessage.warning(`连接服务器失败，${retryCount + 1}秒后重试...`)
      await sleep(RETRY_DELAY)
      return loadDomains(retryCount + 1)
    }
    ElMessage.error(error.message.includes('Network Error') ? 
      '无法连接到服务器，请确保后端服务已启动' : '加载域名列表失败')
  } finally {
    loading.value = false
  }
}

// 计算属性
const filteredDomains = computed(() => {
  if (!filterStatus.value) return domains.value
  
  return domains.value.filter(domain => {
    if (filterStatus.value === 'EXPIRING') {
      return isExpiringSoon(domain.certificateExpiryDate)
    }
    return domain.certificateStatus === filterStatus.value
  })
})

const totalDomains = computed(() => domains.value.length)

const expiringDomains = computed(() => {
  const thirtyDaysFromNow = new Date()
  thirtyDaysFromNow.setDate(thirtyDaysFromNow.getDate() + 30)
  
  return domains.value.filter(domain => {
    if (!domain.certificateExpiryDate) return false
    const expiryDate = new Date(domain.certificateExpiryDate)
    return expiryDate <= thirtyDaysFromNow
  }).length
})

const errorDomains = computed(() => {
  return domains.value.filter(domain => domain.certificateStatus === 'ERROR').length
})

// 域名操作函数
const addDomain = async () => {
  if (!newDomain.value.domainName.trim()) {
    ElMessage.error('域名不能为空')
    return
  }
  
  if (newDomain.value.notificationEmail && !isValidEmail(newDomain.value.notificationEmail)) {
    ElMessage.error('请输入有效的邮箱地址')
    return
  }

  addingDomain.value = true
  try {
    await axios.post(`${API_BASE_URL}/domains`, newDomain.value)
    ElMessage.success('添加域名成功')
    dialogVisible.value = false
    await loadDomains()
  } catch (error) {
    const errorMessage = error.response?.data?.error || '添加域名失败'
    ElMessage.error(errorMessage)
  } finally {
    addingDomain.value = false
  }
}

const isValidEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

const deleteDomain = async (domain) => {
  try {
    await axios.delete(`${API_BASE_URL}/domains/${domain.id}`)
    ElMessage.success('删除域名成功')
    loadDomains()
  } catch (error) {
    ElMessage.error('删除域名失败')
  }
}

const checkCertificate = async (domain) => {
  try {
    await axios.post(`${API_BASE_URL}/domains/${domain.id}/check`)
    ElMessage.success('检查证书成功')
    loadDomains()
  } catch (error) {
    ElMessage.error('检查证书失败')
  }
}

const toggleAutoRenewal = async (domain) => {
  try {
    await axios.put(`${API_BASE_URL}/domains/${domain.id}/auto-renewal`)
    ElMessage.success('更新自动续期设置成功')
    loadDomains()
  } catch (error) {
    ElMessage.error('更新自动续期设置失败')
  }
}

// 辅助函数
const formatDate = (date) => {
  return date ? moment(date).format('YYYY-MM-DD HH:mm:ss') : '-'
}

const getStatusType = (status) => {
  return status === 'VALID' ? 'success' : 'danger'
}

const getStatusText = (status) => {
  return status === 'VALID' ? '正常' : '异常'
}

const isExpiringSoon = (date) => {
  if (!date) return false
  const expiryDate = moment(date)
  const daysUntilExpiry = expiryDate.diff(moment(), 'days')
  return daysUntilExpiry <= 30 && daysUntilExpiry >= 0
}

const getDaysUntilExpiry = (date) => {
  if (!date) return 0
  return moment(date).diff(moment(), 'days')
}

const getExpiryClass = (date) => {
  if (!date) return ''
  return isExpiringSoon(date) ? 'expiring-soon' : ''
}

const getRowClassName = ({ row }) => {
  return row.certificateStatus === 'ERROR' ? 'error-row' : ''
}

const showEditDialog = (domain) => {
  isEditing.value = true
  editingId.value = domain.id
  newDomain.value = {
    domainName: domain.domainName,
    notificationEmail: domain.notificationEmail,
    autoRenewal: domain.autoRenewal
  }
  dialogVisible.value = true
}

const showAddDomainDialog = () => {
  isEditing.value = false
  editingId.value = null
  newDomain.value = {
    domainName: '',
    notificationEmail: '',
    autoRenewal: true
  }
  dialogVisible.value = true
}

const updateDomain = async () => {
  if (!newDomain.value.domainName.trim()) {
    ElMessage.error('域名不能为空')
    return
  }
  
  if (newDomain.value.notificationEmail && !isValidEmail(newDomain.value.notificationEmail)) {
    ElMessage.error('请输入有效的邮箱地址')
    return
  }

  try {
    isUpdating.value = true
    loading.value = true
    await axios.put(`${API_BASE_URL}/domains/${editingId.value}`, newDomain.value)
    ElMessage.success('更新域名成功')
    dialogVisible.value = false
    await loadDomains()
  } catch (error) {
    console.error('Error updating domain:', error)
    const errorMessage = error.response?.data?.error || '更新域名失败'
    ElMessage.error(errorMessage)
  } finally {
    isUpdating.value = false
    loading.value = false
  }
}

const sendNotification = async (domain) => {
  try {
    await axios.post(`${API_BASE_URL}/domains/${domain.id}/send-notification`)
    ElMessage.success('通知邮件发送成功')
  } catch (error) {
    const errorMessage = error.response?.data?.error || '发送通知失败'
    ElMessage.error(errorMessage)
  }
}

onMounted(() => {
  loadDomains()
})
</script>

<template>
  <div class="domain-manager">
    <!-- 统计面板 -->
    <div class="dashboard">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card class="stat-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>总域名数</span>
              </div>
            </template>
            <div class="stat-value">
              <span class="number">{{ totalDomains }}</span>
              <span class="label">个域名</span>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="stat-card warning" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>即将过期</span>
                <el-tooltip
                  content="30天内即将过期的证书数量"
                  placement="top"
                >
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <div class="stat-value">
              <span class="number">{{ expiringDomains }}</span>
              <span class="label">个域名</span>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="stat-card error" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>异常状态</span>
              </div>
            </template>
            <div class="stat-value">
              <span class="number">{{ errorDomains }}</span>
              <span class="label">个域名</span>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 域名列表 -->
    <el-card class="domain-list">
      <template #header>
        <div class="card-header">
          <div class="left">
            <span class="title">域名证书监控</span>
            <el-select
              v-model="filterStatus"
              placeholder="证书状态"
              clearable
              class="filter-select"
            >
              <el-option label="全部" value="" />
              <el-option label="正常" value="VALID" />
              <el-option label="异常" value="ERROR" />
            </el-select>
          </div>
          <el-button type="primary" @click="showAddDomainDialog">添加域名</el-button>
        </div>
      </template>

      <el-table 
        :data="filteredDomains" 
        style="width: 100%" 
        v-loading="loading"
        :row-class-name="getRowClassName"
      >
        <el-table-column prop="domainName" label="域名" min-width="160" />
        <el-table-column prop="notificationEmail" label="通知邮箱" min-width="160" />
        <el-table-column prop="certificateStatus" label="证书状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.certificateStatus)">
              {{ getStatusText(scope.row.certificateStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="certificateExpiryDate" label="到期时间" width="160">
          <template #default="scope">
            <div :class="getExpiryClass(scope.row.certificateExpiryDate)">
              {{ formatDate(scope.row.certificateExpiryDate) }}
              <el-tag 
                v-if="isExpiringSoon(scope.row.certificateExpiryDate)" 
                size="small" 
                type="warning"
                class="expiry-tag"
              >
                {{ getDaysUntilExpiry(scope.row.certificateExpiryDate) }}天后过期
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="lastChecked" label="最后检查时间" width="160">
          <template #default="scope">
            {{ formatDate(scope.row.lastChecked) }}
          </template>
        </el-table-column>
        <el-table-column prop="autoRenewal" label="自动续期" width="80" align="center">
          <template #default="scope">
            <el-switch
              v-model="scope.row.autoRenewal"
              @change="toggleAutoRenewal(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="360" align="center">
          <template #default="scope">
            <div class="operation-buttons">
              <el-button size="small" @click="checkCertificate(scope.row)">
                检查证书
              </el-button>
              <el-button
                size="small"
                type="warning"
                @click="sendNotification(scope.row)"
                :disabled="!scope.row.certificateExpiryDate"
              >
                发送通知
              </el-button>
              <el-button
                size="small"
                type="primary"
                @click="showEditDialog(scope.row)"
              >
                编辑
              </el-button>
              <el-button
                size="small"
                type="danger"
                @click="deleteDomain(scope.row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑域名对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEditing ? '编辑域名' : '添加域名'" 
      :close-on-click-modal="false"
    >
      <el-form :model="newDomain" label-width="120px">
        <el-form-item label="域名" required>
          <el-input 
            v-model="newDomain.domainName" 
            :placeholder="isEditing ? '' : '请输入域名（如：example.com）'"
            :disabled="isUpdating || addingDomain"
          />
        </el-form-item>
        <el-form-item label="通知邮箱">
          <el-input 
            v-model="newDomain.notificationEmail" 
            placeholder="接收证书过期通知的邮箱"
            :disabled="isUpdating || addingDomain"
          />
        </el-form-item>
        <el-form-item label="自动续期">
          <el-switch 
            v-model="newDomain.autoRenewal"
            :disabled="isUpdating || addingDomain"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button 
            @click="dialogVisible = false" 
            :disabled="isUpdating || addingDomain"
          >
            取消
          </el-button>
          <el-button 
            type="primary" 
            @click="isEditing ? updateDomain() : addDomain()"
            :loading="isUpdating || addingDomain"
          >
            {{ isUpdating || addingDomain ? '处理中...' : '确定' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.domain-manager {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.dashboard {
  margin-bottom: 20px;
}

.stat-card {
  height: 160px;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title {
  font-size: 16px;
  font-weight: bold;
}

.filter-select {
  width: 120px;
}

.stat-value {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 80px;
}

.number {
  font-size: 36px;
  font-weight: bold;
  color: #409EFF;
}

.warning .number {
  color: #E6A23C;
}

.error .number {
  color: #F56C6C;
}

.label {
  margin-top: 8px;
  color: #909399;
  font-size: 14px;
}

.domain-list {
  margin-top: 20px;
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: nowrap;
}

.expiring-soon {
  color: #E6A23C;
}

.expiry-tag {
  margin-left: 8px;
}

:deep(.error-row) {
  --el-table-tr-bg-color: var(--el-color-danger-light-9);
}

:deep(.warning-row) {
  --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.dialog-footer {
  margin-top: 20px;
}
</style> 