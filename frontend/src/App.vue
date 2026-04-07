<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import AccountsView from '@/components/accounts/AccountsView.vue'
import HomeView from '@/components/home/HomeView.vue'
import InstancesView from '@/components/instances/InstancesView.vue'
import AppShell from '@/components/layout/AppShell.vue'
import LogsView from '@/components/logs/LogsView.vue'
import { syncCoordinator } from '@/lib/sync'
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
const { statusSummary } = storeToRefs(homeStore)

const currentView = computed(() => {
  if (activeView.value === 'accounts') return AccountsView
  if (activeView.value === 'instances') return InstancesView
  if (activeView.value === 'logs') return LogsView
  return HomeView
})

onMounted(() => {
  syncCoordinator.register('overview', () => homeStore.loadHome(true, true))
  syncCoordinator.register('accounts', () => accountsStore.loadAccounts(true, true))
  syncCoordinator.register('instances', () => instancesStore.loadInstances(true, true))
  syncCoordinator.register('logs', () => logsStore.refreshLogs())
  syncCoordinator.start(['overview', 'accounts', 'instances', 'logs'])
})

onBeforeUnmount(() => {
  syncCoordinator.stop()
})
</script>

<template>
  <AppShell :active-view="activeView" :status="statusSummary" @select="uiStore.setActiveView">
    <component :is="currentView" />
  </AppShell>
</template>
