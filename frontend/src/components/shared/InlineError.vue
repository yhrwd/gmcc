<script setup lang="ts">
defineProps<{
  message: string
  buttonLabel?: string
}>()

const emit = defineEmits<{
  retry: []
}>()
</script>

<template>
  <div class="inline-error">
    <p class="inline-error__title">这一块信号有点飘</p>
    <p class="inline-error__message">{{ message }}</p>
    <button v-if="buttonLabel" type="button" class="inline-error__button" @click="emit('retry')">
      {{ buttonLabel }}
    </button>
  </div>
</template>

<style scoped>
.inline-error {
  display: grid;
  gap: 0.45rem;
  padding: 1rem;
  border-radius: 1.2rem;
  border: 1px solid rgba(206, 101, 126, 0.12);
  background: linear-gradient(135deg, rgba(255, 244, 247, 0.92), rgba(255, 250, 252, 0.96));
}

.inline-error__title {
  margin: 0;
  font-weight: 700;
  color: #8f2942;
  font-size: 0.95rem;
}

.inline-error__message {
  margin: 0;
  color: #8e6472;
  line-height: 1.55;
}

.inline-error__button {
  width: fit-content;
  border: none;
  border-radius: 999px;
  padding: 0.6rem 0.9rem;
  background: rgba(255, 205, 219, 0.96);
  color: #8f2942;
  font-weight: 700;
  cursor: pointer;
}
</style>
