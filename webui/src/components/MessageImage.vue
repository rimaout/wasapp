<script setup>
import { ref, watch, onBeforeUnmount } from 'vue';
import axios from '../services/axios.js';

// Loads a message image as an authenticated blob (the message image endpoint
// requires a Bearer token, so a plain <img src> cannot be used).
const props = defineProps({
	chatId: { type: String, required: true },
	imageId: { type: String, required: true },
	thumbnail: { type: Boolean, default: false },
});

const loading = ref(true);
const src = ref(null);
const failed = ref(false);

async function load() {
	loading.value = true;
	failed.value = false;
	if (src.value) {
		URL.revokeObjectURL(src.value);
		src.value = null;
	}
	try {
		const res = await axios.get('/chats/' + props.chatId + '/images/' + props.imageId, { responseType: 'blob' });
		src.value = URL.createObjectURL(res.data);
	} catch (e) {
		failed.value = true;
	}
	loading.value = false;
}

watch(() => [props.chatId, props.imageId], load, { immediate: true });

onBeforeUnmount(() => {
	if (src.value) URL.revokeObjectURL(src.value);
});
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
