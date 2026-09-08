<script>
/**
 * MessageStatusIcon — compact read-receipt icon for a message status.
 * 'DELIVERED' renders a single check; 'RECEIVED'/'READ' render a double check.
 * The icon inherits `currentColor`, so the parent controls its color (an
 * `is-read` class is added to the root for the read accent).
 *
 * @property {'DELIVERED'|'RECEIVED'|'READ'} status - the message status.
 */
export default {
	props: {
		status: { type: String, required: true }, // Message status to render
	},
};
</script>

<template>
	<span class="status-icon" :class="{ 'is-received': status === 'RECEIVED', 'is-read': status === 'READ' }">
		<i v-if="status === 'DELIVERED'" class="bi bi-check"></i>
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
