<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const currentTime = ref(new Date().toLocaleString('zh-CN'))
let timer = null

onMounted(() => {
  timer = setInterval(() => {
    currentTime.value = new Date().toLocaleString('zh-CN')
  }, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<template>
  <div class="min-h-screen bg-white flex items-center justify-center p-4">
    <div class="text-center">
      <!-- 主标题 -->
      <h1 class="text-6xl md:text-8xl font-bold text-black mb-4">
        兑兑猫
      </h1>

      <!-- 英文副标题 -->
      <p class="text-xl md:text-2xl text-gray-600 mb-12">
        DuiDuiMao CDK Exchange Platform
      </p>

      <!-- 状态信息 -->
      <div class="space-y-4">
        <div class="text-2xl md:text-3xl text-gray-500 font-medium">
          服务即将开通
        </div>

        <!-- 实时时间 -->
        <div class="text-lg md:text-xl text-gray-400 font-mono">
          {{ currentTime }}
        </div>
      </div>

      <!-- 底部描述 -->
      <div class="mt-16 text-sm md:text-base text-gray-400">
        基于 LinuxDo Connect 的公益 CDK 兑换平台
      </div>

      <!-- 装饰线 -->
      <div class="mt-8 flex justify-center">
        <div class="w-32 h-0.5 bg-black"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 确保严格黑白配色，无蓝紫色，无渐变 */
</style>
