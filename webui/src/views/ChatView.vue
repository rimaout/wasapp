<script>
import MessageBubble from '../components/MessageBubble.vue';
import MessageInput from '../components/MessageInput.vue';
import { usePolling } from '../composables/usePolling.js';

export default {
	components: { MessageBubble, MessageInput },
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
	},
	methods: {
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
		<div ref="messagesArea" class="chat-messages flex-grow-1">
			<LoadingSpinner v-if="loading" />
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

			<MessageBubble v-for="msg in messages" :key="msg.id" :message="msg" :isGroup="isGroup" />
		</div>

		<MessageInput v-model="newMsg" :sending="sending" @send="sendMessage" />
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
</style>
