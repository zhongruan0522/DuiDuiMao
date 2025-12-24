<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 via-blue-50 to-purple-50">
    <!-- 顶部导航 -->
    <Navigation
      title="兑兑猫"
      :user-info="userInfo"
      @login="router.push('/login')"
      @navigate="handleNavigation"
    />

    <!-- 主内容区 -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- 欢迎标题 -->
      <div class="text-center mb-12">
        <h1 class="text-4xl font-bold text-gray-800 mb-4">
          欢迎来到兑兑猫 🎉
        </h1>
        <p class="text-lg text-gray-600">
          公益 CDK 兑换平台，选择您需要的额度档位进行兑换
        </p>
      </div>

      <!-- 档位列表 -->
      <div v-if="userInfo" class="max-w-6xl mx-auto">
        <h2 class="text-2xl font-bold text-gray-800 mb-6">可兑换档位 📊</h2>

        <!-- 加载状态 -->
        <div v-if="tiersLoading" class="text-center py-12">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-pink-500"></div>
          <p class="mt-2 text-gray-600">加载中...</p>
        </div>

        <!-- 档位卡片列表 -->
        <div v-else-if="tiers.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="tier in sortedTiers"
            :key="tier.id"
            class="bg-white/80 backdrop-blur-sm rounded-xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition-all"
            :class="{ 'opacity-60': tier.stock === 0 || tier.required_level > userInfo.trust_level }"
          >
            <!-- 档位头部 -->
            <div class="flex justify-between items-start mb-4">
              <div>
                <h3 class="text-xl font-bold text-gray-800">{{ tier.name }}</h3>
                <p class="text-sm text-gray-500 mt-1">所需等级: Lv.{{ tier.required_level }}</p>
              </div>
              <div class="text-right">
                <div class="text-2xl font-bold text-pink-600">{{ tier.quota }}</div>
                <div class="text-xs text-gray-500">额度</div>
              </div>
            </div>

            <!-- 档位信息 -->
            <div class="space-y-2 mb-4">
              <div class="flex justify-between text-sm">
                <span class="text-gray-600">每日限购:</span>
                <span class="font-semibold text-gray-800">
                  {{ tier.daily_limit === 0 ? '不限' : tier.daily_limit }}
                </span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-gray-600">剩余库存:</span>
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

            <!-- 兑换按钮 -->
            <button
              @click="handleRedeem(tier)"
              :disabled="tier.stock === 0 || tier.required_level > userInfo.trust_level"
              class="w-full px-4 py-2 bg-gradient-to-r from-pink-400 to-purple-500 text-white rounded-lg font-medium hover:from-pink-500 hover:to-purple-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <template v-if="tier.stock === 0">
                库存不足
              </template>
              <template v-else-if="tier.required_level > userInfo.trust_level">
                等级不足
              </template>
              <template v-else>
                立即兑换
              </template>
            </button>
          </div>
        </div>

        <!-- 空状态提示 -->
        <div v-else class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg p-12 text-center border border-gray-100">
          <div class="text-6xl mb-4">🎁</div>
          <h3 class="text-xl font-semibold text-gray-800 mb-2">暂无可兑换档位</h3>
          <p class="text-gray-600 mb-6">
            管理员还没有添加任何档位，敬请期待～
          </p>
        </div>
      </div>

      <!-- 未登录提示 -->
      <div v-else class="max-w-2xl mx-auto">
        <div class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg p-12 text-center border border-gray-100">
          <div class="text-6xl mb-4">🔐</div>
          <h3 class="text-xl font-semibold text-gray-800 mb-2">请先登录</h3>
          <p class="text-gray-600 mb-6">
            登录后即可查看和兑换 CDK 档位
          </p>
          <button
            @click="router.push('/login')"
            class="px-8 py-3 bg-gradient-to-r from-pink-400 to-purple-500 text-white rounded-lg font-medium hover:from-pink-500 hover:to-purple-600 transition-all"
          >
            立即登录
          </button>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getUserInfo, clearAuth } from '../utils/auth'
import { logout, getUserTiers } from '../utils/api'
import Navigation from '../components/Navigation.vue'

const router = useRouter()
const userInfo = ref(null)
const tiers = ref([])
const tiersLoading = ref(false)

// 排序后的档位列表（按sort_order降序）
const sortedTiers = computed(() => {
  return [...tiers.value].sort((a, b) => b.sort_order - a.sort_order)
})

onMounted(async () => {
  userInfo.value = getUserInfo()

  // 移除管理员自动跳转逻辑，允许管理员访问用户端首页
  // 管理员可以通过下拉菜单的"管理后台"按钮进入后台

  // 如果用户已登录，加载档位列表
  if (userInfo.value) {
    await loadTiers()
  }
})

// 加载档位列表
async function loadTiers() {
  tiersLoading.value = true
  try {
    tiers.value = await getUserTiers()
  } catch (error) {
    console.error('加载档位列表失败:', error)
  } finally {
    tiersLoading.value = false
  }
}

// 处理导航事件
function handleNavigation(action) {
  switch (action) {
    case 'history':
      // TODO: 实现兑换记录页面
      alert('兑换记录功能即将上线～')
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

// 处理兑换操作
function handleRedeem(tier) {
  // TODO: 实现兑换逻辑
  alert(`兑换功能即将上线！\n档位: ${tier.name}\n额度: ${tier.quota}`)
}
</script>
