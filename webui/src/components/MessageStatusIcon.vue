<script>

/**
 * MessageStatusIcon.vue - the icon changes based on the message status: "DELIVERED", "RECEIVED", or "READ"
 *
 * Used in:
 *  - ChatListView.vue to show the status of the last message in a chat preview.
 *  - MessageBubble.vue to show the status of a message in the chat view.
 * Props:
 *  - status (String): The status of the message, which can be "DELIVERED", "RECEIVED", or "READ".
 *
 * Emits: None
 */
export default {
	props: {
		status: { type: String, required: true }, // Message status to render
	},
};
</script>

<template>
	<span class="status-icon" :class="{ 'is-received': status === 'RECEIVED', 'is-read': status === 'READ' }">
		<!-- Single checkmark icon from Bootstrap Icons for 'DELIVERED' status -->
		<i v-if="status === 'DELIVERED'" class="bi  bi-check"></i>
		<!-- Double checkmark icon from Bootstrap Icons for 'RECEIVED' status -->
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
