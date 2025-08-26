<script setup>
import { reactive, onMounted } from 'vue'
import { LoadClaudeConfig, SaveClaudeConfig, LoadConfigProfiles, SaveConfigProfile, ApplyConfigProfile, DeleteConfigProfile } from '../../wailsjs/go/app/App'

const state = reactive({
  authToken: '',
  baseURL: '',
  loading: false,
  message: '',
  messageType: '',
  configLoaded: false,
  showPassword: false,
  // Profile management
  profiles: [],
  showProfileModal: false,
  newProfile: {
    name: '',
    description: '',
    authToken: '',
    baseURL: ''
  },
  activeView: 'config', // 'config' or 'profiles'
  // Confirmation dialog
  showConfirmDialog: false,
  confirmDialog: {
    title: '',
    message: '',
    onConfirm: null,
    onCancel: null
  }
})

onMounted(async () => {
  await loadConfig()
  await loadProfiles()
})

async function loadConfig(showMessages = true) {
  state.loading = true
  if (showMessages) {
    state.message = ''
  }
  
  try {
    const result = await LoadClaudeConfig()
    
    if (result.success && result.data) {
      state.authToken = result.data.env.ANTHROPIC_AUTH_TOKEN || ''
      state.baseURL = result.data.env.ANTHROPIC_BASE_URL || ''
      state.configLoaded = true
      if (showMessages) {
        showMessage('配置加载成功 ✨', 'success')
      }
    } else {
      state.configLoaded = false
      if (showMessages) {
        showMessage('配置文件不存在，将创建新配置 🚀', 'info')
      }
    }
  } catch (error) {
    if (showMessages) {
      showMessage('加载配置失败: ' + error.toString(), 'error')
    }
  } finally {
    state.loading = false
  }
}

async function saveConfig() {
  if (!state.authToken.trim()) {
    showMessage('请输入 Auth Token 🔑', 'error')
    return
  }
  
  if (!state.baseURL.trim()) {
    showMessage('请输入 Base URL 🌐', 'error')
    return
  }

  state.loading = true
  
  try {
    const result = await SaveClaudeConfig(state.authToken.trim(), state.baseURL.trim())
    
    if (result.success) {
      state.configLoaded = true
      showMessage('配置保存成功！🎉', 'success')
    } else {
      showMessage(result.message || '保存配置失败', 'error')
    }
  } catch (error) {
    showMessage('保存配置失败: ' + error.toString(), 'error')
  } finally {
    state.loading = false
  }
}

function showMessage(message, type) {
  state.message = message
  state.messageType = type
  
  if (type === 'success' || type === 'info') {
    setTimeout(() => {
      state.message = ''
      state.messageType = ''
    }, 4000)
  }
}

function clearMessage() {
  state.message = ''
  state.messageType = ''
}

function togglePasswordVisibility() {
  state.showPassword = !state.showPassword
}

// Profile management functions
async function loadProfiles() {
  try {
    const result = await LoadConfigProfiles()
    if (result.success) {
      state.profiles = result.data || []
    } else {
      state.profiles = []
    }
  } catch (error) {
    console.error('Failed to load profiles:', error)
    state.profiles = []
  }
}

async function saveProfile() {
  if (!state.newProfile.name.trim()) {
    showMessage('请输入配置文件名称 📝', 'error')
    return
  }
  
  if (!state.newProfile.authToken.trim()) {
    showMessage('请输入 Auth Token 🔑', 'error')
    return
  }
  
  if (!state.newProfile.baseURL.trim()) {
    showMessage('请输入 Base URL 🌐', 'error')
    return
  }

  state.loading = true
  
  try {
    const result = await SaveConfigProfile(
      state.newProfile.name.trim(),
      state.newProfile.authToken.trim(),
      state.newProfile.baseURL.trim(),
      state.newProfile.description.trim()
    )
    
    if (result.success) {
      showMessage('配置文件保存成功！🎉', 'success')
      state.showProfileModal = false
      state.newProfile = { name: '', description: '', authToken: '', baseURL: '' }
      await loadProfiles()
    } else {
      showMessage(result.message || '保存配置文件失败', 'error')
    }
  } catch (error) {
    showMessage('保存配置文件失败: ' + error.toString(), 'error')
  } finally {
    state.loading = false
  }
}

