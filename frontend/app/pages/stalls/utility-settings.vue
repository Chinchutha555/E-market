<template><UtilitySettingsMain /></template>
<script setup lang="ts">
import UtilitySettingsMain from '~/core/modules/stall/ui/UtilitySettingsMain.vue'
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'

definePageMeta({
  middleware: (to) => {
    if (!import.meta.client) return
    const auth = useAuthStore()
    auth.hydrate()
    if (!auth.isLoggedIn) return navigateTo('/auth/login')
    if (auth.role === 'staff') {
      const assignedId = auth.currentUser?.assigned_markets?.[0]?.id
      if (!assignedId || String(to.query.market || '') !== String(assignedId)) return navigateTo('/staff', { replace: true })
    }
  }
})
</script>
