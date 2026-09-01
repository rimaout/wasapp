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
				<h2 class="fw-normal" style="color: var(--theme-fg-dark);">Hi 👋 {{ userName }}!!</h2>
				<p class="welcome-hint">Choose a chat or tap + to start a new one</p>
			</div>
		</div>
	</div>
</template>

<style scoped>
.welcome-view {
	display: flex;
	flex-direction: column;
	height: 100vh;
	height: 100dvh;
	background-color: var(--theme-bg-dark);
	background-image: url('/chat-bg.svg');
	background-repeat: repeat;
	background-size: 500px;
}

.welcome-header {
	height: var(--topbar-height);
	flex-shrink: 0;
	background: var(--theme-bg-darker);
}

.welcome-body {
	flex: 1 1 auto;
}

.welcome-hint {
	display: inline-block;
	background: var(--theme-bg-highlight);
	color: var(--theme-fg-dark);
	padding: 4px 10px;
	border-radius: 999px;
}
</style>
