<template>
  <div class="staff-entry-skeleton" aria-label="กำลังเปิดพื้นที่จัดการตลาด">
    <span></span><span></span><span></span>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'
import { AuthService } from '~/core/modules/auth/infrastructure/AuthService'

const auth = useAuthStore()
onMounted(async () => {
  auth.hydrate()
  if (!auth.isLoggedIn || auth.role !== 'staff') return
  if (!auth.currentUser?.assigned_markets?.length) {
    const profile = await new AuthService().getProfile()
    if (profile.status && profile.data) auth.setUser(profile.data)
  }
  const market = auth.currentUser?.assigned_markets?.[0]
  if (market) await navigateTo(`/stalls/map?market=${market.id}`, { replace: true })
})

definePageMeta({
  middleware: () => {
    if (!import.meta.client) return
    const session = useAuthStore()
    session.hydrate()
    if (!session.isLoggedIn) return navigateTo('/auth/login')
    if (session.role === 'admin') return navigateTo('/', { replace: true })
    if (session.role !== 'staff') return navigateTo('/', { replace: true })
    const market = session.currentUser?.assigned_markets?.[0]
    if (!market) return
    return navigateTo(`/stalls/map?market=${market.id}`, { replace: true })
  }
})
</script>

<style scoped>
.staff-entry-skeleton { padding: 40px; display: grid; gap: 16px; }
.staff-entry-skeleton span { display: block; height: 76px; border-radius: 16px; background: linear-gradient(90deg,#eef2f7 25%,#f8fafc 50%,#eef2f7 75%); background-size: 200% 100%; animation: staff-shimmer 1.4s infinite; }
.staff-entry-skeleton span:last-child { height: 520px; }
@keyframes staff-shimmer { to { background-position: -200% 0; } }
</style>
