<template>
  <StallManageMain />
</template>

<script setup lang="ts">
import StallManageMain from '~/core/modules/stall/ui/StallManageMain.vue'
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'

definePageMeta({
  keepalive: true,
  middleware: (to) => {
    const authStore = useAuthStore()
    authStore.hydrate()
    if (import.meta.client && !authStore.isLoggedIn) {
      return navigateTo('/auth/login')
    }
    if (import.meta.client && !['admin', 'staff'].includes(authStore.role)) {
      return navigateTo({ path: '/stalls/book', query: to.query }, { replace: true })
    }
    if (import.meta.client && authStore.role === 'staff') {
      const assignedId = authStore.currentUser?.assigned_markets?.[0]?.id
      if (!assignedId) return navigateTo('/staff', { replace: true })
      if (String(to.query.market || '') !== String(assignedId)) {
        return navigateTo({ path: '/stalls/map', query: { ...to.query, market: assignedId } }, { replace: true })
      }
    }
  }
})
</script>
