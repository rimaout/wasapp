<script setup>
import { RouterView } from 'vue-router'
import { ref } from 'vue'
import ChatsListView from './views/ChatsListView.vue'
import SearchBar from './components/SearchBar.vue'

const searchQuery = ref('')
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
				<button type="button" class="plus-btn" @click.prevent>
					<svg class="feather plus-icon"><use href="/feather-sprite-v4.29.0.svg#plus-square"/></svg>
				</button>
			</a>
			<button class="navbar-toggler d-md-none me-2 collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>
		</header>

		<div class="container-fluid">
			<div class="row">
				<nav id="sidebarMenu" class="col-md-4 col-lg-3 d-md-block bg-light sidebar collapse">
					<div class="position-sticky sidebar-sticky pt-2">
						<SearchBar v-model="searchQuery" />
						<ChatsListView :searchQuery="searchQuery" />
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

.plus-icon {
	width: 26px;
	height: 26px;
}
</style>
