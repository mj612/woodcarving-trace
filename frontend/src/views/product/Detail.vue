<template>
  <div class="product-detail">
    <div class="page-container">
      <div class="card">
        <h2 class="form-title">产品详情</h2>
        
        <div v-if="product">
          <el-tabs type="border-card">
            <el-tab-pane label="基本信息">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="产品ID">{{ product.productId }}</el-descriptions-item>
                <el-descriptions-item label="产品名称">{{ chainData.productName }}</el-descriptions-item>
                <el-descriptions-item label="尺寸">{{ chainData.dimensions }}</el-descriptions-item>
                <el-descriptions-item label="重量">{{ chainData.weight }} kg</el-descriptions-item>
                <el-descriptions-item label="状态">
                  <el-tag :type="getStatusType(chainData.status)">
                    {{ getStatusText(chainData.status) }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="创建时间">{{ product.createdAt }}</el-descriptions-item>
                <el-descriptions-item label="雕刻耗时">{{ chainData.carveTime }} 小时</el-descriptions-item>
                <el-descriptions-item label="原料ID">{{ chainData.materialId }}</el-descriptions-item>
                <el-descriptions-item label="工艺描述" :span="2">{{ chainData.craftDesc }}</el-descriptions-item>
              </el-descriptions>
            </el-tab-pane>
            
            <el-tab-pane label="图片展示">
              <div class="image-gallery">
                <el-image
                  v-for="(img, index) in imageList"
                  :key="index"
                  :src="img"
                  :preview-src-list="imageList"
                  class="gallery-image"
                ></el-image>
                <div v-if="!imageList || imageList.length === 0" class="no-image">
                  暂无图片
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
          
          <div class="action-buttons" style="margin-top: 20px;">
            <el-button @click="$router.back()">返回</el-button>
            <el-button type="primary" @click="showTransferDialog = true" 
                       v-if="chainData && chainData.status === 'produced'">转移产品</el-button>
          </div>
        </div>
        
        <div v-else>
          <el-alert title="产品不存在" type="error" show-icon></el-alert>
        </div>
      </div>
      
      <!-- 转移对话框 -->
      <el-dialog title="转移产品" :visible.sync="showTransferDialog" width="500px">
        <el-form :model="transferForm" label-width="100px">
          <el-form-item label="目标用户ID" required>
            <el-input v-model="transferForm.toUser" placeholder="请输入目标用户ID"></el-input>
          </el-form-item>
          
          <el-form-item label="目标用户名" required>
            <el-input v-model="transferForm.toName" placeholder="请输入目标用户名称"></el-input>
          </el-form-item>
          
          <el-form-item label="转移类型" required>
            <el-select v-model="transferForm.transferType" placeholder="请选择转移类型" style="width: 100%;">
              <el-option label="入库" value="store_in"></el-option>
              <el-option label="出库" value="store_out"></el-option>
              <el-option label="销售" value="sell"></el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="位置">
            <el-input v-model="transferForm.location" placeholder="请输入位置信息"></el-input>
          </el-form-item>
          
          <el-form-item label="备注">
            <el-input 
              v-model="transferForm.remarks" 
              type="textarea" 
              :rows="3"
              placeholder="请输入备注信息"
            ></el-input>
          </el-form-item>
        </el-form>
        
        <div slot="footer">
          <el-button @click="showTransferDialog = false">取消</el-button>
          <el-button type="primary" @click="submitTransfer">确定转移</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import { getProductDetail, transferProduct } from '@/api'
import { getImageUrls } from '@/utils/url'

export default {
  name: 'ProductDetail',
  data() {
    return {
      product: null,
      chainData: null,
      records: [],
      loading: false,
      showTransferDialog: false,
      transferForm: {
        toUser: '',
        toName: '',
        transferType: '',
        location: '',
        remarks: ''
      }
    }
  },
  computed: {
    imageList() {
      if (!this.product || !this.product.images) {
        return []
      }
      return getImageUrls(this.product.images)
    }
  },
  mounted() {
    this.loadProductDetail()
  },
  methods: {
    async loadProductDetail() {
      const productId = this.$route.params.id
      this.loading = true
      
      try {
        const res = await getProductDetail(productId)
        this.product = res.data.product
        this.chainData = res.data.chainData
        this.records = res.data.records || []
      } catch (error) {
        console.error('加载产品详情失败:', error)
        this.$message.error('加载产品详情失败')
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
    
    async submitTransfer() {
      // 验证表单
      if (!this.transferForm.toUser || !this.transferForm.toName || !this.transferForm.transferType) {
        this.$message.warning('请填写完整的转移信息')
        return
      }
      
      try {
        await transferProduct({
          productId: this.product.productId,
          ...this.transferForm
        })
        
        this.$message.success('产品转移成功')
        this.showTransferDialog = false
        this.resetTransferForm()
        
        // 重新加载产品详情
        this.loadProductDetail()
      } catch (error) {
        console.error('产品转移失败:', error)
        this.$message.error('产品转移失败，请重试')
      }
    },
    
    resetTransferForm() {
      this.transferForm = {
        toUser: '',
        toName: '',
        transferType: '',
        location: '',
        remarks: ''
      }
    }
  }
}
</script>

<style scoped>
.image-gallery {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.gallery-image {
  width: 200px;
  height: 150px;
  border-radius: 8px;
  overflow: hidden;
}

.no-image {
  color: #999;
  font-style: italic;
  padding: 20px;
}

.action-buttons {
  display: flex;
  gap: 10px;
}
</style>
