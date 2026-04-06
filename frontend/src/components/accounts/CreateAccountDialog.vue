<script setup lang="ts">
import { reactive, computed } from 'vue'
import { validateCreateAccount } from '@/stores/accounts'

const props = defineProps<{
  open: boolean
  submitting: boolean
}>()

const emit = defineEmits<{
  close: []
  submit: [payload: { id: string; label: string; note: string }]
}>()

const form = reactive({
  id: '',
  label: '',
  note: '',
})

const validationMessage = computed(() => validateCreateAccount(form).message)

function handleSubmit() {
  const validation = validateCreateAccount(form)
  if (!validation.valid) {
    return
  }
  emit('submit', { ...form })
}

function resetAndClose() {
  form.id = ''
  form.label = ''
  form.note = ''
  emit('close')
}
</script>

<template>
  <div v-if="open" class="dialog-shell">
    <div class="dialog-panel">
      <div class="dialog-panel__header">
        <div>
          <p class="dialog-panel__eyebrow">新建账号</p>
          <h3 class="dialog-panel__title">创建账号</h3>
        </div>
      </div>

      <label class="dialog-panel__field">
        <span>账号 ID</span>
        <input v-model="form.id" class="dialog-panel__input" placeholder="acc-main" />
      </label>

      <label class="dialog-panel__field">
        <span>昵称</span>
        <input v-model="form.label" class="dialog-panel__input" placeholder="主账号" />
      </label>

      <label class="dialog-panel__field">
        <span>备注</span>
        <textarea v-model="form.note" class="dialog-panel__textarea" placeholder="可填写用途说明" />
      </label>

      <p v-if="validationMessage" class="dialog-panel__hint">{{ validationMessage }}</p>

      <div class="dialog-panel__actions">
        <button type="button" class="dialog-panel__ghost" :disabled="submitting" @click="resetAndClose">取消</button>
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
  color: var(--text-main);
}

.dialog-panel__input,
.dialog-panel__textarea {
  width: 100%;
  border: 1px solid rgba(108, 127, 149, 0.18);
  border-radius: 1rem;
  padding: 0.85rem 0.95rem;
  background: rgba(250, 252, 255, 0.95);
}

.dialog-panel__textarea {
  min-height: 6rem;
  resize: vertical;
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
