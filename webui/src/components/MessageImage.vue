<script>
/**
 * MessageImage - A component that loads and displays an image from a chat message.
 *
 * Used in:
 *  - MessageBubble.vue to display images in chat messages
 *  - MessageInput.vue to display image previews when replying to a message.
 *
 * Props:
 *  - chatId: The ID of the chat the image belongs to.
 *  - imageId: The ID of the image to load.
 *  - thumbnail: Whether to display the image as a thumbnail (small square, used for image previews in the chat list). Defaults to false.
 *
 * Emits: None
 */
export default {
	props: {
		chatId:    { type: String,  required: true },
		imageId:   { type: String,  required: true },
		thumbnail: { type: Boolean, default: false },
	},

	data() {
		return {
			loading: true, // Whether the image is currently being loaded
			src: null,     // The blob URL of the loaded image (null if not loaded or failed)
			failed: false, // Whether the image failed to load
		};
	},

	watch: {
		// Reload when the chat or image changes
		chatId:  'load',
		imageId: 'load',
	},

	created() {
		this.load(); // Load the image when the component is created
	},

	methods: {
		async load() {
			this.loading = true;
			this.failed = false;

			// Release the previous blob URL before loading a new image
			if (this.src) {
				URL.revokeObjectURL(this.src);
				this.src = null;
			}

			// Request the image from the server as a blob (binary data)
			try {
				const res = await this.$axios.get('/chats/' + this.chatId + '/images/' + this.imageId, { responseType: 'blob' });
				this.src = URL.createObjectURL(res.data);
			} catch (e) {
				this.failed = true;
			}
			this.loading = false;
		},
	},

	beforeUnmount() {
		if (this.src) URL.revokeObjectURL(this.src); // Release the blob URL when the component is destroyed
	},
};
</script>

<template>
	<div class="message-image" :class="{ thumbnail }">

		<!-- Show a loading spinner while the image is being loaded -->
		<LoadingSpinner v-if="loading" :loading="loading" />

		<!-- Show the image if it was successfully loaded -->
		<img v-else-if="src" :src="src" alt="Message image" class="message-image-img" />

		<!-- Show a placeholder if the image failed to load -->
		<span v-else-if="failed" class="message-image-failed">Image unavailable</span>
	</div>
</template>

<style scoped>
.message-image {
	display: flex;
	justify-content: center;
	align-items: center;
	min-width: 100px;
	min-height: 56px;
}

.message-image-img {
	max-width: 100%;
	max-height: 200px;
	border-radius: 8px;
	display: block;
}

.message-image-failed {
	color: var(--theme-fg-dark);
	font-size: 0.85rem;
}

.message-image.thumbnail {
	min-width: 0;
	min-height: 0;
	width: 40px;
	height: 40px;
	flex-shrink: 0;
	overflow: hidden;
}

.message-image.thumbnail .message-image-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
	max-height: none;
}
</style>