async function applyProfile(profileName) {
  state.loading = true
  
  try {
    const result = await ApplyConfigProfile(profileName)
    
    if (result.success) {
      showMessage(`已应用配置文件 "${profileName}" ✨`, 'success')
      await loadConfig(false)
      state.activeView = 'config'
    } else {
      showMessage(result.message || '应用配置文件失败', 'error')
    }
  } catch (error) {
    showMessage('应用配置文件失败: ' + error.toString(), 'error')
  } finally {
    state.loading = false
  }
}

async function deleteProfile(profileName) {
  const confirmed = await showConfirm(
    '删除配置文件', 
    `确定要删除配置文件 "${profileName}" 吗？此操作不可恢复。`
  )
  
  if (!confirmed) {
    return
  }
  
  state.loading = true
  
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
    state.loading = false
  }
}

function openProfileModal() {
  state.showProfileModal = true
  state.newProfile = { 
    name: '', 
    description: '', 
    authToken: '', // 默认为空
    baseURL: '' 
  }
}

function closeProfileModal() {
  state.showProfileModal = false
  state.newProfile = { name: '', description: '', authToken: '', baseURL: '' }
}

function switchView(view) {
  state.activeView = view
}

// Custom confirmation dialog functions
function showConfirm(title, message) {
  return new Promise((resolve) => {
    state.confirmDialog = {
      title,
      message,
      onConfirm: () => {
        state.showConfirmDialog = false
        resolve(true)
      },
      onCancel: () => {
        state.showConfirmDialog = false
        resolve(false)
      }
    }
    state.showConfirmDialog = true
  })
}

function closeConfirmDialog() {
  state.showConfirmDialog = false
  if (state.confirmDialog.onCancel) {
    state.confirmDialog.onCancel()
  }
}
</script>

