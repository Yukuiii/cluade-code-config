<template>
  <div v-if="show" class="modal-overlay" @click="handleClose">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h3>保存配置文件</h3>
        <button @click="handleClose" class="modal-close">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
            <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2"/>
          </svg>
        </button>
      </div>
      <div class="modal-body">
        <form @submit.prevent="handleSave">
          <div class="form-group">
            <label for="profileName" class="form-label">配置文件名称</label>
            <input 
              id="profileName"
              v-model="profileData.name"
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
                v-model="profileData.authToken"
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
                v-model="profileData.baseURL"
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
              v-model="profileData.description"
              class="form-textarea"
              placeholder="描述这个配置文件的用途..."
              rows="3"
              maxlength="200"
            ></textarea>
          </div>
        </form>
      </div>
      <div class="modal-footer">
        <button @click="handleClose" class="btn btn--secondary">取消</button>
        <button @click="handleSave" :disabled="loading" class="btn btn--primary">
          <svg v-if="!loading" width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z" stroke="currentColor" stroke-width="2" fill="none"/>
            <polyline points="17,21 17,13 7,13 7,21" stroke="currentColor" stroke-width="2"/>
            <polyline points="7,3 7,8 15,8" stroke="currentColor" stroke-width="2"/>
          </svg>
          <svg v-else class="spinner" width="16" height="16" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" opacity="0.25"/>
            <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="4"/>
          </svg>
          {{ loading ? '保存中...' : '保存配置文件' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, defineEmits, defineProps } from 'vue'
import { SaveConfigProfile } from '../../wailsjs/go/app/App'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'saved', 'error'])

// 响应式状态
const loading = ref(false)
const profileData = ref({
  name: '',
  authToken: '',
  baseURL: '',
  description: ''
})

// 监听 show 属性变化，重置表单
watch(() => props.show, (newShow) => {
  if (newShow) {
    resetForm()
  }
})

// 重置表单
function resetForm() {
  profileData.value = {
    name: '',
    authToken: '',
    baseURL: '',
    description: ''
  }
}

// 处理保存
async function handleSave() {
  if (!profileData.value.name.trim()) {
    emit('error', '请输入配置文件名称 📝')
    return
  }
  
  if (!profileData.value.authToken.trim()) {
    emit('error', '请输入 Auth Token 🔑')
    return
  }
  
  if (!profileData.value.baseURL.trim()) {
    emit('error', '请输入 Base URL 🌐')
    return
  }

  loading.value = true
  
  try {
    const result = await SaveConfigProfile(
      profileData.value.name.trim(),
      profileData.value.authToken.trim(),
      profileData.value.baseURL.trim(),
      profileData.value.description.trim()
    )
    
    if (result.success) {
      emit('saved', '配置文件保存成功！🎉')
      emit('close')
    } else {
      emit('error', result.message || '保存配置文件失败')
    }
  } catch (error) {
    emit('error', '保存配置文件失败: ' + error.toString())
  } finally {
    loading.value = false
  }
}

// 处理关闭
function handleClose() {
  emit('close')
}
</script>

<style scoped>
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

/* Form Styles */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
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

.form-input:focus {
  outline: none;
  border-color: #1e40af;
  box-shadow: 0 0 0 4px rgba(30, 64, 175, 0.1);
  background-color: #fefefe;
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

/* Button Styles */
.btn {
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
</style>