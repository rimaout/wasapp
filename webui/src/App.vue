<script setup>
import { RouterView } from 'vue-router'
import { ref } from 'vue'
import ChatsListView from './views/ChatsListView.vue'
import UserSearchView from './views/UserSearchView.vue'
import SearchBar from './components/SearchBar.vue'
import ChatAvatar from './components/ChatAvatar.vue'

const searchQuery = ref('')
const showUserSearch = ref(false)
</script>

<script>
export default {}
</script>

<template>

	<div v-if="$route.path === '/login'">
		<RouterView />
	</div>

	<div v-else>
		<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0">
			<a class="navbar-brand col-md-4 col-lg-3 me-0 px-3 fs-5 fw-bold d-flex justify-content-between align-items-center gap-2" href="#/chats">
				<span>WASApp</span>
				<button type="button" class="plus-btn" @click.prevent="showUserSearch = true">
					<svg class="feather plus-icon"><use href="/feather-sprite-v4.29.0.svg#plus-square"/></svg>
				</button>
			</a>
			<div v-if="$route.params.chatId" class="chat-nav d-flex align-items-center flex-grow-1 px-3">
				<ChatAvatar :chatId="$route.params.chatId" :displayName="$route.query.name || 'Chat'" :size="32" :isGroup="$route.query.group === '1'" class="me-2" />
				<span class="fw-semibold text-truncate chat-nav-name">{{ $route.query.name || 'Chat' }}</span>
			</div>
			<button class="navbar-toggler d-md-none me-2 collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>
		</header>

		<div class="container-fluid">
			<div class="row">
				<nav id="sidebarMenu" class="col-md-4 col-lg-3 d-md-block bg-light sidebar collapse">
					<div class="position-sticky sidebar-sticky pt-2">
						<UserSearchView v-if="showUserSearch" @close="showUserSearch = false" />
						<template v-else>
							<SearchBar v-model="searchQuery" />
							<ChatsListView :searchQuery="searchQuery" />
						</template>
					</div>
				</nav>

				<main class="col-md-8 ms-sm-auto col-lg-9 px-md-4">
					<RouterView />
				</main>
			</div>
		</div>
	</div>
</template>

<style>
.plus-btn {
	background: none;
	border: none;
	padding: 0;
	cursor: pointer;
	color: var(--tn-fg-dark);
	line-height: 1;
}

.plus-btn:hover {
	color: var(--tn-blue);
}

.plus-btn:active {
	color: var(--tn-cyan);
}

.plus-btn .plus-icon {
	width: 23px;
	height: 23px;
}

.chat-nav-name {
	color: var(--tn-fg);
	font-size: 1.15rem;
}
</style>
