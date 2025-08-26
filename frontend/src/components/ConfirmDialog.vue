<template>
  <div v-if="show" class="modal-overlay" @click="handleCancel">
    <div class="modal confirmation-modal" @click.stop>
      <div class="modal-header">
        <h3>{{ title }}</h3>
      </div>
      <div class="modal-body">
        <div class="confirmation-content">
          <div class="confirmation-icon">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
              <path d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" fill="none"/>
            </svg>
          </div>
          <p class="confirmation-message">{{ message }}</p>
        </div>
      </div>
      <div class="modal-footer">
        <button @click="handleCancel" class="btn btn--secondary">取消</button>
        <button @click="handleConfirm" class="btn btn--danger">确认删除</button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '确认操作'
  },
  message: {
    type: String,
    default: '您确定要执行此操作吗？'
  }
})

const emit = defineEmits(['confirm', 'cancel'])

function handleConfirm() {
  emit('confirm')
}

function handleCancel() {
  emit('cancel')
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

.modal-body {
  padding: 2rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 0 2rem 2rem;
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

.btn--danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.btn--danger:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(239, 68, 68, 0.4);
}
</style>