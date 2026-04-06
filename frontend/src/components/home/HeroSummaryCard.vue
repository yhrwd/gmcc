<script setup lang="ts">
import BaseCard from '@/components/shared/BaseCard.vue'
import type { HomeHeroView } from '@/types/view'

defineProps<{
  hero: HomeHeroView | null
}>()
</script>

<template>
  <BaseCard eyebrow="系统概览" title="主面板" accent="mint">
    <div v-if="hero" class="hero-card">
      <div>
        <h2 class="hero-card__title">{{ hero.title }}</h2>
        <p class="hero-card__subtitle">{{ hero.subtitle }}</p>
      </div>
      <div class="hero-card__metrics" aria-label="资源摘要">
        <div class="hero-card__metric">
          <span class="hero-card__metric-label">CPU</span>
          <strong class="hero-card__metric-value">{{ hero.cpuLabel }}</strong>
        </div>
        <div class="hero-card__metric">
          <span class="hero-card__metric-label">内存</span>
          <strong class="hero-card__metric-value">{{ hero.memoryLabel }}</strong>
        </div>
        <div class="hero-card__metric hero-card__metric--wide">
          <span class="hero-card__metric-label">最新采样</span>
          <strong class="hero-card__metric-value">{{ hero.sampleLabel }}</strong>
        </div>
      </div>
    </div>
  </BaseCard>
</template>

<style scoped>
.hero-card {
  display: grid;
  gap: 0.55rem;
  position: relative;
  min-height: 8.75rem;
  align-content: start;
  padding-right: 4.5rem;
  animation: hero-fade-in 420ms cubic-bezier(0.22, 1, 0.36, 1);
}

.hero-card::after {
  content: '概览';
  position: absolute;
  top: 0;
  right: 0;
  padding: 0.35rem 0.68rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.86);
  color: var(--text-soft);
  font-size: 0.72rem;
  font-weight: 700;
}

.hero-card__title {
  margin: 0;
  font-size: clamp(1.8rem, 3vw, 2.8rem);
  line-height: 1;
  color: var(--text-main);
  max-width: 12ch;
}

.hero-card__subtitle {
  margin: 0.25rem 0 0;
  color: var(--text-soft);
  line-height: 1.5;
  max-width: 30rem;
}

.hero-card__metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.7rem;
  margin-top: 0.7rem;
}

.hero-card__metric {
  display: grid;
  gap: 0.28rem;
  padding: 0.8rem 0.9rem;
  border-radius: 1rem;
  background: rgba(255, 250, 242, 0.62);
  border: 1px solid rgba(131, 167, 141, 0.14);
  transform: translateY(0);
  transition: transform 220ms cubic-bezier(0.22, 1, 0.36, 1), box-shadow 220ms cubic-bezier(0.22, 1, 0.36, 1), background-color 220ms cubic-bezier(0.22, 1, 0.36, 1);
}

.hero-card__metric:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 22px rgba(162, 132, 82, 0.14);
}

.hero-card__metric:nth-child(1) {
  box-shadow: inset 0 0 0 1px rgba(242, 142, 22, 0.12);
}

.hero-card__metric:nth-child(2) {
  box-shadow: inset 0 0 0 1px rgba(131, 167, 141, 0.12);
}

.hero-card__metric:nth-child(3) {
  box-shadow: inset 0 0 0 1px rgba(237, 195, 174, 0.18);
}

.hero-card__metric-label {
  font-size: 0.76rem;
  color: var(--text-soft);
}

.hero-card__metric-value {
  font-size: 1.02rem;
  color: var(--text-main);
}

.hero-card__metric--wide {
  min-width: 0;
}

@keyframes hero-fade-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 740px) {
  .hero-card {
    padding-right: 0;
  }

  .hero-card__metrics {
    grid-template-columns: 1fr;
  }
}

@media (prefers-reduced-motion: reduce) {
  .hero-card,
  .hero-card__metric {
    animation: none;
    transition: none;
  }
}

</style>
