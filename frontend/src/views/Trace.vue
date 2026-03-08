<template>
  <div class="trace-page">
    <div class="page-container">
      <div class="card">
        <h2 class="form-title">产品溯源查询</h2>
        
        <el-form :inline="true" class="search-form">
          <el-form-item label="产品ID">
            <el-input 
              v-model="searchId" 
              placeholder="请输入产品ID或扫描二维码"
              style="width: 300px;"
            ></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="searchTrace">查询</el-button>
            <el-button @click="clearSearch">清空</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 溯源结果 -->
      <div v-if="traceData" class="card">
        <div class="trace-header">
          <h3>溯源信息</h3>
          <el-button 
            type="primary" 
            size="small" 
            icon="el-icon-link"
            @click="goToBlockchainExplorer"
          >
            在区块链浏览器中查看
          </el-button>
        </div>
        
        <!-- 产品基本信息 -->
        <div class="product-info">
          <h4>产品信息</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="产品ID">{{ traceData.productId }}</el-descriptions-item>
            <el-descriptions-item label="产品名称">{{ traceData.productName }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ traceData.createdAt }}</el-descriptions-item>
            <el-descriptions-item label="当前状态">
              <el-tag :type="getStatusType(traceData.status)">
                {{ getStatusText(traceData.status) }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>
        
        <!-- 溯源时间轴 -->
        <div class="trace-timeline">
          <h4>流转记录</h4>
          <div class="timeline">
            <div 
              v-for="(record, index) in traceHistory" 
              :key="index" 
              class="timeline-item"
            >
              <div class="timeline-dot"></div>
              <div class="timeline-content">
                <div class="timeline-title">{{ record.action }}</div>
                <div class="timeline-desc">{{ record.description }}</div>
                <div class="timeline-time">{{ record.time }}</div>
                <div class="timeline-operator">操作人：{{ record.operator }}</div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 原料信息 -->
        <div v-if="traceData.materialInfo" class="material-info">
          <h4>原料信息</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="原料ID">{{ traceData.materialInfo.materialId }}</el-descriptions-item>
            <el-descriptions-item label="木材种类">{{ traceData.materialInfo.woodType }}</el-descriptions-item>
            <el-descriptions-item label="产地">{{ traceData.materialInfo.origin }}</el-descriptions-item>
            <el-descriptions-item label="供应商">{{ traceData.materialInfo.supplier }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
      
      <!-- 二维码扫描 -->
      <div class="card">
        <h3>扫码查询</h3>
        <div class="qr-scan-area">
          <div class="qr-placeholder">
            <i class="el-icon-camera" style="font-size: 48px; color: #ccc;"></i>
            <p>点击此处打开摄像头扫描二维码</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getTrace } from '@/api'

export default {
  name: 'Trace',
  data() {
    return {
      searchId: '',
      traceData: null,
      traceHistory: []
    }
  },
  mounted() {
    // 检查URL参数，支持扫码直接跳转
    const id = this.$route.query.id
    if (id) {
      this.searchId = id
      this.searchTrace()
    }
  },
  methods: {
    async searchTrace() {
      if (!this.searchId.trim()) {
        this.$message.warning('请输入产品ID')
        return
      }
      
      try {
        const res = await getTrace(this.searchId)
        const data = res.data
        
        // 处理产品基本信息
        this.traceData = {
          productId: data.productId || this.searchId,
          productName: data.productName || '未知产品',
          createdAt: data.createdAt || '',
          status: data.status || 'unknown',
          materialInfo: data.materialInfo || null
        }
        
        // 处理流转记录
        if (data.history && Array.isArray(data.history)) {
          this.traceHistory = data.history.map(record => ({
            action: record.action || '操作',
            description: record.description || '',
            time: record.timestamp || record.time || '',
            operator: record.operator || '系统'
          }))
        } else {
          this.traceHistory = []
        }
        
        this.$message.success('查询成功')
      } catch (error) {
        console.error('查询溯源信息失败:', error)
        this.$message.error('查询失败，请检查产品ID是否正确')
        this.traceData = null
        this.traceHistory = []
      }
    },
    
    clearSearch() {
      this.searchId = ''
      this.traceData = null
      this.traceHistory = []
    },
    
    getStatusType(status) {
      const map = {
        'produced': 'success',
        'in_storage': 'warning',
        'sold': 'info'
      }
      return map[status] || 'info'
    },
    
    getStatusText(status) {
      const map = {
        'produced': '制作完成',
        'in_storage': '库存中',
        'sold': '已售出'
      }
      return map[status] || status
    },
    
    goToBlockchainExplorer() {
      this.$router.push('/blockchain-explorer')
    }
  }
}
</script>

<style scoped>
.trace-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
}

.card {
  background: white;
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card h2, .card h3, .card h4 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
}

.trace-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.trace-header h3 {
  margin: 0;
}

.form-title {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 30px;
}

.product-info,
.material-info {
  margin-bottom: 30px;
}

.search-form {
  margin-bottom: 20px;
  text-align: center;
}

.timeline {
  position: relative;
  padding-left: 30px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 15px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #e0e0e0;
}

.timeline-item {
  position: relative;
  margin-bottom: 30px;
}

.timeline-dot {
  position: absolute;
  left: -34px;
  top: 8px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #409eff;
  border: 3px solid white;
  box-shadow: 0 0 0 2px #409eff;
}

.timeline-content {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  border-left: 3px solid #409eff;
}

.timeline-title {
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.timeline-desc {
  color: #666;
  margin-bottom: 5px;
}

.timeline-time {
  color: #999;
  font-size: 12px;
  margin-bottom: 5px;
}

.timeline-operator {
  color: #409eff;
  font-size: 12px;
}

.qr-scan-area {
  text-align: center;
}

.qr-placeholder {
  border: 2px dashed #ccc;
  border-radius: 8px;
  padding: 40px;
  cursor: pointer;
}

.qr-placeholder:hover {
  border-color: #409eff;
  background-color: #f0f9ff;
}
</style>