<template>
  <div class="profiles-section">
    <!-- Status Message for Profiles -->
    <Transition name="message" mode="out-in">
      <div v-if="message" :class="['status-message', `status-message--${messageType}`]">
        <div class="message-content">
          <span class="message-text">{{ message }}</span>
          <button @click="clearMessage" class="message-close">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" />
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
            <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2z" stroke="currentColor"
              stroke-width="2" fill="none" />
            <path d="M16 9l-4 4-4-4" stroke="currentColor" stroke-width="2" />
          </svg>
        </div>
        <div class="header-text">
          <h2>配置文件管理</h2>
          <p>管理和切换您保存的配置文件</p>
        </div>
      </div>
      <button @click="$emit('openProfileModal')" class="btn btn--primary">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none" />
          <line x1="12" y1="8" x2="12" y2="16" stroke="currentColor" stroke-width="2" />
          <line x1="8" y1="12" x2="16" y2="12" stroke="currentColor" stroke-width="2" />
        </svg>
        新建配置文件
      </button>
    </div>

    <!-- Profiles List -->
    <div class="profiles-list">
      <div v-if="profiles.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
            <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2z" stroke="currentColor"
              stroke-width="1.5" fill="none" />
            <path d="M8 10h8M8 14h8" stroke="currentColor" stroke-width="1.5" />
          </svg>
        </div>
        <h3>暂无配置文件</h3>
        <p>创建您的第一个配置文件来快速切换不同的Claude Code设置</p>
        <button @click="$emit('openProfileModal')" class="btn btn--primary btn--large">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none" />
            <line x1="12" y1="8" x2="12" y2="16" stroke="currentColor" stroke-width="2" />
            <line x1="8" y1="12" x2="16" y2="12" stroke="currentColor" stroke-width="2" />
          </svg>
          创建配置文件
        </button>
      </div>

      <div v-else class="profiles-grid">
        <div v-for="profile in profiles" :key="profile.name" class="profile-card">
          <div class="profile-header">
            <div class="profile-content-row">
              <div class="profile-info">
                <h3 class="profile-name">{{ profile.name }}</h3>
                <p v-if="profile.description" class="profile-description">{{ profile.description }}</p>
                <div v-if="profile.baseURL" class="profile-url">{{ profile.baseURL }}</div>
              </div>
              <div class="profile-timestamp">
                <span class="profile-date">
                  {{ new Date(profile.updatedAt).toLocaleDateString('zh-CN', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric'
                  }) }}
                </span>
              </div>
            </div>
          </div>
          <div class="profile-actions">
            <button @click="applyProfile(profile.name)" :disabled="loading" class="btn btn--primary btn--apply"
              title="应用此配置文件">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <path d="M9 12l2 2 4-4" stroke="currentColor" stroke-width="2" fill="none" />
                <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2" fill="none" />
              </svg>
              应用
            </button>
            <button @click="handleDeleteProfile(profile.name)" :disabled="loading" class="btn btn--ghost btn--danger"
              title="删除此配置文件">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <polyline points="3,6 5,6 21,6" stroke="currentColor" stroke-width="2" />
                <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" stroke="currentColor"
                  stroke-width="2" fill="none" />
              </svg>
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 确认删除对话框 -->
    <ConfirmDialog :show="showDeleteConfirm" :title="deleteConfirm.title" :message="deleteConfirm.message"
      @confirm="confirmDelete" @cancel="cancelDelete" />
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
/* 主容器样式 */
.profiles-section {
  padding: 1.5rem;
  max-width: 1200px;
  margin: 0 auto;
}

/* 状态消息样式 */
.status-message {
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  font-weight: 500;
}

.status-message--success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.status-message--error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.message-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.message-close {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.message-close:hover {
  opacity: 1;
}

/* 头部样式 */
.profiles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px;
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-text h2 {
  margin: 0 0 4px 0;
  font-size: 24px;
  font-weight: 700;
  color: #1a202c;
}

.header-text p {
  margin: 0;
  color: #718096;
  font-size: 14px;
}

/* 配置文件网格布局 */
.profiles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  margin-top: 1rem;
}

/* 配置文件卡片样式 */
.profile-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.profile-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

/* 配置文件头部 */
.profile-header {
  padding: 24px;
  border-bottom: 1px solid #f1f5f9;
}

/* 内容行 - 左侧信息，右侧时间 */
.profile-content-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
}

/* 左侧信息区域 */
.profile-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.profile-name {
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
  line-height: 1.4;
}

.profile-description {
  color: #6b7280;
  font-size: 14px;
  margin: 0;
  line-height: 1.5;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}

.profile-url {
  color: #6366f1;
  font-size: 12px;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', monospace;
  background: #f8fafc;
  padding: 6px 12px;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  max-width: fit-content;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 右侧时间区域 */
.profile-timestamp {
  flex-shrink: 0;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  margin-left: auto;
}

.profile-date {
  color: #6b7280;
  font-size: 13px;
  font-weight: 500;
  background: #f3f4f6;
  padding: 6px 12px;
  border-radius: 6px;
  white-space: nowrap;
  text-align: center;
  min-width: 80px;
  display: block;
}

/* 操作按钮区域 */
.profile-actions {
  padding: 16px 24px 20px 24px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  background: #f8fafc;
}

/* 按钮基础样式 */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  text-decoration: none;
  white-space: nowrap;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 主要按钮样式 */
.btn--primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: 1px solid transparent;
}

.btn--primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn--primary:active:not(:disabled) {
  transform: translateY(0);
}

/* 幽灵按钮样式 */
.btn--ghost {
  background: transparent;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

.btn--ghost:hover:not(:disabled) {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

/* 危险按钮样式 */
.btn--danger {
  color: #dc2626;
  border-color: #fecaca;
}

.btn--danger:hover:not(:disabled) {
  background: #fef2f2;
  border-color: #fca5a5;
  color: #b91c1c;
}

/* 大按钮样式 */
.btn--large {
  padding: 12px 24px;
  font-size: 16px;
}

/* 空状态样式 */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #64748b;
}

.empty-icon {
  margin-bottom: 1rem;
  opacity: 0.6;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
  color: #374151;
}

.empty-state p {
  font-size: 14px;
  margin: 0 0 2rem 0;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
  line-height: 1.6;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .profiles-section {
    padding: 1rem;
  }

  .profiles-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
    text-align: center;
  }

  .profiles-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .profile-actions {
    flex-direction: column;
  }

  .btn {
    justify-content: center;
  }
}

/* 过渡动画 */
.message-enter-active,
.message-leave-active {
  transition: all 0.3s ease;
}

.message-enter-from,
.message-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>