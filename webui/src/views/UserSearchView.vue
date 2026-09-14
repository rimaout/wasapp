<script>
import SearchBar from '../components/ui/SearchBar.vue';
import UserAvatar from '../components/ui/UserAvatar.vue';
import { useUsers } from '../composables/useUsers.js';
import { refreshChats } from '../composables/useChats.js';
import { navigateToChat } from '../services/chatNavigation.js';
import { getErrorMessage } from '../services/utils.js';

/**
 * UserSearchView component allows users to search for other users and initiate a chat with them.
 * When a user is selected, it attempts to create a new chat or navigate to an existing one.
 *
 * Used in: App.vue when the sidebar mode is set to 'users'.
 *
 * Props: None
 *
 * Emits: 'close' - emitted when the user search view should be closed
 */

export default {
	components: { SearchBar, UserAvatar },
	emits: ['close'],
	// Get shared users state and methods from useUsers
	setup() {
		const { users, loading, errormsg, fetchUsers, filteredUsers } = useUsers();
		return { users, loading, errormsg, fetchUsers, filteredUsers };
	},
	data() {
		return {
			query: '',
			creating: null,
		};
	},
	computed: {
		visibleUsers() {
			return this.filteredUsers(this.query);
		},
	},
	created() {
		this.fetchUsers();
	},
	methods: {
		async selectUser(user) {
			if (this.creating) return; // Prevent multiple simultaneous chat creation requests

			// Clean any previus state
			this.creating = user.id;
			this.errormsg = null;

			// Create a new chat with the selected user
			try {
				const response = await this.$axios.post('/users/' + user.id + '/chats');
				refreshChats().catch(() => {});
				this.goToChat(response.data.id, response.data.displayName, false); // "false" indicates that this is not a group chat
			} catch (e) {
				if (e.response && e.response.status === 409 && e.response.data && e.response.data.chatId) {
					// Conflict: a chat with this user already exists, navigates to that chat instead (no error message is shown)
					refreshChats().catch(() => {});
					this.goToChat(e.response.data.chatId, user.name, false);
				} else {
					// Other errors: show an error message
					this.errormsg = getErrorMessage(e);
				}
			} finally {
				this.creating = null;
			}
		},

		goToChat(chatId, displayName, isGroupChat) {
			navigateToChat(this.$router, { id: chatId, displayName, isGroupChat });
			this.$emit('close');
		},
	},
};
</script>

<template>
	<div class="user-search-view">
		<SearchBar v-model:search-text="query" placeholder="Search users..." />

		<div class="user-list fade-bottom flex-grow-1">
			<!-- Error message display -->
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

			<!-- Loading spinner display -->
			<LoadingSpinner v-if="loading && users.length === 0" />

			<!-- No users found message display -->
			<div v-if="!loading && visibleUsers.length === 0" class="text-muted text-center py-5">
				<p class="mb-2 fs-5">No users found</p>
				<p v-if="query">Try a different search</p>
			</div>

			<!-- List of users display -->
			<div v-if="visibleUsers.length > 0" class="list-group list-group-flush">
				<!-- User row for each visible user -->
				<a
					v-for="user in visibleUsers"
					:key="user.id"
					class="list-group-item list-group-item-action d-flex align-items-center user-row"
					:class="{ disabled: creating }"
					@click.prevent="selectUser(user)"
					href="#"
				>
				<!-- Note: by default, clicking on an <a> tag would navigate to the href. So we use
						   @click.prevent to prevent the default navigation behavior and instead call the selectUser method.

					The href is set to a default value of "#" to prevent navigation, but the actual navigation is handled by the selectUser() method when the chat is clicked.
				-->

					<!-- User avatar and name display -->
					<UserAvatar :userId="user.id" :displayName="user.name" :size="48" class="me-3" />

					<!-- User name display -->
					<span class="user-name fw-semibold text-truncate">{{ user.name }}</span>
				</a>
			</div>
		</div>
	</div>
</template>

<style scoped>
.user-search-view {
	display: flex;
	flex-direction: column;
	height: 100%;
}

.user-list {
	overflow-y: auto;
	scrollbar-width: none;
	-ms-overflow-style: none;
	padding-bottom: 48px;
}

.user-list::-webkit-scrollbar {
	display: none;
}

.user-row {
	border: 0 !important;
	border-radius: 12px !important;
	margin: 2px 12px;
	width: auto;
	padding: 8px 16px;
}

.user-row:hover {
	background-color: var(--theme-bg-highlight) !important;
}

.user-name {
	color: var(--theme-fg);
}
</style>
