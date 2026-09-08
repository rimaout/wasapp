<script>
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
			return this.imageFile != null || this.removeRequested;
		},
		previewUrl() {
			return this.removeRequested ? null : this.currentImageUrl;
		},
		canRemove() {
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

		onFileSelected(file) {
			this.imageFile = file;
			if (file) this.removeRequested = false;
		},

		onPickerError(message) {
			this.errormsg = message;
		},

		requestRemove() {
			this.removeRequested = true;
		},

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
		<div v-if="title" class="action-title">{{ title }}</div>

		<ImagePicker :model-value="imageFile" :preview-url="previewUrl" @update:model-value="onFileSelected" @error="onPickerError" />
		<button v-if="canRemove" type="button" class="remove-image-btn" @click="requestRemove">Remove image</button>

		<ErrorMsg v-if="errormsg" :msg="errormsg" />

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
