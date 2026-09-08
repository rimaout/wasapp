<script setup>
import { ref, computed } from 'vue';
import SearchBar from './ui/SearchBar.vue';
import UserAvatar from './ui/UserAvatar.vue';
import MemberChips from './ui/MemberChips.vue';
import { useUsers } from '../composables/useUsers.js';

// Search + select users, shown as removable chips. Excludes the logged-in user
// and any IDs passed in `excludeIds`.
const props = defineProps({
	modelValue:  { type: Array,   default: () => []          },
	excludeIds:  { type: Array,   default: () => []          },
	placeholder: { type: String,  default: 'Search users...' },
	compact:     { type: Boolean, default: false             },
	light:       { type: Boolean, default: false             },
});
const emit = defineEmits(['update:modelValue']);

const query = ref('');
const { users, loading, errormsg, fetchUsers, filteredUsers } = useUsers();
fetchUsers();

const selectedIds = computed(() => new Set(props.modelValue.map(m => m.id)));
const visibleUsers = computed(() => {
	let list = filteredUsers(query.value, props.excludeIds);
	return list.filter(u => !selectedIds.value.has(u.id));
});

function isSelected(user) {
	return selectedIds.value.has(user.id);
}

function toggleUser(user) {
	let list = props.modelValue.slice();
	let idx = list.findIndex(m => m.id === user.id);
	if (idx === -1) {
		list.push({ id: user.id, name: user.name });
	} else {
		list.splice(idx, 1);
	}
	emit('update:modelValue', list);
}

function removeMember(member) {
	emit('update:modelValue', props.modelValue.filter(m => m.id !== member.id));
}
</script>

<template>
	<div class="user-picker">
		<SearchBar v-model="query" :placeholder="placeholder" :compact="compact" :light="light" />

		<div v-if="modelValue.length > 0" class="px-3 py-2">
			<MemberChips :members="modelValue" :light="light" @remove="removeMember" />
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg" />

		<hr class="list-divider" :class="{ compact: compact }" />

		<div class="user-list fade-bottom">
			<LoadingSpinner v-if="loading && users.length === 0" />

			<div v-if="!loading && visibleUsers.length === 0" class="text-muted text-center py-5">
				<p class="mb-2 fs-5">No users found</p>
			</div>

			<div v-if="visibleUsers.length > 0" class="list-group list-group-flush">
				<a
					v-for="user in visibleUsers"
					:key="user.id"
					class="list-group-item list-group-item-action d-flex align-items-center user-row"
					:class="{ selected: isSelected(user), compact: compact }"
					href="#"
					@click.prevent="toggleUser(user)"
				>
					<UserAvatar :userId="user.id" :displayName="user.name" :size="compact ? 32 : 48" class="me-3" />
					<span class="user-name fw-semibold text-truncate flex-grow-1">{{ user.name }}</span>
					<svg v-if="isSelected(user)" class="feather check-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
				</a>
			</div>
		</div>
	</div>
</template>

<style scoped>
.user-picker {
	display: flex;
	flex-direction: column;
	min-height: 0;
	flex: 1 1 auto;
}

.list-divider {
	border: none;
	border-top: 1.5px solid var(--theme-border);
	margin: 6px 20px 4px;
}

.list-divider.compact {
	margin: 6px 12px 10px;
}

.user-list {
	flex: 1 1 auto;
	min-height: 0;
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
	background-color: transparent !important;
}

.user-row.compact {
	padding: 4px 12px;
}

.user-row:hover {
	background-color: var(--theme-overlay) !important;
}

.user-row.selected {
	background-color: var(--theme-overlay) !important;
}

.user-name {
	color: var(--theme-fg);
}

.check-icon {
	width: 20px;
	height: 20px;
	color: var(--theme-blue);
	flex-shrink: 0;
}
</style>
