<template>
  <nav class="bg-white/80 backdrop-blur-sm border-b border-gray-100 relative" style="z-index: 10000;">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo区域 -->
        <div class="flex items-center gap-2">
          <span class="text-2xl">🐱</span>
          <span class="font-bold text-xl bg-gradient-to-r from-pink-400 to-purple-500 bg-clip-text text-transparent">
            {{ title }}
          </span>
        </div>

        <!-- 用户菜单 -->
        <div v-if="userInfo" class="relative">
          <button
            @click="toggleMenu"
            class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <div class="w-8 h-8 rounded-full bg-gradient-to-r from-pink-400 to-purple-500 flex items-center justify-center text-white font-medium">
              {{ userInfo.name?.charAt(0) || '🐱' }}
            </div>
            <div class="text-left hidden sm:block">
              <p class="text-sm font-medium text-gray-700">{{ userInfo.name }}</p>
              <p class="text-xs text-gray-500">
                {{ userInfo.is_admin ? '管理员' : `Lv.${userInfo.trust_level}` }}
              </p>
            </div>
            <svg
              class="w-4 h-4 text-gray-400 transition-transform"
              :class="{ 'rotate-180': showMenu }"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          <!-- 下拉菜单 -->
          <transition
            enter-active-class="transition ease-out duration-100"
            enter-from-class="transform opacity-0 scale-95"
            enter-to-class="transform opacity-100 scale-100"
            leave-active-class="transition ease-in duration-75"
            leave-from-class="transform opacity-100 scale-100"
            leave-to-class="transform opacity-0 scale-95"
          >
            <div
              v-if="showMenu"
              class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-lg border border-gray-100 py-2"
              style="z-index: 9999;"
            >
              <!-- 用户信息卡片 -->
              <div class="px-4 py-3 border-b border-gray-100">
                <div class="flex items-center gap-3 mb-3">
                  <div class="w-12 h-12 rounded-full bg-gradient-to-r from-pink-400 to-purple-500 flex items-center justify-center text-white text-lg font-medium">
                    {{ userInfo.name?.charAt(0) || '🐱' }}
                  </div>
                  <div>
                    <p class="font-semibold text-gray-800">{{ userInfo.name }}</p>
                    <p class="text-xs text-gray-500">@{{ userInfo.username }}</p>
                  </div>
                </div>
                <div class="grid grid-cols-2 gap-2 text-xs">
                  <div class="bg-pink-50 rounded-lg px-2 py-1">
                    <p class="text-gray-600">信任等级</p>
                    <p class="font-semibold text-pink-600">Lv.{{ userInfo.trust_level }}</p>
                  </div>
                  <div class="bg-purple-50 rounded-lg px-2 py-1">
                    <p class="text-gray-600">账户类型</p>
                    <p class="font-semibold text-purple-600">{{ userInfo.is_admin ? '管理员' : '普通用户' }}</p>
                  </div>
                </div>
              </div>

              <!-- 菜单项 -->
              <div class="py-1">
                <!-- 我的兑换记录 -->
                <button
                  @click="handleMenuClick('history')"
                  class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-pink-50 hover:text-pink-600 transition-colors flex items-center gap-2"
                >
                  <span class="text-base">📜</span>
                  <span>我的兑换记录</span>
                </button>

                <!-- 管理后台（仅管理员可见） -->
                <button
                  v-if="userInfo.is_admin"
                  @click="handleMenuClick('admin')"
                  class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-purple-50 hover:text-purple-600 transition-colors flex items-center gap-2"
                >
                  <span class="text-base">🎛️</span>
                  <span>管理后台</span>
                </button>

                <!-- 返回首页（仅在管理后台显示） -->
                <button
                  v-if="showHomeButton"
                  @click="handleMenuClick('home')"
                  class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors flex items-center gap-2"
                >
                  <span class="text-base">🏠</span>
                  <span>返回首页</span>
                </button>

                <!-- 退出登录 -->
                <button
                  @click="handleMenuClick('logout')"
                  class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 hover:text-red-600 transition-colors flex items-center gap-2 border-t border-gray-100 mt-1"
                >
                  <span class="text-base">👋</span>
                  <span>退出登录</span>
                </button>
              </div>
            </div>
          </transition>
        </div>

        <!-- 未登录状态 -->
        <div v-else>
          <button
            @click="$emit('login')"
            class="px-6 py-2 bg-gradient-to-r from-pink-400 to-purple-500 text-white rounded-lg font-medium hover:from-pink-500 hover:to-purple-600 transition-all"
          >
            登录
          </button>
        </div>
      </div>
    </div>

    <!-- 点击外部关闭菜单 -->
    <div
      v-if="showMenu"
      @click="showMenu = false"
      class="fixed inset-0"
      style="z-index: 9998;"
    ></div>
  </nav>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// 组件属性
const props = defineProps({
  title: {
    type: String,
    default: '兑兑猫'
  },
  userInfo: {
    type: Object,
    default: null
  },
  showHomeButton: {
    type: Boolean,
    default: false
  }
})

// 组件事件
const emit = defineEmits(['login', 'logout', 'navigate'])

const showMenu = ref(false)

// 切换菜单显示
function toggleMenu() {
  showMenu.value = !showMenu.value
}

// ESC键关闭菜单
function handleEscape(event) {
  if (event.key === 'Escape' && showMenu.value) {
    showMenu.value = false
  }
}

// 处理菜单点击
function handleMenuClick(action) {
  showMenu.value = false
  emit('navigate', action)
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})
</script>
