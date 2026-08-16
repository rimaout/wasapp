<script setup>
import { ref } from 'vue';
import UserAvatar from './UserAvatar.vue';
import PopupMenu from './PopupMenu.vue';
import { getUserId, getUserName, setUserName, clearAuth } from '../services/auth.js';

const userId = getUserId();
const name = ref(getUserName() || 'User');
const avatarVersion = ref(0);

const profileItems = [
	{ id: 'name', label: 'Change Profile Name', icon: 'edit-2', action: 'rename' },
	{ id: 'image', label: 'Change Profile Image', icon: 'image', action: 'image' },
	{ id: 'logout', label: 'Logout', icon: 'log-out', danger: true },
];

function onSelect(item) {
	if (item.id === 'logout') {
		clearAuth();
		window.location.hash = '#/login';
	}
}

function onDone(payload) {
	if (payload.action === 'rename') {
		setUserName(payload.name);
		name.value = payload.name;
	} else if (payload.action === 'image') {
		avatarVersion.value++;
	}
}
</script>

<template>
	<div class="profile-bar">
		<UserAvatar :userId="userId" :displayName="name" :size="40" :version="avatarVersion" />
		<span class="profile-greeting">Hi 👋 {{ name }}</span>
		<PopupMenu class="ms-auto" icon="settings" label="Profile options" direction="up" :items="profileItems" :target="{ kind: 'me', userId }" :initial-name="name" @select="onSelect" @done="onDone" />
	</div>
</template>

<style scoped>
.profile-bar {
	display: flex;
	align-items: center;
	gap: 12px;
	margin: 4px 12px 12px;
	padding: 8px 12px;
	background: var(--tn-bg-light);
	border-radius: 12px;
	flex-shrink: 0;
}

.profile-greeting {
	color: var(--tn-fg);
	font-weight: 600;
	font-size: 0.95rem;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}
</style>
