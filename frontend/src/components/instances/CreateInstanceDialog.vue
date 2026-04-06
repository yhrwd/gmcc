<script setup lang="ts">
import { reactive, computed } from 'vue'
import { validateCreateInstance } from '@/stores/instances'
import type { ViewAccount } from '@/types/view'

const props = defineProps<{
  open: boolean
  submitting: boolean
  accounts: ViewAccount[]
}>()

const emit = defineEmits<{
  close: []
  submit: [payload: { id: string; accountId: string; serverAddress: string; enabled: boolean; autoStart: boolean }]
}>()

const form = reactive({
  id: '',
  accountId: '',
  serverAddress: '',
  enabled: true,
  autoStart: false,
})

const validationMessage = computed(() => validateCreateInstance(form).message)
const availableAccounts = computed(() => props.accounts.filter((account) => account.authStatusTone === 'ready'))

function handleSubmit() {
  const validation = validateCreateInstance(form)
  if (!validation.valid) {
    return
  }
  emit('submit', { ...form })
}

function handleClose() {
  emit('close')
}
</script>

<template>
  <div v-if="open" class="dialog-shell">
    <div class="dialog-panel">
      <div class="dialog-panel__header">
        <div>
          <p class="dialog-panel__eyebrow">新建实例</p>
          <h3 class="dialog-panel__title">创建实例</h3>
        </div>
      </div>

      <label class="dialog-panel__field">
        <span>实例 ID</span>
        <input v-model="form.id" class="dialog-panel__input" placeholder="bot-1" />
      </label>

      <label class="dialog-panel__field">
        <span>绑定账号</span>
        <select v-model="form.accountId" class="dialog-panel__input">
          <option value="">请选择已登录账号</option>
          <option v-for="account in availableAccounts" :key="account.id" :value="account.id">{{ account.label }}</option>
        </select>
      </label>

      <label class="dialog-panel__field">
        <span>服务器地址</span>
        <input v-model="form.serverAddress" class="dialog-panel__input" placeholder="mc.example.com" />
      </label>

      <label class="dialog-panel__toggle">
        <input v-model="form.enabled" type="checkbox" />
        <span>创建后启用</span>
      </label>

      <label class="dialog-panel__toggle">
        <input v-model="form.autoStart" type="checkbox" />
        <span>创建后自动启动</span>
      </label>

      <p v-if="validationMessage" class="dialog-panel__hint">{{ validationMessage }}</p>

      <div class="dialog-panel__actions">
        <button type="button" class="dialog-panel__ghost" :disabled="submitting" @click="handleClose">取消</button>
        <button type="button" class="dialog-panel__primary" :disabled="submitting" @click="handleSubmit">
          {{ submitting ? '创建中...' : '确认创建' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dialog-shell {
  position: fixed;
  inset: 0;
  z-index: 30;
  display: grid;
  place-items: center;
  padding: 1rem;
  background: rgba(23, 28, 42, 0.24);
  backdrop-filter: blur(10px);
}

.dialog-panel {
  width: min(32rem, 100%);
  display: grid;
  gap: 0.95rem;
  padding: 1.2rem;
  border-radius: 1.65rem;
  background: rgba(255, 255, 255, 0.92);
}

.dialog-panel__header,
.dialog-panel__actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.dialog-panel__eyebrow,
.dialog-panel__title,
.dialog-panel__hint {
  margin: 0;
}

.dialog-panel__eyebrow {
  color: var(--text-soft);
  text-transform: uppercase;
  letter-spacing: 0.18em;
  font-size: 0.72rem;
}

.dialog-panel__field {
  display: grid;
  gap: 0.45rem;
}

.dialog-panel__input {
  width: 100%;
  border: 1px solid rgba(108, 127, 149, 0.18);
  border-radius: 1rem;
  padding: 0.85rem 0.95rem;
  background: rgba(250, 252, 255, 0.95);
}

.dialog-panel__toggle {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  color: var(--text-main);
}

.dialog-panel__hint {
  color: #b14b65;
}

.dialog-panel__ghost,
.dialog-panel__primary {
  border: none;
  border-radius: 999px;
  padding: 0.75rem 1rem;
  cursor: pointer;
}

.dialog-panel__ghost {
  background: rgba(240, 235, 227, 0.96);
  color: var(--text-main);
}

.dialog-panel__primary {
  background: linear-gradient(135deg, #f28e16, #edc3ae);
  color: #4f341b;
  font-weight: 700;
}
</style>
