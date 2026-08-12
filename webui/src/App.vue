<script setup>
import { RouterView } from 'vue-router'
import ChatsListView from './views/ChatsListView.vue'
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
				<button type="button" class="btn btn-sm btn-outline-primary d-flex align-items-center justify-content-center" @click.prevent style="width: 30px; height: 30px;">+</button>
			</a>
			<button class="navbar-toggler d-md-none me-2 collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>
		</header>

		<div class="container-fluid">
			<div class="row">
				<nav id="sidebarMenu" class="col-md-4 col-lg-3 d-md-block bg-light sidebar collapse">
					<div class="position-sticky sidebar-sticky pt-2">
						<ChatsListView />
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
</style>
