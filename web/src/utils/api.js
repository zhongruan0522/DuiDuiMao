import { getToken, doubleEncode, doubleDecode } from './auth'

const API_BASE_URL = '/api'

/**
 * 通用请求函数
 * @param {string} url
 * @param {Object} options
 * @returns {Promise}
 */
async function request(url, options = {}) {
  const token = getToken()
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  }

  // 如果有token，添加到请求头
  if (token && url !== '/auth/admin/login' && url !== '/auth/callback') {
    headers['Authorization'] = `Bearer ${token}`
  }

  const config = {
    ...options,
    headers
  }

  try {
    const response = await fetch(`${API_BASE_URL}${url}`, config)
    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || '请求失败')
    }

    return data
  } catch (error) {
    console.error('请求错误:', error)
    throw error
  }
}

/**
 * GET 请求
 */
export function get(url) {
  return request(url, { method: 'GET' })
}

/**
 * POST 请求
 */
export function post(url, data) {
  return request(url, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

/**
 * PUT 请求
 */
export function put(url, data) {
  return request(url, {
    method: 'PUT',
    body: JSON.stringify(data)
  })
}

/**
 * DELETE 请求
 */
export function del(url) {
  return request(url, { method: 'DELETE' })
}

// ========== 认证相关API ==========

/**
 * 管理员账密登录
 * @param {string} username - 原始用户名
 * @param {string} password - 原始密码
 */
export async function adminLogin(username, password) {
  // 加密用户名和密码
  const encryptedUsername = doubleEncode(username)
  const encryptedPassword = doubleEncode(password)

  const response = await post('/auth/admin/login', {
    username: encryptedUsername,
    password: encryptedPassword
  })

  // 解密响应数据
  if (response.success && response.data) {
    const token = doubleDecode(response.data.token)
    const userStr = doubleDecode(response.data.user)
    const user = JSON.parse(userStr)

    return { token, user }
  }

  throw new Error('登录失败')
}

/**
 * 获取LinuxDo OAuth登录URL
 */
export async function getOAuthURL() {
  const response = await get('/auth/login')
  return response.data.url
}

/**
 * 登出
 */
export function logout() {
  return post('/auth/logout')
}

/**
 * 获取当前用户信息
 */
export function getCurrentUser() {
  return get('/user/me')
}

// ========== 档位管理API ==========

/**
 * 获取用户端档位列表（仅显示启用的档位）
 */
export async function getUserTiers() {
  const response = await get('/tiers')
  if (response.success && response.data) {
    // 解密每个档位的数据
    return response.data.map(tier => ({
      id: parseInt(doubleDecode(tier.id)),
      name: doubleDecode(tier.name),
      quota: parseInt(doubleDecode(tier.quota)),
      required_level: parseInt(doubleDecode(tier.required_level)),
      daily_limit: parseInt(doubleDecode(tier.daily_limit)),
      stock: parseInt(doubleDecode(tier.stock)),
      is_active: doubleDecode(tier.is_active) === 'true',
      sort_order: parseInt(doubleDecode(tier.sort_order)),
      created_at: doubleDecode(tier.created_at),
      updated_at: doubleDecode(tier.updated_at)
    }))
  }
  return []
}

/**
 * 获取所有档位列表（管理端）
 */
export async function getTiers() {
  const response = await get('/admin/tiers')
  if (response.success && response.data) {
    // 解密每个档位的数据
    return response.data.map(tier => ({
      id: parseInt(doubleDecode(tier.id)),
      name: doubleDecode(tier.name),
      quota: parseInt(doubleDecode(tier.quota)),
      required_level: parseInt(doubleDecode(tier.required_level)),
      daily_limit: parseInt(doubleDecode(tier.daily_limit)),
      stock: parseInt(doubleDecode(tier.stock)),
      is_active: doubleDecode(tier.is_active) === 'true',
      sort_order: parseInt(doubleDecode(tier.sort_order)),
      created_at: doubleDecode(tier.created_at),
      updated_at: doubleDecode(tier.updated_at)
    }))
  }
  return []
}

/**
 * 创建档位
 */
export async function createTier(tierData) {
  const response = await post('/admin/tiers', tierData)
  return response
}

/**
 * 更新档位
 */
export async function updateTier(id, tierData) {
  const response = await put(`/admin/tiers/${id}`, tierData)
  return response
}

/**
 * 删除档位
 */
export async function deleteTier(id) {
  const response = await del(`/admin/tiers/${id}`)
  return response
}

// ========== 系统设置API ==========

/**
 * 获取系统设置
 */
export async function getSettings() {
  const response = await get('/admin/settings')
  if (response.success && response.data) {
    // 解密设置数据
    return {
      global_enabled: doubleDecode(response.data.global_enabled) === 'true',
      announcement: doubleDecode(response.data.announcement),
      order_expire_minutes: parseInt(doubleDecode(response.data.order_expire_minutes))
    }
  }
  return null
}

/**
 * 更新系统设置
 * @param {Object} settings - 原始设置数据
 * @param {boolean} settings.global_enabled - 全局开关
 * @param {string} settings.announcement - 公告内容
 * @param {number} settings.order_expire_minutes - 订单超时时间
 */
export async function updateSettings(settings) {
  // 加密设置数据
  const encryptedData = {}

  if (settings.global_enabled !== undefined) {
    encryptedData.global_enabled = doubleEncode(String(settings.global_enabled))
  }

  if (settings.announcement !== undefined) {
    encryptedData.announcement = doubleEncode(settings.announcement)
  }

  if (settings.order_expire_minutes !== undefined) {
    encryptedData.order_expire_minutes = doubleEncode(String(settings.order_expire_minutes))
  }

  const response = await put('/admin/settings', encryptedData)
  return response
}
