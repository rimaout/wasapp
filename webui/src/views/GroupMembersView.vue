<script setup>
import { ref, computed } from 'vue';
import SearchBar from '../components/SearchBar.vue';
import UserAvatar from '../components/UserAvatar.vue';
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

		<div v-if="members.length > 0" class="member-chips px-3 py-2">
			<span v-for="m in members" :key="m.id" class="member-chip">
				<UserAvatar :userId="m.id" :displayName="m.name" :size="24" />
				<span class="member-chip-name">{{ m.name }}</span>
				<button type="button" class="member-chip-remove" @click="removeMember(m)">×</button>
			</span>
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
					class="list-group-item list-group-item-action d-flex align-items-center px-3 py-2 user-row"
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
			<button type="button" class="btn btn-primary w-100 create-btn" :disabled="members.length === 0" @click="create">
				Create
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

.member-chips {
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
	flex-shrink: 0;
}

.member-chip {
	display: inline-flex;
	align-items: center;
	gap: 6px;
	background: var(--tn-bg-highlight);
	border-radius: 999px;
	padding: 3px 8px 3px 3px;
	color: var(--tn-fg);
}

.member-chip-name {
	font-size: 0.85rem;
	max-width: 120px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.member-chip-remove {
	background: none;
	border: none;
	padding: 0;
	cursor: pointer;
	color: var(--tn-fg-dark);
	line-height: 1;
	font-size: 1rem;
}

.member-chip-remove:hover {
	color: var(--tn-red);
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
	font-weight: 600;
}
</style>
