<script setup lang="ts">
defineProps<{
  title: string
  description: string
  buttonLabel?: string
}>()

const emit = defineEmits<{
  action: []
}>()
</script>

<template>
  <div class="empty-state">
    <div class="empty-state__pet-wrap">
      <div class="empty-state__pet">o(>w<)o</div>
      <span class="empty-state__sticker">基地便条</span>
    </div>
    <h4 class="empty-state__title">{{ title }}</h4>
    <p class="empty-state__description">{{ description }}</p>
    <button v-if="buttonLabel" class="empty-state__button" type="button" @click="emit('action')">
      {{ buttonLabel }}
    </button>
  </div>
</template>

<style scoped>
.empty-state {
  display: grid;
  gap: 0.65rem;
  place-items: start;
  padding: 1.2rem;
  border: 1px dashed rgba(135, 154, 176, 0.28);
  border-radius: 1.5rem;
  background: rgba(255, 255, 255, 0.38);
}

.empty-state__pet-wrap {
  display: flex;
  align-items: center;
  gap: 0.55rem;
}

.empty-state__pet {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 3.1rem;
  min-height: 3.1rem;
  padding: 0.55rem 0.8rem;
  border-radius: 1.25rem;
  background: linear-gradient(135deg, rgba(255, 240, 214, 0.9), rgba(220, 245, 239, 0.88));
  font-weight: 700;
  color: var(--text-main);
  box-shadow: 0 10px 18px rgba(115, 128, 151, 0.08);
}

.empty-state__sticker {
  display: inline-flex;
  align-items: center;
  padding: 0.35rem 0.6rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.84);
  color: var(--text-soft);
  font-size: 0.72rem;
  font-weight: 700;
}

.empty-state__title {
  margin: 0;
  font-size: 1rem;
  color: var(--text-main);
}

.empty-state__description {
  margin: 0;
  color: var(--text-soft);
  line-height: 1.6;
}

.empty-state__button {
  border: none;
  border-radius: 999px 999px 999px 0.95rem;
  padding: 0.72rem 1.02rem;
  background: linear-gradient(135deg, #7be0cb, #9bbdff);
  color: #12323c;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 12px 22px rgba(122, 170, 213, 0.18);
}
</style>
