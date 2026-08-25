<script setup>
/**
 * MessageStatusIcon — compact read-receipt icon for a message status.
 * 'delivered' renders a single check; 'received'/'read' render a double check.
 * The icon inherits `currentColor`, so the parent controls its color (an
 * `is-read` class is added to the root for the read accent).
 *
 * @property {'delivered'|'received'|'read'} status - the message status.
 */
defineProps({
	status: { type: String, required: true },
});
</script>

<template>
	<span class="status-icon" :class="{ 'is-received': status === 'received', 'is-read': status === 'read' }">
		<i v-if="status === 'delivered'" class="bi bi-check"></i>
		<i v-else class="bi bi-check-all"></i>
	</span>
</template>

<style scoped>
.status-icon {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	color: currentColor;
	font-size: 15px;
	line-height: 1;
}
</style>
