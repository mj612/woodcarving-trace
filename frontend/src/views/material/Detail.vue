<template>
  <div class="material-detail">
    <div class="page-container">
      <div class="card">
        <h2 class="form-title">原料详情</h2>
        
        <div v-if="material" v-loading="loading">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="原料ID">{{ material.materialId || chainData.materialId }}</el-descriptions-item>
            <el-descriptions-item label="木材种类">{{ chainData.woodType }}</el-descriptions-item>
            <el-descriptions-item label="产地">{{ chainData.origin }}</el-descriptions-item>
            <el-descriptions-item label="采伐证编号">{{ chainData.harvestCert }}</el-descriptions-item>
            <el-descriptions-item label="数量">{{ chainData.quantity }}</el-descriptions-item>
            <el-descriptions-item label="质量等级">
              <el-tag :type="getQualityType(chainData.quality)">
                {{ chainData.quality }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(chainData.status)">
                {{ getStatusText(chainData.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="供应商">{{ chainData.supplierName }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ material.createdAt }}</el-descriptions-item>
            <el-descriptions-item label="交易ID">{{ material.txId }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ material.description }}</el-descriptions-item>
          </el-descriptions>
          
          <!-- 流转记录 -->
          <div v-if="records && records.length > 0" style="margin-top: 30px;">
            <h3>流转记录</h3>
            <el-timeline>
              <el-timeline-item 
                v-for="(record, index) in records" 
                :key="index"
                :timestamp="record.timestamp"
                placement="top">
                <el-card>
                  <p><strong>从:</strong> {{ record.fromOwner }}</p>
                  <p><strong>到:</strong> {{ record.toName }} ({{ record.toOwner }})</p>
                  <p><strong>位置:</strong> {{ record.location }}</p>
                  <p v-if="record.remarks"><strong>备注:</strong> {{ record.remarks }}</p>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </div>
          
          <div class="action-buttons" style="margin-top: 20px;">
            <el-button @click="$router.back()">返回</el-button>
            <el-button type="primary" @click="transferMaterial" 
                       v-if="chainData && chainData.status === 'available'">转移原料</el-button>
          </div>
        </div>
        
        <div v-else>
          <el-alert title="原料不存在" type="error" show-icon></el-alert>
        </div>
      </div>
      
      <!-- 转移对话框 -->
      <el-dialog title="转移原料" :visible.sync="showTransferDialog" width="500px">
        <el-form :model="transferForm" label-width="100px">
          <el-form-item label="目标用户ID" required>
            <el-input v-model="transferForm.toUser" placeholder="请输入目标用户ID"></el-input>
            <div style="color: #999; font-size: 12px; margin-top: 5px;">
              提示：可以是工匠的用户ID
            </div>
          </el-form-item>
          
          <el-form-item label="目标用户名" required>
            <el-input v-model="transferForm.toName" placeholder="请输入目标用户名称"></el-input>
          </el-form-item>
          
          <el-form-item label="位置">
            <el-input v-model="transferForm.location" placeholder="请输入位置信息（如：工作台A区）"></el-input>
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
import { getMaterialDetail, transferMaterial } from '@/api'

export default {
  name: 'MaterialDetail',
  data() {
    return {
      material: null,
      chainData: null,
      records: [],
      loading: false,
      showTransferDialog: false,
      transferForm: {
        toUser: '',
        toName: '',
        location: '',
        remarks: ''
      }
    }
  },
  mounted() {
    this.loadMaterialDetail()
  },
  methods: {
    async loadMaterialDetail() {
      const materialId = this.$route.params.id
      this.loading = true
      
      try {
        const res = await getMaterialDetail(materialId)
        this.material = res.data.material
        this.chainData = res.data.chainData
        this.records = res.data.records || []
      } catch (error) {
        console.error('加载原料详情失败:', error)
        this.$message.error('加载原料详情失败')
      } finally {
        this.loading = false
      }
    },
    
    getQualityType(quality) {
      const map = {
        '优': 'success',
        '良': 'warning',
        '一般': 'info'
      }
      return map[quality] || 'info'
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
    
    transferMaterial() {
      this.showTransferDialog = true
    },
    
    async submitTransfer() {
      // 验证表单
      if (!this.transferForm.toUser || !this.transferForm.toName) {
        this.$message.warning('请填写完整的转移信息')
        return
      }
      
      try {
        await transferMaterial({
          materialId: this.material.materialId || this.chainData.materialId,
          ...this.transferForm
        })
        
        this.$message.success('原料转移成功')
        this.showTransferDialog = false
        this.resetTransferForm()
        
        // 重新加载原料详情
        this.loadMaterialDetail()
      } catch (error) {
        console.error('原料转移失败:', error)
        this.$message.error('原料转移失败，请重试')
      }
    },
    
    resetTransferForm() {
      this.transferForm = {
        toUser: '',
        toName: '',
        location: '',
        remarks: ''
      }
    }
  }
}
</script>