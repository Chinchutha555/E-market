<template>
	<ApplicationMain />
</template>

<script setup lang="ts">
import ApplicationMain from "~/core/modules/application/ui/ApplicationMain.vue"

definePageMeta({
	middleware: "auth"
})
</script>
