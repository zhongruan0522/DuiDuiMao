<template>
  <div class="min-h-screen flex items-center justify-center bg-white dark:bg-gray-900 transition-colors duration-300">
    <div class="text-center px-6">
      <!-- 404 大标题 -->
      <div class="relative">
        <h1 class="text-[12rem] font-bold text-gray-200 dark:text-gray-800 select-none animate-pulse">
          404
        </h1>

        <!-- 猫咪图标 (使用 CSS 画一个简单的猫头) -->
        <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
          <div class="cat-face">
            <div class="cat-ear cat-ear-left"></div>
            <div class="cat-ear cat-ear-right"></div>
            <div class="cat-head">
              <div class="cat-eyes">
                <div class="cat-eye"></div>
                <div class="cat-eye"></div>
              </div>
              <div class="cat-nose"></div>
              <div class="cat-whiskers">
                <div class="whisker whisker-left"></div>
                <div class="whisker whisker-right"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 提示文字 -->
      <h2 class="text-2xl md:text-3xl font-semibold text-gray-800 dark:text-gray-200 mb-4 mt-8">
        页面走丢啦
      </h2>
      <p class="text-gray-600 dark:text-gray-400 mb-8 max-w-md mx-auto">
        你访问的页面不存在，可能是链接错误或页面已被移除。
      </p>

      <!-- 返回按钮 -->
      <button
        @click="goHome"
        class="inline-flex items-center px-6 py-3 bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900 font-medium rounded-lg hover:bg-gray-800 dark:hover:bg-gray-200 transition-all duration-200 transform hover:scale-105 active:scale-95"
      >
        <svg
          class="w-5 h-5 mr-2"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
          />
        </svg>
        返回首页
      </button>

      <!-- 当前路径提示 -->
      <p class="mt-6 text-sm text-gray-500 dark:text-gray-500">
        当前路径: <code class="bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">{{ currentPath }}</code>
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const currentPath = computed(() => route.fullPath)

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
/* 猫咪图标样式 */
.cat-face {
  position: relative;
  width: 120px;
  height: 120px;
}

.cat-head {
  position: relative;
  width: 80px;
  height: 80px;
  background-color: currentColor;
  border-radius: 50%;
  margin: 20px auto 0;
  animation: float 3s ease-in-out infinite;
}

.cat-ear {
  position: absolute;
  width: 0;
  height: 0;
  border-left: 20px solid transparent;
  border-right: 20px solid transparent;
  border-bottom: 30px solid currentColor;
}

.cat-ear-left {
  top: 0;
  left: 10px;
  transform: rotate(-20deg);
}

.cat-ear-right {
  top: 0;
  right: 10px;
  transform: rotate(20deg);
}

.cat-eyes {
  display: flex;
  justify-content: center;
  gap: 20px;
  padding-top: 25px;
}

.cat-eye {
  width: 8px;
  height: 8px;
  background-color: white;
  border-radius: 50%;
  animation: blink 3s infinite;
}

.cat-nose {
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-top: 8px solid white;
  margin: 8px auto 0;
}

.cat-whiskers {
  position: relative;
  width: 100%;
  height: 20px;
}

.whisker {
  position: absolute;
  width: 30px;
  height: 1px;
  background-color: currentColor;
  opacity: 0.6;
}

.whisker-left {
  left: -15px;
  top: 10px;
  transform: rotate(-15deg);
}

.whisker-right {
  right: -15px;
  top: 10px;
  transform: rotate(15deg);
}

/* 暗黑模式下的猫咪颜色 */
@media (prefers-color-scheme: dark) {
  .cat-head {
    color: #f3f4f6;
  }

  .cat-eye {
    background-color: #1f2937;
  }

  .cat-nose {
    border-top-color: #1f2937;
  }
}

@media (prefers-color-scheme: light) {
  .cat-head {
    color: #1f2937;
  }
}

/* 动画效果 */
@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

@keyframes blink {
  0%, 90%, 100% {
    opacity: 1;
  }
  95% {
    opacity: 0;
  }
}

/* 暗黑模式支持 */
@media (prefers-color-scheme: dark) {
  .dark\:bg-gray-900 {
    background-color: #111827;
  }

  .dark\:bg-gray-800 {
    background-color: #1f2937;
  }

  .dark\:bg-gray-100 {
    background-color: #f3f4f6;
  }

  .dark\:text-gray-200 {
    color: #e5e7eb;
  }

  .dark\:text-gray-400 {
    color: #9ca3af;
  }

  .dark\:text-gray-500 {
    color: #6b7280;
  }

  .dark\:text-gray-800 {
    color: #1f2937;
  }

  .dark\:text-gray-900 {
    color: #111827;
  }

  .dark\:hover\:bg-gray-200:hover {
    background-color: #e5e7eb;
  }
}
</style>
