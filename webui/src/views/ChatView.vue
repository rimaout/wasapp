<script>
import MessageBubble from '../components/MessageBubble.vue';
import MessageInput from '../components/MessageInput.vue';
import ChatAvatar from '../components/ChatAvatar.vue';
import { usePolling } from '../composables/usePolling.js';
import { formatDay, isNewDay } from '../services/utils.js';

export default {
	components: { MessageBubble, MessageInput, ChatAvatar },
	data() {
		return {
			messages: [],
			newMsg: '',
			sending: false,
			loading: true,
			errormsg: null,
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
		chatName() {
			return this.$route.query.name || 'Chat';
		},
	},
	methods: {
		formatDay,
		isNewDay,

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
		this.stopPolling = usePolling(() => this.fetchMessages(false), 10000);
	},
	beforeUnmount() {
		if (this.stopPolling) this.stopPolling();
	},
};
</script>

<template>
	<div class="chat-view">
		<div class="chat-header">
			<ChatAvatar :chatId="chatId" :displayName="chatName" :size="40" :isGroup="isGroup" />
			<span class="chat-header-name">{{ chatName }}</span>
		</div>

		<div ref="messagesArea" class="chat-messages flex-grow-1">
			<LoadingSpinner v-if="loading" />
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

			<template v-for="(msg, i) in messages" :key="msg.id">
				<div v-if="i === 0 || isNewDay(messages[i - 1], msg)" class="date-divider">
					{{ formatDay(msg.sendTime) }}
				</div>
				<MessageBubble :message="msg" :isGroup="isGroup" />
			</template>
		</div>

		<MessageInput v-model="newMsg" :sending="sending" @send="sendMessage" />
	</div>
</template>

<style scoped>
.chat-view {
	display: flex;
	flex-direction: column;
	height: 100vh;
}

.chat-header {
	display: flex;
	align-items: center;
	gap: 12px;
	height: 56px;
	padding: 0 16px;
	border-bottom: 1px solid var(--tn-border);
	background: var(--tn-bg-darker);
	flex-shrink: 0;
}

.chat-header-name {
	color: var(--tn-fg);
	font-size: 1.15rem;
	font-weight: 600;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.chat-messages {
	overflow-y: auto;
	padding: 16px 16px 0;
	display: flex;
	flex-direction: column;
	gap: 6px;
	scrollbar-width: none;
	-ms-overflow-style: none;
}

.chat-messages::-webkit-scrollbar {
	display: none;
}

.date-divider {
	align-self: center;
	font-size: 0.75rem;
	color: var(--tn-fg-dark);
	background: var(--tn-bg-highlight);
	padding: 4px 12px;
	border-radius: 999px;
	margin: 8px 0;
}
</style>
