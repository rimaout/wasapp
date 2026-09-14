<script>
import SearchBar from './ui/SearchBar.vue';
import UserAvatar from './ui/UserAvatar.vue';
import MemberChips from './ui/MemberChips.vue';
import { useUsers } from '../composables/useUsers.js';

/**
 * UserPicker component allows users to search and select multiple users from a list.
 * Selected users are displayed, and users can be added or removed from the selection.
 *
 * Props:
 *  - userList (Array): The currently selected users (array of objects with id and name).
 *  - excludeIds (Array): List of user IDs to exclude from the selection list.
 *  - placeholder (String): Placeholder text for the search bar.
 *  - compact (Boolean): Whether to use a compact layout for the user list.
 *  - light (Boolean): Whether to use a light theme for the component (useful for darker backgrounds).
 *
 * Emits:
 *  - 'update:userList': Emitted when the selected users are updated.
 */
export default {
	components: { SearchBar, UserAvatar, MemberChips },

	props: {
		userList:    { type: Array,   default: () => []          },
		excludeIds:  { type: Array,   default: () => []          },
		placeholder: { type: String,  default: 'Search users...' },
		compact:     { type: Boolean, default: false             },
		light:       { type: Boolean, default: false             },
	},

	emits: ['update:userList'],

	// Get shared users state and methods from useUsers
	setup() {
		const { users, loading, errormsg, fetchUsers, filteredUsers } = useUsers();
		return { users, loading, errormsg, fetchUsers, filteredUsers };
	},

	data() {
		return {
			query: '',
		};
	},

	computed: {
		// Ids of the currently selected users
		selectedIds() {
			return new Set(this.userList.map(m => m.id));
		},
		// List of all the non selected users (can be filtered by search query and excludeIds)
		visibleUsers() {
			let list = this.filteredUsers(this.query, this.excludeIds);
			return list.filter(u => !this.selectedIds.has(u.id));
		},
	},

	created() {
		this.fetchUsers();
	},

	methods: {

		isSelected(user) {
			return this.selectedIds.has(user.id);
		},


		// ----- ACTION HANDLERS -----

		// Toggle a user in the selected list, used to add or remove a user from the selection
		toggleUser(user) {
			let list = this.userList.slice();
			let idx = list.findIndex(m => m.id === user.id);

			// Toggle the user in the selected list
			if (idx === -1) {
				list.push({ id: user.id, name: user.name });
			} else {
				list.splice(idx, 1); // Note: 1 is the number of elements to remove
			}

			// Emit the updated list of selected users
			this.$emit('update:userList', list);
		},

		// Remove a user from the selected list, used when clicking the remove button on a chip
		removeMember(member) {
			this.$emit('update:userList', this.userList.filter(m => m.id !== member.id));
		},
	},
};
</script>

<template>
	<div class="user-picker">
		<!-- Search bar for filtering users -->
		<SearchBar v-model:search-text="query" :placeholder="placeholder" :compact="compact" :light="light" />

		<!-- Display selected users as chips -->
		<div v-if="userList.length > 0" class="px-3 py-2">
			<MemberChips :members="userList" :light="light" @remove="removeMember" />
		</div>

		<!-- Display error message if there is an error fetching users -->
		<ErrorMsg v-if="errormsg" :msg="errormsg" />

		<hr class="list-divider" :class="{ compact: compact }" />

		<!-- List of users that can be selected -->
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
					@click.prevent="toggleUser(user)"
					href="#"
				>
				<!-- Note: by default, clicking on an <a> tag would navigate to the href. So we use
						   @click.prevent to prevent the default navigation behavior and instead call the toggleUser() method.

					The href is set to a default value of "#" to prevent navigation, but the actual navigation is handled by the toggleUser() method when the chat is clicked.
				-->
					<!-- User avatar and name display -->
					<UserAvatar :userId="user.id" :displayName="user.name" :size="compact ? 32 : 48" class="me-3" />

					<!-- User name display -->
					<span class="user-name fw-semibold text-truncate flex-grow-1">{{ user.name }}</span>

					<!-- Check icon to indicate selected users -->
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
