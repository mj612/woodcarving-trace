<template>
  <div class="material-list">
    <div class="page-container">
      <div class="action-bar">
        <h2>原料列表</h2>
        <el-button type="primary" @click="$router.push('/dashboard/materials/create')">
          添加原料
        </el-button>
      </div>
      
      <div class="card">
        <el-table :data="materials" style="width: 100%" v-loading="loading">
          <el-table-column prop="materialId" label="原料ID" width="150"></el-table-column>
          <el-table-column prop="woodType" label="木材种类" width="120"></el-table-column>
          <el-table-column prop="origin" label="产地" width="120"></el-table-column>
          <el-table-column prop="quantity" label="数量" width="100"></el-table-column>
          <el-table-column prop="quality" label="质量等级" width="100"></el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template slot-scope="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="180"></el-table-column>
          <el-table-column label="操作" width="150">
            <template slot-scope="scope">
              <el-button size="mini" @click="viewDetail(scope.row)">查看</el-button>
              <el-button size="mini" type="primary" @click="transfer(scope.row)" 
                         v-if="scope.row.status === 'available'">转移</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script>
import { getMaterialList } from '@/api'

export default {
  name: 'MaterialList',
  data() {
    return {
      materials: [],
      loading: false
    }
  },
  
  mounted() {
    this.loadMaterials()
  },
  methods: {
    async loadMaterials() {
      this.loading = true
      try {
        const res = await getMaterialList()
        console.log('原料列表API响应:', res)
        
        // API返回的是 { list: [], total: 0 } 格式
        const list = (res.data && res.data.list) ? res.data.list : []
        console.log('原料列表原始数据:', list)
        
        // 从链上数据中提取字段到顶层
        this.materials = list.map(item => {
          const chainData = item.chainData || {}
          console.log('原料项:', item.materialId, 'chainData:', chainData)
          
          return {
            ...item,
            materialId: item.materialId || chainData.materialId,
            woodType: chainData.woodType || '未知',
            origin: chainData.origin || '未知',
            quantity: chainData.quantity || 0,
            quality: chainData.quality || '未知',
            status: chainData.status || 'unknown'
          }
        })
        
        console.log('处理后的原料列表:', this.materials)
      } catch (error) {
        console.error('加载原料列表失败:', error)
        this.$message.error('加载原料列表失败')
        this.materials = []
      } finally {
        this.loading = false
      }
    },
    
    getStatusType(status) {
      const map = {
        'available': 'success',
        'used': 'info',
        'transferred': 'warning'
      }
      return map[status] || 'info'
    },
    
    getStatusText(status) {
      const map = {
        'available': '可用',
        'used': '已使用',
        'transferred': '已转移'
      }
      return map[status] || status
    },
    
    viewDetail(row) {
      this.$router.push(`/dashboard/materials/${row.materialId}`)
    },
    
    transfer(row) {
      // 跳转到转移页面
      console.log('转移原料:', row)
    }
  }
}
</script>