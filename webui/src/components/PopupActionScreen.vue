<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import ImagePicker from './ImagePicker.vue';
import { fetchAvatar, updateName, updateAvatar, deleteAvatar } from '../services/api.js';
import { getErrorMessage } from '../services/utils.js';

// Inline action panel (rename or image change) shown in place of the popup menu.
const props = defineProps({
	mode: { type: String, required: true },
	title: { type: String, default: '' },
	initialName: { type: String, default: '' },
	target: { type: Object, required: true },
});
const emit = defineEmits(['done', 'cancel']);

const name = ref(props.initialName || '');
const imageFile = ref(null);
const currentImageUrl = ref(null);
const removeRequested = ref(false);
const errormsg = ref(null);
const busy = ref(false);

const isValid = computed(() => {
	if (props.mode === 'rename') {
		const trimmed = name.value.trim();
		const current = (props.initialName || '').trim();
		return trimmed.length > 0 && trimmed !== current;
	}
	return imageFile.value != null || removeRequested.value;
});

const previewUrl = computed(() => {
	return removeRequested.value ? null : currentImageUrl.value;
});

const canRemove = computed(() => {
	return currentImageUrl.value != null && imageFile.value == null && !removeRequested.value;
});

async function loadCurrentImage() {
	try {
		currentImageUrl.value = URL.createObjectURL(await fetchAvatar(props.target));
	} catch (e) {
		currentImageUrl.value = null;
	}
}

function onFileSelected(file) {
	imageFile.value = file;
	if (file) removeRequested.value = false;
}

function onPickerError(message) {
	errormsg.value = message;
}

function requestRemove() {
	removeRequested.value = true;
}

async function confirm() {
	if (busy.value || !isValid.value) return;

	busy.value = true;
	errormsg.value = null;
	try {
		if (props.mode === 'rename') {
			const newName = await updateName(props.target, name.value.trim());
			emit('done', { action: 'rename', name: newName });
		} else if (imageFile.value) {
			await updateAvatar(props.target, imageFile.value);
			emit('done', { action: 'image' });
		} else if (removeRequested.value) {
			await deleteAvatar(props.target);
			emit('done', { action: 'image' });
		}
	} catch (e) {
		errormsg.value = getErrorMessage(e);
	} finally {
		busy.value = false;
	}
}

onMounted(() => {
	if (props.mode === 'image') loadCurrentImage();
});

onBeforeUnmount(() => {
	if (currentImageUrl.value) URL.revokeObjectURL(currentImageUrl.value);
});
</script>

<template>
	<div class="action-screen" :class="mode">
		<div v-if="title" class="action-title">{{ title }}</div>

		<template v-if="mode === 'rename'">
			<input type="text" class="form-control name-input" v-model="name" maxlength="24" placeholder="Enter name" @keyup.enter="confirm" />
		</template>

		<template v-else>
			<ImagePicker :model-value="imageFile" :preview-url="previewUrl" @update:model-value="onFileSelected" @error="onPickerError" />
			<button v-if="canRemove" type="button" class="remove-image-btn" @click="requestRemove">Remove image</button>
		</template>

		<ErrorMsg v-if="errormsg" :msg="errormsg" />

		<div class="action-footer">
			<button type="button" class="action-btn cancel" :disabled="busy" aria-label="Cancel" title="Cancel" @click="emit('cancel')">
				<svg class="feather action-icon"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
			</button>
			<button type="button" class="action-btn confirm" :disabled="busy || !isValid" @click="confirm">
				<svg class="feather confirm-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
				<span>Confirm</span>
			</button>
		</div>
	</div>
</template>

<style scoped>
.action-screen {
	display: flex;
	flex-direction: column;
	gap: 12px;
	padding: 8px;
}

.action-screen.rename {
	width: 260px;
}

.action-screen.image {
	width: 200px;
}

.action-title {
	color: var(--tn-fg);
	font-weight: 600;
	font-size: 1rem;
}

.field-label {
	color: var(--tn-fg);
	font-weight: 600;
	font-size: 0.85rem;
}

.name-input {
	background: var(--tn-bg-highlight);
	border-color: var(--tn-border);
	color: var(--tn-fg);
}

.remove-image-btn {
	align-self: center;
	background: none;
	border: none;
	color: var(--tn-red);
	font-size: 0.85rem;
	cursor: pointer;
	padding: 0;
}

.remove-image-btn:hover {
	text-decoration: underline;
}

.action-footer {
	display: flex;
	gap: 8px;
	width: 100%;
	margin-top: 4px;
}

.action-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	height: 40px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	flex-shrink: 0;
}

.action-btn:disabled {
	cursor: default;
}

.action-btn.cancel:disabled {
	opacity: 0.5;
}

.action-btn.cancel {
	width: 40px;
	background: rgba(247, 118, 142, 0.12);
	color: var(--tn-red);
}

.action-btn.cancel:hover:not(:disabled) {
	background: rgba(247, 118, 142, 0.22);
}

.action-btn.confirm {
	flex: 1 1 auto;
	background: rgba(158, 206, 106, 0.15);
	color: var(--tn-green);
	font-weight: 600;
}

.action-btn.confirm:hover:not(:disabled) {
	background: rgba(158, 206, 106, 0.25);
}

.action-btn.confirm:disabled {
	background: rgba(255, 255, 255, 0.08);
	color: var(--tn-fg-dark);
}

.action-icon {
	width: 22px;
	height: 22px;
}

.confirm-icon {
	width: 18px;
	height: 18px;
}
</style>
