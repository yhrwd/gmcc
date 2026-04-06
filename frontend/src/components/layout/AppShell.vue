<script setup lang="ts">
import MobileNav from '@/components/layout/MobileNav.vue'
import SidebarNav from '@/components/layout/SidebarNav.vue'
import TopStatusBar from '@/components/layout/TopStatusBar.vue'
import ToastHost from '@/components/shared/ToastHost.vue'
import type { ActiveView, ViewStatusSummary } from '@/types/view'

defineProps<{
  activeView: ActiveView
  status: ViewStatusSummary | null
}>()

const emit = defineEmits<{
  select: [view: ActiveView]
}>()
</script>

<template>
  <div class="app-shell">
    <aside class="app-shell__sidebar">
      <SidebarNav :active-view="activeView" @select="emit('select', $event)" />
    </aside>

    <div class="app-shell__content">
      <TopStatusBar v-if="activeView === 'home'" :status="status" />
      <main class="app-shell__main">
        <slot />
      </main>
      <div class="app-shell__mobile">
        <MobileNav :active-view="activeView" @select="emit('select', $event)" />
      </div>
    </div>

    <ToastHost />
  </div>
</template>

<style scoped>
.app-shell {
  display: grid;
  grid-template-columns: 17.2rem minmax(0, 1fr);
  gap: 1rem;
  min-height: 100vh;
  padding: 1rem;
}

.app-shell__sidebar {
  position: sticky;
  top: 1rem;
  align-self: start;
  min-height: calc(100vh - 2rem);
  padding: 1rem;
  border-radius: 1.8rem;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.74), rgba(255, 255, 255, 0.62));
  border: 1px solid rgba(104, 128, 145, 0.14);
  backdrop-filter: blur(18px);
}

.app-shell__content {
  display: grid;
  gap: 0.85rem;
  min-width: 0;
}

.app-shell__main {
  min-width: 0;
  display: grid;
  gap: 0.85rem;
}

.app-shell__mobile {
  display: none;
}

@media (max-width: 960px) {
  .app-shell {
    grid-template-columns: 1fr;
    padding: 0.8rem 0.8rem 5.8rem;
  }

  .app-shell__sidebar {
    display: none;
  }

  .app-shell__mobile {
    display: block;
  }
}
</style>
