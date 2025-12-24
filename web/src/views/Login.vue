<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-pink-50 via-blue-50 to-purple-50">
    <div class="w-full max-w-md px-6">
      <!-- 标题 -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold bg-gradient-to-r from-pink-400 to-purple-500 bg-clip-text text-transparent mb-2">
          🐱 兑兑猫
        </h1>
        <p class="text-gray-500 text-sm">公益 CDK 兑换平台</p>
      </div>

      <!-- 登录卡片 -->
      <div class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-xl p-8 border border-gray-100">
        <!-- 管理员账密登录表单 -->
        <div class="space-y-4 mb-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">用户名</label>
            <input
              v-model="adminForm.username"
              type="text"
              placeholder="请输入用户名"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-pink-400 focus:border-transparent outline-none transition-all"
              @keyup.enter="handleAdminLogin"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">密码</label>
            <input
              v-model="adminForm.password"
              type="password"
              placeholder="请输入密码"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-pink-400 focus:border-transparent outline-none transition-all"
              @keyup.enter="handleAdminLogin"
            />
          </div>
          <button
            @click="handleAdminLogin"
            :disabled="loading || !adminForm.username || !adminForm.password"
            class="w-full bg-gradient-to-r from-pink-400 to-purple-500 text-white py-3 px-4 rounded-lg font-medium hover:from-pink-500 hover:to-purple-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-md hover:shadow-lg"
          >
            <span v-if="!loading">🔑 登录</span>
            <span v-else>登录中...</span>
          </button>
        </div>

        <!-- 分割线 -->
        <div class="relative my-6">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-gray-200"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-4 bg-white text-gray-500">或</span>
          </div>
        </div>

        <!-- LinuxDo OAuth 登录 -->
        <div class="space-y-3">
          <button
            @click="handleOAuthLogin"
            :disabled="loading"
            class="w-full bg-white border-2 border-blue-500 text-blue-600 py-3 px-4 rounded-lg font-medium hover:bg-blue-50 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span v-if="!loading">🔗 使用 LinuxDo 登录</span>
            <span v-else>跳转中...</span>
          </button>
          <p class="text-xs text-gray-500 text-center">
            使用 LinuxDo 账号登录，安全便捷
          </p>
        </div>

        <!-- 错误提示 -->
        <div v-if="errorMessage" class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-sm text-red-600">{{ errorMessage }}</p>
        </div>
      </div>

      <!-- 底部说明 -->
      <p class="text-center text-xs text-gray-400 mt-6">
        登录即代表您同意遵守平台使用规则
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { adminLogin, getOAuthURL } from '../utils/api'
import { setToken, setUserInfo } from '../utils/auth'

const router = useRouter()

const loading = ref(false)
const errorMessage = ref('')

// 管理员登录表单
const adminForm = ref({
  username: '',
  password: ''
})

// LinuxDo OAuth 登录
async function handleOAuthLogin() {
  try {
    loading.value = true
    errorMessage.value = ''

    const url = await getOAuthURL()
    // 跳转到 LinuxDo 登录页面
    window.location.href = url
  } catch (error) {
    errorMessage.value = error.message || 'OAuth登录失败'
  } finally {
    loading.value = false
  }
}

// 管理员账密登录
async function handleAdminLogin() {
  if (!adminForm.value.username || !adminForm.value.password) {
    errorMessage.value = '请输入用户名和密码'
    return
  }

  try {
    loading.value = true
    errorMessage.value = ''

    const { token, user } = await adminLogin(adminForm.value.username, adminForm.value.password)

    // 保存登录信息
    setToken(token)
    setUserInfo(user)

    // 跳转到首页或管理页
    router.push(user.is_admin ? '/admin' : '/')
  } catch (error) {
    errorMessage.value = error.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 自定义样式 */
</style>
