<template>
  <div class="storage-index">
    <div class="page-container">
      <div class="action-bar">
        <h2>仓储管理</h2>
        <div>
          <el-button type="primary" @click="showInboundDialog = true">入库登记</el-button>
          <el-button type="warning" @click="showOutboundDialog = true">出库登记</el-button>
        </div>
      </div>
      
      <div class="card">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="库存列表" name="inventory">
            <el-table :data="inventoryList" style="width: 100%" v-loading="loading">
              <el-table-column prop="productId" label="产品ID" width="150"></el-table-column>
              <el-table-column prop="productName" label="产品名称" width="180"></el-table-column>
              <el-table-column prop="warehouse" label="仓库" width="120"></el-table-column>
              <el-table-column prop="location" label="库位" width="120"></el-table-column>
              <el-table-column prop="temperature" label="温度(℃)" width="100"></el-table-column>
              <el-table-column prop="humidity" label="湿度(%)" width="100"></el-table-column>
              <el-table-column prop="entryTime" label="入库时间" width="180"></el-table-column>
              <el-table-column label="操作" width="150">
                <template slot-scope="scope">
                  <el-button size="mini" @click="viewDetail(scope.row)">详情</el-button>
                  <el-button size="mini" type="warning" @click="outbound(scope.row)">出库</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          
          <el-tab-pane label="入库记录" name="inbound">
            <el-table :data="inboundRecords" style="width: 100%">
              <el-table-column prop="recordId" label="记录ID" width="150"></el-table-column>
              <el-table-column prop="productId" label="产品ID" width="150"></el-table-column>
              <el-table-column prop="warehouse" label="目标仓库" width="120"></el-table-column>
              <el-table-column prop="operator" label="操作员" width="120"></el-table-column>
              <el-table-column prop="entryTime" label="入库时间" width="180"></el-table-column>
              <el-table-column prop="remarks" label="备注"></el-table-column>
            </el-table>
          </el-tab-pane>
          
          <el-tab-pane label="出库记录" name="outbound">
            <el-table :data="outboundRecords" style="width: 100%">
              <el-table-column prop="recordId" label="记录ID" width="150"></el-table-column>
              <el-table-column prop="productId" label="产品ID" width="150"></el-table-column>
              <el-table-column prop="target" label="目标方" width="150"></el-table-column>
              <el-table-column prop="trackingNo" label="物流单号" width="150"></el-table-column>
              <el-table-column prop="operator" label="操作员" width="120"></el-table-column>
              <el-table-column prop="exitTime" label="出库时间" width="180"></el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>
      
      <!-- 入库对话框 -->
      <el-dialog title="入库登记" :visible.sync="showInboundDialog" width="500px">
        <el-form :model="inboundForm" label-width="100px">
          <el-form-item label="产品ID">
            <el-input v-model="inboundForm.productId" placeholder="请输入产品ID"></el-input>
          </el-form-item>
          <el-form-item label="仓库">
            <el-select v-model="inboundForm.warehouseId" placeholder="请选择仓库">
              <el-option label="一号仓库" value="WH001"></el-option>
              <el-option label="二号仓库" value="WH002"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="库位">
            <el-input v-model="inboundForm.location" placeholder="请输入库位"></el-input>
          </el-form-item>
          <el-form-item label="温湿度">
            <el-col :span="11">
              <el-input-number v-model="inboundForm.temperature" :min="-10" :max="50" :step="0.1"></el-input-number>
              <span style="margin-left: 5px;">℃</span>
            </el-col>
            <el-col class="line" :span="2" style="text-align: center;">-</el-col>
            <el-col :span="11">
              <el-input-number v-model="inboundForm.humidity" :min="0" :max="100"></el-input-number>
              <span style="margin-left: 5px;">%</span>
            </el-col>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="inboundForm.remarks" type="textarea"></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button @click="showInboundDialog = false">取消</el-button>
          <el-button type="primary" @click="submitInbound">确定</el-button>
        </div>
      </el-dialog>
      
      <!-- 出库对话框 -->
      <el-dialog title="出库登记" :visible.sync="showOutboundDialog" width="500px">
        <el-form :model="outboundForm" label-width="100px">
          <el-form-item label="产品ID">
            <el-input v-model="outboundForm.productId" placeholder="请输入产品ID"></el-input>
          </el-form-item>
          <el-form-item label="仓库">
            <el-select v-model="outboundForm.warehouseId" placeholder="请选择仓库">
              <el-option label="一号仓库" value="WH001"></el-option>
              <el-option label="二号仓库" value="WH002"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="库位">
            <el-input v-model="outboundForm.location" placeholder="请输入库位"></el-input>
          </el-form-item>
          <el-form-item label="温湿度">
            <el-col :span="11">
              <el-input-number v-model="outboundForm.temperature" :min="-10" :max="50" :step="0.1"></el-input-number>
              <span style="margin-left: 5px;">℃</span>
            </el-col>
            <el-col class="line" :span="2" style="text-align: center;">-</el-col>
            <el-col :span="11">
              <el-input-number v-model="outboundForm.humidity" :min="0" :max="100"></el-input-number>
              <span style="margin-left: 5px;">%</span>
            </el-col>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="outboundForm.remarks" type="textarea"></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button @click="showOutboundDialog = false">取消</el-button>
          <el-button type="primary" @click="submitOutbound">确定</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import { recordStorage, getProductList } from '@/api'

