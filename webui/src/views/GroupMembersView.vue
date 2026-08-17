<script setup>
import { ref, computed } from 'vue';
import SearchBar from '../components/SearchBar.vue';
import UserAvatar from '../components/UserAvatar.vue';
import MemberChips from '../components/MemberChips.vue';
import { useUsers } from '../composables/useUsers.js';

const props = defineProps({
	members: { type: Array, default: () => [] },
});
const emit = defineEmits(['create', 'update:members']);

const query = ref('');

const { users, loading, errormsg, fetchUsers, filteredUsers } = useUsers();
fetchUsers();

const selectedIds = computed(() => new Set(props.members.map(m => m.id)));
const visibleUsers = computed(() => filteredUsers(query.value));

function isSelected(user) {
	return selectedIds.value.has(user.id);
}

function toggleUser(user) {
	let list = props.members.slice();
	let idx = list.findIndex(m => m.id === user.id);
	if (idx === -1) {
		list.push({ id: user.id, name: user.name });
	} else {
		list.splice(idx, 1);
	}
	emit('update:members', list);
}

function removeMember(user) {
	emit('update:members', props.members.filter(m => m.id !== user.id));
}

function create() {
	if (props.members.length === 0) return;
	emit('create');
}
</script>

<template>
	<div class="group-members-view">
		<SearchBar v-model="query" placeholder="Search users..." />

		<div v-if="members.length > 0" class="px-3 py-2">
			<MemberChips :members="members" @remove="removeMember" />
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="user-list flex-grow-1">
			<LoadingSpinner v-if="loading && users.length === 0" />

			<div v-if="!loading && visibleUsers.length === 0" class="text-muted text-center py-5">
				<p class="mb-2 fs-5">No users found</p>
			</div>

			<div v-if="visibleUsers.length > 0" class="list-group list-group-flush">
				<a
					v-for="user in visibleUsers"
					:key="user.id"
					class="list-group-item list-group-item-action d-flex align-items-center user-row"
					:class="{ selected: isSelected(user) }"
					href="#"
					@click.prevent="toggleUser(user)"
				>
					<UserAvatar :userId="user.id" :displayName="user.name" :size="48" class="me-3" />
					<span class="user-name fw-semibold text-truncate flex-grow-1">{{ user.name }}</span>
					<svg v-if="isSelected(user)" class="feather check-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
				</a>
			</div>
		</div>

		<div class="footer px-3 py-3">
			<button type="button" class="create-btn" :disabled="members.length === 0" @click="create">
				Continue
			</button>
		</div>
	</div>
</template>

<style scoped>
.group-members-view {
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

.user-row.selected {
	background-color: var(--tn-bg-highlight) !important;
}

.user-name {
	color: var(--tn-fg);
}

.check-icon {
	width: 20px;
	height: 20px;
	color: var(--tn-blue);
	flex-shrink: 0;
}

.footer {
	flex-shrink: 0;
}

.create-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 100%;
	height: 40px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	background: rgba(158, 206, 106, 0.15);
	color: var(--tn-green);
	font-weight: 600;
}

.create-btn:hover:not(:disabled) {
	background: rgba(158, 206, 106, 0.25);
}

.create-btn:disabled {
	background: rgba(255, 255, 255, 0.08);
	color: var(--tn-fg-dark);
	cursor: default;
}
</style>
