<script>
import { getUserName } from '../services/auth.js';

export default {
	computed: {
		userName() {
			return getUserName();
		},
	},
};
</script>

<template>
	<div class="d-flex align-items-center justify-content-center min-vh-100">
		<div class="text-center">
			<h2 class="fw-normal" style="color: var(--tn-fg-dark);">Hi 👋 {{ userName }}!!</h2>
			<p class="text-muted">Choose a chat or tap + to start a new one</p>
		</div>
	</div>
</template>

<style scoped>
</style>
