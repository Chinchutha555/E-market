<template>
  <AppImageUpload
    :model-value="preview"
    :label="title"
    :hint="hint"
    accept="image/png,image/jpeg"
    :max-size-mb="1"
    select-text="เพิ่มรูปภาพ"
    :preview-alt="title"
    @change="onChange"
  />
</template>
<script setup lang="ts">
import AppImageUpload from '~/components/AppImageUpload.vue'

defineProps<{ title: string; hint: string; preview?: string }>()
const emit = defineEmits<{ (e: 'file', file: File): void }>()
const onChange = (file: File | File[] | null) => { if (file instanceof File) emit('file', file) }
</script>
