<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import AccountSummaryPanel from '@/components/home/AccountSummaryPanel.vue'
import HeroSummaryCard from '@/components/home/HeroSummaryCard.vue'
import InstanceSummaryPanel from '@/components/home/InstanceSummaryPanel.vue'
import LogSummaryPanel from '@/components/home/LogSummaryPanel.vue'
import { useAccountsStore } from '@/stores/accounts'
import { useHomeStore } from '@/stores/home'
import { useInstancesStore } from '@/stores/instances'
import { useLogsStore } from '@/stores/logs'
import { buildHeroView } from '@/lib/mappers'

const homeStore = useHomeStore()
const accountsStore = useAccountsStore()
const instancesStore = useInstancesStore()
const logsStore = useLogsStore()

const {
  statusState,
  resourcesState,
  statusError,
  resourcesError,
  statusSummary,
  resourceSnapshot,
} = storeToRefs(homeStore)
const { state: accountsState, errorMessage: accountsError, items: accountItems } = storeToRefs(accountsStore)
const { state: instancesState, errorMessage: instancesError, items: instanceItems } = storeToRefs(instancesStore)
const { state: logsState, errorMessage: logsError, items: logItems } = storeToRefs(logsStore)

const hero = computed(() => buildHeroView({
  status: statusSummary.value,
  resources: resourceSnapshot.value,
  accounts: accountItems.value,
}))

const summaryLogs = computed(() => logItems.value.slice(0, 10))
</script>

<template>
  <section class="home-view">
    <div class="home-view__hero">
      <HeroSummaryCard :hero="hero" />
    </div>

    <section class="home-view__guide" aria-label="总览使用提示">
      <div class="home-view__guide-card">
        <span class="home-view__guide-label">使用提示</span>
        <p class="home-view__guide-text">
          先在“账号”中完成账号创建和登录，再到“实例”里创建并启动实例；当前页主要用于快速查看系统状态、账号状态、实例状态和最近操作记录。
        </p>
      </div>
    </section>

    <div v-if="statusState === 'error' || resourcesState === 'error'" class="home-view__system-error">
      <span v-if="statusState === 'error'">{{ statusError }}</span>
      <span v-if="resourcesState === 'error'">{{ resourcesError }}</span>
    </div>

    <div class="home-view__summary">
      <InstanceSummaryPanel
        :state="instancesState"
        :error-message="instancesError"
        :instances="instanceItems"
        @retry="instancesStore.loadInstances(true)"
      />
      <AccountSummaryPanel
        :state="accountsState"
        :error-message="accountsError"
        :accounts="accountItems"
        @retry="accountsStore.loadAccounts(true)"
      />
    </div>

    <LogSummaryPanel :state="logsState" :error-message="logsError" :logs="summaryLogs" @retry="logsStore.loadLogs(true)" />
  </section>
</template>

<style scoped>
.home-view {
  display: grid;
  gap: 0.9rem;
}

.home-view__guide {
  display: grid;
}

.home-view__guide-card {
  display: grid;
  gap: 0.45rem;
  padding: 0.9rem 1rem;
  border-radius: 1rem;
  background: rgba(255, 250, 242, 0.7);
  border: 1px solid rgba(242, 142, 22, 0.14);
}

.home-view__guide-label {
  font-size: 0.76rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #9a6024;
}

.home-view__guide-text {
  margin: 0;
  color: var(--text-soft);
  line-height: 1.6;
}

.home-view__summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(min(100%, 22rem), 1fr));
  gap: 0.9rem;
  align-items: start;
}

.home-view__system-error {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  color: #a14b5e;
}

@media (min-width: 1100px) {
  .home-view__hero {
    min-height: 15rem;
  }
}

@media (max-width: 900px) {
  .home-view__summary {
    grid-template-columns: 1fr;
  }
}
</style>
