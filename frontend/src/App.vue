<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import AccountsView from '@/components/accounts/AccountsView.vue'
import HomeView from '@/components/home/HomeView.vue'
import InstancesView from '@/components/instances/InstancesView.vue'
import AppShell from '@/components/layout/AppShell.vue'
import LogsView from '@/components/logs/LogsView.vue'
import { useAccountsStore } from '@/stores/accounts'
import { useHomeStore } from '@/stores/home'
import { useInstancesStore } from '@/stores/instances'
import { useLogsStore } from '@/stores/logs'
import { useUiStore } from '@/stores/ui'

const uiStore = useUiStore()
const homeStore = useHomeStore()
const accountsStore = useAccountsStore()
const instancesStore = useInstancesStore()
const logsStore = useLogsStore()

const { activeView } = storeToRefs(uiStore)
const { statusSummary, resourceSnapshot } = storeToRefs(homeStore)

const currentView = computed(() => {
  if (activeView.value === 'accounts') return AccountsView
  if (activeView.value === 'instances') return InstancesView
  if (activeView.value === 'logs') return LogsView
  return HomeView
})

let fastRefreshId: number | null = null
let slowRefreshId: number | null = null

onMounted(() => {
  void Promise.all([
    homeStore.loadHome(),
    accountsStore.loadAccounts(),
    instancesStore.loadInstances(),
    logsStore.loadLogs(),
  ])

  fastRefreshId = window.setInterval(() => {
    void homeStore.loadHome(true, true)
    void accountsStore.loadAccounts(true, true)
    void instancesStore.loadInstances(true, true)
  }, 8000)

  slowRefreshId = window.setInterval(() => {
    void logsStore.refreshLogs()
  }, 20000)
})

onBeforeUnmount(() => {
  if (fastRefreshId !== null) window.clearInterval(fastRefreshId)
  if (slowRefreshId !== null) window.clearInterval(slowRefreshId)
})
</script>

<template>
  <AppShell :active-view="activeView" :status="statusSummary" @select="uiStore.setActiveView">
    <component :is="currentView" />
  </AppShell>
</template>
