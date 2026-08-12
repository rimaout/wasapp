<script>
import { formatTime, getStatusIcon, getStatusColor, isMyMessageById } from '../services/utils.js';

export default {
	data() {
		return {
			messages: [],
			newMsg: '',
			sending: false,
			loading: true,
			errormsg: null,
			intervalId: null,
			markedRead: false,
		};
	},
	computed: {
		chatId() {
			return this.$route.params.chatId;
		},
		isGroup() {
			return this.$route.query.group === '1';
		},
	},
	methods: {
		formatTime,
		getStatusIcon,
		getStatusColor,
		isMyMessageById,

		async fetchMessages(scroll) {
			try {
				let response = await this.$axios.get('/chats/' + this.chatId + '/messages');
				this.messages = response.data.messages.slice().reverse();

				if (!this.markedRead) {
					this.markedRead = true;
					this.$axios.post('/chats/' + this.chatId + '/read').catch(() => {});
				}
			} catch (e) {
				this.errormsg = e.response?.data?.message || e.toString();
			}
			this.loading = false;
			if (scroll) this.scrollToBottom();
		},

		async sendMessage() {
			let text = this.newMsg.trim();
			if (!text || this.sending) return;

			this.sending = true;
			try {
				let formData = new FormData();
				formData.append('text', text);

				let response = await this.$axios.post('/chats/' + this.chatId + '/messages', formData);
				this.messages.push(response.data);
				this.newMsg = '';
				this.scrollToBottom();
			} catch (e) {
				this.errormsg = e.response?.data?.message || e.toString();
			}
			this.sending = false;
		},

		scrollToBottom() {
			this.$nextTick(() => {
				let el = this.$refs.messagesArea;
				if (el) el.scrollTop = el.scrollHeight;
			});
		},

		initMessageText(msg) {
			let who = this.isMyMessageById(msg.sender.id) ? 'You' : msg.sender.name;
			return this.isGroup ? 'Group created by ' + who : who + ' started this chat';
		},
	},
	watch: {
		chatId() {
			this.messages = [];
			this.loading = true;
			this.markedRead = false;
			this.errormsg = null;
			this.fetchMessages(true);
		},
	},
	mounted() {
		this.fetchMessages(true);
		this.intervalId = setInterval(() => this.fetchMessages(false), 10000);
	},
	beforeUnmount() {
		if (this.intervalId) clearInterval(this.intervalId);
	},
};
</script>

<template>
	<div class="chat-view">
		<!-- Messages -->
		<div ref="messagesArea" class="chat-messages flex-grow-1">
			<LoadingSpinner v-if="loading" />
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

			<div v-for="msg in messages" :key="msg.id" class="message-wrapper" :class="isMyMessageById(msg.sender.id) ? 'me' : 'other'">
				<div v-if="msg.isInitMessage" class="message-init text-muted">{{ initMessageText(msg) }}</div>

				<div v-else-if="msg.isDeleted" class="message-deleted text-muted fst-italic">
					{{ isMyMessageById(msg.sender.id) ? 'You deleted this message' : 'Message deleted' }}
				</div>

				<div v-else class="message-bubble" :class="isMyMessageById(msg.sender.id) ? 'me' : 'other'">
					<div v-if="isGroup && !isMyMessageById(msg.sender.id)" class="message-sender">{{ msg.sender.name }}</div>
					<div class="message-text">{{ msg.content?.text || '' }}</div>
					<div class="message-meta">
						<span class="message-time">{{ formatTime(msg.sendTime) }}</span>
						<span v-if="isMyMessageById(msg.sender.id)" class="message-check" :class="getStatusColor(msg.status)">
							{{ getStatusIcon(msg.status) }}
						</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Input bar -->
		<div class="chat-input d-flex align-items-center px-3">
			<input
				type="text"
				class="form-control"
				placeholder="Type a message..."
				v-model="newMsg"
				@keyup.enter="sendMessage"
				:disabled="sending"
			/>
			<button class="send-btn ms-2" @click="sendMessage" :disabled="sending || !newMsg.trim()">
				<svg class="feather" style="width: 20px; height: 20px;"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
			</button>
		</div>
	</div>
</template>

<style scoped>
.chat-view {
	display: flex;
	flex-direction: column;
	height: calc(100vh - 48px);
}

.chat-messages {
	overflow-y: auto;
	padding: 16px;
	display: flex;
	flex-direction: column;
	gap: 6px;
	scrollbar-width: none;
	-ms-overflow-style: none;
}

.chat-messages::-webkit-scrollbar {
	display: none;
}

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

.chat-input {
	height: 56px;
	flex-shrink: 0;
	border-top: 1px solid var(--tn-border);
}

.send-btn {
	background: none;
	border: none;
	padding: 8px;
	color: var(--tn-blue);
	cursor: pointer;
	line-height: 1;
	flex-shrink: 0;
}

.send-btn:hover {
	color: var(--tn-cyan);
}

.send-btn:disabled {
	color: var(--tn-comment);
	cursor: default;
}
</style>
