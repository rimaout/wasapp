<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import ImagePicker from './ui/ImagePicker.vue';
import ConfirmBar from './ui/ConfirmBar.vue';
import { fetchAvatar, updateAvatar, deleteAvatar } from '../services/api.js';
import { getErrorMessage } from '../services/utils.js';

/**
 * ImageActionScreen — an inline panel shown in place of the popup menu, used to
 * change or remove the image of a profile or group chat. Supply it as the
 * `component` of an ActionMenu item.
 *
 * @property {string} [title=''] - heading shown at the top of the panel.
 * @property {Object} target - what is being edited; either { kind:'me', userId }
 *   or { kind:'group', chatId }. Used to fetch and update the avatar.
 */
const props = defineProps({
	title:  { type: String, default: ''    },
	target: { type: Object, required: true },
});
/**
 * Events:
 *   done({ action:'image' }) - fired after the image is changed or removed.
 *   cancel - fired when the user cancels.
 */
const emit = defineEmits(['done', 'cancel']);

const imageFile = ref(null);
const currentImageUrl = ref(null);
const removeRequested = ref(false);
const errormsg = ref(null);
const busy = ref(false);

const isValid = computed(() => imageFile.value != null || removeRequested.value);

const previewUrl = computed(() => removeRequested.value ? null : currentImageUrl.value);

const canRemove = computed(() => currentImageUrl.value != null && imageFile.value == null && !removeRequested.value);

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
		if (imageFile.value) {
			await updateAvatar(props.target, imageFile.value);
		} else if (removeRequested.value) {
			await deleteAvatar(props.target);
		}
		emit('done', { action: 'image' });
	} catch (e) {
		errormsg.value = getErrorMessage(e);
	} finally {
		busy.value = false;
	}
}

onMounted(loadCurrentImage);

onBeforeUnmount(() => {
	if (currentImageUrl.value) URL.revokeObjectURL(currentImageUrl.value);
});
</script>

<template>
	<div class="action-screen image">
		<div v-if="title" class="action-title">{{ title }}</div>

		<ImagePicker :model-value="imageFile" :preview-url="previewUrl" @update:model-value="onFileSelected" @error="onPickerError" />
		<button v-if="canRemove" type="button" class="remove-image-btn" @click="requestRemove">Remove image</button>

		<ErrorMsg v-if="errormsg" :msg="errormsg" />

		<ConfirmBar
			:cancel-disabled="busy"
			:confirm-disabled="busy || !isValid"
			confirm-text="Confirm"
			confirm-icon="check"
			@cancel="emit('cancel')"
			@confirm="confirm"
		/>
	</div>
</template>

<style scoped>
.action-screen.image {
	display: flex;
	flex-direction: column;
	gap: 12px;
	padding: 4px 8px 8px;
	width: 200px;
}

.action-title {
	color: var(--theme-fg);
	font-weight: 600;
	font-size: 1rem;
}

.remove-image-btn {
	align-self: center;
	background: none;
	border: none;
	color: var(--theme-red);
	font-size: 0.85rem;
	cursor: pointer;
	padding: 0;
}

.remove-image-btn:hover {
	text-decoration: underline;
}
</style>
