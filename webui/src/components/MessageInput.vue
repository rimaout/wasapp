<script>
import IconButton from './ui/IconButton.vue';
import MessageImage from './MessageImage.vue';

/**
 * MessageInput component provides an input area for users to type messages, attach images, and send messages in a chat.
 *
 * Used in: ChatView.vue as the message input area at the bottom of the chat.
 *
 * Props:
 *  - messageText (String): The current text of the message being written.
 *  - images (Array): An array of image files attached to the message.
 *  - sending (Boolean): Indicates whether a message is currently being sent.
 *  - chatId (String): The ID of the current chat.
 *  - replyTo (Object): The message object that the user is replying to, if any. Object structure: { sender: { name }, content: { text, msgImageId }, forwardedFrom: { text, msgImageId } }.
 *
 * Emits:
 *  - 'update:messageText' - emitted when the message text is updated.
 *  - 'update:images' - emitted when the attached images are updated.
 *  - 'send' - emitted when the user sends a message.
 *  - 'error' - emitted when there is an error (e.g., file size limit exceeded).
 *  - 'clearReply' - emitted when the user clears the reply preview.
 */
export default {
	components: { IconButton, MessageImage },

	props: {
		messageText: { type: String,  default: ''       },
		images:      { type: Array,   default: () => [] },
		sending:     { type: Boolean, default: false    },
		chatId:      { type: String,  default: ''       },
		replyTo:     { type: Object,  default: null     },
	},

	emits: ['update:messageText', 'update:images', 'send', 'error', 'clearReply'],

	data() {
		return {
			previewUrls: [], // Array of temporary URLs for previewing attached images
		};
	},

	computed: {
		text: {
			get() { return this.messageText; },
			set(v) { this.$emit('update:messageText', v); },
		},

		canSend() {
			// The user can send a message if there is text or at least one image attached, and no message is currently being sent.
			return (this.text.trim() || this.images.length > 0) && !this.sending;
		},

		replyPreview() {
			const r = this.replyTo;
			if (!r) return null;

			// Determine the content to display in the reply preview. If the message is forwarded, use the forwarded content; otherwise, use the original content.
			const fwd = r.forwardedFrom || {};
			const content = r.content || {};

			// Return an object containing the sender's name, text, and image information for the reply preview.
			return {
				name: r.sender.name || '',
				text: fwd.text || content.text || '',
				hasImage: !!(fwd.msgImageId || content.msgImageId),
				imageId: fwd.msgImageId || content.msgImageId,
			};
		},
	},

	watch: {
		// Watch for changes in the images to update the preview URLs
		images(files) {
			this.previewUrls.forEach(u => URL.revokeObjectURL(u));				// Revoke previous object URLs to free memory
			this.previewUrls = (files || []).map(f => URL.createObjectURL(f));  // Create new object URLs for the current images for preview
		},
	},

	methods: {
		// Simulates a click on the hidden <input type="file"> element so the user can upload images using the custom image button instead of default browser styles.
		triggerFileInput() {
			this.$refs.fileInput.click();
		},
		// Handler for the file input change event - Validates selected files and emits an update event with the new images array.
		onFileChange(e) {
			const files = Array.from(e.target.files || []);

			// If no files were selected, exit
			if (files.length === 0) return;

			// Validate each selected file's size (max 5 MB) and the total number of images (max 10)
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

			// Emit the updated images array to the parent component
			this.$emit('update:images', this.images.concat(files));
			e.target.value = '';
		},
		// Handler for the remove image button click event - Removes the selected image from the images array and emits an update event.
		removeImage(idx) {
			this.$emit('update:images', this.images.filter((_, i) => i !== idx));
		},
		// Handler for the send button click event - Emits a 'send' event if the user can send a message.
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

		<!-- Reply Preview Bar: Displays the message being replied to -->
		<div v-if="replyPreview" class="reply-preview-bar">

			<!-- Thumbnail Image Displays -->
			<MessageImage v-if="replyPreview.hasImage" :chat-id="chatId" :image-id="replyPreview.imageId" thumbnail />

			<!-- Reply Information: Displays the sender's name and message text -->
			<div class="reply-info">
				<span class="reply-name">{{ replyPreview.name }}</span>
				<span class="reply-text">{{ replyPreview.text || 'Image' }}</span>
			</div>

			<!-- Remove Reply Button: Allows the user to clear the replys -->
			<button type="button" class="reply-remove" @click="$emit('clearReply')" title="Remove reply">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
			</button>
		</div>

		<!-- Image Preview Bar: Displays thumbnails of attached images -->
		<div v-if="images.length > 0" class="image-preview-bar">

			<!-- For each attached image, display a thumbnail preview and a remove button -->
			<div v-for="(img, idx) in images" :key="idx" class="image-thumb">
				<img :src="previewUrls[idx]" alt="Attachment preview" />
				<button type="button" class="remove-image" @click="removeImage(idx)" title="Remove image">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
				</button>
			</div>
		</div>

		<!-- Message Input Section: Contains the image upload button, text input field, and send button -->
		<div class="message-input-group">
			<IconButton class="me-3" icon="image" @click="triggerFileInput" />
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

		<!-- Hidden file input for image uploads, triggered by the custom image button -->
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
	background: var(--theme-bg-highlight);
	border: 1px solid var(--theme-border);
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
	color: var(--theme-blue);
}

.reply-text {
	font-size: 0.85rem;
	color: var(--theme-fg);
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
	color: var(--theme-fg-dark);
	cursor: pointer;
}

.reply-remove:hover {
	background: var(--theme-overlay);
	color: var(--theme-red);
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
	background: var(--theme-bg-highlight);
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
	color: var(--theme-fg);
	cursor: pointer;
}

.remove-image .feather {
	width: 14px;
	height: 14px;
}

.message-input-group {
	display: flex;
	align-items: center;
	background: var(--theme-bg-highlight);
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
	color: var(--theme-fg);
	font-size: 1rem;
	min-width: 0;
}

.message-field::placeholder {
	color: var(--theme-comment);
}

.send-btn {
	width: 38px;
	height: 38px;
	display: flex;
	align-items: center;
	justify-content: center;
	background-color: var(--theme-blue);
	border: none;
	border-radius: 50%;
	color: var(--theme-bg-darker);
	cursor: pointer;
	flex-shrink: 0;
}

.send-btn:hover {
	filter: brightness(1.1);
}

.send-btn:disabled {
	background-color: var(--theme-comment);
	color: var(--theme-bg-darker);
	cursor: default;
}
</style>
