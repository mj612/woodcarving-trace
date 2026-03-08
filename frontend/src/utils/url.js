/**
 * URL处理工具函数
 */

// 后端服务器地址（用于静态文件访问，不包含/api/v1）
const API_BASE_URL = process.env.VUE_APP_SERVER_URL || 'http://localhost:3000'

/**
 * 获取完整的文件URL
 * @param {string} path - 文件路径
 * @returns {string} 完整的URL
 */
export function getFileUrl(path) {
  if (!path) {
    return ''
  }
  
  // 如果已经是完整URL，直接返回
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path
  }
  
  // 如果是相对路径，添加后端服务器地址
  if (path.startsWith('/uploads/')) {
    return `${API_BASE_URL}${path}`
  }
  
  // 如果只是文件名，添加完整路径
  return `${API_BASE_URL}/uploads/${path}`
}

/**
 * 获取图片URL列表
 * @param {string|array} images - 图片路径（字符串或数组）
 * @returns {array} 完整URL数组
 */
export function getImageUrls(images) {
  if (!images) {
    return []
  }
  
  // 如果是字符串，按逗号分割
  let imageList = []
  if (typeof images === 'string') {
    imageList = images.split(',').filter(img => img.trim())
  } else if (Array.isArray(images)) {
    imageList = images
  } else {
    return []
  }
  
  // 转换为完整URL
  return imageList.map(img => getFileUrl(img))
}

export default {
  getFileUrl,
  getImageUrls
}