<template>
  <div class="config-manager">
    <!-- Navigation Tabs -->
    <div class="nav-tabs">
      <button 
        @click="switchView('config')"
        :class="['nav-tab', { 'nav-tab--active': state.activeView === 'config' }]"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M12 1L3 5V11C3 16.55 6.84 21.74 12 23C17.16 21.74 21 16.55 21 11V5L12 1Z" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <span>当前配置</span>
      </button>
      <button 
        @click="switchView('profiles')"
        :class="['nav-tab', { 'nav-tab--active': state.activeView === 'profiles' }]"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2z" stroke="currentColor" stroke-width="2" fill="none"/>
          <path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <span>配置文件 ({{ state.profiles.length }})</span>
      </button>
    </div>

    <!-- Current Configuration View -->
    <div v-if="state.activeView === 'config'">
    <!-- Main Configuration Card -->
    <div class="main-card">
      <div class="card-glow"></div>
      
      <div class="card-header">
        <div class="header-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" fill="currentColor"/>
          </svg>
        </div>
        <h2>Claude Code 配置</h2>
        <p>安全管理您的 API 密钥和服务器配置</p>
      </div>
      
      <div class="card-body">
        <!-- Status Message -->
        <Transition name="message" mode="out-in">
          <div v-if="state.message" :class="['status-message', `status-message--${state.messageType}`]">
            <div class="message-content">
              <span class="message-text">{{ state.message }}</span>
              <button @click="clearMessage" class="message-close">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2"/>
                </svg>
              </button>
            </div>
          </div>
        </Transition>

        <!-- Configuration Form -->
        <form @submit.prevent="saveConfig" class="config-form">
          <div class="form-group">
            <label for="authToken" class="form-label">
              <div class="label-content">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                  <path d="M12 1L3 5V11C3 16.55 6.84 21.74 12 23C17.16 21.74 21 16.55 21 11V5L12 1Z" stroke="currentColor" stroke-width="2" fill="none"/>
                </svg>
                <span>Auth Token</span>
              </div>
              <span class="required">必填</span>
            </label>
            
            <div class="input-wrapper">
              <input 
                id="authToken"
                v-model="state.authToken"
                :type="state.showPassword ? 'text' : 'password'"
                class="form-input"
                placeholder="sk-ant-api03-xxxxxxxxxxxx"
                :disabled="state.loading"
              />
              <button 
                type="button" 
                @click="togglePasswordVisibility"
                class="password-toggle"
                :class="{ active: state.showPassword }"
              >
                <svg v-if="state.showPassword" width="20" height="20" viewBox="0 0 24 24" fill="none">
                  <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24M1 1l22 22" stroke="currentColor" stroke-width="2"/>
                </svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2" fill="none"/>
                  <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2" fill="none"/>
                </svg>
              </button>
            </div>
          </div>

          <div class="form-group">
            <label for="baseURL" class="form-label">
              <div class="label-content">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                  <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
                  <path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" stroke="currentColor" stroke-width="2" fill="none"/>
                  <path d="M2 12h20" stroke="currentColor" stroke-width="2"/>
                </svg>
                <span>Base URL</span>
              </div>
              <span class="required">必填</span>
            </label>
            
            <div class="input-wrapper">
              <input 
                id="baseURL"
                v-model="state.baseURL"
                type="url"
                class="form-input"
                placeholder="https://api.anthropic.com"
                :disabled="state.loading"
              />
            </div>
          </div>

          <div class="form-actions">
            <button 
              type="button" 
              @click="loadConfig" 
              :disabled="state.loading"
              class="btn btn--secondary"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <path d="M1 4v6h6M23 20v-6h-6" stroke="currentColor" stroke-width="2"/>
                <path d="M20.49 9A9 9 0 005.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 003.51 15" stroke="currentColor" stroke-width="2" fill="none"/>
              </svg>
              {{ state.loading ? '加载中...' : '加载' }}
            </button>
            
            <button 
              type="submit" 
              :disabled="state.loading"
              class="btn btn--primary"
            >
              <svg v-if="!state.loading" width="16" height="16" viewBox="0 0 24 24" fill="none">
                <path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z" stroke="currentColor" stroke-width="2" fill="none"/>
                <polyline points="17,21 17,13 7,13 7,21" stroke="currentColor" stroke-width="2"/>
                <polyline points="7,3 7,8 15,8" stroke="currentColor" stroke-width="2"/>
              </svg>
              <svg v-else class="spinner" width="16" height="16" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" opacity="0.25"/>
                <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="4"/>
              </svg>
              {{ state.loading ? '保存中...' : '保存' }}
            </button>
            
            <button 
              type="button" 
              @click="openProfileModal"
              :disabled="state.loading"
              class="btn btn--accent"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
                <line x1="12" y1="8" x2="12" y2="16" stroke="currentColor" stroke-width="2"/>
                <line x1="8" y1="12" x2="16" y2="12" stroke="currentColor" stroke-width="2"/>
              </svg>
              另存为
            </button>
          </div>
        </form>
      </div>
    </div>
    </div>

    <!-- Profiles View -->
    <div v-else-if="state.activeView === 'profiles'">
      <!-- Status Message for Profiles -->
      <Transition name="message" mode="out-in">
        <div v-if="state.message" :class="['status-message', `status-message--${state.messageType}`]">
          <div class="message-content">
            <span class="message-text">{{ state.message }}</span>
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
        <button @click="openProfileModal" class="btn btn--primary">
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
        <div v-if="state.profiles.length === 0" class="empty-state">
          <div class="empty-icon">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
              <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2z" stroke="currentColor" stroke-width="1.5" fill="none"/>
              <path d="M8 10h8M8 14h8" stroke="currentColor" stroke-width="1.5"/>
            </svg>
          </div>
          <h3>暂无配置文件</h3>
          <p>创建您的第一个配置文件来快速切换不同的Claude Code设置</p>
          <button @click="openProfileModal" class="btn btn--primary btn--large">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
              <line x1="12" y1="8" x2="12" y2="16" stroke="currentColor" stroke-width="2"/>
              <line x1="8" y1="12" x2="16" y2="12" stroke="currentColor" stroke-width="2"/>
            </svg>
            创建配置文件
          </button>
        </div>
        
        <div v-else class="profiles-grid">
          <div v-for="profile in state.profiles" :key="profile.name" class="profile-card">
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
                  :disabled="state.loading"
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
                  @click="deleteProfile(profile.name)"
                  :disabled="state.loading"
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
    </div>

    <!-- Profile Modal -->
    <div v-if="state.showProfileModal" class="modal-overlay" @click="closeProfileModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>保存配置文件</h3>
          <button @click="closeProfileModal" class="modal-close">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveProfile">
            <div class="form-group">
              <label for="profileName" class="form-label">配置文件名称</label>
              <input 
                id="profileName"
                v-model="state.newProfile.name"
                type="text"
                class="form-input"
                placeholder="例如：开发环境、生产环境"
                required
                maxlength="50"
              />
            </div>
            
            <div class="form-group">
              <label for="profileAuthToken" class="form-label">
                <div class="label-content">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                    <path d="M12 1L3 5V11C3 16.55 6.84 21.74 12 23C17.16 21.74 21 16.55 21 11V5L12 1Z" stroke="currentColor" stroke-width="2" fill="none"/>
                  </svg>
                  <span>Auth Token</span>
                </div>
                <span class="required">必填</span>
              </label>
              <div class="input-wrapper">
                <input 
                  id="profileAuthToken"
                  v-model="state.newProfile.authToken"
                  type="password"
                  class="form-input"
                  placeholder="sk-ant-api03-xxxxxxxxxxxx"
                  required
                />
              </div>
            </div>
            
            <div class="form-group">
              <label for="profileBaseURL" class="form-label">
                <div class="label-content">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
                    <path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" stroke="currentColor" stroke-width="2" fill="none"/>
                    <path d="M2 12h20" stroke="currentColor" stroke-width="2"/>
                  </svg>
                  <span>Base URL</span>
                </div>
                <span class="required">必填</span>
              </label>
              <div class="input-wrapper">
                <input 
                  id="profileBaseURL"
                  v-model="state.newProfile.baseURL"
                  type="url"
                  class="form-input"
                  placeholder="https://api.anthropic.com"
                  required
                />
              </div>
            </div>
            
            <div class="form-group">
              <label for="profileDescription" class="form-label">描述（可选）</label>
              <textarea 
                id="profileDescription"
                v-model="state.newProfile.description"
                class="form-textarea"
                placeholder="描述这个配置文件的用途..."
                rows="3"
                maxlength="200"
              ></textarea>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button @click="closeProfileModal" class="btn btn--secondary">取消</button>
          <button @click="saveProfile" :disabled="state.loading" class="btn btn--primary">
            <svg v-if="!state.loading" width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z" stroke="currentColor" stroke-width="2" fill="none"/>
              <polyline points="17,21 17,13 7,13 7,21" stroke="currentColor" stroke-width="2"/>
              <polyline points="7,3 7,8 15,8" stroke="currentColor" stroke-width="2"/>
            </svg>
            <svg v-else class="spinner" width="16" height="16" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" opacity="0.25"/>
              <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="4"/>
            </svg>
            {{ state.loading ? '保存中...' : '保存配置文件' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Status Panel for Config View -->
    <div v-if="state.activeView === 'config'" class="status-panel">
      <div class="status-header">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2" fill="none"/>
          <path d="M12 1v6m0 6v6m11-7h-6m-6 0H1" stroke="currentColor" stroke-width="2"/>
        </svg>
        <h3>系统状态</h3>
      </div>
      
      <div class="status-items">
        <div class="status-item">
          <div class="status-indicator">
            <div :class="['indicator-dot', state.configLoaded ? 'dot--success' : 'dot--warning']"></div>
            <span class="status-label">配置状态</span>
          </div>
          <span :class="['status-value', state.configLoaded ? 'value--success' : 'value--warning']">
            {{ state.configLoaded ? '✅ 已加载' : '⚠️ 未找到' }}
          </span>
        </div>
        
        <div class="status-item">
          <div class="status-indicator">
            <div class="indicator-dot dot--info"></div>
            <span class="status-label">配置文件</span>
          </div>
          <span class="status-value value--info">~/.claude/settings.json</span>
        </div>
      </div>
    </div>
  </div>

  <!-- Confirmation Dialog -->
  <div v-if="state.showConfirmDialog" class="modal-overlay" @click="closeConfirmDialog">
    <div class="modal confirmation-modal" @click.stop>
      <div class="modal-header">
        <h3>{{ state.confirmDialog.title }}</h3>
      </div>
      <div class="modal-body">
        <div class="confirmation-content">
          <div class="confirmation-icon">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
              <path d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" fill="none"/>
            </svg>
          </div>
          <p class="confirmation-message">{{ state.confirmDialog.message }}</p>
        </div>
      </div>
      <div class="modal-footer">
        <button @click="state.confirmDialog.onCancel && state.confirmDialog.onCancel()" class="btn btn--secondary">取消</button>
        <button @click="state.confirmDialog.onConfirm && state.confirmDialog.onConfirm()" class="btn btn--danger">确认删除</button>
      </div>
    </div>
  </div>
</template>

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

/* Profiles Header */
.profiles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 16px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  margin-bottom: 2rem;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, #1e3a8a, #3b82f6);
  border-radius: 16px;
  color: white;
  box-shadow: 0 8px 16px rgba(30, 58, 138, 0.3);
}

.header-text h2 {
  margin: 0 0 0.25rem;
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a202c;
}

.header-text p {
  margin: 0;
  color: #64748b;
  font-size: 0.9rem;
}

/* Profiles List */
.profiles-list {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 16px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #64748b;
}

.empty-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 96px;
  height: 96px;
  background: linear-gradient(135deg, #f8fafc, #e2e8f0);
  border-radius: 24px;
  color: #94a3b8;
  margin-bottom: 1.5rem;
}

.empty-state h3 {
  margin: 0 0 0.5rem;
  font-size: 1.25rem;
  font-weight: 600;
  color: #374151;
}

.empty-state p {
  margin: 0 0 2rem;
  font-size: 0.95rem;
  line-height: 1.5;
}

.profiles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
}

