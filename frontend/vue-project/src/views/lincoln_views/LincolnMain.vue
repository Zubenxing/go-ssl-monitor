<template>
  <div class="lincoln-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="left-section">
            <h2>备份状态监控</h2>
            <div class="search-section">
              <el-input
                v-model="searchQuery"
                placeholder="搜索IP/服务器/版本"
                class="search-input"
                clearable
                @clear="handleSearchClear"
                @input="handleSearch">
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-date-picker
                v-model="dateRange"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                value-format="YYYY-MM-DD HH:mm:ss"
                :shortcuts="dateShortcuts"
                @change="handleDateChange"
              />
            </div>
          </div>
          <el-button type="primary" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-table 
        v-loading="loading"
        :data="paginatedData"
        style="width: 100%">
        <el-table-column 
          prop="id" 
          label="ID" 
          width="100"
          sortable
          :sort-method="sortById" />
        <el-table-column prop="ip" label="IP地址" width="160" />
        <el-table-column prop="server_name" label="服务器名称" min-width="200" />
        <el-table-column label="开始时间" width="180" sortable :sort-method="sortByStartTime">
          <template #default="scope">
            {{ formatTime(scope.row.start_time) }}
          </template>
        </el-table-column>
        <el-table-column label="结束时间" width="180" sortable :sort-method="sortByEndTime">
          <template #default="scope">
            {{ formatTime(scope.row.end_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="script_version" label="脚本版本" width="120" align="center" />
        <el-table-column 
          label="备份状态" 
          width="100" 
          align="center"
          :filters="[
            { text: '正常', value: 0 },
            { text: '异常', value: 1 }
          ]"
          :filter-method="filterBackupStatus"
          filter-placement="bottom">
          <template #default="scope">
            <el-tag :type="scope.row.backup_status === 0 ? 'success' : 'danger'">
              {{ scope.row.backup_status === 0 ? '正常' : '异常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column 
          label="告警状态" 
          width="120" 
          align="center"
          :filters="[
            { text: '正常', value: 0 },
            { text: '告警已触发', value: 1 },
            { text: '告警未触发', value: 2 }
          ]"
          :filter-method="filterAlertStatus"
          filter-placement="bottom">
          <template #default="scope">
            <el-tag :type="getAlertStatusType(scope.row.alert_status, scope.row.backup_status)">
              {{ getAlertStatusText(scope.row.alert_status, scope.row.backup_status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Search } from '@element-plus/icons-vue'
import axios from 'axios'

const loading = ref(false)
const backupLogs = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')
const dateRange = ref(null)
const backupStatusFilter = ref([])
const alertStatusFilter = ref([])

// 日期快捷选项
const dateShortcuts = [
  {
    text: '最近一天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24)
      // 格式化日期为字符串
      return [
        formatDateToString(start),
        formatDateToString(end)
      ]
    },
  },
  {
    text: '最近三天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 3)
      return [
        formatDateToString(start),
        formatDateToString(end)
      ]
    },
  },
  {
    text: '最近一周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [
        formatDateToString(start),
        formatDateToString(end)
      ]
    },
  },
]

// 添加日期格式化辅助函数
const formatDateToString = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 过滤后的数据
const filteredData = computed(() => {
  let result = [...backupLogs.value]
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(item => 
      item.ip.toLowerCase().includes(query) ||
      item.server_name.toLowerCase().includes(query) ||
      item.script_version.toLowerCase().includes(query)
    )
  }
  
  // 修改日期范围过滤逻辑
  if (dateRange.value && dateRange.value[0] && dateRange.value[1]) {
    const selectedStartDate = new Date(dateRange.value[0])
    const selectedEndDate = new Date(dateRange.value[1])
    
    result = result.filter(item => {
      // 如果记录没有开始时间，则跳过
      if (!item.start_time) return false
      
      const recordStartTime = new Date(item.start_time)
      const recordEndTime = item.end_time ? new Date(item.end_time) : null
      
      // 如果记录有结束时间
      if (recordEndTime) {
        // 检查记录的时间段是否与选择的时间段有重叠
        return (
          // 记录的开始时间在选择的时间范围内
          (recordStartTime >= selectedStartDate && recordStartTime <= selectedEndDate) ||
          // 记录的结束时间在选择的时间范围内
          (recordEndTime >= selectedStartDate && recordEndTime <= selectedEndDate) ||
          // 记录的时间段完全包含选择的时间范围
          (recordStartTime <= selectedStartDate && recordEndTime >= selectedEndDate)
        )
      }
      
      // 如果记录没有结束时间，只检查开始时间是否在范围内
      return recordStartTime >= selectedStartDate && recordStartTime <= selectedEndDate
    })
  }

  // 备份状态筛选
  if (backupStatusFilter.value.length > 0) {
    result = result.filter(item => backupStatusFilter.value.includes(item.backup_status))
  }

  // 告警状态筛选
  if (alertStatusFilter.value.length > 0) {
    result = result.filter(item => alertStatusFilter.value.includes(item.alert_status))
  }
  
  return result
})

