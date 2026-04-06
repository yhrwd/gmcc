<script setup lang="ts">
import StatusBadge from '@/components/shared/StatusBadge.vue'
import type { ViewAccount } from '@/types/view'

defineProps<{
  account: ViewAccount
  deleting: boolean
  loginStarting: boolean
}>()

const emit = defineEmits<{
  login: [accountId: string]
  delete: [accountId: string]
}>()
</script>

<template>
  <article class="account-card">
    <div class="account-card__header">
      <div>
        <h3 class="account-card__title">{{ account.label }}</h3>
        <p class="account-card__id">{{ account.id }}</p>
      </div>
      <StatusBadge :label="account.authStatusLabel" :tone="account.authStatusTone" />
    </div>
    <p class="account-card__note">{{ account.note || '当前没有额外备注。' }}</p>
    <div class="account-card__actions">
      <button
        v-if="account.authStatusTone !== 'ready' && account.authStatusTone !== 'disabled'"
        type="button"
        class="account-card__button"
        :disabled="loginStarting || deleting"
        @click="emit('login', account.id)"
      >
        {{ loginStarting ? '登录中...' : '去登录' }}
      </button>
      <button type="button" class="account-card__button account-card__button--ghost" :disabled="deleting || loginStarting" @click="emit('delete', account.id)">
        {{ deleting ? '删除中...' : '删除账号' }}
      </button>
    </div>
  </article>
</template>

<style scoped>
.account-card {
  display: grid;
  gap: 0.75rem;
  padding: 1rem;
  border-radius: 1.35rem;
  background: rgba(255, 255, 255, 0.62);
  border: 1px solid rgba(107, 131, 147, 0.14);
}

.account-card__header {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 1rem;
}

.account-card__title,
.account-card__id,
.account-card__note {
  margin: 0;
}

.account-card__id,
.account-card__note {
  color: var(--text-soft);
}

.account-card__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.55rem;
}

.account-card__button {
  border: none;
  border-radius: 999px 999px 999px 0.95rem;
  padding: 0.7rem 0.95rem;
  background: linear-gradient(135deg, rgba(255, 231, 200, 0.96), rgba(255, 214, 196, 0.92));
  color: #85512e;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 12px 22px rgba(219, 160, 127, 0.16);
}

.account-card__button--ghost {
  background: rgba(236, 241, 248, 0.95);
  color: #576980;
  box-shadow: none;
}
</style>
