<script>
import { RouterView } from 'vue-router'
import ChatsListView from './views/ChatsListView.vue'
import UserSearchView from './views/UserSearchView.vue'
import GroupMembersView from './views/GroupMembersView.vue'
import GroupDetailsView from './views/GroupDetailsView.vue'
import SidebarHeader from './components/SidebarHeader.vue'
import SearchBar from './components/ui/SearchBar.vue'
import ProfileBar from './components/ProfileBar.vue'

// Main app component. It contains the sidebar and the main content area.
//
// The main constent area shows by default the welcome screen, or the chat view if a chat is selected
//
// The sidebar can be in different modes:
//  - 'chats': shows the chat list and search bar
//  - 'users': shows the user search view for starting a direct chat
//  - 'group-members': shows the group member selection view for creating a group chat

export default {
	components: { RouterView, ChatsListView, UserSearchView, GroupMembersView, GroupDetailsView, SidebarHeader, SearchBar, ProfileBar },

	data() {
		return {
			sidebarMode: 'chats',   // current mode of the sidebar: 'chats', 'users', 'group-members', 'group-details'
			groupMembers: [],       // selected members for creating a group chat (used in 'group-members' and 'group-details' modes)
			searchQuery: '',        // search query for search bar (used in 'chats' mode)
		};
	},

	computed: {
		isChatRoute() {
			return this.$route.path.startsWith('/chats/');
		},
	},

	methods: {
		startGroupCreation() {
			this.groupMembers = [];
			this.sidebarMode = 'group-members';
		},

		finishGroupCreation() {
			this.groupMembers = [];
			this.sidebarMode = 'chats';
		},

		cancelGroupCreation() {
			this.groupMembers = [];
			this.sidebarMode = 'chats';
		},

		handleBackButton() {
			if (this.sidebarMode === 'users') {
				this.sidebarMode = 'chats';
			} else if (this.sidebarMode === 'group-members') {
				this.sidebarMode = 'chats';
			} else if (this.sidebarMode === 'group-details') {
				this.sidebarMode = 'group-members';
			}
		},
	},
};
</script>

<template>

	<!-- If the current route is '/login', show only the login view -->
	<div v-if="$route.path === '/login'">
		<RouterView /> <!-- RouterView automatically renders the component for the current route (in this case /login) -->
	</div>

	<!-- Otherwise, show the main app layout with sidebar and main content area -->
	<div v-else class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-4 col-lg-3 d-md-block sidebar p-0" :class="{ 'd-none d-md-block': isChatRoute }">
				<div class="sidebar-inner">

					<!-- Sidebar header: shows the title and buttons for creating direct chats or group chats -->
					<SidebarHeader :mode="sidebarMode" @create-direct="sidebarMode = 'users'" @create-group="startGroupCreation" @back="handleBackButton" @home="sidebarMode = 'chats'" />

					<!-- Sidebar content: shows different views based on the current sidebar mode -->
					<div class="sidebar-sticky" :class="{ 'fade-bottom': sidebarMode === 'chats' }">

						<!-- Chats list view: shows the search bar and the list of chats -->
						<template v-if="sidebarMode === 'chats'">
							<SearchBar v-model:search-text="searchQuery" />
							<ChatsListView :searchQuery="searchQuery" />
						</template>

						<!-- DIRECT CHAT CREATION -->
						<!-- User search view: shows the user search view for starting a direct chat -->
						<UserSearchView v-else-if="sidebarMode === 'users'" @close="sidebarMode = 'chats'" />

						<!-- GROUP CHAT CREATION -->
						<!-- Group members view: shows the group member selection view for creating a group chat -->
						<GroupMembersView v-else-if="sidebarMode === 'group-members'" v-model:members="groupMembers" @create="sidebarMode = 'group-details'" @cancel="cancelGroupCreation" />
						<!-- Group details view: shows the group details view for finalizing the group chat creation -->
						<GroupDetailsView v-else-if="sidebarMode === 'group-details'" v-model:members="groupMembers" @done="finishGroupCreation" @add-members="sidebarMode = 'group-members'" @cancel="cancelGroupCreation" />
					</div>

					<!-- Profile bar: shows the user's profile information at the bottom of the sidebar -->
					<ProfileBar v-if="sidebarMode === 'chats'" />
				</div>
			</nav>

			<!-- Main content area: shows the chat view or welcome screen -->
			<main class="col-md-8 ms-sm-auto col-lg-9 p-0" :class="{ 'd-none d-md-block': !isChatRoute }">
				<RouterView /> <!-- RouterView automatically renders the component for the current route (in this case /chats/:id) -->
			</main>
		</div>
	</div>
</template>
