<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 via-blue-50 to-purple-50">
    <!-- 顶部导航 -->
    <Navigation
      title="兑兑猫 - 管理后台"
      :user-info="userInfo"
      :show-home-button="true"
      @login="router.push('/login')"
      @navigate="handleNavigation"
    />

    <!-- 主内容区 -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="text-center mb-12">
        <h1 class="text-4xl font-bold text-gray-800 mb-4">
          管理后台 🎛️
        </h1>
        <p class="text-lg text-gray-600">
          欢迎回来，{{ userInfo?.name }}！
        </p>
      </div>

      <!-- 功能模块 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- 档位管理 -->
        <div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition-shadow">
          <div class="text-4xl mb-4">📊</div>
          <h3 class="text-xl font-bold text-gray-800 mb-2">档位管理</h3>
          <p class="text-gray-600 text-sm mb-4">管理额度档位，设置库存、限额等</p>
          <button
            @click="router.push('/admin/tiers')"
            class="w-full px-4 py-2 bg-gradient-to-r from-pink-500 to-rose-500 text-white rounded-lg font-medium hover:from-pink-600 hover:to-rose-600 transition-all"
          >
            进入管理
          </button>
        </div>

        <!-- CDK导入管理 -->
        <div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition-shadow">
          <div class="text-4xl mb-4">📦</div>
          <h3 class="text-xl font-bold text-gray-800 mb-2">CDK导入管理</h3>
          <p class="text-gray-600 text-sm mb-4">批量导入CDK，管理CDK库存</p>
          <button class="w-full px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg font-medium hover:from-blue-600 hover:to-cyan-600 transition-all">
            进入管理
          </button>
        </div>

        <!-- 兑换记录 -->
        <div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition-shadow">
          <div class="text-4xl mb-4">📜</div>
          <h3 class="text-xl font-bold text-gray-800 mb-2">兑换记录</h3>
          <p class="text-gray-600 text-sm mb-4">查看所有用户的CDK兑换记录</p>
          <button class="w-full px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-500 text-white rounded-lg font-medium hover:from-green-600 hover:to-emerald-600 transition-all">
            查看记录
          </button>
        </div>

        <!-- 系统设置 -->
        <div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition-shadow">
          <div class="text-4xl mb-4">⚙️</div>
          <h3 class="text-xl font-bold text-gray-800 mb-2">系统设置</h3>
          <p class="text-gray-600 text-sm mb-4">配置系统参数，管理平台设置</p>
          <button
            @click="router.push('/admin/settings')"
            class="w-full px-4 py-2 bg-gradient-to-r from-orange-500 to-red-500 text-white rounded-lg font-medium hover:from-orange-600 hover:to-red-600 transition-all"
          >
            进入设置
          </button>
        </div>

        <!-- 返回首页 -->
        <div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition-shadow">
          <div class="text-4xl mb-4">🏠</div>
          <h3 class="text-xl font-bold text-gray-800 mb-2">返回首页</h3>
          <p class="text-gray-600 text-sm mb-4">回到用户首页查看平台信息</p>
          <button
            @click="router.push('/')"
            class="w-full px-4 py-2 bg-gradient-to-r from-gray-500 to-gray-600 text-white rounded-lg font-medium hover:from-gray-600 hover:to-gray-700 transition-all"
          >
            返回首页
          </button>
        </div>
      </div>

      <!-- 提示信息 -->
      <div class="mt-8 text-center">
        <p class="text-sm text-gray-500">
          💡 提示：具体管理功能正在开发中，敬请期待
        </p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getUserInfo, clearAuth } from '../utils/auth'
import { logout } from '../utils/api'
import Navigation from '../components/Navigation.vue'

const router = useRouter()
const userInfo = ref(null)

onMounted(() => {
  userInfo.value = getUserInfo()

  // 如果不是管理员，跳转到首页
  if (!userInfo.value || !userInfo.value.is_admin) {
    router.push('/')
  }
})

// 处理导航事件
function handleNavigation(action) {
  switch (action) {
    case 'history':
      // TODO: 实现兑换记录页面
      alert('兑换记录功能即将上线～')
      break
    case 'home':
      router.push('/')
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
</script>
