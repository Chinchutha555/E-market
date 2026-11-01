<template>
  <div class="body-breadcrumb-wrapper">
    <div class="body-breadcrumb-left">
      <!-- Back Button: Square rounded button with left arrow icon matching user mockup -->
      <button
        class="body-breadcrumb-back-btn"
        type="button"
        @click="handleBack"
        title="ย้อนกลับ"
        aria-label="ย้อนกลับ"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="15 18 9 12 15 6"></polyline>
        </svg>
      </button>

      <!-- Breadcrumb Items -->
      <nav class="body-breadcrumb-trail" aria-label="Breadcrumb">
        <NuxtLink v-if="!isStaffUser" to="/" class="body-breadcrumb-link">หน้าหลัก</NuxtLink>

        <!-- Optional Market Parent (e.g. ตลาดธนบุรี (สนามหลวง 2)) -->
        <template v-if="shouldShowMarketParent">
          <span v-if="!isStaffUser" class="body-breadcrumb-sep">&gt;</span>
          <NuxtLink v-if="current || (items && items.length)" :to="resolvedMarketId ? `/stalls?market=${resolvedMarketId}` : '/stalls'" class="body-breadcrumb-link">
            {{ resolvedMarketName }}
          </NuxtLink>
          <span v-else class="body-breadcrumb-current">{{ resolvedMarketName }}</span>
        </template>

        <!-- Optional Intermediate Items List -->
        <template v-for="(item, idx) in items" :key="idx">
          <span v-if="!isStaffUser || shouldShowMarketParent || idx > 0" class="body-breadcrumb-sep">&gt;</span>
          <button
            v-if="item.onClick"
            type="button"
            class="body-breadcrumb-link body-breadcrumb-btn"
            @click="item.onClick()"
          >
            {{ item.label }}
          </button>
          <NuxtLink v-else-if="item.to" :to="item.to" class="body-breadcrumb-link">{{ item.label }}</NuxtLink>
          <span v-else class="body-breadcrumb-current">{{ item.label }}</span>
        </template>

        <!-- Active / Current Page Label -->
        <template v-if="current">
          <span v-if="!isStaffUser || shouldShowMarketParent || (items && items.length)" class="body-breadcrumb-sep">&gt;</span>
          <span class="body-breadcrumb-current">{{ current }}</span>
        </template>
      </nav>
    </div>

    <!-- Optional Right Slot (e.g. Action buttons) -->
    <div v-if="$slots.actions" class="body-breadcrumb-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { getMarketGisConfig } from '~/core/modules/stall/data/markets_registry'

const props = defineProps<{
  current?: string
  market?: string
  marketId?: string | number
  items?: Array<{ label: string; to?: string; onClick?: () => void }>
  backTo?: string
  onBack?: () => void
  forceNavigation?: boolean
}>()

const route = useRoute()
const { applyMarketOverride } = useMarketOverrides()
const { user } = useAuth()
const isStaffUser = computed(() => user.value?.role === 'staff')

const resolvedMarketId = computed(() => {
  const raw = props.marketId ?? route.query.market ?? (route.path.startsWith('/markets/') ? route.params.id : undefined)
  return raw !== undefined && raw !== null && String(raw).trim() !== '' ? String(raw).trim() : undefined
})

const resolvedMarketName = computed(() => {
  if (!props.market && !resolvedMarketId.value) return ''

  if (resolvedMarketId.value) {
    const cfg = getMarketGisConfig(resolvedMarketId.value)
    const overridden = applyMarketOverride({ id: cfg.id, slug: cfg.slug, name: cfg.name })
    const expectedName = overridden?.name || cfg.name
    // If parent passed default 'ตลาดนัดจตุจักร' before async load finished for another market, use the resolved market name immediately
    if (!props.market || (props.market === 'ตลาดนัดจตุจักร' && String(cfg.id) !== '1')) {
      return expectedName
    }
  }

  return props.market || ''
})

const shouldShowMarketParent = computed(() => {
  if (!resolvedMarketName.value) return false
  if (route.path === '/stalls') return false
  if (props.current && resolvedMarketName.value.trim() === props.current.trim()) return false
  return true
})

const handleBack = async () => {
  if (props.onBack) {
    props.onBack()
    return
  }

  if (props.backTo) {
    await navigateTo(props.backTo)
    return
  }

  if (resolvedMarketId.value && route.path !== '/stalls') {
    await navigateTo(`/stalls?market=${resolvedMarketId.value}`)
    return
  }

  if (import.meta.client && window.history.state && typeof window.history.state.back === 'string') {
    const backUrl = window.history.state.back
    const backPathOnly = backUrl.split('?')[0]?.split('#')[0] || ''
    const isAuthPage = backPathOnly === '/login' || backPathOnly === '/auth/login'
    const isSameFullUrl = backUrl === route.fullPath
    const isRedirectLoop =
      (route.path === '/stalls/manage' && backPathOnly === '/stalls/book') ||
      (route.path === '/stalls/book' && backPathOnly === '/stalls/manage')

    if (!isAuthPage && !isSameFullUrl && !isRedirectLoop) {
      await navigateTo(backUrl)
      return
    }
  }

  await navigateTo('/')
}
</script>

<style scoped>
.body-breadcrumb-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
  padding: 2px 0;
  flex-wrap: wrap;
}

.body-breadcrumb-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.body-breadcrumb-back-btn {
  width: 28px;
  height: 28px;
  border-radius: 7px;
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--tx-475569);
  transition: all 0.2s ease;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
  flex-shrink: 0;
}

.body-breadcrumb-back-btn:hover {
  background: var(--bg-f8fafc);
  border-color: var(--bd-cbd5e1);
  color: var(--primary);
  transform: translateX(-1px);
}

.body-breadcrumb-back-btn:active {
  transform: translateX(0);
}

.body-breadcrumb-trail {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 0.85rem;
  font-family: 'Bai Jamjuree', sans-serif;
  flex-wrap: wrap;
  line-height: 1.4;
}

.body-breadcrumb-link {
  color: var(--tx-64748b);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.15s ease;
}

.body-breadcrumb-btn {
  border: 0;
  background: transparent;
  padding: 0;
  font: inherit;
  cursor: pointer;
}

.body-breadcrumb-link:hover {
  color: var(--primary);
  text-decoration: underline;
}

.body-breadcrumb-sep {
  color: var(--tx-cbd5e1);
  font-size: 0.75rem;
  font-weight: 300;
  user-select: none;
}

.body-breadcrumb-current {
  color: var(--tx-0f172a);
  font-weight: 600;
  letter-spacing: -0.01em;
}

.body-breadcrumb-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
