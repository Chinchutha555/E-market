<template>
  <div class="booking-skeleton" role="status" aria-live="polite" aria-label="กำลังโหลดรายการจอง">
    <article
      v-for="index in count"
      :key="index"
      class="booking-skeleton-card"
      :class="{ 'is-expanded': index === 1 }"
    >
      <div class="skeleton-card-header">
        <span class="skeleton-shape skeleton-icon"></span>
        <div class="skeleton-heading">
          <span class="skeleton-shape skeleton-label"></span>
          <span class="skeleton-shape skeleton-title"></span>
        </div>
        <span class="skeleton-shape skeleton-chip"></span>
        <span class="skeleton-shape skeleton-chip is-wide"></span>
        <span class="skeleton-shape skeleton-time"></span>
        <span class="skeleton-shape skeleton-status"></span>
      </div>

      <template v-if="index === 1">
        <div class="skeleton-progress">
          <span v-for="step in 3" :key="step" class="skeleton-step">
            <span class="skeleton-shape skeleton-step-dot"></span>
            <span class="skeleton-shape skeleton-step-label"></span>
          </span>
        </div>
        <div class="skeleton-detail-grid">
          <span v-for="detail in 3" :key="detail" class="skeleton-detail">
            <span class="skeleton-shape skeleton-detail-label"></span>
            <span class="skeleton-shape skeleton-detail-value"></span>
          </span>
        </div>
      </template>
    </article>
    <span class="sr-only">กำลังโหลดข้อมูลรายการจอง</span>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  count?: number
}>(), {
  count: 2
})
</script>

<style scoped>
.booking-skeleton {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding-top: 4px;
}

.booking-skeleton-card {
  padding: 20px 22px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-left: 4px solid var(--bd-dbeafe);
  border-radius: 12px;
  overflow: hidden;
}

.skeleton-card-header {
  display: flex;
  align-items: center;
  gap: 14px;
}

.skeleton-heading {
  width: 150px;
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.skeleton-shape {
  display: block;
  background: linear-gradient(90deg, var(--bg-f1f5f9) 25%, var(--bg-e2e8f0) 45%, var(--bg-f1f5f9) 65%);
  background-size: 240% 100%;
  animation: bookingSkeletonShimmer 1.35s ease-in-out infinite;
}

.skeleton-icon { width: 42px; height: 42px; border-radius: 50%; flex-shrink: 0; }
.skeleton-label { width: 92px; height: 9px; border-radius: 5px; }
.skeleton-title { width: 118px; height: 18px; border-radius: 6px; }
.skeleton-chip { width: 92px; height: 30px; border-radius: 7px; }
.skeleton-chip.is-wide { width: 150px; }
.skeleton-time { width: 210px; height: 13px; border-radius: 6px; }
.skeleton-status { width: 106px; height: 30px; border-radius: 999px; margin-left: auto; }

.skeleton-progress {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-top: 22px;
  padding: 20px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-f1f5f9);
  border-radius: 10px;
}

.skeleton-step,
.skeleton-detail {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 9px;
}

.skeleton-step-dot { width: 30px; height: 30px; border-radius: 50%; }
.skeleton-step-label { width: 84px; height: 10px; border-radius: 5px; }

.skeleton-detail-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 28px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--bd-f1f5f9);
}

.skeleton-detail { align-items: flex-start; }
.skeleton-detail-label { width: 86px; height: 9px; border-radius: 5px; }
.skeleton-detail-value { width: 142px; max-width: 85%; height: 16px; border-radius: 5px; }

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

@keyframes bookingSkeletonShimmer {
  0% { background-position: 100% 0; }
  100% { background-position: -100% 0; }
}

@media (prefers-reduced-motion: reduce) {
  .skeleton-shape { animation: none; }
}

@media (max-width: 768px) {
  .skeleton-card-header { flex-wrap: wrap; }
  .skeleton-time { width: 160px; }
  .skeleton-status { margin-left: 0; }
  .skeleton-progress,
  .skeleton-detail-grid { grid-template-columns: 1fr; }
  .skeleton-step { flex-direction: row; justify-content: flex-start; }
}
</style>
