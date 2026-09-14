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

/**
 * ChatView component displays the chat interface for a specific chat, allowing users to view messages, send new messages, and perform actions like replying or deleting messages.
 *
 * Used in: App.vue when the user navigates to a specific chat via the /chats/:chatId route.
 *
 * Props: None
 *
 * Emits: None
 */
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
		// Cheat option list used for the ActionMenu component in the chat header
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
				// Fetch messages for the current chat from the server
				let response = await this.$axios.get('/chats/' + this.chatId + '/messages');
				this.messages = response.data.messages.slice().reverse();

				// Mark the chat as read on the server and refresh the chat list
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

			// Prevent sending empty messages or sending while a message is already being sent
			if ((!text && this.newImages.length === 0) || this.sending) return;

			this.sending = true;
			try {
				// Determine the appropriate URL for sending the message, depending on whether it's a reply or a new message
				let base = '/chats/' + this.chatId + '/messages';
				let url = this.replyTo ? base + '/' + this.replyTo.id + '/replies' : base;

				if (this.newImages.length === 0) {
					// If there are no images, send a text message
					let formData = new FormData();
					formData.append('text', text);
					let response = await this.$axios.post(url, formData);
					this.messages.push(response.data);
					this.newMsg = '';
				} else {
					// If there are images, send each image as a separate message
					while (this.newImages.length > 0) {
						let image = this.newImages[0];
						let isLast = this.newImages.length === 1;
						let formData = new FormData();

						// If it's the last image and there's text, include the text in the last message
						if (isLast && text) formData.append('text', text);

						// Append the image file to the form data and send the message
						formData.append('imageFile', image);
						let response = await this.$axios.post(url, formData);

						// Add the sent message to the messages array
						this.messages.push(response.data);

						// Remove the sent image from the newImages array and clear the text if it was included in the last message
						this.newImages = this.newImages.slice(1);
						if (isLast) this.newMsg = '';
					}
				}

				// Clear the replyTo state and scroll to the bottom of the chat after sending the message
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
			// Reset the state when the chatId changes
			this.messages = [];
			this.loading = true;
			this.errormsg = null;
			this.nameOverride = null;
			this.replyTo = null;
			this.fetchMessages(true);
		},
	},
	mounted() {
		// Fetch messages when the component is mounted and start polling for new messages every 10 seconds
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

		<!-- CHAT HEADER -->
		<div class="chat-header">
			<!-- Back button for mobile view -->
			<button class="chat-back-btn d-md-none" @click="$router.push('/chats')" title="Back">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#arrow-left"/></svg>
			</button>

			<!-- Chat avatar and name -->
			<ChatAvatar :chatId="chatId" :displayName="chatName" :size="40" :isGroup="isGroup" :version="avatarVersion" />
			<span class="chat-header-name">{{ chatName }}</span>

			<!-- Members Menu: button to open popmenu for view and export members in group chat -->
			<ChatMembersPopup v-if="isGroup" :chatId="chatId" />

			<!-- Action Menu: button to open menu for editing group name, image or leave group -->
			<ActionMenu v-if="isGroup" class="ms-auto" trigger-button-icon="more-vertical" placement="down-right" :items="chatOptions" trigger-button-filled @done="onChatAction" @select="onChatSelect" />
		</div>

		<!-- CHAT MESSAGES AREA -->
		<div ref="messagesArea" class="chat-messages flex-grow-1">

			<!-- Loading spinner and error message -->
			<LoadingSpinner v-if="loading" />
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

			<!-- Render each message in the chat -->
			<template v-for="(msg, i) in messages" :key="msg.id">

				<!-- Initial message rendering: shows the first message in the chat with a date divider -->
				<template v-if="msg.isInitMessage">
					<MessageBubble :message="msg" :isGroup="isGroup" :showAvatar="isGroup" @action="onMessageAction" />
					<div class="date-divider">{{ formatDay(msg.sendTime) }}</div>
				</template>

				<!-- Join/Leave message rendering: shows a message when a user joins or leaves the chat, with a date divider if it's a new day -->
				<template v-else-if="msg.isJoinMessage || msg.isLeaveMessage">
					<div v-if="i === 0 || isNewDay(messages[i - 1], msg)" class="date-divider">
						{{ formatDay(msg.sendTime) }}
					</div>
					<MessageBubble :message="msg" :isGroup="isGroup" @action="onMessageAction" />
				</template>

				<!-- Regular message rendering: shows a standard message with a date divider if it's a new day -->
				<template v-else>
					<div v-if="i === 0 || isNewDay(messages[i - 1], msg)" class="date-divider">
						{{ formatDay(msg.sendTime) }}
					</div>
					<MessageBubble :message="msg" :isGroup="isGroup" :showAvatar="isGroup" :previous-message="messages[i - 1]" @action="onMessageAction" />
				</template>
			</template>
		</div>

		<!-- MESSAGE INPUT AREA -->
		<MessageInput v-model:message-text="newMsg" v-model:images="newImages" :sending="sending" :chat-id="chatId" :reply-to="replyTo" @send="sendMessage" @error="onImageError" @clear-reply="replyTo = null" />
	</div>
</template>

<style scoped>
.chat-view {
	display: flex;
	flex-direction: column;
	height: 100vh;
	height: 100dvh;
	position: relative;
	background-color: var(--theme-bg-dark);
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
	color: var(--theme-fg-dark);
	cursor: pointer;
	flex-shrink: 0;
	margin-left: -8px;
}

.chat-back-btn:hover {
	background: var(--theme-overlay);
}

.chat-header {
	display: flex;
	align-items: center;
	gap: 12px;
	height: var(--topbar-height);
	padding: 0 16px;
	background: var(--theme-bg-darker);
	flex-shrink: 0;
}

.chat-header-name {
	color: var(--theme-fg);
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
	color: var(--theme-fg-dark);
	background: var(--theme-bg-highlight);
	padding: 4px 12px;
	border-radius: 999px;
	margin: 8px 0;
}
</style>
