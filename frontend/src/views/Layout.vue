<template>
  <div class="layout">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside width="200px">
        <el-menu
          :default-active="$route.path"
          class="el-menu-vertical"
          @select="handleMenuSelect"
          router
        >
          <el-menu-item index="/dashboard" router>
            <i class="el-icon-house"></i>
            <span slot="title">工作台</span>
          </el-menu-item>
          
          <el-submenu index="material">
            <template slot="title">
              <i class="el-icon-box"></i>
              <span>原料管理</span>
            </template>
            <el-menu-item index="/dashboard/materials">原料列表</el-menu-item>
            <el-menu-item index="/dashboard/materials/create">添加原料</el-menu-item>
          </el-submenu>
          
          <el-submenu index="product">
            <template slot="title">
              <i class="el-icon-goods"></i>
              <span>产品管理</span>
            </template>
            <el-menu-item index="/dashboard/products">产品列表</el-menu-item>
            <el-menu-item index="/dashboard/products/create">创建产品</el-menu-item>
          </el-submenu>
          
          <el-menu-item index="/dashboard/storage">
            <i class="el-icon-office-building"></i>
            <span slot="title">仓储管理</span>
          </el-menu-item>
          
          <el-menu-item index="/dashboard/sales">
            <i class="el-icon-shopping-cart-full"></i>
            <span slot="title">销售管理</span>
          </el-menu-item>
          
          <el-menu-item index="/trace">
            <i class="el-icon-search"></i>
            <span slot="title">溯源查询</span>
          </el-menu-item>
          
          <el-menu-item index="/blockchain-explorer">
            <i class="el-icon-coin"></i>
            <span slot="title">区块链浏览器</span>
          </el-menu-item>
          
          <el-menu-item index="/dashboard/profile">
            <i class="el-icon-user"></i>
            <span slot="title">个人中心</span>
          </el-menu-item>
          
          <!-- 管理员菜单 -->
          <el-submenu index="admin" v-if="userInfo && userInfo.role === 'supervisor'">
            <template slot="title">
              <i class="el-icon-setting"></i>
              <span>系统管理</span>
            </template>
            <el-menu-item index="/dashboard/admin/users">用户管理</el-menu-item>
          </el-submenu>
        </el-menu>
      </el-aside>
      
      <!-- 主体内容 -->
      <el-container>
        <!-- 头部 -->
        <el-header height="60px" class="header">
          <div class="header-content">
            <h2>木雕工艺品溯源系统</h2>
            <div class="user-info">
              <span>欢迎，{{ userInfo ? (userInfo.realName || userInfo.username) : '用户' }}</span>
              <el-dropdown @command="handleUserCommand">
                <span class="el-dropdown-link">
                  <i class="el-icon-arrow-down el-icon--right"></i>
                </span>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
            </div>
          </div>
        </el-header>
        
        <!-- 内容区域 -->
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'Layout',
  computed: {
    ...mapState(['userInfo'])
  },
  methods: {
    handleMenuSelect(index) {
      console.log('菜单选择:', index)
    },
    
    handleMenuClick(path) {
      this.$router.push(path)
    },
    
    handleUserCommand(command) {
      if (command === 'logout') {
        this.$store.dispatch('logout')
        this.$router.push('/login')
      } else if (command === 'profile') {
        this.$router.push('/profile')
      }
    }
  }
}
</script>

<style scoped>
.layout {
  height: 100vh;
}

.el-menu-vertical:not(.el-menu--collapse) {
  width: 200px;
  min-height: 400px;
}

.header {
  background-color: #409eff;
  color: white;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0,0,0,.1);
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h2 {
  margin: 0;
  font-size: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.el-dropdown-link {
  cursor: pointer;
  color: white;
}

.el-main {
  padding: 0;
  background-color: #f5f5f5;
}
</style>