<template>
  <div class="profiles-section">
    <!-- Status Message for Profiles -->
    <Transition name="message" mode="out-in">
      <div v-if="message" :class="['status-message', `status-message--${messageType}`]">
        <div class="message-content">
          <span class="message-text">{{ message }}</span>
          <button @click="clearMessage" class="message-close">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2"/>
            </svg>
          </button>
        </div>
      </div>
    </Transition>

    <!-- Profiles Header -->
    <div class="profiles-header">
      <div class="header-content">
        <div class="header-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
            <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2z" stroke="currentColor" stroke-width="2" fill="none"/>
            <path d="M16 9l-4 4-4-4" stroke="currentColor" stroke-width="2"/>
          </svg>
        </div>
        <div class="header-text">
          <h2>配置文件管理</h2>
          <p>管理和切换您保存的配置文件</p>
        </div>
      </div>
      <button @click="$emit('openProfileModal')" class="btn btn--primary">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
          <line x1="12" y1="8" x2="12" y2="16" stroke="currentColor" stroke-width="2"/>
          <line x1="8" y1="12" x2="16" y2="12" stroke="currentColor" stroke-width="2"/>
        </svg>
        新建配置文件
      </button>
    </div>

    <!-- Profiles List -->
    <div class="profiles-list">
      <div v-if="profiles.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
            <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2z" stroke="currentColor" stroke-width="1.5" fill="none"/>
            <path d="M8 10h8M8 14h8" stroke="currentColor" stroke-width="1.5"/>
          </svg>
        </div>
        <h3>暂无配置文件</h3>
        <p>创建您的第一个配置文件来快速切换不同的Claude Code设置</p>
        <button @click="$emit('openProfileModal')" class="btn btn--primary btn--large">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
            <line x1="12" y1="8" x2="12" y2="16" stroke="currentColor" stroke-width="2"/>
            <line x1="8" y1="12" x2="16" y2="12" stroke="currentColor" stroke-width="2"/>
          </svg>
          创建配置文件
        </button>
      </div>
      
      <div v-else class="profiles-grid">
        <div v-for="profile in profiles" :key="profile.name" class="profile-card">
          <div class="profile-header">
            <div class="profile-info">
              <h3 class="profile-name">{{ profile.name }}</h3>
              <p v-if="profile.description" class="profile-description">{{ profile.description }}</p>
              <div class="profile-meta">
                <span class="profile-date">{{ new Date(profile.updatedAt).toLocaleDateString('zh-CN') }}</span>
              </div>
            </div>
            <div class="profile-actions">
              <button 
                @click="applyProfile(profile.name)"
                :disabled="loading"
                class="btn btn--small btn--primary"
                title="应用此配置文件"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <path d="M9 12l2 2 4-4" stroke="currentColor" stroke-width="2" fill="none"/>
                  <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2" fill="none"/>
                </svg>
                应用
              </button>
              <button 
                @click="handleDeleteProfile(profile.name)"
                :disabled="loading"
                class="btn btn--small btn--danger"
                title="删除此配置文件"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <polyline points="3,6 5,6 21,6" stroke="currentColor" stroke-width="2"/>
                  <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" stroke="currentColor" stroke-width="2" fill="none"/>
                </svg>
                删除
              </button>
            </div>
          </div>
          <div class="profile-details">
            <div class="profile-field">
              <span class="field-label">Base URL:</span>
              <span class="field-value">{{ profile.baseURL || 'N/A' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 确认删除对话框 -->
    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="deleteConfirm.title"
      :message="deleteConfirm.message"
      @confirm="confirmDelete"
      @cancel="cancelDelete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, defineEmits } from 'vue'
import { LoadConfigProfiles, ApplyConfigProfile, DeleteConfigProfile } from '../../wailsjs/go/app/App'
import ConfirmDialog from './ConfirmDialog.vue'

const emit = defineEmits(['openProfileModal', 'profileApplied'])

// 响应式状态
const profiles = ref([])
const loading = ref(false)
const message = ref('')
const messageType = ref('')
const showDeleteConfirm = ref(false)
const deleteConfirm = ref({
  title: '',
  message: '',
  profileName: ''
})

// 初始化时加载配置文件
onMounted(async () => {
  await loadProfiles()
})

// 加载所有配置文件
async function loadProfiles() {
  try {
    const result = await LoadConfigProfiles()
    if (result.success) {
      profiles.value = result.data || []
    } else {
      profiles.value = []
    }
  } catch (error) {
    console.error('Failed to load profiles:', error)
    profiles.value = []
  }
}

// 应用配置文件
async function applyProfile(profileName) {
  loading.value = true
  
  try {
    const result = await ApplyConfigProfile(profileName)
    
    if (result.success) {
      showMessage(`已应用配置文件 "${profileName}" ✨`, 'success')
      emit('profileApplied', profileName)
    } else {
      showMessage(result.message || '应用配置文件失败', 'error')
    }
  } catch (error) {
    showMessage('应用配置文件失败: ' + error.toString(), 'error')
  } finally {
    loading.value = false
  }
}

// 处理删除配置文件请求
function handleDeleteProfile(profileName) {
  deleteConfirm.value = {
    title: '删除配置文件',
    message: `确定要删除配置文件 "${profileName}" 吗？此操作不可恢复。`,
    profileName
  }
  showDeleteConfirm.value = true
}

// 确认删除
async function confirmDelete() {
  const profileName = deleteConfirm.value.profileName
  showDeleteConfirm.value = false
  loading.value = true
  
  try {
    const result = await DeleteConfigProfile(profileName)
    
    if (result.success) {
      showMessage(`配置文件 "${profileName}" 已删除 🗑️`, 'success')
      await loadProfiles()
    } else {
      showMessage(result.message || '删除配置文件失败', 'error')
    }
  } catch (error) {
    showMessage('删除配置文件失败: ' + error.toString(), 'error')
  } finally {
    loading.value = false
  }
}

// 取消删除
function cancelDelete() {
  showDeleteConfirm.value = false
}

// 显示消息
function showMessage(msg, type) {
  message.value = msg
  messageType.value = type
  
  if (type === 'success' || type === 'info') {
    setTimeout(() => {
      message.value = ''
      messageType.value = ''
    }, 4000)
  }
}

// 清除消息
function clearMessage() {
  message.value = ''
  messageType.value = ''
}

// 暴露加载方法供父组件调用
defineExpose({
  loadProfiles
})
</script>

<style scoped>
/* 样式将在后面的步骤中统一整理 */
</style>