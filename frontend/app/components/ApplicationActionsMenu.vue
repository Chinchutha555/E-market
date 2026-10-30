<template>
  <div v-if="hasActions" ref="root" class="more-root" @click.stop>
    <AppButton :size="size" variant="outline" :aria-expanded="open" aria-haspopup="menu" @click.stop="toggle">
      <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="5" cy="12" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/></svg></template>เพิ่มเติม
    </AppButton>
    <Teleport to="body"><div v-if="open" ref="menu" class="more-menu" role="menu" :style="menuStyle">
      <button v-if="showDocument" type="button" role="menuitem" @click="choose('document')"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg>ตัวอย่างคำร้อง</button>
      <button v-if="showIdCopy" type="button" role="menuitem" @click="choose('idCopy')"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="16" rx="2"/><circle cx="9" cy="10" r="2"/><path d="M6 16c.5-2 5.5-2 6 0M15 9h3M15 13h3"/></svg>ดูสำเนาบัตรประชาชน</button>
      <div v-if="canReview && (showDocument || showIdCopy)" class="more-divider"></div>
      <button v-if="canReview" type="button" role="menuitem" @click="choose('revision')"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 4v6h6"/><path d="M3.5 15a9 9 0 1 0 2-9.4L1 10"/></svg>ส่งกลับแก้ไข</button>
      <button v-if="canReview && showReject" type="button" role="menuitem" class="danger" @click="choose('reject')"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"/><path d="m15 9-6 6M9 9l6 6"/></svg>ไม่อนุมัติ</button>
    </div></Teleport>
  </div>
</template>
<script setup lang="ts">
const props = withDefaults(defineProps<{ showDocument?: boolean; showIdCopy?: boolean; canReview?: boolean; showReject?: boolean; size?: 'sm' | 'md' }>(), { size: 'sm', showReject: true })
const emit = defineEmits<{ document: []; idCopy: []; revision: []; reject: [] }>()
const root = ref<HTMLElement | null>(null)
const menu = ref<HTMLElement | null>(null)
const open = ref(false)
const menuStyle = ref<Record<string, string>>({})
const hasActions = computed(() => props.showDocument || props.showIdCopy || props.canReview)
const position = () => {
  const rect = root.value?.getBoundingClientRect()
  if (!rect) return
  const width = 230
  const height = (Number(props.showDocument) + Number(props.showIdCopy) + (props.canReview ? (props.showReject ? 2 : 1) : 0)) * 42 + 20
  menuStyle.value = { left: Math.max(8, Math.min(window.innerWidth - width - 8, rect.right - width)) + 'px', top: (rect.bottom + height + 8 > window.innerHeight ? Math.max(8, rect.top - height - 6) : rect.bottom + 6) + 'px' }
}
const toggle = () => { if (!open.value) position(); open.value = !open.value }
const choose = (action: 'document' | 'idCopy' | 'revision' | 'reject') => { open.value = false; if (action === 'document') emit('document'); else if (action === 'idCopy') emit('idCopy'); else if (action === 'revision') emit('revision'); else emit('reject') }
const dismiss = (event: Event) => { const target = event.target as Node; if (!root.value?.contains(target) && !menu.value?.contains(target)) open.value = false }
const keydown = (event: KeyboardEvent) => { if (event.key === 'Escape' && open.value) { open.value = false; root.value?.querySelector('button')?.focus() } }
const closeOnScroll = () => { open.value = false }
onMounted(() => { document.addEventListener('pointerdown', dismiss); document.addEventListener('keydown', keydown); window.addEventListener('resize', closeOnScroll); window.addEventListener('scroll', closeOnScroll, true) })
onBeforeUnmount(() => { document.removeEventListener('pointerdown', dismiss); document.removeEventListener('keydown', keydown); window.removeEventListener('resize', closeOnScroll); window.removeEventListener('scroll', closeOnScroll, true) })
</script>
<style scoped>
.more-root { display: inline-flex; }
.more-menu { position: fixed; z-index: 1500; width: 230px; padding: 5px; border: 1px solid #e2e8f0; border-radius: 10px; background: #fff; box-shadow: 0 8px 24px rgba(15,23,42,.14); }
.more-menu button { display: flex; align-items: center; gap: 10px; width: 100%; min-height: 40px; padding: 8px 10px; border: 0; border-radius: 7px; background: transparent; color: #334155; font: inherit; font-size: .875rem; text-align: left; cursor: pointer; }
.more-menu button:hover, .more-menu button:focus-visible { background: #f1f5f9; outline: none; }
.more-menu button.danger { color: #dc2626; } .more-menu button.danger:hover { background: #fef2f2; }
.more-menu svg { width: 16px; height: 16px; flex: 0 0 16px; } .more-divider { height: 1px; margin: 4px 0; background: #f1f5f9; }
</style>
