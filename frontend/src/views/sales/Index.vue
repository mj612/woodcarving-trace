<template>
  <div class="sales-index">
    <div class="page-container">
      <div class="action-bar">
        <h2>销售管理</h2>
        <el-button type="primary" @click="showSaleDialog = true">记录销售</el-button>
      </div>
      
      <div class="card">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="销售订单" name="orders">
            <el-table :data="salesOrders" style="width: 100%" v-loading="loading">
              <el-table-column prop="orderId" label="订单号" width="150"></el-table-column>
              <el-table-column prop="productId" label="产品ID" width="150"></el-table-column>
              <el-table-column prop="productName" label="产品名称" width="180"></el-table-column>
              <el-table-column prop="buyerName" label="购买方" width="120"></el-table-column>
              <el-table-column prop="price" label="价格(元)" width="100"></el-table-column>
              <el-table-column prop="saleDate" label="销售日期" width="180"></el-table-column>
              <el-table-column prop="status" label="状态" width="100">
                <template slot-scope="scope">
                  <el-tag :type="getStatusType(scope.row.status)">
                    {{ getStatusText(scope.row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150">
                <template slot-scope="scope">
                  <el-button size="mini" @click="viewOrder(scope.row)">详情</el-button>
                  <el-button size="mini" type="primary" @click="updateTracking(scope.row)" 
                             v-if="scope.row.status === 'paid'">物流</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          
          <el-tab-pane label="销售统计" name="statistics">
            <div class="stats-grid">
              <div class="stat-card">
                <div class="stat-value">¥128,500</div>
                <div class="stat-label">本月销售额</div>
              </div>
              <div class="stat-card">
                <div class="stat-value">24</div>
                <div class="stat-label">本月订单数</div>
              </div>
              <div class="stat-card">
                <div class="stat-value">¥5,354</div>
                <div class="stat-label">平均客单价</div>
              </div>
              <div class="stat-card">
                <div class="stat-value">92%</div>
                <div class="stat-label">客户满意度</div>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
      
      <!-- 销售记录对话框 -->
      <el-dialog title="记录销售" :visible.sync="showSaleDialog" width="600px">
        <el-form :model="saleForm" :rules="saleRules" ref="saleForm" label-width="100px">
          <el-form-item label="产品ID" prop="productId">
            <el-input v-model="saleForm.productId" placeholder="请输入产品ID"></el-input>
          </el-form-item>
          <el-form-item label="订单号" prop="orderId">
            <el-input v-model="saleForm.orderId" placeholder="请输入订单号"></el-input>
          </el-form-item>
          <el-form-item label="购买方姓名" prop="buyerName">
            <el-input v-model="saleForm.buyerName" placeholder="请输入购买方姓名"></el-input>
          </el-form-item>
          <el-form-item label="联系方式" prop="buyerContact">
            <el-input v-model="saleForm.buyerContact" placeholder="请输入联系方式"></el-input>
          </el-form-item>
          <el-form-item label="价格(元)" prop="price">
            <el-input-number v-model="saleForm.price" :min="0" :step="100"></el-input-number>
          </el-form-item>
          <el-form-item label="销售日期">
            <el-date-picker
              v-model="saleForm.saleDate"
              type="datetime"
              placeholder="选择销售日期"
              format="yyyy-MM-dd HH:mm:ss"
              value-format="yyyy-MM-dd HH:mm:ss">
            </el-date-picker>
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button @click="showSaleDialog = false">取消</el-button>
          <el-button type="primary" @click="submitSale">确定</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import { recordSales, getProductList } from '@/api'

export default {
  name: 'SalesIndex',
  data() {
    return {
      activeTab: 'orders',
      showSaleDialog: false,
      loading: false,
      saleForm: {
        productId: '',
        orderId: '',
        buyerName: '',
        buyerContact: '',
        price: 0,
        saleDate: ''
      },
      saleRules: {
        productId: [
          { required: true, message: '请输入产品ID', trigger: 'blur' }
        ],
        orderId: [
          { required: true, message: '请输入订单号', trigger: 'blur' }
        ],
        buyerName: [
          { required: true, message: '请输入购买方姓名', trigger: 'blur' }
        ],
        buyerContact: [
          { required: true, message: '请输入联系方式', trigger: 'blur' }
        ],
        price: [
          { required: true, message: '请输入价格', trigger: 'blur' }
        ]
      },
      salesOrders: []
    }
  },
  
  mounted() {
    this.loadSalesOrders()
  },
  
  methods: {
    async loadSalesOrders() {
      this.loading = true
      try {
        // 获取已售产品列表
        const res = await getProductList({ status: 'sold' })
        // 修复：API返回的是 { list: [], total: 0 } 格式
        this.salesOrders = (res.data && res.data.list) ? res.data.list : []
      } catch (error) {
        console.error('加载销售订单失败:', error)
        this.salesOrders = []
      } finally {
        this.loading = false
      }
    },
    getStatusType(status) {
      const map = {
        'pending': 'info',
        'paid': 'warning',
        'shipped': 'success',
        'completed': 'success',
        'cancelled': 'danger'
      }
      return map[status] || 'info'
    },
    
    getStatusText(status) {
      const map = {
        'pending': '待付款',
        'paid': '已付款',
        'shipped': '已发货',
        'completed': '已完成',
        'cancelled': '已取消'
      }
      return map[status] || status
    },
    
    viewOrder(row) {
      this.$message.info(`查看订单 ${row.orderId} 详情`)
    },
    
    updateTracking(row) {
      this.$prompt('请输入物流单号', '更新物流信息', {
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }).then(({ value }) => {
        this.$message.success(`物流单号已更新为: ${value}`)
      })
    },
    
    submitSale() {
      this.$refs.saleForm.validate(async (valid) => {
        if (valid) {
          try {
            await recordSales(this.saleForm)
            this.$message.success('销售记录添加成功')
            this.showSaleDialog = false
            this.resetSaleForm()
            this.loadSalesOrders()
          } catch (error) {
            console.error('添加销售记录失败:', error)
            this.$message.error('添加销售记录失败')
          }
        }
      })
    },
    
    resetSaleForm() {
      this.$refs.saleForm.resetFields()
      this.saleForm.saleDate = ''
    }
  }
}
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  color: #666;
  font-size: 14px;
}
</style>