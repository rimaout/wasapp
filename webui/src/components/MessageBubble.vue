<script>
import { formatTime, getStatusIcon, getStatusColor, isMyMessageById } from '../services/utils.js';

export default {
	props: {
		message: { type: Object, required: true },
		isGroup: { type: Boolean, default: false },
	},
	computed: {
		isMine() {
			return isMyMessageById(this.message.sender.id);
		},
	},
	methods: {
		formatTime,
		getStatusIcon,
		getStatusColor,

		initMessageText() {
			let who = this.isMine ? 'You' : this.message.sender.name;
			return this.isGroup ? 'Group created by ' + who : who + ' started this chat';
		},
	},
};
</script>

<template>
	<div class="message-wrapper" :class="isMine ? 'me' : 'other'">
		<div v-if="message.isInitMessage" class="message-init text-muted">{{ initMessageText() }}</div>

		<div v-else-if="message.isDeleted" class="message-deleted text-muted fst-italic">
			{{ isMine ? 'You deleted this message' : 'Message deleted' }}
		</div>

		<div v-else class="message-bubble" :class="isMine ? 'me' : 'other'">
			<div v-if="isGroup && !isMine" class="message-sender">{{ message.sender.name }}</div>
			<div class="message-text">{{ message.content?.text || '' }}</div>
			<div class="message-meta">
				<span class="message-time">{{ formatTime(message.sendTime) }}</span>
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

.message-sender {
	font-size: 0.75rem;
	font-weight: 600;
	color: var(--tn-cyan);
	margin-bottom: 2px;
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
