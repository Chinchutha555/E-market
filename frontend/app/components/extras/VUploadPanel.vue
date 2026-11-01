<template>
  <AppFileUpload
    v-model="selectedFile"
    :label="title"
    :hint="hint"
    :accept="accept"
    @change="handleChange"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AppFileUpload from '~/components/AppFileUpload.vue'

withDefaults(defineProps<{
  title?: string
  hint?: string
  accept?: string
}>(), {
  title: 'แนบไฟล์',
  hint: 'รองรับไฟล์ PDF, PNG, JPG ขนาดไม่เกิน 10MB',
  accept: '.pdf,.png,.jpg,.jpeg'
})

const emit = defineEmits<{ 'file-selected': [file: File] }>()
const selectedFile = ref<File | null>(null)
const handleChange = (file: File | null) => { if (file) emit('file-selected', file) }
</script>