.profile-card {
  background: linear-gradient(135deg, #ffffff, #f8fafc);
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  transition: all 0.3s ease;
}

.profile-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
  border-color: #cbd5e0;
}

.profile-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 1.5rem;
  gap: 1rem;
}

.profile-info {
  flex: 1;
  min-width: 0;
}

.profile-name {
  margin: 0 0 0.5rem;
  font-size: 1.125rem;
  font-weight: 600;
  color: #1a202c;
  word-break: break-word;
}

.profile-description {
  margin: 0 0 0.75rem;
  font-size: 0.875rem;
  color: #64748b;
  line-height: 1.4;
  word-break: break-word;
}

.profile-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.profile-date {
  font-size: 0.75rem;
  color: #94a3b8;
  font-weight: 500;
}

.profile-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex-shrink: 0;
}

.profile-details {
  padding: 0 1.5rem 1.5rem;
  border-top: 1px solid #f1f5f9;
  background: #fafbfc;
}

.profile-field {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 0;
  font-size: 0.875rem;
}

.field-label {
  font-weight: 500;
  color: #6b7280;
}

.field-value {
  font-weight: 500;
  color: #374151;
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  word-break: break-all;
}

/* Enhanced Button Styles */
.btn--accent {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  box-shadow: 0 8px 16px rgba(59, 130, 246, 0.3);
}

