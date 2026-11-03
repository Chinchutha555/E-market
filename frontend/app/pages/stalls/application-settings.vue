<template>
  <TradeContractDraft />
</template>
<script setup lang="ts">
import TradeContractDraft from '~/components/TradeContractDraft.vue'
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'

definePageMeta({
  middleware: (to) => {
    if (!import.meta.client) return
    const auth = useAuthStore()
    auth.hydrate()
    if (!auth.isLoggedIn) return navigateTo('/auth/login')
    if (auth.role === 'staff' && to.query.market) {
      const assigned = auth.currentUser?.assigned_markets?.some((m: any) => String(m.id) === String(to.query.market || ''))
      if (!assigned) return navigateTo('/staff', { replace: true })
    }
  }
})
</script>
