<script setup>
import { RouterView, useRoute } from 'vue-router'
import { ref, computed } from 'vue'
import ChatsListView from './views/ChatsListView.vue'
import UserSearchView from './views/UserSearchView.vue'
import GroupMembersView from './views/GroupMembersView.vue'
import GroupDetailsView from './views/GroupDetailsView.vue'
import SidebarHeader from './components/SidebarHeader.vue'
import SearchBar from './components/SearchBar.vue'
import ProfileBar from './components/ProfileBar.vue'

const route = useRoute();
const isChatRoute = computed(() => route.path.startsWith('/chats/'));
const searchQuery = ref('')
const sidebarMode = ref('chats')
const groupMembers = ref([])

function startGroupCreation() {
	groupMembers.value = [];
	sidebarMode.value = 'group-members';
}

function finishGroupCreation() {
	groupMembers.value = [];
	sidebarMode.value = 'chats';
}

function handleBack() {
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
			<nav id="sidebarMenu" class="col-md-5 col-lg-4 d-md-block sidebar p-0" :class="{ 'd-none d-md-block': isChatRoute }">
				<div class="sidebar-inner">
					<SidebarHeader :mode="sidebarMode" @create-direct="sidebarMode = 'users'" @create-group="startGroupCreation" @back="handleBack" @home="sidebarMode = 'chats'" />
					<div class="sidebar-sticky" :class="{ 'fade-bottom': sidebarMode === 'chats' }">
						<template v-if="sidebarMode === 'chats'">
							<SearchBar v-model="searchQuery" />
							<ChatsListView :searchQuery="searchQuery" />
						</template>
						<UserSearchView v-else-if="sidebarMode === 'users'" @close="sidebarMode = 'chats'" />
						<GroupMembersView v-else-if="sidebarMode === 'group-members'" v-model:members="groupMembers" @create="sidebarMode = 'group-details'" />
						<GroupDetailsView v-else-if="sidebarMode === 'group-details'" :members="groupMembers" @done="finishGroupCreation" />
					</div>
					<ProfileBar v-if="sidebarMode === 'chats'" />
				</div>
			</nav>

			<main class="col-md-7 ms-sm-auto col-lg-8 p-0" :class="{ 'd-none d-md-block': !isChatRoute }">
				<RouterView />
			</main>
		</div>
	</div>
</template>