export default {
  name: 'StorageIndex',
  data() {
    return {
      activeTab: 'inventory',
      showInboundDialog: false,
      showOutboundDialog: false,
      loading: false,
      inboundForm: {
        productId: '',
        operationType: 'in',
        warehouseId: '',
        location: '',
        temperature: 20,
        humidity: 60,
        remarks: ''
      },
      outboundForm: {
        productId: '',
        operationType: 'out',
        warehouseId: '',
        location: '',
        temperature: 20,
        humidity: 60,
        remarks: ''
      },
      inventoryList: [],
      inboundRecords: [],
      outboundRecords: []
    }
  },
  
  mounted() {
    this.loadInventory()
  },
  
  methods: {
    async loadInventory() {
      this.loading = true
      try {
        // 获取在库产品列表
        const res = await getProductList({ status: 'in_storage' })
        this.inventoryList = res.data.list || []
      } catch (error) {
        console.error('加载库存列表失败:', error)
        this.inventoryList = []
      } finally {
        this.loading = false
      }
    },
    viewDetail(row) {
      this.$message.info(`查看产品 ${row.productId} 详情`)
    },
    
    outbound(row) {
      this.outboundForm.productId = row.productId
      this.showOutboundDialog = true
    },
    
    async submitInbound() {
      try {
        await recordStorage(this.inboundForm)
        this.$message.success('入库登记成功')
        this.showInboundDialog = false
        this.resetInboundForm()
        this.loadInventory()
      } catch (error) {
        console.error('入库登记失败:', error)
        this.$message.error('入库登记失败')
      }
    },
    
    async submitOutbound() {
      try {
        await recordStorage(this.outboundForm)
        this.$message.success('出库登记成功')
        this.showOutboundDialog = false
        this.resetOutboundForm()
        this.loadInventory()
      } catch (error) {
        console.error('出库登记失败:', error)
        this.$message.error('出库登记失败')
      }
    },
    
    resetInboundForm() {
      this.inboundForm = {
        productId: '',
        operationType: 'in',
        warehouseId: '',
        location: '',
        temperature: 20,
        humidity: 60,
        remarks: ''
      }
    },
    
    resetOutboundForm() {
      this.outboundForm = {
        productId: '',
        operationType: 'out',
        warehouseId: '',
        location: '',
        temperature: 20,
        humidity: 60,
        remarks: ''
      }
    }
  }
}
</script>