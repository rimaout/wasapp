<script setup>
import { ref, watch, computed, onBeforeUnmount } from 'vue';

// Circular/rounded image picker: handles file selection, 5MB check and preview.
const props = defineProps({
	modelValue: { type: File, default: null },
	previewUrl: { type: String, default: null },
	radius: { type: String, default: '50%' },
	size: { type: Number, default: 96 },
});
const emit = defineEmits(['update:modelValue', 'error']);

const objectUrl = ref(null);

const iconSize = computed(() => Math.round(props.size * 28 / 96));

watch(() => props.modelValue, (file) => {
	if (objectUrl.value) URL.revokeObjectURL(objectUrl.value);
	objectUrl.value = file ? URL.createObjectURL(file) : null;
});

const displaySrc = computed(() => objectUrl.value || props.previewUrl || null);

function onFileChange(e) {
	const file = e.target.files && e.target.files[0];
	if (!file) return;

	if (file.size > 5 * 1024 * 1024) {
		emit('error', 'Image must be at most 5 MB');
		e.target.value = '';
		return;
	}

	emit('update:modelValue', file);
}

onBeforeUnmount(() => {
	if (objectUrl.value) URL.revokeObjectURL(objectUrl.value);
});
</script>

<template>
	<label class="image-picker" :style="{ borderRadius: radius, width: size + 'px', height: size + 'px' }">
		<img v-if="displaySrc" :src="displaySrc" class="image-preview" />
		<span v-else class="image-placeholder">
			<svg class="feather camera-icon" :style="{ width: iconSize + 'px', height: iconSize + 'px' }"><use href="/feather-sprite-v4.29.0.svg#camera"/></svg>
		</span>
		<input type="file" accept="image/jpeg,image/png,image/webp" class="d-none" @change="onFileChange" />
	</label>
</template>

<style scoped>
.image-picker {
	overflow: hidden;
	cursor: pointer;
	flex-shrink: 0;
	align-self: center;
	display: flex;
	align-items: center;
	justify-content: center;
}

.image-preview {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.image-placeholder {
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	background: var(--theme-bg-highlight);
	color: var(--theme-fg-dark);
}
</style>
