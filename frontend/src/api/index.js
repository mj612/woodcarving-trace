import request from '@/utils/request'

// 用户相关API
export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

export function register(data) {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

export function updateUserInfo(data) {
  return request({
    url: '/user/info',
    method: 'put',
    data
  })
}

export function changePassword(data) {
  return request({
    url: '/user/password',
    method: 'put',
    data
  })
}

// 原料相关API
export function createMaterial(data) {
  return request({
    url: '/materials',
    method: 'post',
    data
  })
}

export function getMaterialList(params) {
  return request({
    url: '/materials',
    method: 'get',
    params
  })
}

export function getMaterialDetail(id) {
  return request({
    url: `/materials/${id}`,
    method: 'get'
  })
}

export function transferMaterial(data) {
  return request({
    url: '/materials/transfer',
    method: 'post',
    data
  })
}

// 产品相关API
export function createProduct(data) {
  return request({
    url: '/products',
    method: 'post',
    data
  })
}

export function getProductList(params) {
  return request({
    url: '/products',
    method: 'get',
    params
  })
}

export function getProductDetail(id) {
  return request({
    url: `/products/${id}`,
    method: 'get'
  })
}

export function transferProduct(data) {
  return request({
    url: '/products/transfer',
    method: 'post',
    data
  })
}

// 仓储相关API
export function recordStorage(data) {
  return request({
    url: '/storage/record',
    method: 'post',
    data
  })
}

// 销售相关API
export function recordSales(data) {
  return request({
    url: '/sales/record',
    method: 'post',
    data
  })
}

// 溯源查询API（公开）
export function getTrace(id) {
  return request({
    url: `/trace/${id}`,
    method: 'get'
  })
}

export function getHistory(id) {
  return request({
    url: `/history/${id}`,
    method: 'get'
  })
}

// 管理后台API
export function getUserList(params) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  })
}

export function updateUserStatus(data) {
  return request({
    url: '/admin/users/status',
    method: 'put',
    data
  })
}

// 统计数据API
export function getStats() {
  return request({
    url: '/stats',
    method: 'get'
  })
}

export function getRecentActivities(params) {
  return request({
    url: '/activities/recent',
    method: 'get',
    params
  })
}

// 文件上传API
export function uploadFile(file) {
  const formData = new FormData()
  formData.append('file', file)
  
  return request({
    url: '/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