.btn--accent:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 12px 24px rgba(59, 130, 246, 0.4);
}

.btn--small {
  padding: 0.5rem 0.875rem;
  font-size: 0.875rem;
  border-radius: 12px;
}

.btn--large {
  padding: 1.125rem 2rem;
  font-size: 1.125rem;
  border-radius: 18px;
}

.btn--danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.btn--danger:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(239, 68, 68, 0.4);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal {
  background: white;
  border-radius: 24px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 32px 64px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  transform: translateY(0);
  transition: all 0.3s ease;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2rem 2rem 0;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a202c;
}

.modal-close {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 12px;
  transition: background-color 0.2s ease;
  color: #6b7280;
}

.modal-close:hover {
  background-color: #f3f4f6;
  color: #374151;
}

.modal-body {
  padding: 2rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 0 2rem 2rem;
}

.form-textarea {
  width: 100%;
  padding: 1rem 1.25rem;
  border: 2px solid #e5e7eb;
  border-radius: 16px;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.3s ease;
  background-color: #ffffff;
  color: #374151;
  font-family: inherit;
  resize: vertical;
  min-height: 80px;
}

.form-textarea:focus {
  outline: none;
  border-color: #1e40af;
  box-shadow: 0 0 0 4px rgba(30, 64, 175, 0.1);
  background-color: #fefefe;
}

