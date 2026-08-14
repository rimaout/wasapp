<script>
import SearchBar from '../components/SearchBar.vue';
import UserAvatar from '../components/UserAvatar.vue';
import { getUserId } from '../services/auth.js';

export default {
	components: { SearchBar, UserAvatar },
	emits: ['close'],
	data() {
		return {
			query: '',
			users: [],
			errormsg: null,
			loading: false,
			creating: null,
		};
	},
	computed: {
		filteredUsers() {
			let selfId = getUserId();
			let list = this.users.filter(u => u.id !== selfId);
			if (!this.query) return list;
			let q = this.query.toLowerCase();
			return list.filter(u => u.name.toLowerCase().includes(q));
		},
	},
	methods: {
		async fetchUsers() {
			this.errormsg = null;
			try {
				let response = await this.$axios.get('/users');
				this.users = response.data.usersList;
			} catch (e) {
				this.errormsg = e.response?.data?.message || e.toString();
			}
			this.loading = false;
		},

		async selectUser(user) {
			if (this.creating) return;
			this.creating = user.id;
			this.errormsg = null;
			try {
				let response = await this.$axios.post('/user/' + user.id + '/chats');
				this.openChat(response.data.id, response.data.displayName, false);
			} catch (e) {
				if (e.response && e.response.status === 409 && e.response.data && e.response.data.chatId) {
					this.openChat(e.response.data.chatId, user.name, false);
				} else {
					this.errormsg = e.response?.data?.message || e.toString();
				}
			} finally {
				this.creating = null;
			}
		},

		openChat(chatId, displayName, isGroup) {
			this.$router.push('/chats/' + chatId +
				'?name=' + encodeURIComponent(displayName) +
				'&group=' + (isGroup ? '1' : '0'));
			let sidebar = document.getElementById('sidebarMenu');
			if (sidebar && window.innerWidth < 768) {
				let bsCollapse = bootstrap.Collapse.getOrCreateInstance(sidebar);
				bsCollapse.hide();
			}
			this.$emit('close');
		},
	},
	mounted() {
		this.loading = true;
		this.fetchUsers();
	},
};
</script>

<template>
	<div>
		<div class="d-flex align-items-center px-3 pt-3 pb-2 user-search-header">
			<button type="button" class="back-btn" @click="$emit('close')">
				<svg class="feather back-icon"><use href="/feather-sprite-v4.29.0.svg#arrow-left"/></svg>
			</button>
			<span class="fw-semibold user-search-title">New chat</span>
		</div>

		<SearchBar v-model="query" placeholder="Search users..." />

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<LoadingSpinner v-if="loading && users.length === 0" />

		<div v-if="!loading && query && filteredUsers.length === 0" class="text-muted text-center py-5">
			<p class="mb-2 fs-5">No users found</p>
			<p>Try a different search</p>
		</div>

		<div v-if="!loading && !query && filteredUsers.length === 0" class="text-muted text-center py-5">
			<p class="mb-2 fs-5">No users found</p>
		</div>

		<div v-if="filteredUsers.length > 0" class="list-group list-group-flush">
			<a
				v-for="user in filteredUsers"
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
