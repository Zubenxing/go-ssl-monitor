<template>
  <div class="lincoln-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <h2>备份状态监控</h2>
          <el-button type="primary" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-table 
        v-loading="loading"
        :data="backupLogs"
        style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="server_name" label="服务器名称" width="180" />
        <el-table-column label="开始时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.start_time) }}
          </template>
        </el-table-column>
        <el-table-column label="结束时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.end_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="script_version" label="脚本版本" width="120" />
        <el-table-column label="备份状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.backup_status === 0 ? 'success' : 'danger'">
              {{ scope.row.backup_status === 0 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="告警状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.alert_status === 0 ? 'success' : 'warning'">
              {{ scope.row.alert_status === 0 ? '正常' : '告警' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import axios from 'axios'

const loading = ref(false)
const backupLogs = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

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
    total.value = response.data.length
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
  fetchData()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchData()
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  return time
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.lincoln-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
  font-size: 18px;
  color: var(--text-color);
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 