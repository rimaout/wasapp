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
	<div class="welcome-view">
		<div class="welcome-header"></div>
		<div class="welcome-body d-flex align-items-center justify-content-center">
			<div class="text-center">
				<h2 class="fw-normal" style="color: var(--tn-fg-dark);">Hi 👋 {{ userName }}!!</h2>
				<p class="text-muted">Choose a chat or tap + to start a new one</p>
			</div>
		</div>
	</div>
</template>

<style scoped>
.welcome-view {
	display: flex;
	flex-direction: column;
	height: 100vh;
}

.welcome-header {
	height: var(--topbar-height);
	flex-shrink: 0;
	background: var(--tn-bg-darker);
}

.welcome-body {
	flex: 1 1 auto;
}
</style>
