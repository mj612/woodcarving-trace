<template>
  <div class="blockchain-explorer">
    <div class="page-container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1>区块链浏览器</h1>
        <p class="subtitle">查看链上交易记录和区块信息</p>
      </div>

      <!-- 搜索框 -->
      <div class="card search-card">
        <el-input
          v-model="searchQuery"
          placeholder="搜索交易ID、产品ID或原料ID"
          class="search-input"
          @keyup.enter="handleSearch"
        >
          <el-button slot="append" icon="el-icon-search" @click="handleSearch"></el-button>
        </el-input>
      </div>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stats-row">
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: #409eff;">
              <i class="el-icon-files"></i>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ blockCount }}</div>
              <div class="stat-label">区块总数</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: #67c23a;">
              <i class="el-icon-document"></i>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ txCount }}</div>
              <div class="stat-label">交易总数</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: #e6a23c;">
              <i class="el-icon-box"></i>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ assetCount }}</div>
              <div class="stat-label">资产总数</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: #f56c6c;">
              <i class="el-icon-time"></i>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ avgBlockTime }}s</div>
              <div class="stat-label">平均出块时间</div>
            </div>
          </div>
        </el-col>
      </el-row>

      <!-- Tab切换 -->
      <div class="card">
        <el-tabs v-model="activeTab">
          <!-- 最新交易 -->
          <el-tab-pane label="最新交易" name="transactions">
            <el-table :data="transactions" style="width: 100%">
              <el-table-column prop="txId" label="交易ID" width="180">
                <template slot-scope="scope">
                  <el-link type="primary" @click="viewTransaction(scope.row)">
                    {{ scope.row.txId.substring(0, 16) }}...
                  </el-link>
                </template>
              </el-table-column>
              <el-table-column prop="action" label="操作类型" width="120">
                <template slot-scope="scope">
                  <el-tag :type="getActionType(scope.row.action)" size="small">
                    {{ getActionText(scope.row.action) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="assetId" label="资产ID" width="150" />
              <el-table-column prop="operator" label="操作人" width="120" />
              <el-table-column prop="timestamp" label="时间" width="180">
                <template slot-scope="scope">
                  {{ formatTime(scope.row.timestamp) }}
                </template>
              </el-table-column>
              <el-table-column label="状态" width="100">
                <template>
                  <el-tag type="success" size="small">已确认</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100">
                <template slot-scope="scope">
                  <el-button type="text" size="small" @click="viewTransaction(scope.row)">
                    详情
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            
            <div class="pagination">
              <el-pagination
                @current-change="handlePageChange"
                :current-page="currentPage"
                :page-size="pageSize"
                layout="total, prev, pager, next"
                :total="totalTransactions"
              />
            </div>
          </el-tab-pane>

          <!-- 区块列表 -->
          <el-tab-pane label="区块列表" name="blocks">
            <el-table :data="blocks" style="width: 100%">
              <el-table-column prop="blockNumber" label="区块高度" width="120" />
              <el-table-column prop="blockHash" label="区块哈希" width="200">
                <template slot-scope="scope">
                  {{ scope.row.blockHash.substring(0, 20) }}...
                </template>
              </el-table-column>
              <el-table-column prop="txCount" label="交易数" width="100" />
              <el-table-column prop="timestamp" label="时间" width="180">
                <template slot-scope="scope">
                  {{ formatTime(scope.row.timestamp) }}
                </template>
              </el-table-column>
              <el-table-column prop="dataSize" label="数据大小" width="120">
                <template slot-scope="scope">
                  {{ scope.row.dataSize }} KB
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100">
                <template slot-scope="scope">
                  <el-button type="text" size="small" @click="viewBlock(scope.row)">
                    详情
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <!-- 资产列表 -->
          <el-tab-pane label="资产列表" name="assets">
            <el-table :data="assets" style="width: 100%">
              <el-table-column prop="assetId" label="资产ID" width="150" />
              <el-table-column prop="assetType" label="类型" width="100">
                <template slot-scope="scope">
                  <el-tag :type="scope.row.assetType === 'product' ? 'success' : 'warning'" size="small">
                    {{ scope.row.assetType === 'product' ? '产品' : '原料' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="name" label="名称" width="180" />
              <el-table-column prop="owner" label="当前所有者" width="120" />
              <el-table-column prop="status" label="状态" width="100">
                <template slot-scope="scope">
                  <el-tag size="small">{{ scope.row.status }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="createdAt" label="创建时间" width="180">
                <template slot-scope="scope">
                  {{ formatTime(scope.row.createdAt) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150">
                <template slot-scope="scope">
                  <el-button type="text" size="small" @click="viewTrace(scope.row.assetId)">
                    查看溯源
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>

    <!-- 交易详情对话框 -->
    <el-dialog title="交易详情" :visible.sync="txDialogVisible" width="600px">
      <div v-if="selectedTx" class="tx-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="交易ID">{{ selectedTx.txId }}</el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <el-tag :type="getActionType(selectedTx.action)">
              {{ getActionText(selectedTx.action) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="资产ID">{{ selectedTx.assetId }}</el-descriptions-item>
          <el-descriptions-item label="操作人">{{ selectedTx.operator }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ formatTime(selectedTx.timestamp) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag type="success">已确认</el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="tx-data">
          <h4>交易数据</h4>
          <pre>{{ JSON.stringify(selectedTx.data, null, 2) }}</pre>
        </div>
      </div>
    </el-dialog>

    <!-- 区块详情对话框 -->
    <el-dialog title="区块详情" :visible.sync="blockDialogVisible" width="600px">
      <div v-if="selectedBlock" class="block-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="区块高度">{{ selectedBlock.blockNumber }}</el-descriptions-item>
          <el-descriptions-item label="区块哈希">{{ selectedBlock.blockHash }}</el-descriptions-item>
          <el-descriptions-item label="前一区块哈希">{{ selectedBlock.prevHash }}</el-descriptions-item>
          <el-descriptions-item label="交易数量">{{ selectedBlock.txCount }}</el-descriptions-item>
          <el-descriptions-item label="数据大小">{{ selectedBlock.dataSize }} KB</el-descriptions-item>
          <el-descriptions-item label="时间戳">{{ formatTime(selectedBlock.timestamp) }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getStats } from '@/api'

export default {
  name: 'BlockchainExplorer',
  data() {
    return {
      searchQuery: '',
      activeTab: 'transactions',
      
      // 统计数据
      blockCount: 0,
      txCount: 0,
      assetCount: 0,
      avgBlockTime: 2.5,
      
      // 交易列表
      transactions: [],
      currentPage: 1,
      pageSize: 10,
      totalTransactions: 0,
      
      // 区块列表
      blocks: [],
      
      // 资产列表
      assets: [],
      
      // 对话框
      txDialogVisible: false,
      blockDialogVisible: false,
      selectedTx: null,
      selectedBlock: null
    }
  },
  
  mounted() {
    this.loadData()
  },
  
  methods: {
    async loadData() {
      // 模拟加载区块链数据
      this.loadStats()
      this.loadTransactions()
      this.loadBlocks()
      this.loadAssets()
    },
    
    loadStats() {
      // 模拟统计数据
      this.blockCount = 1247
      this.txCount = 3856
      this.assetCount = 156
      this.avgBlockTime = 2.5
    },
    
    loadTransactions() {
      // 模拟交易数据
      const mockTxs = []
      const actions = ['create', 'transfer', 'storage', 'sales']
      const operators = ['供应商A', '工匠B', '仓管C', '销售商D']
      
      for (let i = 0; i < 50; i++) {
        mockTxs.push({
          txId: this.generateTxId(),
          action: actions[Math.floor(Math.random() * actions.length)],
          assetId: `ASSET${1000 + i}`,
          operator: operators[Math.floor(Math.random() * operators.length)],
          timestamp: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000),
          data: {
            key1: 'value1',
            key2: 'value2'
          }
        })
      }
      
      this.transactions = mockTxs.slice(0, this.pageSize)
      this.totalTransactions = mockTxs.length
    },
    
    loadBlocks() {
      // 模拟区块数据
      const mockBlocks = []
      for (let i = 1247; i > 1237; i--) {
        mockBlocks.push({
          blockNumber: i,
          blockHash: this.generateHash(),
          prevHash: this.generateHash(),
          txCount: Math.floor(Math.random() * 10) + 1,
          dataSize: (Math.random() * 50 + 10).toFixed(2),
          timestamp: new Date(Date.now() - (1247 - i) * 2.5 * 1000)
        })
      }
      this.blocks = mockBlocks
    },
    
    loadAssets() {
      // 模拟资产数据
      const mockAssets = []
      const types = ['product', 'material']
      const names = ['龙凤呈祥木雕', '花鸟屏风', '红木家具', '黄花梨原木', '紫檀木料']
      const owners = ['工匠A', '仓库B', '销售商C', '供应商D']
      const statuses = ['available', 'in_storage', 'sold', 'transferred']
      
      for (let i = 0; i < 20; i++) {
        mockAssets.push({
          assetId: `ASSET${1000 + i}`,
          assetType: types[Math.floor(Math.random() * types.length)],
          name: names[Math.floor(Math.random() * names.length)],
          owner: owners[Math.floor(Math.random() * owners.length)],
          status: statuses[Math.floor(Math.random() * statuses.length)],
          createdAt: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000)
        })
      }
      this.assets = mockAssets
    },
    
    generateTxId() {
      return 'TX' + Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
    },
    
    generateHash() {
      return '0x' + Array.from({length: 64}, () => 
        Math.floor(Math.random() * 16).toString(16)
      ).join('')
    },
    
    formatTime(time) {
      if (!time) return ''
      const date = new Date(time)
      return date.toLocaleString('zh-CN')
    },
    
    getActionType(action) {
      const map = {
        'create': 'success',
        'transfer': 'warning',
        'storage': 'info',
        'sales': 'danger'
      }
      return map[action] || 'info'
    },
    
    getActionText(action) {
      const map = {
        'create': '创建',
        'transfer': '转移',
        'storage': '仓储',
        'sales': '销售'
      }
      return map[action] || action
    },
    
    handleSearch() {
      if (!this.searchQuery.trim()) {
        this.$message.warning('请输入搜索内容')
        return
      }
      this.$message.info('搜索功能开发中...')
    },
    
    handlePageChange(page) {
      this.currentPage = page
      this.loadTransactions()
    },
    
    viewTransaction(tx) {
      this.selectedTx = tx
      this.txDialogVisible = true
    },
    
    viewBlock(block) {
      this.selectedBlock = block
      this.blockDialogVisible = true
    },
    
    viewTrace(assetId) {
      this.$router.push(`/trace?id=${assetId}`)
    }
  }
}
</script>

<style scoped>
.blockchain-explorer {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px;
}

.page-container {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 30px;
}

.page-header h1 {
  font-size: 32px;
  color: #333;
  margin-bottom: 10px;
}

.subtitle {
  color: #666;
  font-size: 14px;
}

.card {
  background: white;
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.search-card {
  padding: 20px;
}

.search-input {
  max-width: 600px;
  margin: 0 auto;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
}

.stat-icon i {
  font-size: 28px;
  color: white;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.pagination {
  margin-top: 20px;
  text-align: center;
}

.tx-detail, .block-detail {
  padding: 10px 0;
}

.tx-data {
  margin-top: 20px;
}

.tx-data h4 {
  margin-bottom: 10px;
  color: #333;
}

.tx-data pre {
  background: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 12px;
}
</style>
