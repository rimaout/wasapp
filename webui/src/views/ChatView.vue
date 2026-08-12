<script>
export default {
	computed: {
		chatId() {
			return this.$route.params.chatId;
		},
	},
};
</script>

<template>
	<div>
		<div class="d-flex align-items-center pt-3 pb-2 mb-3 border-bottom">
			<button class="btn btn-sm btn-outline-secondary me-2" @click="$router.back()">Back</button>
			<h1 class="h3 mb-0">Chat {{ chatId }}</h1>
		</div>
		<p class="text-muted">Chat view coming soon...</p>
	</div>
</template>

<style>
</style>
