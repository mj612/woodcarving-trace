<template>
  <div class="admin-users">
    <div class="page-container">
      <div class="action-bar">
        <h2>用户管理</h2>
        <el-button type="primary" @click="loadUsers">刷新列表</el-button>
      </div>
      
      <div class="card">
        <!-- 筛选条件 -->
        <el-form :inline="true" class="filter-form">
          <el-form-item label="角色">
            <el-select v-model="filters.role" placeholder="请选择角色" clearable>
              <el-option label="原料供应商" value="supplier"></el-option>
              <el-option label="雕刻工匠" value="artisan"></el-option>
              <el-option label="仓库管理员" value="warehouse"></el-option>
              <el-option label="销售商" value="seller"></el-option>
              <el-option label="消费者" value="consumer"></el-option>
              <el-option label="管理员" value="supervisor"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="filters.status" placeholder="请选择状态" clearable>
              <el-option label="启用" value="1"></el-option>
              <el-option label="禁用" value="0"></el-option>
              <el-option label="待审核" value="2"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadUsers">搜索</el-button>
            <el-button @click="resetFilters">重置</el-button>
          </el-form-item>
        </el-form>
        
        <!-- 用户表格 -->
        <el-table :data="users" style="width: 100%" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80"></el-table-column>
          <el-table-column prop="username" label="用户名" width="120"></el-table-column>
          <el-table-column prop="realName" label="真实姓名" width="120"></el-table-column>
          <el-table-column prop="role" label="角色" width="120">
            <template slot-scope="scope">
              <el-tag :type="getRoleType(scope.row.role)">
                {{ getRoleText(scope.row.role) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="company" label="公司" width="150"></el-table-column>
          <el-table-column prop="phone" label="电话" width="120"></el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template slot-scope="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="注册时间" width="180"></el-table-column>
          <el-table-column label="操作" width="200">
            <template slot-scope="scope">
              <el-button size="mini" @click="viewUser(scope.row)">详情</el-button>
              <el-button 
                size="mini" 
                type="warning" 
                @click="toggleStatus(scope.row)"
                v-if="scope.row.status !== 2"
              >
                {{ scope.row.status === 1 ? '禁用' : '启用' }}
              </el-button>
              <el-button 
                size="mini" 
                type="success" 
                @click="approveUser(scope.row)"
                v-if="scope.row.status === 2"
              >
                审核通过
              </el-button>
              <el-button size="mini" type="danger" @click="deleteUser(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination">
          <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="pagination.currentPage"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="pagination.pageSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total">
          </el-pagination>
        </div>
      </div>
      
      <!-- 用户详情对话框 -->
      <el-dialog title="用户详情" :visible.sync="showDetailDialog" width="600px">
        <div v-if="selectedUser">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="用户ID">{{ selectedUser.id }}</el-descriptions-item>
            <el-descriptions-item label="用户名">{{ selectedUser.username }}</el-descriptions-item>
            <el-descriptions-item label="真实姓名">{{ selectedUser.realName }}</el-descriptions-item>
            <el-descriptions-item label="角色">{{ getRoleText(selectedUser.role) }}</el-descriptions-item>
            <el-descriptions-item label="公司">{{ selectedUser.company }}</el-descriptions-item>
            <el-descriptions-item label="电话">{{ selectedUser.phone }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ selectedUser.email }}</el-descriptions-item>
            <el-descriptions-item label="地址">{{ selectedUser.address }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(selectedUser.status)">
                {{ getStatusText(selectedUser.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="注册时间">{{ selectedUser.createdAt }}</el-descriptions-item>
          </el-descriptions>
        </div>
        <div slot="footer">
          <el-button @click="showDetailDialog = false">关闭</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import { getUserList, updateUserStatus } from '@/api'

export default {
  name: 'AdminUsers',
  data() {
    return {
      loading: false,
      showDetailDialog: false,
      selectedUser: null,
      filters: {
        role: '',
        status: ''
      },
      users: [],
      pagination: {
        currentPage: 1,
        pageSize: 10,
        total: 0
      }
    }
  },
  mounted() {
    this.loadUsers()
  },
  methods: {
    async loadUsers() {
      this.loading = true
      try {
        const params = {
          page: this.pagination.currentPage,
          pageSize: this.pagination.pageSize,
          ...this.filters
        }
        const res = await getUserList(params)
        this.users = res.data.list || []
        this.pagination.total = res.data.total || 0
      } catch (error) {
        console.error('加载用户列表失败:', error)
        this.$message.error('加载用户列表失败')
      } finally {
        this.loading = false
      }
    },
    
    resetFilters() {
      this.filters = {
        role: '',
        status: ''
      }
      this.loadUsers()
    },
    
    getRoleType(role) {
      const map = {
        'supplier': 'success',
        'artisan': 'warning',
        'warehouse': 'info',
        'seller': 'primary',
        'consumer': '',
        'supervisor': 'danger'
      }
      return map[role] || ''
    },
    
    getRoleText(role) {
      const map = {
        'supplier': '原料供应商',
        'artisan': '雕刻工匠',
        'warehouse': '仓库管理员',
        'seller': '销售商',
        'consumer': '消费者',
        'supervisor': '管理员'
      }
      return map[role] || role
    },
    
    getStatusType(status) {
      const map = {
        0: 'danger',
        1: 'success',
        2: 'warning'
      }
      return map[status] || ''
    },
    
    getStatusText(status) {
      const map = {
        0: '禁用',
        1: '启用',
        2: '待审核'
      }
      return map[status] || status
    },
    
    viewUser(user) {
      this.selectedUser = user
      this.showDetailDialog = true
    },
    
    toggleStatus(user) {
      const action = user.status === 1 ? '禁用' : '启用'
      const newStatus = user.status === 1 ? 0 : 1
      
      this.$confirm(`确定要${action}用户 ${user.username} 吗？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await updateUserStatus({
            userId: user.id,
            status: newStatus
          })
          user.status = newStatus
          this.$message.success(`${action}成功`)
        } catch (error) {
          console.error('更新用户状态失败:', error)
          this.$message.error('操作失败')
        }
      })
    },
    
    approveUser(user) {
      this.$confirm(`确定要审核通过用户 ${user.username} 吗？`, '审核确认', {
        confirmButtonText: '通过',
        cancelButtonText: '拒绝',
        type: 'warning'
      }).then(async () => {
        try {
          await updateUserStatus({
            userId: user.id,
            status: 1
          })
          user.status = 1
          this.$message.success('审核通过')
        } catch (error) {
          console.error('审核失败:', error)
          this.$message.error('审核失败')
        }
      }).catch(async () => {
        try {
          await updateUserStatus({
            userId: user.id,
            status: 0
          })
          user.status = 0
          this.$message.info('已拒绝')
        } catch (error) {
          console.error('操作失败:', error)
        }
      })
    },
    
    deleteUser(user) {
      this.$confirm(`确定要删除用户 ${user.username} 吗？此操作不可恢复！`, '删除确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      }).then(() => {
        // TODO: 实现删除用户API
        this.$message.warning('删除功能暂未实现')
      })
    },
    
    handleSizeChange(val) {
      this.pagination.pageSize = val
      this.loadUsers()
    },
    
    handleCurrentChange(val) {
      this.pagination.currentPage = val
      this.loadUsers()
    }
  }
}
</script>

<style scoped>
.filter-form {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>