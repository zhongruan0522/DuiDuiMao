<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-pink-50 via-blue-50 to-purple-50">
    <div class="text-center">
      <div class="inline-block animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-pink-500 mb-4"></div>
      <p class="text-lg text-gray-600">{{ message }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { setToken, setUserInfo } from '../utils/auth'

const router = useRouter()
const route = useRoute()
const message = ref('正在处理登录...')

onMounted(async () => {
  try {
    // 获取URL参数
    const code = route.query.code
    const state = route.query.state

    if (!code) {
      throw new Error('缺少授权码')
    }

    // 调用回调接口
    const response = await fetch(`/api/auth/callback?code=${code}&state=${state}`)
    const data = await response.json()

    if (!response.ok || !data.success) {
      throw new Error(data.error || '登录失败')
    }

    // 解密响应数据
    const token = doubleDecode(data.data.token)
    const userStr = doubleDecode(data.data.user)
    const user = JSON.parse(userStr)

    // 保存登录信息
    setToken(token)
    setUserInfo(user)

    message.value = '登录成功！正在跳转...'

    // 延迟跳转
    setTimeout(() => {
      router.push(user.is_admin ? '/admin' : '/')
    }, 1000)
  } catch (error) {
    message.value = '登录失败: ' + error.message
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  }
})

// 双重Base64解密
function doubleDecode(encoded) {
  const first = decodeURIComponent(escape(atob(encoded)))
  const second = decodeURIComponent(escape(atob(first)))
  return second
}
</script>
