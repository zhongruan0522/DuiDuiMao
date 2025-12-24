<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 via-blue-50 to-purple-50">
    <!-- 顶部导航 -->
    <Navigation
      title="系统设置"
      :user-info="userInfo"
      :show-home-button="true"
      @login="router.push('/login')"
      @navigate="handleNavigation"
    />

    <!-- 主内容区 -->
    <main class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- 页面标题 -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800 mb-2">系统设置 ⚙️</h1>
        <p class="text-gray-600">管理平台的全局设置，修改后自动保存并立即生效</p>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="text-center py-12">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-pink-500"></div>
        <p class="mt-2 text-gray-600">加载中...</p>
      </div>

      <!-- 设置表单 -->
      <div v-else class="bg-white/80 backdrop-blur-sm rounded-xl shadow-lg p-8 border border-gray-100">
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- 全局开关 -->
          <div class="border-b border-gray-200 pb-6">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <label class="text-lg font-semibold text-gray-800 mb-1 block">全局购买开关</label>
                <p class="text-sm text-gray-600">关闭后，用户端将无法进行购买操作</p>
              </div>
              <div class="flex items-center">
                <button
                  type="button"
                  @click="formData.global_enabled = !formData.global_enabled"
                  :class="[
                    'relative inline-flex h-8 w-16 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-pink-500 focus:ring-offset-2',
                    formData.global_enabled ? 'bg-gradient-to-r from-pink-500 to-rose-500' : 'bg-gray-300'
                  ]"
                >
                  <span
                    :class="[
                      'inline-block h-6 w-6 transform rounded-full bg-white shadow-lg transition-transform',
                      formData.global_enabled ? 'translate-x-9' : 'translate-x-1'
                    ]"
                  ></span>
                </button>
                <span class="ml-3 text-sm font-medium" :class="formData.global_enabled ? 'text-green-600' : 'text-gray-500'">
                  {{ formData.global_enabled ? '已启用' : '已关闭' }}
                </span>
              </div>
            </div>
          </div>

          <!-- 公告内容 -->
          <div class="border-b border-gray-200 pb-6">
            <label class="text-lg font-semibold text-gray-800 mb-1 block">系统公告</label>
            <p class="text-sm text-gray-600 mb-3">在用户端首页顶部显示的公告内容</p>
            <textarea
              v-model="formData.announcement"
              rows="4"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-none"
              placeholder="输入公告内容，支持 Emoji 😊"
            ></textarea>
            <p class="text-xs text-gray-500 mt-2">当前字数: {{ formData.announcement.length }}</p>
          </div>

          <!-- 订单超时时间 -->
          <div class="pb-6">
            <label class="text-lg font-semibold text-gray-800 mb-1 block">订单超时时间</label>
            <p class="text-sm text-gray-600 mb-3">用户下单后的有效时间，超时未支付将自动取消</p>
            <div class="flex items-center gap-3">
              <input
                v-model.number="formData.order_expire_minutes"
                type="number"
                required
                min="1"
                max="1440"
                class="w-32 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
              />
              <span class="text-gray-700 font-medium">分钟</span>
            </div>
            <p class="text-xs text-gray-500 mt-2">建议范围: 5-60 分钟（最大 1440 分钟 = 24 小时）</p>
          </div>

          <!-- 按钮组 -->
          <div class="flex gap-4 pt-4">
            <button
              type="button"
              @click="resetForm"
              class="flex-1 px-6 py-3 bg-gray-200 text-gray-700 rounded-lg font-medium hover:bg-gray-300 transition-colors"
            >
              重置
            </button>
            <button
              type="submit"
              :disabled="submitting"
              class="flex-1 px-6 py-3 bg-gradient-to-r from-pink-500 to-rose-500 text-white rounded-lg font-medium hover:from-pink-600 hover:to-rose-600 transition-all disabled:opacity-50 shadow-md hover:shadow-lg"
            >
              {{ submitting ? '保存中...' : '💾 保存设置' }}
            </button>
          </div>
        </form>

        <!-- 配置文件提示 -->
        <div class="mt-8 p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <div class="flex items-start gap-3">
            <span class="text-2xl">ℹ️</span>
            <div class="flex-1">
              <h3 class="text-sm font-semibold text-blue-800 mb-1">配置热更新</h3>
              <p class="text-xs text-blue-700">
                修改后的设置会立即生效并自动保存到 <code class="bg-blue-100 px-1 py-0.5 rounded">config.yaml</code> 文件中，无需重启服务器～
              </p>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 提示消息 -->
    <div
      v-if="message"
      :class="[
        'fixed bottom-4 right-4 px-6 py-3 rounded-lg shadow-lg transition-all z-50',
        message.type === 'success' ? 'bg-green-500 text-white' : 'bg-red-500 text-white'
      ]"
    >
      {{ message.text }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getSettings, updateSettings, logout } from '../utils/api'
import { getUserInfo, clearAuth } from '../utils/auth'
import Navigation from '../components/Navigation.vue'

const router = useRouter()

// 用户信息
const userInfo = ref(null)

// 数据状态
const loading = ref(false)
const submitting = ref(false)
const message = ref(null)

// 原始数据（用于重置）
const originalData = ref(null)

// 表单数据
const formData = ref({
  global_enabled: true,
  announcement: '',
  order_expire_minutes: 15
})

// 处理导航事件
function handleNavigation(action) {
  switch (action) {
    case 'history':
      alert('兑换记录功能即将上线～')
      break
    case 'home':
      router.push('/')
      break
    case 'admin':
      router.push('/admin')
      break
    case 'logout':
      handleLogout()
      break
  }
}

// 退出登录
async function handleLogout() {
  try {
    await logout()
  } catch (error) {
    console.error('登出失败:', error)
  } finally {
    clearAuth()
    userInfo.value = null
    router.push('/login')
  }
}

// 加载系统设置
async function loadSettings() {
  loading.value = true
  try {
    const settings = await getSettings()
    if (settings) {
      formData.value = { ...settings }
      originalData.value = { ...settings }
    }
  } catch (error) {
    showMessage('加载设置失败: ' + error.message, 'error')
  } finally {
    loading.value = false
  }
}

// 重置表单
function resetForm() {
  if (originalData.value) {
    formData.value = { ...originalData.value }
    showMessage('已重置为原始设置', 'success')
  }
}

// 提交表单
async function handleSubmit() {
  // 验证数据
  if (formData.value.order_expire_minutes < 1 || formData.value.order_expire_minutes > 1440) {
    showMessage('订单超时时间必须在 1-1440 分钟之间', 'error')
    return
  }

  submitting.value = true
  try {
    await updateSettings(formData.value)
    originalData.value = { ...formData.value }
    showMessage('设置保存成功！已自动生效～', 'success')
  } catch (error) {
    showMessage('保存失败: ' + error.message, 'error')
  } finally {
    submitting.value = false
  }
}

// 显示消息
function showMessage(text, type = 'success') {
  message.value = { text, type }
  setTimeout(() => {
    message.value = null
  }, 3000)
}

// 组件挂载时加载数据
onMounted(() => {
  userInfo.value = getUserInfo()

  // 如果不是管理员,跳转到首页
  if (!userInfo.value || !userInfo.value.is_admin) {
    router.push('/')
    return
  }

  loadSettings()
})
</script>
