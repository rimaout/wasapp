<script>
import ChatAvatar from '../components/ui/ChatAvatar.vue';
import MessageStatusIcon from '../components/MessageStatusIcon.vue';
import { isMyMessage, formatPreviewTime, getChatSnippet } from '../services/utils.js';
import { usePolling } from '../composables/usePolling.js';
import { useChats } from '../composables/useChats.js';
import { navigateToChat } from '../services/chatNavigation.js';

/**
 * ChatsListView component displays a list of chat conversations, allowing users to search, view, and navigate to individual chats.
 *
 * Props:
 *  - searchQuery (String): The current search query used to filter the chat list.
 *  - beforeUnmount: Stops polling when the component is destroyed.
 */
export default {
	components: { ChatAvatar, MessageStatusIcon },
	props: {
		searchQuery: { type: String, default: '' },
	},
	// Get shared chats state and methods from useChats
	setup() {
		const { chats, errormsg, version, refreshChats } = useChats();
		return { chats, errormsg, version, refreshChats };
	},
	data() {
		return {
			loading: false,
		}
	},
	computed: {
		filteredChats() {
			if (!this.searchQuery) return this.chats;
			let q = this.searchQuery.toLowerCase();
			return this.chats.filter(c => c.displayName.toLowerCase().includes(q));
		},
	},
	methods: {
		isMyMessage,
		formatPreviewTime,
		getChatSnippet,

		openChat(chat) {
			navigateToChat(this.$router, chat);
		},

		isActive(chatId) {
			return this.$route && this.$route.params.chatId === chatId;
		},
	},
	mounted() {
		this.loading = true;
		this.refreshChats().finally(() => this.loading = false);
		this.stopPolling = usePolling(() => this.refreshChats(), 8000);
	},
	beforeUnmount() {
		if (this.stopPolling) this.stopPolling();
	},
};
</script>

<template>
	<div>
		<!-- Error message display -->
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<!-- Loading spinner display when loading and no chats are present -->
		<LoadingSpinner v-if="loading && chats.length === 0" />

		<!-- Display message when no search results are found after a search-->
		<div v-if="!loading && searchQuery && filteredChats.length === 0 && chats.length > 0" class="text-muted text-center py-5">
			<p class="mb-2 fs-5">No results found</p>
			<p>Try a different search or start a new chat</p>
		</div>

		<!-- Display message when there are no conversations yet -->
		<div v-if="!loading && !searchQuery && chats.length === 0" class="text-muted text-center py-5">
			<p class="mb-2 fs-5">No conversations yet</p>
			<p>Start a new chat to begin messaging</p>
		</div>

		<!-- List of filtered chats -->
		<div v-if="filteredChats.length > 0" class="list-group list-group-flush">
			<a
				v-for="chat in filteredChats"
				:key="chat.id"
				class="list-group-item list-group-item-action d-flex align-items-center chat-row"
				:class="{ active: isActive(chat.id) }"
				@click.prevent="openChat(chat)"
				href="#"
			>
			<!-- Note: by default, clicking on an <a> tag would navigate to the href. So we use
						   @click.prevent to prevent the default navigation behavior and instead call the openChat() method.

				 The href is set to a default value of "#" to prevent navigation, but the actual navigation is handled by the openChat() method when the chat is clicked.
				-->

				<!-- ChatAvatar Display -->
				<ChatAvatar :chatId="chat.id" :displayName="chat.displayName" :size="48" :isGroup="chat.isGroupChat" :version="version" class="me-3" />

				<div class="chat-info flex-grow-1 min-w-0">
					<div class="d-flex justify-content-between align-items-baseline">

						<!-- Chat Name Display -->
						<span class="chat-name fw-semibold text-truncate">
							{{ chat.displayName }}
						</span>

						<div class="d-flex align-items-center flex-shrink-0 ms-2">
							<!-- Message Status Icon Display - for user's own messages (status icon is for sent/received/read badge) -->
							<span v-if="isMyMessage(chat.lastMessage.senderName)" class="me-1 chat-check">
								<MessageStatusIcon :status="chat.lastMessage.status" />
							</span>

							<!-- Last Message Time Display -->
							<small class="chat-time">{{ formatPreviewTime(chat.lastMessage.sendTime) }}</small>
						</div>
					</div>

					<!-- Last Message Snippet Display -->
					<div class="d-flex justify-content-between align-items-center">
						<small class="chat-snippet text-truncate" :class="{ 'text-muted fst-italic': chat.lastMessage.isDeleted }">
							<!-- Display sender name for group chats, unless the last message is an init/join/leave message -->
							<span v-if="chat.isGroupChat && !chat.lastMessage.isInitMessage && !chat.lastMessage.isJoinMessage && !chat.lastMessage.isLeaveMessage" class="text-muted">
								{{ isMyMessage(chat.lastMessage.senderName) ? 'You' : chat.lastMessage.senderName }}:
							</span>
							<!-- Display the last message snippet -->
							{{ getChatSnippet(chat) }}
						</small>
						<!-- Unread message count badge display -->
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
	border: 0 !important;
	border-radius: 12px !important;
	margin: 4px 12px !important;
	width: auto;
	padding: 10px 16px;
}

.chat-row:hover {
	background-color: var(--theme-bg-highlight) !important;
}

.chat-row.active {
	background-color: var(--theme-bg-highlight) !important;
}

.chat-name {
	color: var(--theme-fg);
}

.chat-snippet {
	color: var(--theme-fg-dark);
	max-width: 100%;
}

.chat-time {
	color: var(--theme-comment);
	white-space: nowrap;
}

.chat-check {
	display: inline-flex;
	color: var(--theme-fg-dark);
}

.chat-check .is-received {
	opacity: 0.7;
}

.chat-check .is-read {
	color: var(--theme-cyan);
}

.unread-badge {
	font-size: 0.7rem;
	padding: 0.25em 0.55em;
}

.min-w-0 {
	min-width: 0;
}
</style>
