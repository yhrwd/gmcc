<script setup lang="ts">
defineProps<{
  label: string
  tone: 'ready' | 'pending' | 'disabled' | 'invalid' | 'unknown' | 'starting' | 'running' | 'reconnecting' | 'stopped' | 'success' | 'error'
}>()
</script>

<template>
  <span class="status-badge" :data-tone="tone">
    <span class="status-badge__dot"></span>
    <span>{{ label }}</span>
  </span>
</template>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.44rem 0.82rem;
  border-radius: 999px 999px 999px 0.9rem;
  font-size: 0.78rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  border: 1px solid transparent;
  box-shadow: 0 8px 16px rgba(85, 98, 120, 0.08);
}

.status-badge__dot {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 999px;
  background: currentColor;
  box-shadow: 0 0 0 0.18rem rgba(255, 255, 255, 0.44);
}

.status-badge[data-tone='ready'],
.status-badge[data-tone='running'],
.status-badge[data-tone='success'] {
  color: #456255;
  background: linear-gradient(135deg, rgba(188, 206, 192, 0.86), rgba(238, 247, 242, 0.92));
  border-color: rgba(131, 167, 141, 0.28);
}

.status-badge[data-tone='pending'] {
  color: #9a6024;
  background: linear-gradient(135deg, rgba(242, 210, 167, 0.9), rgba(248, 239, 223, 0.95));
  border-color: rgba(242, 142, 22, 0.22);
}

.status-badge[data-tone='starting'],
.status-badge[data-tone='reconnecting'] {
  color: #546c76;
  background: linear-gradient(135deg, rgba(208, 223, 230, 0.9), rgba(238, 244, 246, 0.96));
  border-color: rgba(158, 181, 191, 0.24);
}

.status-badge[data-tone='disabled'],
.status-badge[data-tone='stopped'] {
  color: #746f67;
  background: linear-gradient(135deg, rgba(234, 231, 224, 0.88), rgba(245, 242, 236, 0.94));
  border-color: rgba(170, 161, 148, 0.22);
}

.status-badge[data-tone='invalid'],
.status-badge[data-tone='error'] {
  color: #8c4b44;
  background: linear-gradient(135deg, rgba(237, 195, 174, 0.88), rgba(250, 234, 225, 0.94));
  border-color: rgba(195, 132, 103, 0.24);
}

.status-badge[data-tone='unknown'] {
  color: #5e6b66;
  background: linear-gradient(135deg, rgba(222, 231, 228, 0.9), rgba(243, 246, 244, 0.96));
  border-color: rgba(140, 157, 151, 0.24);
}
</style>
