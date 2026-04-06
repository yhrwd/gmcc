<script setup lang="ts">
defineProps<{
  title?: string
  eyebrow?: string
  accent?: 'mint' | 'peach' | 'sky' | 'gold'
}>()
</script>

<template>
  <section class="base-card" :data-accent="accent || 'mint'">
    <div class="base-card__corner"></div>
    <header v-if="title || eyebrow" class="base-card__header">
      <div>
        <p v-if="eyebrow" class="base-card__eyebrow">{{ eyebrow }}</p>
        <h3 v-if="title" class="base-card__title">{{ title }}</h3>
      </div>
    </header>
    <div class="base-card__body">
      <slot />
    </div>
  </section>
</template>

<style scoped>
.base-card {
  position: relative;
  overflow: hidden;
  border: 1px solid var(--line-soft);
  border-radius: 28px;
  padding: 1.15rem;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.92), rgba(250, 246, 239, 0.76));
  box-shadow: 0 20px 40px rgba(115, 110, 93, 0.09);
  backdrop-filter: blur(18px);
  transition: transform 260ms cubic-bezier(0.22, 1, 0.36, 1), box-shadow 260ms cubic-bezier(0.22, 1, 0.36, 1);
}

.base-card:hover {
  transform: translateY(-2px);
}

.base-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.34), transparent 24%),
    linear-gradient(0deg, rgba(255, 255, 255, 0.12), transparent 38%);
  pointer-events: none;
}

.base-card::after {
  content: '';
  position: absolute;
  inset: auto -12% -55% auto;
  width: 9rem;
  height: 9rem;
  border-radius: 2rem;
  opacity: 0.45;
  transform: rotate(18deg);
  background: rgba(255, 255, 255, 0.56);
}

.base-card__corner {
  position: absolute;
  top: -0.55rem;
  right: 1.1rem;
  width: 3.5rem;
  height: 1.2rem;
  border-radius: 999px;
  background: rgba(242, 142, 22, 0.3);
  transform: rotate(6deg);
  box-shadow: 0 6px 14px rgba(190, 128, 45, 0.18);
}

.base-card[data-accent='mint'] {
  box-shadow: 0 18px 45px rgba(131, 167, 141, 0.12), inset 0 0 0 1px rgba(131, 167, 141, 0.1);
}

.base-card[data-accent='peach'] {
  box-shadow: 0 18px 45px rgba(237, 195, 174, 0.16), inset 0 0 0 1px rgba(237, 195, 174, 0.14);
}

.base-card[data-accent='sky'] {
  box-shadow: 0 18px 45px rgba(208, 223, 230, 0.18), inset 0 0 0 1px rgba(208, 223, 230, 0.16);
}

.base-card[data-accent='gold'] {
  box-shadow: 0 18px 45px rgba(242, 142, 22, 0.12), inset 0 0 0 1px rgba(242, 142, 22, 0.12);
}

.base-card__header {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 1rem;
  position: relative;
  z-index: 1;
  margin-bottom: 0.9rem;
}

.base-card__eyebrow {
  margin: 0 0 0.25rem;
  font-size: 0.73rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--text-soft);
}

.base-card__title {
  margin: 0;
  font-size: 1.05rem;
  color: var(--text-main);
}

.base-card__body {
  position: relative;
  z-index: 1;
}

@media (prefers-reduced-motion: reduce) {
  .base-card {
    transition: none;
  }

  .base-card:hover {
    transform: none;
  }
}

@media (max-width: 640px) {
  .base-card__header {
    flex-direction: column;
  }
}
</style>
