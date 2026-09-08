<script>
import SearchBar from '../components/ui/SearchBar.vue';
import UserAvatar from '../components/ui/UserAvatar.vue';
import { useUsers } from '../composables/useUsers.js';
import { refreshChats } from '../composables/useChats.js';
import { navigateToChat } from '../services/chatNavigation.js';
import { getErrorMessage } from '../services/utils.js';

export default {
	components: { SearchBar, UserAvatar },
	emits: ['close'],

	// Shared users state comes from the useUsers composable (module-level refs);
	// expose it here so the rest of the Options API can use it via `this`.
	setup(props, { emit }) {
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
			if (this.creating) return;
			this.creating = user.id;
			this.errormsg = null;
			try {
				const response = await this.$axios.post('/users/' + user.id + '/chats');
				refreshChats().catch(() => {});
				this.goToChat(response.data.id, response.data.displayName, false);
			} catch (e) {
				if (e.response && e.response.status === 409 && e.response.data && e.response.data.chatId) {
					refreshChats().catch(() => {});
					this.goToChat(e.response.data.chatId, user.name, false);
				} else {
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
		<SearchBar v-model="query" placeholder="Search users..." />

		<div class="user-list fade-bottom flex-grow-1">
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
			<LoadingSpinner v-if="loading && users.length === 0" />

			<div v-if="!loading && query && visibleUsers.length === 0" class="text-muted text-center py-5">
				<p class="mb-2 fs-5">No users found</p>
				<p>Try a different search</p>
			</div>

			<div v-if="!loading && !query && visibleUsers.length === 0" class="text-muted text-center py-5">
				<p class="mb-2 fs-5">No users found</p>
			</div>

			<div v-if="visibleUsers.length > 0" class="list-group list-group-flush">
				<a
					v-for="user in visibleUsers"
					:key="user.id"
					class="list-group-item list-group-item-action d-flex align-items-center user-row"
					:class="{ disabled: creating }"
					href="#"
					@click.prevent="selectUser(user)"
				>
					<UserAvatar :userId="user.id" :displayName="user.name" :size="48" class="me-3" />
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
