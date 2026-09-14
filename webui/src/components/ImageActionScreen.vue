<script>
import ImagePicker from './ui/ImagePicker.vue';
import ConfirmBar from './ui/ConfirmBar.vue';
import { fetchAvatar, updateAvatar, deleteAvatar } from '../services/api.js';
import { getErrorMessage } from '../services/utils.js';

/**
 * ImageActionScreen component allows users to change or remove the avatar image for a user or group.
 *
 * Used in:
 *  - ProfileBar.vue to change the current user's avatar.
 *  - ChatView.vue to change a group's avatar.
 *
 * Props:
 *  - title (String): The title displayed at the top of the screen.
 *  - target (Object): The target entity for which the avatar is being changed, structure: { kind: 'me' | 'group', userId | chatId }.
 *
 * Emits:
 *  - 'done' - emitted after the image is changed or removed, with an object containing the action.
 *  - 'cancel' - emitted when the user cancels the action.
 */
export default {
	components: { ImagePicker, ConfirmBar },
	props: {
		title:  { type: String, default: ''    },
		target: { type: Object, required: true },
	},

	/**
	 * Events:
	 *   done({ action:'image' }) - fired after the image is changed or removed.
	 *   cancel - fired when the user cancels.
     *
	 * Note: action: 'image' is needed to distinguish this event from other 'done' events in the parent component.
	 */
	emits: ['done', 'cancel'],

	data() {
		return {
			imageFile: null,
			currentImageUrl: null,
			removeRequested: false,
			errormsg: null,
			busy: false,
		};
	},

	computed: {
		isValid() {
			// The action is valid if there is a new image file selected or if the user has requested to remove the current image.
			return this.imageFile != null || this.removeRequested;
		},
		previewUrl() {
			// If a new image file is selected, create a temporary URL for preview.
			return this.removeRequested ? null : this.currentImageUrl;
		},
		canRemove() {
			// The user can remove the image if there is a current image, no new image file is selected, and the remove action has not been requested yet
			return this.currentImageUrl != null && this.imageFile == null && !this.removeRequested;
		},
	},

	methods: {
		async loadCurrentImage() {
			try {
				this.currentImageUrl = URL.createObjectURL(await fetchAvatar(this.target));
			} catch (e) {
				this.currentImageUrl = null;
			}
		},

		// Event handler for <ImagePicker>'s 'update:selected-image' event.
		onFileSelected(file) {
			this.imageFile = file;
			if (file) this.removeRequested = false;
		},

		// Event handler for <ImagePicker>'s 'error' event.
		onPickerError(message) {
			this.errormsg = message;
		},

		// Hadler for the "Remove image" button click event.
		requestRemove() {
			this.removeRequested = true;
		},

		// Confirm button click handler: updates or removes the avatar image based on user selection.
		async confirm() {
			if (this.busy || !this.isValid) return;

			this.busy = true;
			this.errormsg = null;

			try {
				if (this.imageFile) {
					await updateAvatar(this.target, this.imageFile);
				} else if (this.removeRequested) {
					await deleteAvatar(this.target);
				}
				this.$emit('done', { action: 'image' });
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			} finally {
				this.busy = false;
			}
		},
	},

	mounted() {
		this.loadCurrentImage();
	},

	beforeUnmount() {
		if (this.currentImageUrl) URL.revokeObjectURL(this.currentImageUrl);
	},
};
</script>

<template>
	<div class="action-screen image">

		<!-- Title Display -->
		<div v-if="title" class="action-title">{{ title }}</div>

		<!-- Image Picker Component (clickable image to change the image) -->
		<ImagePicker :selected-image="imageFile" :preview-url="previewUrl" @update:selected-image="onFileSelected" @error="onPickerError" />

		<!-- Remove Image Button -->
		<button v-if="canRemove" type="button" class="remove-image-btn" @click="requestRemove">Remove image</button>

		<!-- Error Message Display -->
		<ErrorMsg v-if="errormsg" :msg="errormsg" />

		<!-- ConfirmBar Component for confirming or canceling the action -->
		<ConfirmBar
			:cancel-disabled="busy"
			:confirm-disabled="busy || !isValid"
			confirm-text="Confirm"
			confirm-icon="check"
			@cancel="$emit('cancel')"
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