.current-config-preview {
  margin-top: 1.5rem;
  padding: 1.25rem;
  background: linear-gradient(135deg, #f8fafc, #f1f5f9);
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.current-config-preview h4 {
  margin: 0 0 1rem;
  font-size: 0.95rem;
  font-weight: 600;
  color: #374151;
}

.config-preview-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
  font-size: 0.875rem;
}

.config-preview-item:not(:last-child) {
  border-bottom: 1px solid #e5e7eb;
}

.preview-label {
  font-weight: 500;
  color: #6b7280;
}

.preview-value {
  font-weight: 500;
  color: #374151;
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  word-break: break-all;
}

.main-card {
  position: relative;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  box-shadow: 0 32px 64px rgba(0, 0, 0, 0.12);
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-glow {
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(45deg, #1e3a8a, #3b82f6, #60a5fa);
  border-radius: 26px;
  z-index: -1;
  opacity: 0.6;
  filter: blur(8px);
}

.card-header {
  padding: 2.5rem 2rem 1.5rem;
  text-align: center;
  background: linear-gradient(135deg, rgba(30, 58, 138, 0.05), rgba(59, 130, 246, 0.05));
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #1e3a8a, #3b82f6);
  border-radius: 20px;
  color: white;
  margin-bottom: 1rem;
  box-shadow: 0 12px 24px rgba(30, 58, 138, 0.3);
}

.card-header h2 {
  margin: 0 0 0.5rem;
  font-size: 2rem;
  font-weight: 700;
  color: #1a202c;
  letter-spacing: -0.025em;
}

.card-header p {
  margin: 0;
  color: #64748b;
  font-size: 1rem;
  font-weight: 400;
}

.card-body {
  padding: 2rem;
}

/* Message Styles */
.status-message {
  border-radius: 16px;
  margin-bottom: 2rem;
  overflow: hidden;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.message-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
}

.message-text {
  font-weight: 500;
  font-size: 0.95rem;
}

.message-close {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 8px;
  transition: background-color 0.2s ease;
  color: inherit;
}

.message-close:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

.status-message--success {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
}

.status-message--error {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.status-message--info {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
}

/* Message Transitions */
.message-enter-active, .message-leave-active {
  transition: all 0.3s ease;
}

.message-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.message-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* Form Styles */
.config-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.form-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: #374151;
  font-size: 0.95rem;
}

.label-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #1e40af;
}

.required {
  color: #ef4444;
  font-size: 0.85rem;
  font-weight: 500;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.form-input {
  width: 100%;
  padding: 1rem 1.25rem;
  border: 2px solid #e5e7eb;
  border-radius: 16px;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.3s ease;
  background-color: #ffffff;
  color: #374151;
}

.input-wrapper .form-input {
  padding-right: 3.5rem; /* Add more right padding for password toggle button */
}

.form-input:focus {
  outline: none;
  border-color: #1e40af;
  box-shadow: 0 0 0 4px rgba(30, 64, 175, 0.1);
  background-color: #fefefe;
}

.form-input:disabled {
  background-color: #f9fafb;
  cursor: not-allowed;
  opacity: 0.7;
}

.password-toggle {
  position: absolute;
  right: 1rem;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 8px;
  color: #9ca3af;
  transition: all 0.2s ease;
}

.password-toggle:hover,
.password-toggle.active {
  color: #1e40af;
  background-color: rgba(30, 64, 175, 0.1);
}

/* Button Styles */
.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  border: none;
  border-radius: 16px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.7;
  transform: none !important;
}

.btn--primary {
  background: linear-gradient(135deg, #1e40af, #2563eb);
  color: white;
  box-shadow: 0 8px 16px rgba(30, 64, 175, 0.3);
}

.btn--primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 12px 24px rgba(30, 64, 175, 0.4);
}

.btn--secondary {
  background: linear-gradient(135deg, #f8fafc, #f1f5f9);
  color: #475569;
  border: 2px solid #e2e8f0;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
}

.btn--secondary:hover:not(:disabled) {
  background: linear-gradient(135deg, #f1f5f9, #e2e8f0);
  border-color: #cbd5e0;
  transform: translateY(-1px);
}

.spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Status Panel */
.status-panel {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 1.5rem;
  box-shadow: 0 16px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.status-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.25rem;
  color: #1e40af;
}

.status-header h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: #374151;
}

.status-items {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.indicator-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  position: relative;
}

.indicator-dot::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  animation: pulse-dot 2s ease-in-out infinite;
}

.dot--success {
  background-color: #3b82f6;
}

.dot--success::after {
  background-color: #3b82f6;
}

.dot--warning {
  background-color: #f59e0b;
}

.dot--warning::after {
  background-color: #f59e0b;
}

.dot--info {
  background-color: #3b82f6;
}

.dot--info::after {
  background-color: #3b82f6;
}

@keyframes pulse-dot {
  0% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(1.2); }
  100% { opacity: 1; transform: scale(1); }
}

.status-label {
  font-weight: 500;
  color: #6b7280;
  font-size: 0.9rem;
}

.status-value {
  font-weight: 600;
  font-size: 0.9rem;
}

.value--success {
  color: #3b82f6;
}

.value--warning {
  color: #f59e0b;
}

.value--info {
  color: #6b7280;
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  font-size: 0.8rem;
}

/* Responsive Design */
@media (max-width: 640px) {
  .card-header,
  .card-body {
    padding-left: 1.5rem;
    padding-right: 1.5rem;
  }
  
  .card-header h2 {
    font-size: 1.75rem;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .status-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .status-indicator {
    width: 100%;
  }
}

/* Confirmation Dialog */
.confirmation-modal {
  max-width: 480px;
}

.confirmation-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 1.5rem;
}

.confirmation-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #fef3c7, #fed7aa);
  border-radius: 20px;
  color: #d97706;
}

.confirmation-message {
  margin: 0;
  font-size: 1.1rem;
  line-height: 1.5;
  color: #374151;
}

.modal-footer .btn--danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.modal-footer .btn--danger:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(239, 68, 68, 0.4);
}
</style>