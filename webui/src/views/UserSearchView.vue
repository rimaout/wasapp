<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import axios from '../services/axios.js';
import SearchBar from '../components/SearchBar.vue';
import UserAvatar from '../components/UserAvatar.vue';
import { useUsers } from '../composables/useUsers.js';

const emit = defineEmits(['close', 'create-group']);
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
		openChat(response.data.id, response.data.displayName, false);
	} catch (e) {
		if (e.response && e.response.status === 409 && e.response.data && e.response.data.chatId) {
			openChat(e.response.data.chatId, user.name, false);
		} else {
			errormsg.value = e.response?.data?.message || e.toString();
		}
	} finally {
		creating.value = null;
	}
}

function openChat(chatId, displayName, isGroup) {
	router.push('/chats/' + chatId +
		'?name=' + encodeURIComponent(displayName) +
		'&group=' + (isGroup ? '1' : '0'));
	let sidebar = document.getElementById('sidebarMenu');
	if (sidebar && window.innerWidth < 768) {
		let bsCollapse = bootstrap.Collapse.getOrCreateInstance(sidebar);
		bsCollapse.hide();
	}
	emit('close');
}
</script>

<template>
	<div>
		<div class="d-flex align-items-center px-3 pt-3 pb-2 user-search-header">
			<button type="button" class="back-btn" @click="emit('close')">
				<svg class="feather back-icon"><use href="/feather-sprite-v4.29.0.svg#arrow-left"/></svg>
			</button>
			<span class="fw-semibold user-search-title">New chat</span>
		</div>

		<SearchBar v-model="query" placeholder="Search users..." />

		<button type="button" class="new-group-btn" @click="emit('create-group')">
			<svg class="feather new-group-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
			<span>Create new group</span>
		</button>

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
				class="list-group-item list-group-item-action d-flex align-items-center px-3 py-2 user-row"
				:class="{ disabled: creating }"
				href="#"
				@click.prevent="selectUser(user)"
			>
				<UserAvatar :userId="user.id" :displayName="user.name" :size="48" class="me-3" />
				<span class="user-name fw-semibold text-truncate">{{ user.name }}</span>
			</a>
		</div>
	</div>
</template>

<style scoped>
.user-search-header {
	color: var(--tn-fg);
}

.user-search-title {
	font-size: 1.05rem;
}

.back-btn {
	background: none;
	border: none;
	padding: 0;
	margin-right: 12px;
	cursor: pointer;
	color: var(--tn-fg-dark);
	line-height: 1;
	display: flex;
	align-items: center;
}

.back-btn:hover {
	color: var(--tn-blue);
}

.back-btn .back-icon {
	width: 22px;
	height: 22px;
}

.new-group-btn {
	display: flex;
	align-items: center;
	gap: 10px;
	width: 100%;
	border: none;
	background: none;
	padding: 10px 16px;
	color: var(--tn-fg);
	cursor: pointer;
	text-align: left;
	border-bottom: 1px solid var(--tn-border);
}

.new-group-btn:hover {
	background-color: var(--tn-bg-highlight);
}

.new-group-btn .new-group-icon {
	width: 20px;
	height: 20px;
	color: var(--tn-blue);
}

.user-row {
	border-color: var(--tn-border) !important;
}

.user-row:hover {
	background-color: var(--tn-bg-highlight) !important;
}

.user-name {
	color: var(--tn-fg);
}
</style>
