<script>
// Loads a message image as an authenticated blob (the message image endpoint
// requires a Bearer token, so a plain <img src> cannot be used).
export default {
	props: {
		chatId:    { type: String, required: true  },
		imageId:   { type: String, required: true  },
		thumbnail: { type: Boolean, default: false },
	},

	data() {
		return {
			loading: true,
			src: null,
			failed: false,
		};
	},

	watch: {
		// Reload when the chat or image changes
		chatId:  'load',
		imageId: 'load',
	},

	created() {
		// Load the image when the component is created (equivalent of the previous immediate watcher)
		this.load();
	},

	methods: {
		async load() {
			this.loading = true;
			this.failed = false;
			if (this.src) {
				URL.revokeObjectURL(this.src);
				this.src = null;
			}
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
		if (this.src) URL.revokeObjectURL(this.src);
	},
};
</script>

<template>
	<div class="message-image" :class="{ thumbnail }">
		<LoadingSpinner v-if="loading" :loading="loading" />
		<img v-else-if="src" :src="src" alt="Message image" class="message-image-img" />
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
