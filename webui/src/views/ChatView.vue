<script>
import MessageBubble from '../components/MessageBubble.vue';
import MessageInput from '../components/MessageInput.vue';
import ChatAvatar from '../components/ui/ChatAvatar.vue';
import ActionMenu from '../components/ui/ActionMenu.vue';
import RenameActionScreen from '../components/RenameActionScreen.vue';
import ImageActionScreen from '../components/ImageActionScreen.vue';
import ChatMembersPopup from '../components/ChatMembersPopup.vue';
import { usePolling } from '../composables/usePolling.js';
import { refreshChats } from '../composables/useChats.js';
import { leaveGroup } from '../services/api.js';
import { formatDay, isNewDay, getErrorMessage } from '../services/utils.js';

export default {
	components: { MessageBubble, MessageInput, ChatAvatar, ActionMenu, ChatMembersPopup },
	data() {
		return {
			messages: [],
			newMsg: '',
			newImages: [],
			sending: false,
			loading: true,
			errormsg: null,
			nameOverride: null,
			avatarVersion: 0,
			replyTo: null,
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
			return this.nameOverride || this.$route.query.name || 'Chat';
		},
		chatOptions() {
			return [
				{ id: 'rename', label: 'Change Group Name', icon: 'type', component: RenameActionScreen, props: { title: 'Change Group Name', target: { kind: 'group', chatId: this.chatId }, initialName: this.chatName } },
				{ id: 'image', label: 'Change Group Image', icon: 'image', component: ImageActionScreen, props: { title: 'Change Group Image', target: { kind: 'group', chatId: this.chatId } } },
				{ id: 'leave', label: 'Leave Group', icon: 'log-out', danger: true, dangerText: 'Are you sure you want to leave this group?' },
			];
		},
	},
	methods: {
		formatDay,
		isNewDay,

		async fetchMessages(scroll) {
			try {
				let response = await this.$axios.get('/chats/' + this.chatId + '/messages');
				this.messages = response.data.messages.slice().reverse();

				this.$axios.post('/chats/' + this.chatId + '/read')
				.then(() => refreshChats())
				.catch(() => {});
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			}
			this.loading = false;
			if (scroll) this.scrollToBottom();
		},

		async sendMessage() {
			let text = this.newMsg.trim();
			if ((!text && this.newImages.length === 0) || this.sending) return;

			this.sending = true;
			try {
				let base = '/chats/' + this.chatId + '/messages';
				let url = this.replyTo ? base + '/' + this.replyTo.id + '/reply' : base;

				if (this.newImages.length === 0) {
					let formData = new FormData();
					formData.append('text', text);
					let response = await this.$axios.post(url, formData);
					this.messages.push(response.data);
					this.newMsg = '';
				} else {
					while (this.newImages.length > 0) {
						let image = this.newImages[0];
						let isLast = this.newImages.length === 1;
						let formData = new FormData();
						if (isLast && text) formData.append('text', text);
						formData.append('imageFile', image);

						let response = await this.$axios.post(url, formData);
						this.messages.push(response.data);
						this.newImages = this.newImages.slice(1);
						if (isLast) this.newMsg = '';
					}
				}
				this.replyTo = null;
				this.scrollToBottom();
				refreshChats().catch(() => {});
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			}
			this.sending = false;
		},

		scrollToBottom() {
			this.$nextTick(() => {
				let el = this.$refs.messagesArea;
				if (el) el.scrollTop = el.scrollHeight;
			});
		},

		onImageError(message) {
			this.errormsg = message;
		},

		onChatAction(payload) {
			if (payload.action === 'rename') {
				this.nameOverride = payload.name;
			} else if (payload.action === 'image') {
				this.avatarVersion++;
			}
			refreshChats().catch(() => {});
		},

		onChatSelect(item) {
			if (item.id === 'leave') {
				this.leaveGroup();
			}
		},

		onMessageAction(action) {
			if (action.type === 'reply') {
				this.replyTo = action.message;
			} else if (action.type === 'delete') {
				this.deleteMessage(action.message);
			} else if (action.type === 'react' && action.updatedMessage) {
				const idx = this.messages.findIndex(m => m.id === action.message.id);
				if (idx !== -1) this.messages.splice(idx, 1, action.updatedMessage);
			}
		},

		async deleteMessage(message) {
			try {
				const res = await this.$axios.delete('/chats/' + this.chatId + '/messages/' + message.id);
				const idx = this.messages.findIndex(m => m.id === message.id);
				if (idx !== -1) this.messages.splice(idx, 1, res.data);
				refreshChats().catch(() => {});
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			}
		},

		async leaveGroup() {
			try {
				await leaveGroup(this.chatId);
				refreshChats().catch(() => {});
				this.$router.push('/chats');
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			}
		},
	},
	watch: {
		chatId() {
			this.messages = [];
			this.loading = true;
			this.errormsg = null;
			this.nameOverride = null;
			this.replyTo = null;
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
			<button class="chat-back-btn d-md-none" @click="$router.push('/chats')" aria-label="Back" title="Back">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#arrow-left"/></svg>
			</button>
			<ChatAvatar :chatId="chatId" :displayName="chatName" :size="40" :isGroup="isGroup" :version="avatarVersion" />
			<span class="chat-header-name">{{ chatName }}</span>
			<ChatMembersPopup v-if="isGroup" :chatId="chatId" />
			<ActionMenu v-if="isGroup" class="ms-auto" trigger-button-icon="more-vertical" placement="down-right" :items="chatOptions" trigger-button-filled @done="onChatAction" @select="onChatSelect" />
		</div>

		<div ref="messagesArea" class="chat-messages flex-grow-1">
			<LoadingSpinner v-if="loading" />
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

			<template v-for="(msg, i) in messages" :key="msg.id">
				<template v-if="msg.isInitMessage">
					<MessageBubble :message="msg" :isGroup="isGroup" :showAvatar="isGroup" @action="onMessageAction" />
					<div class="date-divider">{{ formatDay(msg.sendTime) }}</div>
				</template>
				<template v-else-if="msg.isJoinMessage || msg.isLeaveMessage">
					<div v-if="i === 0 || isNewDay(messages[i - 1], msg)" class="date-divider">
						{{ formatDay(msg.sendTime) }}
					</div>
					<MessageBubble :message="msg" :isGroup="isGroup" @action="onMessageAction" />
				</template>
				<template v-else>
					<div v-if="i === 0 || isNewDay(messages[i - 1], msg)" class="date-divider">
						{{ formatDay(msg.sendTime) }}
					</div>
					<MessageBubble :message="msg" :isGroup="isGroup" :showAvatar="isGroup" :previous-message="messages[i - 1]" @action="onMessageAction" />
				</template>
			</template>
		</div>

		<MessageInput v-model="newMsg" v-model:images="newImages" :sending="sending" :chat-id="chatId" :reply-to="replyTo" @send="sendMessage" @error="onImageError" @clear-reply="replyTo = null" />
	</div>
</template>

<style scoped>
.chat-view {
	display: flex;
	flex-direction: column;
	height: 100vh;
	position: relative;
	background-color: var(--tn-bg-dark);
	background-image: url('/chat-bg.svg');
	background-repeat: repeat;
	background-size: 500px;
}

.chat-back-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	background: none;
	border: none;
	border-radius: 50%;
	width: 40px;
	height: 40px;
	color: var(--tn-fg-dark);
	cursor: pointer;
	flex-shrink: 0;
	margin-left: -8px;
}

.chat-back-btn:hover {
	background: rgba(255, 255, 255, 0.08);
}

.chat-header {
	display: flex;
	align-items: center;
	gap: 12px;
	height: var(--topbar-height);
	padding: 0 16px;
	background: var(--tn-bg-darker);
	flex-shrink: 0;
}

.chat-header-name {
	color: var(--tn-fg);
	font-size: 1.15rem;
	font-weight: 600;
	flex: 1;
	min-width: 0;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.chat-messages {
	overflow-y: auto;
	padding: 16px 16px 110px;
	display: flex;
	flex-direction: column;
	gap: 10px;
	scrollbar-width: none;
	-ms-overflow-style: none;
	-webkit-mask-image: linear-gradient(to bottom, black 0%, black calc(100% - 64px), transparent 100%);
	mask-image: linear-gradient(to bottom, black 0%, black calc(100% - 64px), transparent 100%);
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
