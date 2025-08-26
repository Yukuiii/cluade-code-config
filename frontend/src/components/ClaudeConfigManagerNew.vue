<template>
  <div class="config-manager">
    <!-- Navigation Tabs -->
    <div class="nav-tabs">
      <button 
        @click="activeView = 'config'"
        :class="['nav-tab', { 'nav-tab--active': activeView === 'config' }]"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M12 1L3 5V11C3 16.55 6.84 21.74 12 23C17.16 21.74 21 16.55 21 11V5L12 1Z" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <span>当前配置</span>
      </button>
      <button 
        @click="activeView = 'profiles'"
        :class="['nav-tab', { 'nav-tab--active': activeView === 'profiles' }]"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2z" stroke="currentColor" stroke-width="2" fill="none"/>
          <path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <span>配置文件 ({{ profileCount }})</span>
      </button>
    </div>

    <!-- 当前配置视图 -->
    <ConfigManager 
      v-if="activeView === 'config'"
      @openProfileModal="showProfileModal = true"
    />

    <!-- 配置文件管理视图 -->
    <ProfileManager
      v-else-if="activeView === 'profiles'"
      ref="profileManagerRef"
      @openProfileModal="showProfileModal = true"
      @profileApplied="handleProfileApplied"
    />

    <!-- 配置文件保存模态框 -->
    <ProfileModal
      :show="showProfileModal"
      @close="showProfileModal = false"
      @saved="handleProfileSaved"
      @error="handleModalError"
    />

    <!-- 全局消息显示 -->
    <Transition name="message" mode="out-in">
      <div v-if="globalMessage" :class="['status-message', `status-message--${globalMessageType}`]">
        <div class="message-content">
          <span class="message-text">{{ globalMessage }}</span>
          <button @click="clearGlobalMessage" class="message-close">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2"/>
            </svg>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import ConfigManager from './ConfigManager.vue'
import ProfileManager from './ProfileManager.vue'
import ProfileModal from './ProfileModal.vue'
import '../assets/common.css'

// 响应式状态
const activeView = ref('config')
const showProfileModal = ref(false)
const profileManagerRef = ref(null)
const globalMessage = ref('')
const globalMessageType = ref('')

// 计算属性
const profileCount = computed(() => {
  return profileManagerRef.value?.profiles?.length || 0
})

// 处理配置文件应用成功
function handleProfileApplied(profileName) {
  // 切换到配置视图并显示成功消息
  activeView.value = 'config'
}

// 处理配置文件保存成功
function handleProfileSaved(message) {
  showGlobalMessage(message, 'success')
  // 刷新配置文件列表
  if (profileManagerRef.value) {
    profileManagerRef.value.loadProfiles()
  }
}

// 处理模态框错误
function handleModalError(message) {
  showGlobalMessage(message, 'error')
}

// 显示全局消息
function showGlobalMessage(message, type) {
  globalMessage.value = message
  globalMessageType.value = type
  
  if (type === 'success' || type === 'info') {
    setTimeout(() => {
      globalMessage.value = ''
      globalMessageType.value = ''
    }, 4000)
  }
}

// 清除全局消息
function clearGlobalMessage() {
  globalMessage.value = ''
  globalMessageType.value = ''
}
</script>

<style scoped>
.config-manager {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Navigation Tabs */
.nav-tabs {
  display: flex;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 0.5rem;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  margin-bottom: 1rem;
}

.nav-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.875rem 1.25rem;
  border: none;
  border-radius: 12px;
  font-weight: 500;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
  background: transparent;
  color: #6b7280;
}

.nav-tab:hover {
  color: #1e40af;
  background: rgba(30, 64, 175, 0.05);
}

.nav-tab--active {
  background: linear-gradient(135deg, #1e3a8a, #1e40af);
  color: white;
  box-shadow: 0 4px 12px rgba(30, 58, 138, 0.4);
}

.nav-tab--active:hover {
  color: white;
  background: linear-gradient(135deg, #1e3a8a, #1e40af);
}

/* 全局消息样式 */
.status-message {
  position: fixed;
  top: 2rem;
  right: 2rem;
  z-index: 2000;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  min-width: 320px;
  max-width: 480px;
}

@media (max-width: 640px) {
  .status-message {
    left: 1rem;
    right: 1rem;
    top: 1rem;
    min-width: unset;
    max-width: unset;
  }
}
</style>