// 分页数据
const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})

// 总数（只定义一次）
const total = computed(() => filteredData.value.length)

// 排序方法
const sortById = (a, b) => a.id - b.id
const sortByStartTime = (a, b) => new Date(a.start_time) - new Date(b.start_time)
const sortByEndTime = (a, b) => new Date(a.end_time || 0) - new Date(b.end_time || 0)

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
}

const handleSearchClear = () => {
  searchQuery.value = ''
  handleSearch()
}

// 日期范围变化处理
const handleDateChange = () => {
  currentPage.value = 1
}

const fetchData = async () => {
  loading.value = true
  try {
    console.log('开始获取备份日志...')
    const token = localStorage.getItem('token')
    const response = await axios.get('/api/backupLogs', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    console.log('获取到的数据:', response.data)
    backupLogs.value = response.data
  } catch (error) {
    console.error('获取备份日志详细错误:', error)
    ElMessage.error(`获取备份日志失败: ${error.message}`)
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  // 确保时间格式一致
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).replace(/\//g, '-')
}

// 修改筛选方法
const filterBackupStatus = (value, row, column) => {
  backupStatusFilter.value = column.filteredValue || []
  return true  // 实际筛选由 filteredData 处理
}

// 获取告警状态类型
const getAlertStatusType = (alertStatus, backupStatus) => {
  // 如果备份正常，显示成功状态
  if (backupStatus === 0) {
    return 'success'
  }
  
  // 备份异常时，根据告警状态显示不同类型
  switch (alertStatus) {
    case 0: // 正常
      return 'success'
    case 1: // 告警已触发（邮件发送成功）
      return 'warning'
    case 2: // 告警未触发（邮件发送失败）
      return 'danger'
    default:
      return 'info'
  }
}

// 获取告警状态文本
const getAlertStatusText = (alertStatus, backupStatus) => {
  // 如果备份正常，显示正常状态
  if (backupStatus === 0) {
    return '正常'
  }
  
  // 备份异常时，根据告警状态显示不同文本
  switch (alertStatus) {
    case 0:
      return '正常'
    case 1:
      return '告警已触发'
    case 2:
      return '告警未触发'
    default:
      return '未知'
  }
}

// 告警状态筛选方法
const filterAlertStatus = (value, row, column) => {
  alertStatusFilter.value = column.filteredValue || []
  return true  // 实际筛选由 filteredData 处理
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.lincoln-container {
  padding: 20px;
  max-width: 1800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.left-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.search-section {
  display: flex;
  gap: 10px;
  align-items: center;
}

.search-input {
  width: 220px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  padding: 10px 0;
}

:deep(.el-table) {
  margin-bottom: 20px;
}

:deep(.el-table th) {
  background-color: #f5f7fa;
  font-weight: bold;
}

:deep(.el-tag) {
  min-width: 90px;  /* 增加宽度以适应更长的文本 */
  text-align: center;
  font-weight: bold;
}

:deep(.el-tag--success) {
  background-color: #f0f9eb;
  border-color: #e1f3d8;
  color: #67c23a;
}

:deep(.el-tag--warning) {
  background-color: #fdf6ec;
  border-color: #faecd8;
  color: #e6a23c;
}

:deep(.el-tag--danger) {
  background-color: #fef0f0;
  border-color: #fde2e2;
  color: #f56c6c;
}

/* 优化筛选图标样式 */
:deep(.el-table__column-filter-trigger) {
  margin-left: 4px;
}

:deep(.el-table__column-filter-trigger i) {
  color: #909399;
}

:deep(.el-table__column-filter-trigger:hover i) {
  color: #409EFF;
}
</style> 