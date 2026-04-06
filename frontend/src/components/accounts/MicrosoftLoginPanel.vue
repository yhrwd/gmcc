<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import StatusBadge from '@/components/shared/StatusBadge.vue'
import { useAccountsStore } from '@/stores/accounts'
import { useUiStore } from '@/stores/ui'

const accountsStore = useAccountsStore()
const uiStore = useUiStore()
const { loginTask } = storeToRefs(uiStore)

let timerId: number | null = null
let countdownId: number | null = null
const now = ref(Date.now())

const tone = computed(() => {
  if (loginTask.value.taskStatus === 'succeeded') return 'success'
  if (loginTask.value.taskStatus === 'failed' || loginTask.value.taskStatus === 'expired' || loginTask.value.taskStatus === 'cancelled' || loginTask.value.taskStatus === 'replaced') return 'error'
  if (loginTask.value.taskStatus === 'polling' || loginTask.value.taskStatus === 'initializing') return 'pending'
  return 'unknown'
})

function clearTimer() {
  if (timerId !== null) {
    window.clearTimeout(timerId)
    timerId = null
  }
}

function clearCountdown() {
  if (countdownId !== null) {
    window.clearInterval(countdownId)
    countdownId = null
  }
}

function ensureCountdown() {
  clearCountdown()
  if (!loginTask.value.panelOpen || !loginTask.value.expiresAt) {
    return
  }
  now.value = Date.now()
  countdownId = window.setInterval(() => {
    now.value = Date.now()
  }, 1000)
}

const verificationLink = computed(() => loginTask.value.verificationUriComplete || loginTask.value.verificationUri)

const expiresLabel = computed(() => {
  if (!loginTask.value.expiresAt) {
    return ''
  }
  const remainingMs = Math.max(loginTask.value.expiresAt - now.value, 0)
  const totalSeconds = Math.ceil(remainingMs / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
})

async function tickPoll() {
  clearTimer()
  await accountsStore.pollLogin()
  if (loginTask.value.taskStatus === 'polling' && loginTask.value.panelOpen) {
    timerId = window.setTimeout(() => {
      void tickPoll()
    }, loginTask.value.intervalSeconds * 1000)
  }
}

watch(
  () => [loginTask.value.taskStatus, loginTask.value.panelOpen] as const,
  ([taskStatus, panelOpen]) => {
    clearTimer()
    if (taskStatus === 'polling' && panelOpen) {
      timerId = window.setTimeout(() => {
        void tickPoll()
      }, loginTask.value.intervalSeconds * 1000)
    }
  },
  { immediate: true },
)

watch(
  () => [loginTask.value.panelOpen, loginTask.value.expiresAt] as const,
  () => {
    ensureCountdown()
  },
  { immediate: true },
)

onMounted(() => {
  ensureCountdown()
})

onBeforeUnmount(() => {
  clearTimer()
  clearCountdown()
})
</script>

<template>
  <div v-if="loginTask.panelOpen" class="login-panel">
    <div class="login-panel__header">
      <div>
        <p class="login-panel__eyebrow">登录流程</p>
        <h3 class="login-panel__title">账号登录 {{ loginTask.accountId }}</h3>
      </div>
      <button type="button" class="login-panel__ghost" @click="uiStore.closeLoginPanel">收起</button>
    </div>

    <StatusBadge :label="loginTask.taskStatus" :tone="tone" />

    <div v-if="loginTask.userCode" class="login-panel__code">{{ loginTask.userCode }}</div>
    <p v-if="expiresLabel" class="login-panel__meta">设备码剩余时间 {{ expiresLabel }}</p>
    <p v-if="loginTask.verificationUri" class="login-panel__meta login-panel__meta--break">授权地址 {{ loginTask.verificationUri }}</p>
    <p class="login-panel__message">{{ loginTask.lastMessage || '准备开始 Microsoft 设备码登录' }}</p>

    <div class="login-panel__actions">
      <a v-if="verificationLink" class="login-panel__primary" :href="verificationLink" target="_blank" rel="noreferrer">打开验证页</a>
      <button
        v-if="loginTask.taskStatus === 'failed' || loginTask.taskStatus === 'expired' || loginTask.taskStatus === 'cancelled'"
        type="button"
        class="login-panel__ghost"
        @click="accountsStore.startLogin(loginTask.accountId)"
      >
        重新尝试
      </button>
      <button
        v-if="loginTask.taskStatus === 'polling' || loginTask.taskStatus === 'initializing'"
        type="button"
        class="login-panel__ghost"
        disabled
      >
        处理中...
      </button>
    </div>
  </div>
</template>

<style scoped>
.login-panel {
  display: grid;
  gap: 0.9rem;
  padding: 1rem;
  border-radius: 1.35rem;
  background: linear-gradient(135deg, rgba(242, 255, 249, 0.9), rgba(244, 243, 255, 0.9));
  border: 1px solid rgba(112, 141, 160, 0.14);
}

.login-panel__header,
.login-panel__actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.8rem;
}

.login-panel__eyebrow,
.login-panel__title,
.login-panel__meta,
.login-panel__message {
  margin: 0;
}

.login-panel__eyebrow {
  font-size: 0.72rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--text-soft);
}

.login-panel__code {
  padding: 0.95rem 1rem;
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.78);
  font-size: 1.35rem;
  font-weight: 800;
  letter-spacing: 0.16em;
  color: var(--text-main);
}

.login-panel__message {
  color: var(--text-soft);
  line-height: 1.6;
}

.login-panel__meta {
  color: var(--text-main);
  font-size: 0.92rem;
}

.login-panel__meta--break {
  word-break: break-word;
}

.login-panel__ghost,
.login-panel__primary {
  border: none;
  border-radius: 999px;
  padding: 0.72rem 1rem;
  text-decoration: none;
  cursor: pointer;
}

.login-panel__ghost:focus-visible,
.login-panel__primary:focus-visible {
  outline: 3px solid rgba(72, 122, 170, 0.32);
  outline-offset: 2px;
}

.login-panel__ghost {
  background: rgba(227, 234, 245, 0.92);
  color: var(--text-main);
}

.login-panel__primary {
  background: linear-gradient(135deg, #7be0cb, #9cbcff);
  color: #143640;
  font-weight: 700;
}
</style>
