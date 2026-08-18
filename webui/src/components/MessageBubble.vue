<script>
import { formatMessageTime, getStatusIcon, getStatusColor, isMyMessageById } from '../services/utils.js';
import UserAvatar from './UserAvatar.vue';
import MessageImage from './MessageImage.vue';

export default {
	components: { UserAvatar, MessageImage },
	props: {
		message: { type: Object, required: true },
		isGroup: { type: Boolean, default: false },
		showAvatar: { type: Boolean, default: false },
		previousMessage: { type: Object, default: null },
	},
	computed: {
		isMine() {
			return isMyMessageById(this.message.sender.id);
		},
		showSenderHeader() {
			if (!this.isGroup || this.isMine) return false;
			if (this.message.isInitMessage || this.message.isDeleted) return false;
			let prev = this.previousMessage;
			if (!prev) return true;
			if (prev.isInitMessage || prev.isDeleted) return true;
			if (prev.sender.id !== this.message.sender.id) return true;
			if (this.timeGapExceeded(prev.sendTime, this.message.sendTime)) return true;
			return false;
		},
		hasAvatar() {
			return this.showSenderHeader && this.showAvatar;
		},
		hasImage() {
			return !!this.message.content?.msgImageId;
		},
	},
	methods: {
		formatMessageTime,
		getStatusIcon,
		getStatusColor,

		timeGapExceeded(a, b) {
			let diff = Math.abs(new Date(b) - new Date(a));
			return diff > 5 * 60 * 1000;
		},

		initMessageText() {
			let who = this.isMine ? 'You' : this.message.sender.name;
			return this.isGroup ? 'Group created by ' + who : who + ' started this chat';
		},

		systemMessageText(action) {
			let who = this.isMine ? 'You' : this.message.sender.name;
			return who + ' ' + action;
		},
	},
};
</script>

<template>
	<div class="message-wrapper" :class="isMine ? 'me' : 'other'">
		<div v-if="message.isInitMessage" class="message-init text-muted">{{ initMessageText() }}</div>

		<div v-else-if="message.isJoinMessage" class="message-init text-muted">{{ systemMessageText('joined the group') }}</div>

		<div v-else-if="message.isLeaveMessage" class="message-init text-muted">{{ systemMessageText('left the group') }}</div>

		<div v-else-if="message.isDeleted" class="message-deleted text-muted fst-italic">
			{{ isMine ? 'You deleted this message' : 'Message deleted' }}
		</div>

		<div v-else class="message-bubble" :class="[isMine ? 'me' : 'other', { 'has-avatar': hasAvatar, 'has-image': hasImage }]">
			<div v-if="showSenderHeader" class="message-sender">
				<UserAvatar v-if="showAvatar" :userId="message.sender.id" :displayName="message.sender.name" :size="18" />
				<span class="message-sender-name">{{ message.sender.name }}</span>
			</div>

			<div v-if="hasImage" class="message-image-wrap">
				<MessageImage :chat-id="message.chatId" :image-id="message.content.msgImageId" />
				<div class="message-meta message-meta-overlay">
					<span class="message-time">{{ formatMessageTime(message.sendTime) }}</span>
					<span v-if="isMine" class="message-check" :class="getStatusColor(message.status)">
						{{ getStatusIcon(message.status) }}
					</span>
				</div>
			</div>

			<div v-if="message.content?.text" class="message-text" :class="{ 'message-caption': hasImage }">{{ message.content.text }}</div>

			<div v-if="!hasImage" class="message-meta">
				<span class="message-time">{{ formatMessageTime(message.sendTime) }}</span>
				<span v-if="isMine" class="message-check" :class="getStatusColor(message.status)">
					{{ getStatusIcon(message.status) }}
				</span>
			</div>
		</div>
	</div>
</template>

<style scoped>
.message-wrapper {
	display: flex;
	flex-direction: column;
	align-items: flex-start;
}

.message-wrapper.me {
	align-items: flex-end;
}

.message-init {
	font-size: 0.8rem;
	text-align: center;
	width: 100%;
	padding: 4px 0;
}

.message-deleted {
	font-size: 0.85rem;
	padding: 4px 0;
}

.message-bubble {
	max-width: 70%;
	padding: 8px 12px;
	border-radius: 14px;
	word-wrap: break-word;
}

.message-bubble.me {
	background-color: var(--tn-blue);
	color: var(--tn-bg-darker);
	border-bottom-right-radius: 4px;
}

.message-bubble.other {
	background-color: var(--tn-bg-highlight);
	color: var(--tn-fg);
	border-bottom-left-radius: 4px;
}

.message-bubble.has-avatar {
	padding-left: 8px;
}

.message-bubble.has-avatar .message-sender {
	margin-bottom: 6px;
}

.message-bubble.has-image {
	padding: 3px;
	overflow: hidden;
}

.message-bubble.has-image .message-sender {
	padding: 4px 6px 2px;
	margin-bottom: 0;
}

.message-image-wrap {
	position: relative;
	border-radius: 10px;
	overflow: hidden;
}

.message-caption {
	padding: 6px 8px 4px;
}

.message-meta-overlay {
	position: absolute;
	right: 6px;
	bottom: 6px;
	margin: 0;
	padding: 2px 7px;
	border-radius: 999px;
	background: rgba(0, 0, 0, 0.55);
}

.message-meta-overlay .message-time {
	color: #fff;
	opacity: 0.9;
}

.message-sender {
	display: flex;
	align-items: center;
	gap: 6px;
	margin-bottom: 2px;
}

.message-sender-name {
	font-size: 0.75rem;
	font-weight: 600;
	color: var(--tn-cyan);
}

.message-text {
	font-size: 0.95rem;
	white-space: pre-wrap;
}

.message-meta {
	display: flex;
	align-items: center;
	justify-content: flex-end;
	gap: 4px;
	margin-top: 2px;
}

.message-time {
	font-size: 0.7rem;
	opacity: 0.7;
}

.message-check {
	font-size: 0.75rem;
	font-weight: bold;
}
</style>
