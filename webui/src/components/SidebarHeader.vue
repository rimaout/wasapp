<script>
import ActionMenu from './ui/ActionMenu.vue';

const TITLES = {
	users: 'New direct chat',
	'group-members': 'Add group members',
	'group-details': 'New group details',
};

/**
 * SidebarHeader component displays the header of the sidebar, which can change based on the current mode of the sidebar.
 * It shows different titles and actions depending on whether the user is in chats, users, group-members, or group-details mode.
 *
 * Used in: App.vue as the header of the sidebar.
 *
 * Props:
 *  - mode (String): The current mode of the sidebar ('chats', 'users', 'group-members', 'group-details').
 *
 * Emits:
 *  - 'back' - emitted when the back button is clicked (for user, group-members, or group-details modes)
 *  - 'home' - emitted when the home link is clicked (for chats mode).
 *  - 'create-direct' - emitted when the user selects to create a direct chat.
 *  - 'create-group' - emitted when the user selects to create a group chat.
 */
export default {
	components: { ActionMenu },
	props: {
		mode: { type: String, required: true },
	},
	emits: ['back', 'home', 'create-direct', 'create-group'],

	data() {
		return {
			// Itesm for the popup actionMenu component
			newChatItems: [
				{ id: 'direct', label: 'Create Direct Chat', icon: 'message-circle' },
				{ id: 'group',  label: 'Create Group',       icon: 'users'          },
			],
		};
	},

	computed: {
		isChats() {
			return this.mode === 'chats';
		},
		title() {
			return TITLES[this.mode] || '';
		},
	},

	methods: {
		// ---- ACTION HANDLERS ----
		// Handle selection of a new chat type from the action menu
		onNewChatSelect(item) {
			this.$emit(item.id === 'group' ? 'create-group' : 'create-direct');
		},
	},
};
</script>

<template>
	<div class="sidebar-header">

		<!-- CHATS MODE: Show the app logo and the new chat action menu -->
		<template v-if="isChats">
			<!-- Logo (link to home) -->
			<a href="#/chats" class="sidebar-logo fw-bold" @click="$emit('home')">WASApp</a>
			<!-- ActionMenu for creating new chats (direct or group) -->
			<ActionMenu class="ms-auto" trigger-button-icon="message-square-plus" placement="down-right" :items="newChatItems" trigger-button-filled @select="onNewChatSelect" />
		</template>

		<!-- OTHER MODES: Show a back button and the title of the current mode -->
		<template v-else>
			<button type="button" class="back-btn" @click="$emit('back')">
				<svg class="feather back-icon"><use href="/feather-sprite-v4.29.0.svg#arrow-left"/></svg>
			</button>
			<span class="fw-semibold header-title">{{ title }}</span>
		</template>
	</div>
</template>

<style scoped>
.sidebar-header {
	height: var(--topbar-height);
	flex-shrink: 0;
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 0 16px;
	background: var(--theme-bg-darker);
	color: var(--theme-fg);
}

.sidebar-logo {
	font-size: 1.25rem;
	color: var(--theme-fg);
	text-decoration: none;
}

.sidebar-logo:hover {
	color: var(--theme-blue);
}

.header-title {
	font-size: 1.05rem;
}

.back-btn {
	background: none;
	border: none;
	padding: 0;
	cursor: pointer;
	color: var(--theme-fg-dark);
	line-height: 1;
	display: flex;
	align-items: center;
}

.back-btn:hover {
	color: var(--theme-blue);
}

.back-btn .back-icon {
	width: 22px;
	height: 22px;
}
</style>
