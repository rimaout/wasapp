<script>
import IconButton from './ui/IconButton.vue';
import MessageImage from './MessageImage.vue';

export default {
	components: { IconButton, MessageImage },
	props: {
		modelValue: { type: String, default: '' },
		images: { type: Array, default: () => [] },
		sending: { type: Boolean, default: false },
		chatId: { type: String, default: '' },
		replyTo: { type: Object, default: null },
	},
	emits: ['update:modelValue', 'update:images', 'send', 'error', 'clearReply'],
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
		replyPreview() {
			const r = this.replyTo;
			if (!r) return null;
			const fwd = r.forwardedFrom || {};
			const content = r.content || {};
			return {
				name: (r.sender && r.sender.name) || '',
				text: fwd.text || content.text || '',
				hasImage: !!(fwd.msgImageId || content.msgImageId),
				imageId: fwd.msgImageId || content.msgImageId,
			};
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
		<div v-if="replyPreview" class="reply-preview-bar">
			<MessageImage v-if="replyPreview.hasImage" :chat-id="chatId" :image-id="replyPreview.imageId" thumbnail />
			<div class="reply-info">
				<span class="reply-name">{{ replyPreview.name }}</span>
				<span class="reply-text">{{ replyPreview.text || 'Image' }}</span>
			</div>
			<button type="button" class="reply-remove" @click="$emit('clearReply')" aria-label="Remove reply" title="Remove reply">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
			</button>
		</div>

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

.reply-preview-bar {
	display: flex;
	align-items: center;
	gap: 10px;
	background: var(--tn-bg-highlight);
	border: 1px solid var(--tn-border);
	border-radius: 12px;
	padding: 6px 8px;
	margin-bottom: 8px;
}

.reply-info {
	display: flex;
	flex-direction: column;
	min-width: 0;
	flex: 1 1 auto;
}

.reply-name {
	font-size: 0.75rem;
	font-weight: 600;
	color: var(--tn-blue);
}

.reply-text {
	font-size: 0.85rem;
	color: var(--tn-fg);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.reply-remove {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 24px;
	height: 24px;
	flex-shrink: 0;
	border: none;
	border-radius: 50%;
	background: none;
	color: var(--tn-fg-dark);
	cursor: pointer;
}

.reply-remove:hover {
	background: var(--tn-overlay);
	color: var(--tn-red);
}

.reply-remove .feather {
	width: 16px;
	height: 16px;
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
