<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import axios from '../services/axios.js';
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
const imagePreview = ref(null);
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

const previewSrc = computed(() => {
	if (removeRequested.value) return null;
	return imagePreview.value || currentImageUrl.value;
});

const canRemove = computed(() => {
	return currentImageUrl.value != null && imageFile.value == null && !removeRequested.value;
});

function currentAvatarUrl() {
	if (props.target.kind === 'me') {
		return '/users/' + props.target.userId + '/avatar';
	}
	return '/chats/' + props.target.chatId + '/avatar';
}

async function loadCurrentImage() {
	try {
		const response = await axios.get(currentAvatarUrl(), { responseType: 'blob' });
		currentImageUrl.value = URL.createObjectURL(response.data);
	} catch (e) {
		currentImageUrl.value = null;
	}
}

function onFileChange(e) {
	const file = e.target.files && e.target.files[0];
	if (!file) return;

	if (file.size > 5 * 1024 * 1024) {
		errormsg.value = 'Image must be at most 5 MB';
		e.target.value = '';
		return;
	}

	imageFile.value = file;
	removeRequested.value = false;
	if (imagePreview.value) URL.revokeObjectURL(imagePreview.value);
	imagePreview.value = URL.createObjectURL(file);
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
			const newName = name.value.trim();
			if (props.target.kind === 'me') {
				const res = await axios.patch('/me/name', { userName: newName });
				emit('done', { action: 'rename', name: res.data.name });
			} else {
				const res = await axios.patch('/chats/' + props.target.chatId + '/name', { groupName: newName });
				emit('done', { action: 'rename', name: res.data.groupName });
			}
		} else if (imageFile.value) {
			const formData = new FormData();
			formData.append('binaryImage', imageFile.value);
			if (props.target.kind === 'me') {
				await axios.put('/me/avatar', formData);
			} else {
				await axios.put('/chats/' + props.target.chatId + '/avatar', formData);
			}
			emit('done', { action: 'image' });
		} else if (removeRequested.value) {
			if (props.target.kind === 'me') {
				await axios.delete('/me/avatar');
			} else {
				await axios.delete('/chats/' + props.target.chatId + '/avatar');
			}
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
	if (imagePreview.value) URL.revokeObjectURL(imagePreview.value);
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
			<label class="image-picker">
				<img v-if="previewSrc" :src="previewSrc" class="image-preview" />
				<span v-else class="image-placeholder">
					<svg class="feather camera-icon"><use href="/feather-sprite-v4.29.0.svg#camera"/></svg>
				</span>
				<input type="file" accept="image/jpeg,image/png,image/webp" class="d-none" @change="onFileChange" />
			</label>
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

.image-picker {
	width: 96px;
	height: 96px;
	border-radius: 50%;
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
	background: var(--tn-bg-highlight);
	color: var(--tn-fg-dark);
}

.camera-icon {
	width: 28px;
	height: 28px;
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
