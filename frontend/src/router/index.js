import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/store'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/trace/:id',
    name: 'Trace',
    component: () => import('@/views/Trace.vue'),
    meta: { title: '溯源查询', public: true }
  },
  {
    path: '/blockchain-explorer',
    name: 'BlockchainExplorer',
    component: () => import('@/views/BlockchainExplorer.vue'),
    meta: { title: '区块链浏览器', public: true }
  },
  {
    path: '/dashboard',
    component: () => import('@/views/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '工作台' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人中心' }
      },
      // 原料管理
      {
        path: 'materials',
        name: 'Materials',
        component: () => import('@/views/material/List.vue'),
        meta: { title: '原料列表', roles: ['supplier', 'supervisor'] }
      },
      {
        path: 'materials/create',
        name: 'MaterialCreate',
        component: () => import('@/views/material/Create.vue'),
        meta: { title: '添加原料', roles: ['supplier', 'supervisor'] }
      },
      {
        path: 'materials/:id',
        name: 'MaterialDetail',
        component: () => import('@/views/material/Detail.vue'),
        meta: { title: '原料详情', roles: ['supplier', 'supervisor'] }
      },
      // 产品管理
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/product/List.vue'),
        meta: { title: '产品列表' }
      },
      {
        path: 'products/create',
        name: 'ProductCreate',
        component: () => import('@/views/product/Create.vue'),
        meta: { title: '创建产品', roles: ['artisan', 'supervisor'] }
      },
      {
        path: 'products/:id',
        name: 'ProductDetail',
        component: () => import('@/views/product/Detail.vue'),
        meta: { title: '产品详情' }
      },
      // 仓储管理
      {
        path: 'storage',
        name: 'Storage',
        component: () => import('@/views/storage/Index.vue'),
        meta: { title: '仓储管理', roles: ['warehouse', 'supervisor'] }
      },
      // 销售管理
      {
        path: 'sales',
        name: 'Sales',
        component: () => import('@/views/sales/Index.vue'),
        meta: { title: '销售管理', roles: ['seller', 'supervisor'] }
      },
      // 管理后台
      {
        path: 'admin/users',
        name: 'UserManagement',
        component: () => import('@/views/admin/Users.vue'),
        meta: { title: '用户管理', roles: ['supervisor'] }
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  document.title = to.meta.title || '木雕溯源系统'
  
  const isLoggedIn = store.getters.isLoggedIn
  const userRole = store.getters.role
  
  // 公开页面直接通过
  if (to.meta.public) {
    next()
    return
  }
  
  // 需要登录的页面
  if (to.meta.requiresAuth) {
    if (!isLoggedIn) {
      next('/login')
      return
    }
    
    // 如果有token但没有用户信息，尝试获取用户信息
    if (isLoggedIn && !store.state.userInfo) {
      try {
        await store.dispatch('getUserInfo')
      } catch (error) {
        // 获取用户信息失败，清除登录状态
        store.dispatch('logout')
        next('/login')
        return
      }
    }
    
    // 检查角色权限
    if (to.meta.roles && !to.meta.roles.includes(store.getters.role)) {
      next('/dashboard')
      return
    }
  }
  
  // 已登录访问登录页，跳转到工作台
  if (isLoggedIn && (to.path === '/login' || to.path === '/register')) {
    next('/dashboard')
    return
  }
  
  next()
})

export default router
