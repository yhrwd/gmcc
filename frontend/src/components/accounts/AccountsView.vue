<script setup lang="ts">
import { shallowRef, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import AccountCard from '@/components/accounts/AccountCard.vue'
import CreateAccountDialog from '@/components/accounts/CreateAccountDialog.vue'
import MicrosoftLoginPanel from '@/components/accounts/MicrosoftLoginPanel.vue'
import BaseCard from '@/components/shared/BaseCard.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import InlineError from '@/components/shared/InlineError.vue'
import { useAccountsStore } from '@/stores/accounts'

const accountsStore = useAccountsStore()
const { state, errorMessage, items, submitting, deletingId, loginStartingId } = storeToRefs(accountsStore)
const createOpen = shallowRef(false)

async function handleCreate(payload: { id: string; label: string; note: string }) {
  const result = await accountsStore.createAccount(payload)
  if (result.valid) {
    createOpen.value = false
  }
}

onMounted(() => {
  if (state.value === 'idle') {
    void accountsStore.loadAccounts()
  }
})
</script>

<template>
  <section class="accounts-view">
    <BaseCard eyebrow="账号管理" title="账号管理" accent="peach">
      <div class="accounts-view__toolbar">
        <div>
          <h2 class="accounts-view__headline">管理账号</h2>
          <p class="accounts-view__copy">创建账号、查看认证状态并执行登录操作。</p>
        </div>
        <button type="button" class="accounts-view__button" @click="createOpen = true">新增账号</button>
      </div>

      <InlineError v-if="state === 'error'" :message="errorMessage" button-label="重试" @retry="accountsStore.loadAccounts" />
      <EmptyState
        v-else-if="state === 'success' && !items.length"
        title="暂无账号"
        description="请先创建账号，再继续进行登录配置。"
        button-label="创建账号"
        @action="createOpen = true"
      />
      <div v-else class="accounts-view__grid">
        <AccountCard
          v-for="account in items"
          :key="account.id"
          :account="account"
          :deleting="deletingId === account.id"
          :login-starting="loginStartingId === account.id"
          @login="accountsStore.startLogin"
          @delete="accountsStore.deleteAccount"
        />
      </div>
    </BaseCard>

    <MicrosoftLoginPanel />
    <CreateAccountDialog :open="createOpen" :submitting="submitting" @close="createOpen = false" @submit="handleCreate" />
  </section>
</template>

<style scoped>
.accounts-view {
  display: grid;
  gap: 1rem;
}

.accounts-view__toolbar {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1rem;
}

.accounts-view__headline,
.accounts-view__copy {
  margin: 0;
}

.accounts-view__copy {
  margin-top: 0.45rem;
  color: var(--text-soft);
}

.accounts-view__button {
  align-self: start;
  border: none;
  border-radius: 999px;
  padding: 0.8rem 1rem;
  background: linear-gradient(135deg, #ffdcaf, #ffd0c2);
  color: #71472f;
  font-weight: 700;
  cursor: pointer;
}

.accounts-view__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.85rem;
}

@media (max-width: 900px) {
  .accounts-view__toolbar {
    flex-direction: column;
  }

  .accounts-view__grid {
    grid-template-columns: 1fr;
  }
}
</style>
