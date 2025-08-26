<template>
  <div class="config-section">
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
                v-model="authToken"
                :type="showPassword ? 'text' : 'password'"
                class="form-input"
                placeholder="sk-ant-api03-xxxxxxxxxxxx"
                :disabled="loading"
              />
              <button 
                type="button" 
                @click="togglePasswordVisibility"
                class="password-toggle"
                :class="{ active: showPassword }"
              >
                <svg v-if="showPassword" width="20" height="20" viewBox="0 0 24 24" fill="none">
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
                v-model="baseURL"
                type="url"
                class="form-input"
                placeholder="https://api.anthropic.com"
                :disabled="loading"
              />
            </div>
          </div>

          <div class="form-actions">
            <button 
              type="button" 
              @click="loadConfig" 
              :disabled="loading"
              class="btn btn--secondary"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <path d="M1 4v6h6M23 20v-6h-6" stroke="currentColor" stroke-width="2"/>
                <path d="M20.49 9A9 9 0 005.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 003.51 15" stroke="currentColor" stroke-width="2" fill="none"/>
              </svg>
              {{ loading ? '加载中...' : '加载' }}
            </button>
            
            <button 
              type="submit" 
              :disabled="loading"
              class="btn btn--primary"
            >
              <svg v-if="!loading" width="16" height="16" viewBox="0 0 24 24" fill="none">
                <path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z" stroke="currentColor" stroke-width="2" fill="none"/>
                <polyline points="17,21 17,13 7,13 7,21" stroke="currentColor" stroke-width="2"/>
                <polyline points="7,3 7,8 15,8" stroke="currentColor" stroke-width="2"/>
              </svg>
              <svg v-else class="spinner" width="16" height="16" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" opacity="0.25"/>
                <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="4"/>
              </svg>
              {{ loading ? '保存中...' : '保存' }}
            </button>
            
            <button 
              type="button" 
              @click="$emit('openProfileModal')"
              :disabled="loading"
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

    <!-- Status Panel -->
    <div class="status-panel">
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
            <div :class="['indicator-dot', configLoaded ? 'dot--success' : 'dot--warning']"></div>
            <span class="status-label">配置状态</span>
          </div>
          <span :class="['status-value', configLoaded ? 'value--success' : 'value--warning']">
            {{ configLoaded ? '✅ 已加载' : '⚠️ 未找到' }}
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
</template>

<script setup>
import { ref, defineEmits } from 'vue'
import { LoadClaudeConfig, SaveClaudeConfig } from '../../wailsjs/go/app/App'

const emit = defineEmits(['openProfileModal'])

// 响应式状态
const authToken = ref('')
const baseURL = ref('')
const loading = ref(false)
const message = ref('')
const messageType = ref('')
const configLoaded = ref(false)
const showPassword = ref(false)

// 加载配置
async function loadConfig() {
  loading.value = true
  message.value = ''
  
  try {
    const result = await LoadClaudeConfig()
    
    if (result.success && result.data) {
      authToken.value = result.data.env.ANTHROPIC_AUTH_TOKEN || ''
      baseURL.value = result.data.env.ANTHROPIC_BASE_URL || ''
      configLoaded.value = true
      showMessage('配置加载成功 ✨', 'success')
    } else {
      configLoaded.value = false
      showMessage('配置文件不存在，将创建新配置 🚀', 'info')
    }
  } catch (error) {
    showMessage('加载配置失败: ' + error.toString(), 'error')
  } finally {
    loading.value = false
  }
}

// 保存配置
async function saveConfig() {
  if (!authToken.value.trim()) {
    showMessage('请输入 Auth Token 🔑', 'error')
    return
  }
  
  if (!baseURL.value.trim()) {
    showMessage('请输入 Base URL 🌐', 'error')
    return
  }

  loading.value = true
  
  try {
    const result = await SaveClaudeConfig(authToken.value.trim(), baseURL.value.trim())
    
    if (result.success) {
      configLoaded.value = true
      showMessage('配置保存成功！🎉', 'success')
    } else {
      showMessage(result.message || '保存配置失败', 'error')
    }
  } catch (error) {
    showMessage('保存配置失败: ' + error.toString(), 'error')
  } finally {
    loading.value = false
  }
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

// 切换密码可见性
function togglePasswordVisibility() {
  showPassword.value = !showPassword.value
}

// 初始化时加载配置
loadConfig()
</script>

<style scoped>
@import '../assets/common.css';

.config-section {
  display: flex;
  flex-direction: column;
  gap: 2rem;
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

.config-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.form-actions .btn {
  flex: 1;
}

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

/* 响应式设计 */
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
</style>