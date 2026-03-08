<template>
  <div class="product-list">
    <div class="page-container">
      <div class="action-bar">
        <h2>产品列表</h2>
        <el-button type="primary" @click="$router.push('/dashboard/products/create')">
          创建产品
        </el-button>
      </div>
      
      <div class="card">
        <el-table :data="products" style="width: 100%" v-loading="loading">
          <el-table-column prop="productId" label="产品ID" width="150"></el-table-column>
          <el-table-column prop="productName" label="产品名称" width="180"></el-table-column>
          <el-table-column prop="dimensions" label="尺寸" width="120"></el-table-column>
          <el-table-column prop="weight" label="重量(kg)" width="100"></el-table-column>
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
                         v-if="scope.row.status === 'produced'">转移</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script>
import { getProductList } from '@/api'

export default {
  name: 'ProductList',
  data() {
    return {
      products: [],
      loading: false
    }
  },
  
  mounted() {
    this.loadProducts()
  },
  methods: {
    async loadProducts() {
      this.loading = true
      try {
        const res = await getProductList()
        // 修复：API返回的是 { list: [], total: 0 } 格式
        const list = (res.data && res.data.list) ? res.data.list : []
        // 从链上数据中提取字段到顶层
        this.products = list.map(item => {
          const chainData = item.chainData || {}
          return {
            ...item,
            productId: item.productId || chainData.productId,
            productName: chainData.productName || '',
            dimensions: chainData.dimensions || '',
            weight: chainData.weight || 0,
            status: chainData.status || 'unknown'
          }
        })
      } catch (error) {
        console.error('加载产品列表失败:', error)
        this.$message.error('加载产品列表失败')
        this.products = []
      } finally {
        this.loading = false
      }
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
    
    viewDetail(row) {
      this.$router.push(`/dashboard/products/${row.productId}`)
    },
    
    transfer(row) {
      console.log('转移产品:', row)
    }
  }
}
</script>