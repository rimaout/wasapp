<script setup>
import { computed } from 'vue';

const props = defineProps({
	mode: { type: String, required: true },
});
const emit = defineEmits(['new-chat', 'back', 'home']);

const TITLES = {
	users: 'New chat',
	'group-members': 'Add group members',
	'group-details': 'New group',
};

const isChats = computed(() => props.mode === 'chats');
const title = computed(() => TITLES[props.mode] || '');
</script>

<template>
	<div class="sidebar-header">
		<template v-if="isChats">
			<a href="#/chats" class="sidebar-logo fw-bold" @click="emit('home')">WASApp</a>
			<button type="button" class="plus-btn ms-auto" @click="emit('new-chat')">
				<svg class="feather plus-icon"><use href="/feather-sprite-v4.29.0.svg#plus-square"/></svg>
			</button>
		</template>
		<template v-else>
			<button type="button" class="back-btn" @click="emit('back')">
				<svg class="feather back-icon"><use href="/feather-sprite-v4.29.0.svg#arrow-left"/></svg>
			</button>
			<span class="fw-semibold header-title">{{ title }}</span>
		</template>
	</div>
</template>

<style scoped>
.sidebar-header {
	height: var(--topbar-height);
	flex-shrink: 0;
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 0 16px;
	border-bottom: 1px solid var(--tn-border);
	background: var(--tn-bg-darker);
	color: var(--tn-fg);
}

.sidebar-logo {
	font-size: 1.25rem;
	color: var(--tn-fg);
	text-decoration: none;
}

.sidebar-logo:hover {
	color: var(--tn-blue);
}

.header-title {
	font-size: 1.05rem;
}

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

.back-btn {
	background: none;
	border: none;
	padding: 0;
	cursor: pointer;
	color: var(--tn-fg-dark);
	line-height: 1;
	display: flex;
	align-items: center;
}

.back-btn:hover {
	color: var(--tn-blue);
}

.back-btn .back-icon {
	width: 22px;
	height: 22px;
}
</style>
