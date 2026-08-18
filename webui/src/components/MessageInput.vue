<script>
import IconButton from './IconButton.vue';

export default {
	components: { IconButton },
	props: {
		modelValue: { type: String, default: '' },
		images: { type: Array, default: () => [] },
		sending: { type: Boolean, default: false },
	},
	emits: ['update:modelValue', 'update:images', 'send', 'error'],
	data() {
		return {
			previewUrls: [],
		};
	},
	computed: {
		text: {
			get() { return this.modelValue; },
			set(v) { this.$emit('update:modelValue', v); },
		},
		canSend() {
			return (this.text.trim() || this.images.length > 0) && !this.sending;
		},
	},
	watch: {
		images(files) {
			this.previewUrls.forEach(u => URL.revokeObjectURL(u));
			this.previewUrls = (files || []).map(f => URL.createObjectURL(f));
		},
	},
	methods: {
		triggerFileInput() {
			this.$refs.fileInput.click();
		},
		onFileChange(e) {
			const files = Array.from(e.target.files || []);
			if (files.length === 0) return;
			for (const file of files) {
				if (file.size > 5 * 1024 * 1024) {
					this.$emit('error', 'Each image must be at most 5 MB');
					e.target.value = '';
					return;
				}
			}
			if (this.images.length + files.length > 10) {
				this.$emit('error', 'You can attach at most 10 images');
				e.target.value = '';
				return;
			}
			this.$emit('update:images', this.images.concat(files));
			e.target.value = '';
		},
		removeImage(idx) {
			this.$emit('update:images', this.images.filter((_, i) => i !== idx));
		},
		handleSend() {
			if (!this.canSend) return;
			this.$emit('send');
		},
	},
	beforeUnmount() {
		this.previewUrls.forEach(u => URL.revokeObjectURL(u));
	},
};
</script>

<template>
	<div class="chat-input">
		<div v-if="images.length > 0" class="image-preview-bar">
			<div v-for="(img, idx) in images" :key="idx" class="image-thumb">
				<img :src="previewUrls[idx]" alt="Attachment preview" />
				<button type="button" class="remove-image" @click="removeImage(idx)" aria-label="Remove image" title="Remove image">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
				</button>
			</div>
		</div>

		<div class="message-input-group">
			<IconButton class="me-3" icon="image" label="Attach image" @click="triggerFileInput" />
			<input
				type="text"
				class="message-field"
				placeholder="Type a message..."
				v-model="text"
				@keyup.enter="handleSend"
				:disabled="sending"
			/>
			<button v-if="canSend" class="send-btn" @click="handleSend" :disabled="sending">
				<svg class="feather" style="width: 20px; height: 20px;"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
			</button>
		</div>

		<input ref="fileInput" type="file" multiple accept="image/jpeg,image/png,image/webp" class="d-none" @change="onFileChange" />
	</div>
</template>

<style scoped>
.chat-input {
	position: absolute;
	left: 0;
	right: 0;
	bottom: 0;
	padding: 0 24px 14px;
}

.image-preview-bar {
	display: flex;
	gap: 8px;
	margin-bottom: 8px;
	overflow-x: auto;
	padding-bottom: 2px;
	scrollbar-width: none;
	-ms-overflow-style: none;
}

.image-preview-bar::-webkit-scrollbar {
	display: none;
}

.image-thumb {
	position: relative;
	width: 64px;
	height: 64px;
	border-radius: 10px;
	overflow: hidden;
	background: var(--tn-bg-highlight);
	box-shadow: 0 3px 16px rgba(0, 0, 0, 0.3);
}

.image-thumb img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.remove-image {
	position: absolute;
	top: 4px;
	right: 4px;
	width: 20px;
	height: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
	border: none;
	border-radius: 50%;
	background: rgba(31, 35, 53, 0.85);
	color: var(--tn-fg);
	cursor: pointer;
}

.remove-image .feather {
	width: 14px;
	height: 14px;
}

.message-input-group {
	display: flex;
	align-items: center;
	background: var(--tn-bg-highlight);
	border-radius: 999px;
	height: 52px;
	padding: 6px 6px 6px 8px;
	box-shadow: 0 3px 16px rgba(0, 0, 0, 0.3);
}

.message-field {
	flex: 1;
	background: none;
	border: none;
	outline: none;
	color: var(--tn-fg);
	font-size: 1rem;
	min-width: 0;
}

.message-field::placeholder {
	color: var(--tn-comment);
}

.send-btn {
	width: 38px;
	height: 38px;
	display: flex;
	align-items: center;
	justify-content: center;
	background-color: var(--tn-blue);
	border: none;
	border-radius: 50%;
	color: var(--tn-bg-darker);
	cursor: pointer;
	flex-shrink: 0;
}

.send-btn:hover {
	filter: brightness(1.1);
}

.send-btn:disabled {
	background-color: var(--tn-comment);
	color: var(--tn-bg-darker);
	cursor: default;
}
</style>
