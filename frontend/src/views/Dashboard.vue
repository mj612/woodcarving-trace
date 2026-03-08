<template>
  <div class="dashboard">
    <div class="page-container">
      <div class="welcome-card card">
        <h2>欢迎使用木雕工艺品溯源系统</h2>
        <p>当前登录用户：{{ userInfo.realName || userInfo.username }} ({{ roleMap[userInfo.role] }})</p>
      </div>
      
      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stats-row" v-loading="loading">
        <el-col :span="6">
          <div class="stat-card card">
            <div class="stat-number">{{ stats.materialCount || 0 }}</div>
            <div class="stat-label">原料总数</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card card">
            <div class="stat-number">{{ stats.productCount || 0 }}</div>
            <div class="stat-label">产品总数</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card card">
            <div class="stat-number">{{ stats.transactionCount || 0 }}</div>
            <div class="stat-label">交易记录</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card card">
            <div class="stat-number">{{ stats.userCount || 0 }}</div>
            <div class="stat-label">用户数量</div>
          </div>
        </el-col>
      </el-row>
      
      <!-- 快捷操作 -->
      <div class="quick-actions card">
        <h3>快捷操作</h3>
        <el-row :gutter="20">
          <el-col :span="8" v-for="action in quickActions" :key="action.path">
            <el-button 
              type="primary" 
              plain 
              @click="$router.push(action.path)"
              class="action-btn"
            >
              {{ action.label }}
            </el-button>
          </el-col>
        </el-row>
      </div>
      
      <!-- 最近活动 -->
      <div class="recent-activity card">
        <h3>最近活动</h3>
        <el-table :data="recentActivities" style="width: 100%">
          <el-table-column prop="time" label="时间" width="180"></el-table-column>
          <el-table-column prop="action" label="操作"></el-table-column>
          <el-table-column prop="user" label="操作人" width="120"></el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import { getStats, getRecentActivities } from '@/api'

export default {
  name: 'Dashboard',
  data() {
    return {
      stats: {
        materialCount: 0,
        productCount: 0,
        transactionCount: 0,
        userCount: 0
      },
      recentActivities: [],
      loading: false
    }
  },
  computed: {
    ...mapState(['userInfo']),
    roleMap() {
      return {
        'supplier': '原料供应商',
        'artisan': '雕刻工匠',
        'warehouse': '仓库管理员',
        'seller': '销售商',
        'consumer': '消费者',
        'supervisor': '管理员'
      }
    },
    quickActions() {
      const actions = []
      
      switch(this.userInfo.role) {
        case 'supplier':
          actions.push(
            { label: '添加原料', path: '/dashboard/materials/create' }
          )
          break
        case 'artisan':
          actions.push(
            { label: '创建产品', path: '/dashboard/products/create' },
            { label: '查看原料', path: '/dashboard/materials' }
          )
          break
        case 'warehouse':
          actions.push(
            { label: '仓储管理', path: '/dashboard/storage' }
          )
          break
        case 'seller':
          actions.push(
            { label: '销售管理', path: '/dashboard/sales' }
          )
          break
        case 'supervisor':
          actions.push(
            { label: '用户管理', path: '/dashboard/admin/users' },
            { label: '溯源查询', path: '/trace' }
          )
          break
      }
      
      return actions
    }
  },
  mounted() {
    this.loadStats()
    this.loadRecentActivities()
  },
  methods: {
    async loadStats() {
      this.loading = true
      try {
        const res = await getStats()
        this.stats = res.data
      } catch (error) {
        console.error('加载统计信息失败:', error)
        this.$message.error('加载统计信息失败')
      } finally {
        this.loading = false
      }
    },
    
    async loadRecentActivities() {
      try {
        const res = await getRecentActivities({ limit: 10 })
        this.recentActivities = res.data.activities || []
      } catch (error) {
        console.error('加载最近活动失败:', error)
      }
    }
  }
}
</script>

<style scoped>
.dashboard {
  height: 100%;
}

.welcome-card h2 {
  color: #409eff;
  margin-bottom: 10px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
  padding: 20px;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  color: #666;
  font-size: 14px;
}

.quick-actions h3,
.recent-activity h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
}

.action-btn {
  width: 100%;
  height: 60px;
  font-size: 16px;
}

.recent-activity ::v-deep .el-table th {
  background-color: #f5f7fa;
}
</style>