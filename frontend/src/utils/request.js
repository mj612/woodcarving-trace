import axios from 'axios'
import { Message } from 'element-ui'
import store from '@/store'
import router from '@/router'

// 创建axios实例
const service = axios.create({
  baseURL: process.env.VUE_APP_API_BASE_URL || 'http://localhost:3000/api/v1',
  timeout: 30000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 在请求头中添加token
    const token = store.getters.token
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data

    // 根据业务状态码处理
    if (res.code !== 200) {
      Message({
        message: res.msg || '请求失败',
        type: 'error',
        duration: 5000
      })

      // 401: Token过期或无效
      if (res.code === 401) {
        store.dispatch('logout').then(() => {
          router.push('/login')
        })
      }

      return Promise.reject(new Error(res.msg || '请求失败'))
    } else {
      return res
    }
  },
  error => {
    console.error('响应错误:', error)
    
    // 处理401未授权错误
    if (error.response && error.response.status === 401) {
      Message({
        message: '登录已过期，请重新登录',
        type: 'warning',
        duration: 3000
      })
      
      // 清除token并跳转到登录页
      store.dispatch('logout').then(() => {
        router.push('/login')
      })
      
      return Promise.reject(error)
    }
    
    // 显示详细错误信息
    let errorMsg = '网络错误'
    if (error.response && error.response.data) {
      errorMsg = error.response.data.msg || error.response.data.message || error.message
    } else {
      errorMsg = error.message
    }
    
    Message({
      message: errorMsg,
      type: 'error',
      duration: 5000
    })
    return Promise.reject(error)
  }
)

export default service
