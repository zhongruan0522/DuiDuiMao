/**
 * 双重Base64加密
 * @param {string} data - 原始数据
 * @returns {string} 双重加密后的字符串
 */
export function doubleEncode(data) {
  // 第一次编码
  const first = btoa(unescape(encodeURIComponent(data)))
  // 第二次编码
  const second = btoa(unescape(encodeURIComponent(first)))
  return second
}

/**
 * 双重Base64解密
 * @param {string} encoded - 加密后的字符串
 * @returns {string} 解密后的原始数据
 */
export function doubleDecode(encoded) {
  // 第一次解码
  const first = decodeURIComponent(escape(atob(encoded)))
  // 第二次解码
  const second = decodeURIComponent(escape(atob(first)))
  return second
}

/**
 * 获取Token
 * @returns {string|null}
 */
export function getToken() {
  return localStorage.getItem('token')
}

/**
 * 保存Token
 * @param {string} token
 */
export function setToken(token) {
  localStorage.setItem('token', token)
}

/**
 * 删除Token
 */
export function removeToken() {
  localStorage.removeItem('token')
}

/**
 * 获取用户信息
 * @returns {Object|null}
 */
export function getUserInfo() {
  const userStr = localStorage.getItem('userInfo')
  return userStr ? JSON.parse(userStr) : null
}

/**
 * 保存用户信息
 * @param {Object} user
 */
export function setUserInfo(user) {
  localStorage.setItem('userInfo', JSON.stringify(user))
}

/**
 * 删除用户信息
 */
export function removeUserInfo() {
  localStorage.removeItem('userInfo')
}

/**
 * 清除所有登录信息
 */
export function clearAuth() {
  removeToken()
  removeUserInfo()
}
