<script>
import ChatAvatar from '../components/ChatAvatar.vue';
import { isMyMessage, formatTime, getStatusIcon, getStatusColor, getChatSnippet } from '../services/utils.js';

export default {
	components: { ChatAvatar },
	data() {
		return {
			chats: [],
			errormsg: null,
			loading: false,
			intervalId: null,
			pollVersion: 0,
		}
	},
	methods: {
		isMyMessage,
		formatTime,
		getStatusIcon,
		getStatusColor,
		getChatSnippet,

		async fetchChats() {
			this.errormsg = null;
			try {
				let response = await this.$axios.get('/chats');
				this.chats = response.data.chatsPreviewList;
				this.pollVersion++;
			} catch (e) {
				this.errormsg = e.response?.data?.message || e.toString();
			}
			this.loading = false;
		},

		openChat(chatId) {
			this.$router.push('/chats/' + chatId);
			let sidebar = document.getElementById('sidebarMenu');
			if (sidebar && window.innerWidth < 768) {
				let bsCollapse = bootstrap.Collapse.getOrCreateInstance(sidebar);
				bsCollapse.hide();
			}
		},

		isActive(chatId) {
			return this.$route && this.$route.params.chatId === chatId;
		},
	},
	mounted() {
		this.loading = true;
		this.fetchChats();
		this.intervalId = setInterval(() => this.fetchChats(), 10000);
	},
	beforeUnmount() {
		if (this.intervalId) clearInterval(this.intervalId);
	},
};
</script>

<template>
	<div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<LoadingSpinner v-if="loading && chats.length === 0" />

		<div v-if="!loading && chats.length === 0" class="text-muted text-center py-5">
			<p class="mb-2 fs-5">No conversations yet</p>
			<p>Start a new chat to begin messaging</p>
		</div>

		<div v-if="chats.length > 0" class="list-group list-group-flush">
			<a
				v-for="chat in chats"
				:key="chat.id"
				class="list-group-item list-group-item-action d-flex align-items-center px-3 py-3 chat-row"
				:class="{ active: isActive(chat.id) }"
				@click.prevent="openChat(chat.id)"
				href="#"
			>
				<ChatAvatar :chatId="chat.id" :displayName="chat.displayName" :size="48" :isGroup="chat.isGroupChat" :version="pollVersion" class="me-3" />

				<div class="chat-info flex-grow-1 min-w-0">
					<div class="d-flex justify-content-between align-items-baseline">
						<span class="chat-name fw-semibold text-truncate">
							{{ chat.displayName }}
						</span>
						<div class="d-flex align-items-center flex-shrink-0 ms-2">
							<span
								v-if="isMyMessage(chat.lastMessage.senderName)"
								:class="getStatusColor(chat.lastMessage.status)"
								class="me-1 chat-check"
							>{{ getStatusIcon(chat.lastMessage.status) }}</span>
							<small class="chat-time">{{ formatTime(chat.lastMessage.sendTime) }}</small>
						</div>
					</div>
					<div class="d-flex justify-content-between align-items-center">
						<small
							class="chat-snippet text-truncate"
							:class="{ 'text-muted fst-italic': chat.lastMessage.isDeleted }"
						>
							<span v-if="chat.isGroupChat && !chat.lastMessage.isInitMessage" class="text-muted">
								{{ chat.lastMessage.senderName }}:
							</span>
							{{ getChatSnippet(chat) }}
						</small>
						<span
							v-if="chat.unreadCount > 0"
							class="badge rounded-pill bg-primary unread-badge flex-shrink-0 ms-2"
						>{{ chat.unreadCount }}</span>
					</div>
				</div>
			</a>
		</div>
	</div>
</template>

<style scoped>
.chat-row {
	border-color: var(--tn-border) !important;
}

.chat-row:hover {
	background-color: var(--tn-bg-highlight) !important;
}

.chat-row.active {
	background-color: var(--tn-bg-highlight) !important;
	border-color: var(--tn-border) !important;
}

.chat-name {
	color: var(--tn-fg);
}

.chat-snippet {
	color: var(--tn-fg-dark);
	max-width: 100%;
}

.chat-time {
	color: var(--tn-comment);
	white-space: nowrap;
}

.chat-check {
	font-size: 0.8rem;
	font-weight: bold;
}

.unread-badge {
	font-size: 0.7rem;
	padding: 0.25em 0.55em;
}

.min-w-0 {
	min-width: 0;
}
</style>
