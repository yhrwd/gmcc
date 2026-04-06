<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useUiStore } from '@/stores/ui'

const uiStore = useUiStore()
const { toasts } = storeToRefs(uiStore)
</script>

<template>
  <div class="toast-host" aria-live="polite">
    <transition-group name="toast-list">
      <div v-for="toast in toasts" :key="toast.id" class="toast-item" :data-tone="toast.tone">
        <span>{{ toast.message }}</span>
        <button type="button" class="toast-item__button" @click="uiStore.dismissToast(toast.id)">x</button>
      </div>
    </transition-group>
  </div>
</template>

<style scoped>
.toast-host {
  position: fixed;
  right: 1.25rem;
  bottom: 1.25rem;
  z-index: 40;
  display: grid;
  gap: 0.75rem;
  max-width: min(24rem, calc(100vw - 2rem));
}

.toast-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.9rem;
  padding: 0.95rem 1rem;
  border-radius: 1.15rem;
  color: var(--text-main);
  backdrop-filter: blur(18px);
  box-shadow: 0 18px 36px rgba(73, 83, 104, 0.16);
}

.toast-item[data-tone='success'] {
  background: linear-gradient(135deg, rgba(238, 247, 242, 0.94), rgba(224, 236, 228, 0.92));
}

.toast-item[data-tone='error'] {
  background: linear-gradient(135deg, rgba(250, 234, 225, 0.94), rgba(244, 223, 208, 0.92));
}

.toast-item__button {
  border: none;
  background: transparent;
  color: inherit;
  cursor: pointer;
}

.toast-item__button:focus-visible {
  outline: 3px solid rgba(72, 122, 170, 0.3);
  outline-offset: 2px;
}

.toast-list-enter-active,
.toast-list-leave-active {
  transition: opacity 0.24s ease, transform 0.24s ease;
}

.toast-list-enter-from,
.toast-list-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
