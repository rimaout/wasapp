<script setup>
import { computed } from 'vue';
import ActionMenu from './ActionMenu.vue';

const props = defineProps({
	mode: { type: String, required: true },
});
const emit = defineEmits(['back', 'home', 'create-direct', 'create-group']);

const TITLES = {
	users: 'New direct chat',
	'group-members': 'Add group members',
	'group-details': 'New group details',
};

const isChats = computed(() => props.mode === 'chats');
const title = computed(() => TITLES[props.mode] || '');

const newChatItems = [
	{ id: 'direct', label: 'Create Direct Chat', icon: 'message-circle' },
	{ id: 'group', label: 'Create Group', icon: 'users' },
];

function onNewChatSelect(item) {
	emit(item.id === 'group' ? 'create-group' : 'create-direct');
}
</script>

<template>
	<div class="sidebar-header">
		<template v-if="isChats">
			<a href="#/chats" class="sidebar-logo fw-bold" @click="emit('home')">WASApp</a>
			<ActionMenu class="ms-auto" icon="plus-square" label="New chat" placement="down-right" :items="newChatItems" @select="onNewChatSelect" />
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
