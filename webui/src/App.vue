<script setup>

import { RouterView, useRoute } from 'vue-router'
import { ref, computed } from 'vue'
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

const route        = useRoute();
const isChatRoute  = computed(() => route.path.startsWith('/chats/'));

const sidebarMode  = ref('chats') // current mode of the sidebar: 'chats', 'users', 'group-members', 'group-details'
const groupMembers = ref([])	  // selected members for creating a group chat (used in 'group-members' and 'group-details' modes)
const searchQuery  = ref('')      // search query for search bar (used in 'chats' mode)

function startGroupCreation() {
	groupMembers.value = [];
	sidebarMode.value = 'group-members';
}

function finishGroupCreation() {
	groupMembers.value = [];
	sidebarMode.value = 'chats';
}

function cancelGroupCreation() {
	groupMembers.value = [];
	sidebarMode.value = 'chats';
}

function handleBackButton() {
	if (sidebarMode.value === 'users') {
		sidebarMode.value = 'chats';
	} else if (sidebarMode.value === 'group-members') {
		sidebarMode.value = 'chats';
	} else if (sidebarMode.value === 'group-details') {
		sidebarMode.value = 'group-members';
	}
}
</script>

<template>
	<div v-if="$route.path === '/login'">
		<RouterView />
	</div>

	<div v-else class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-4 col-lg-3 d-md-block sidebar p-0" :class="{ 'd-none d-md-block': isChatRoute }">
				<div class="sidebar-inner">
					<SidebarHeader :mode="sidebarMode" @create-direct="sidebarMode = 'users'" @create-group="startGroupCreation" @back="handleBackButton" @home="sidebarMode = 'chats'" />
					<div class="sidebar-sticky" :class="{ 'fade-bottom': sidebarMode === 'chats' }">
						<template v-if="sidebarMode === 'chats'">
							<SearchBar v-model="searchQuery" />
							<ChatsListView :searchQuery="searchQuery" />
						</template>
						<UserSearchView v-else-if="sidebarMode === 'users'" @close="sidebarMode = 'chats'" />
						<GroupMembersView v-else-if="sidebarMode === 'group-members'" v-model:members="groupMembers" @create="sidebarMode = 'group-details'" @cancel="cancelGroupCreation" />
						<GroupDetailsView v-else-if="sidebarMode === 'group-details'" v-model:members="groupMembers" @done="finishGroupCreation" @add-members="sidebarMode = 'group-members'" @cancel="cancelGroupCreation" />
					</div>
					<ProfileBar v-if="sidebarMode === 'chats'" />
				</div>
			</nav>

			<main class="col-md-8 ms-sm-auto col-lg-9 p-0" :class="{ 'd-none d-md-block': !isChatRoute }">
				<RouterView />
			</main>
		</div>
	</div>
</template>
