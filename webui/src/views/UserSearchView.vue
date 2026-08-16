<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import axios from '../services/axios.js';
import SearchBar from '../components/SearchBar.vue';
import UserAvatar from '../components/UserAvatar.vue';
import { useUsers } from '../composables/useUsers.js';
import { refreshChats } from '../composables/useChats.js';
import { navigateToChat } from '../services/chatNavigation.js';
import { getErrorMessage } from '../services/utils.js';

const emit = defineEmits(['close']);
const router = useRouter();

const query = ref('');
const creating = ref(null);

const { users, loading, errormsg, fetchUsers, filteredUsers } = useUsers();
fetchUsers();

const visibleUsers = computed(() => filteredUsers(query.value));

async function selectUser(user) {
	if (creating.value) return;
	creating.value = user.id;
	errormsg.value = null;
	try {
		const response = await axios.post('/user/' + user.id + '/chats');
		refreshChats().catch(() => {});
		goToChat(response.data.id, response.data.displayName, false);
	} catch (e) {
		if (e.response && e.response.status === 409 && e.response.data && e.response.data.chatId) {
			refreshChats().catch(() => {});
			goToChat(e.response.data.chatId, user.name, false);
		} else {
			errormsg.value = getErrorMessage(e);
		}
	} finally {
		creating.value = null;
	}
}

function goToChat(chatId, displayName, isGroupChat) {
	navigateToChat(router, { id: chatId, displayName, isGroupChat });
	emit('close');
}
</script>

<template>
	<div class="user-search-view">
		<SearchBar v-model="query" placeholder="Search users..." />

		<div class="user-list flex-grow-1">
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
	background-color: var(--tn-bg-highlight) !important;
}

.user-name {
	color: var(--tn-fg);
}
</style>
