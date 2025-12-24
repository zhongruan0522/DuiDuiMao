<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 via-blue-50 to-purple-50">
    <!-- 顶部导航 -->
    <Navigation
      title="档位管理"
      :user-info="userInfo"
      :show-home-button="true"
      @login="router.push('/login')"
      @navigate="handleNavigation"
    />

    <!-- 主内容区 -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- 页面标题 -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800 mb-2">档位管理 📊</h1>
        <p class="text-gray-600">管理平台的所有额度档位，设置库存、价格、限额等</p>
      </div>

      <!-- 操作栏 -->
      <div class="mb-6 flex justify-between items-center">
        <div class="text-sm text-gray-600">
          共 <span class="font-semibold text-pink-600">{{ tiers.length }}</span> 个档位
        </div>
        <button
          @click="showCreateDialog"
          class="px-4 py-2 bg-gradient-to-r from-pink-500 to-rose-500 text-white rounded-lg font-medium hover:from-pink-600 hover:to-rose-600 transition-all shadow-md hover:shadow-lg"
        >
          ➕ 创建新档位
        </button>
      </div>

      <!-- 档位列表 -->
      <div v-if="loading" class="text-center py-12">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-pink-500"></div>
        <p class="mt-2 text-gray-600">加载中...</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="tier in sortedTiers"
          :key="tier.id"
          class="bg-white/80 backdrop-blur-sm rounded-xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition-shadow"
        >
          <!-- 档位头部 -->
          <div class="flex justify-between items-start mb-4">
            <div>
              <h3 class="text-xl font-bold text-gray-800">{{ tier.name }}</h3>
              <p class="text-sm text-gray-500">排序: {{ tier.sort_order }}</p>
            </div>
            <span
              :class="[
                'px-2 py-1 text-xs rounded-full font-medium',
                tier.is_active
                  ? 'bg-green-100 text-green-700'
                  : 'bg-gray-100 text-gray-600'
              ]"
            >
              {{ tier.is_active ? '启用' : '禁用' }}
            </span>
          </div>

          <!-- 档位信息 -->
          <div class="space-y-2 mb-4">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">额度值:</span>
              <span class="font-semibold text-pink-600">{{ tier.quota }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">所需等级:</span>
              <span class="font-semibold text-gray-800">Lv.{{ tier.required_level }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">每日限购:</span>
              <span class="font-semibold text-gray-800">
                {{ tier.daily_limit === 0 ? '不限' : tier.daily_limit }}
              </span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">当前库存:</span>
              <span
                :class="[
                  'font-semibold',
                  tier.stock > 10 ? 'text-green-600' : tier.stock > 0 ? 'text-orange-600' : 'text-red-600'
                ]"
              >
                {{ tier.stock }}
              </span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex gap-2">
            <button
              @click="showEditDialog(tier)"
              class="flex-1 px-3 py-2 bg-blue-500 text-white rounded-lg text-sm font-medium hover:bg-blue-600 transition-colors"
            >
              ✏️ 编辑
            </button>
            <button
              @click="confirmDelete(tier)"
              class="flex-1 px-3 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors"
            >
              🗑️ 删除
            </button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && tiers.length === 0" class="text-center py-12">
        <div class="text-6xl mb-4">📦</div>
        <p class="text-gray-600 mb-4">还没有任何档位</p>
        <button
          @click="showCreateDialog"
          class="px-6 py-3 bg-gradient-to-r from-pink-500 to-rose-500 text-white rounded-lg font-medium hover:from-pink-600 hover:to-rose-600 transition-all"
        >
          创建第一个档位
        </button>
      </div>
    </main>

    <!-- 创建/编辑对话框 -->
    <div
      v-if="showDialog"
      class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
      @click.self="closeDialog"
    >
      <div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6">
        <h2 class="text-2xl font-bold text-gray-800 mb-6">
          {{ editingTier ? '编辑档位' : '创建档位' }}
        </h2>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- 档位名称 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">档位名称 *</label>
            <input
              v-model="formData.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
              placeholder="例如: 标准档位"
            />
          </div>

          <!-- 额度值 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">额度值 *</label>
            <input
              v-model.number="formData.quota"
              type="number"
              required
              min="1"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
              placeholder="例如: 100"
            />
          </div>

          <!-- 所需信任等级 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">所需信任等级 (0-4)</label>
            <input
              v-model.number="formData.required_level"
              type="number"
              min="0"
              max="4"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
              placeholder="0"
            />
          </div>

          <!-- 每日限购 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">每日限购 (0=不限)</label>
            <input
              v-model.number="formData.daily_limit"
              type="number"
              min="0"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
              placeholder="0"
            />
          </div>

          <!-- 排序权重 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">排序权重</label>
            <input
              v-model.number="formData.sort_order"
              type="number"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
              placeholder="0"
            />
            <p class="text-xs text-gray-500 mt-1">数值越大越靠前</p>
          </div>

          <!-- 是否启用 -->
          <div class="flex items-center">
            <input
              v-model="formData.is_active"
              type="checkbox"
              id="is_active"
              class="w-4 h-4 text-pink-600 border-gray-300 rounded focus:ring-pink-500"
            />
            <label for="is_active" class="ml-2 text-sm font-medium text-gray-700">启用此档位</label>
          </div>

          <!-- 按钮组 -->
          <div class="flex gap-3 pt-4">
            <button
              type="button"
              @click="closeDialog"
              class="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg font-medium hover:bg-gray-300 transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="submitting"
              class="flex-1 px-4 py-2 bg-gradient-to-r from-pink-500 to-rose-500 text-white rounded-lg font-medium hover:from-pink-600 hover:to-rose-600 transition-all disabled:opacity-50"
            >
              {{ submitting ? '提交中...' : '确定' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <div
      v-if="showDeleteDialog"
      class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
      @click.self="showDeleteDialog = false"
    >
      <div class="bg-white rounded-2xl shadow-2xl max-w-sm w-full p-6">
        <h2 class="text-xl font-bold text-gray-800 mb-4">确认删除</h2>
        <p class="text-gray-600 mb-6">
          确定要删除档位 <span class="font-semibold text-pink-600">{{ deletingTier?.name }}</span> 吗？此操作不可恢复。
        </p>
        <div class="flex gap-3">
          <button
            @click="showDeleteDialog = false"
            class="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg font-medium hover:bg-gray-300 transition-colors"
          >
            取消
          </button>
          <button
            @click="handleDelete"
            :disabled="submitting"
            class="flex-1 px-4 py-2 bg-red-500 text-white rounded-lg font-medium hover:bg-red-600 transition-colors disabled:opacity-50"
          >
            {{ submitting ? '删除中...' : '确认删除' }}
          </button>
        </div>
      </div>
    </div>

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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getTiers, createTier, updateTier, deleteTier, logout } from '../utils/api'
import { getUserInfo, clearAuth } from '../utils/auth'
import Navigation from '../components/Navigation.vue'

const router = useRouter()

// 用户信息
const userInfo = ref(null)

// 数据状态
const tiers = ref([])
const loading = ref(false)
const showDialog = ref(false)
const showDeleteDialog = ref(false)
const editingTier = ref(null)
const deletingTier = ref(null)
const submitting = ref(false)
const message = ref(null)

// 表单数据
const formData = ref({
  name: '',
  quota: 1,
  required_level: 0,
  daily_limit: 0,
  sort_order: 0,
  is_active: true
})

// 排序后的档位列表
const sortedTiers = computed(() => {
  return [...tiers.value].sort((a, b) => b.sort_order - a.sort_order)
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

// 加载档位列表
async function loadTiers() {
  loading.value = true
  try {
    tiers.value = await getTiers()
  } catch (error) {
    showMessage('加载档位列表失败: ' + error.message, 'error')
  } finally {
    loading.value = false
  }
}

// 显示创建对话框
function showCreateDialog() {
  editingTier.value = null
  formData.value = {
    name: '',
    quota: 1,
    required_level: 0,
    daily_limit: 0,
    sort_order: 0,
    is_active: true
  }
  showDialog.value = true
}

// 显示编辑对话框
function showEditDialog(tier) {
  editingTier.value = tier
  formData.value = {
    name: tier.name,
    quota: tier.quota,
    required_level: tier.required_level,
    daily_limit: tier.daily_limit,
    sort_order: tier.sort_order,
    is_active: tier.is_active
  }
  showDialog.value = true
}

// 关闭对话框
function closeDialog() {
  showDialog.value = false
  editingTier.value = null
}

// 提交表单
async function handleSubmit() {
  submitting.value = true
  try {
    if (editingTier.value) {
      // 更新档位
      await updateTier(editingTier.value.id, formData.value)
      showMessage('档位更新成功', 'success')
    } else {
      // 创建档位
      await createTier(formData.value)
      showMessage('档位创建成功', 'success')
    }
    closeDialog()
    await loadTiers()
  } catch (error) {
    showMessage((editingTier.value ? '更新' : '创建') + '失败: ' + error.message, 'error')
  } finally {
    submitting.value = false
  }
}

// 确认删除
function confirmDelete(tier) {
  deletingTier.value = tier
  showDeleteDialog.value = true
}

// 执行删除
async function handleDelete() {
  if (!deletingTier.value) return

  submitting.value = true
  try {
    await deleteTier(deletingTier.value.id)
    showMessage('档位删除成功', 'success')
    showDeleteDialog.value = false
    deletingTier.value = null
    await loadTiers()
  } catch (error) {
    showMessage('删除失败: ' + error.message, 'error')
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

  // 如果不是管理员，跳转到首页
  if (!userInfo.value || !userInfo.value.is_admin) {
    router.push('/')
    return
  }

  loadTiers()
})
</script>
