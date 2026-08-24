<script setup>
import { ref, computed } from 'vue';
import UserAvatar from './UserAvatar.vue';
import ActionMenu from './ActionMenu.vue';
import RenameActionScreen from './RenameActionScreen.vue';
import ImageActionScreen from './ImageActionScreen.vue';
import { getUserId, getUserName, setUserName, clearAuth } from '../services/auth.js';

const userId = getUserId();
const name = ref(getUserName() || 'User');
const avatarVersion = ref(0);

const profileItems = computed(() => [
	{ id: 'name', label: 'Change Profile Name', icon: 'edit-2', component: RenameActionScreen, props: { title: 'Change Profile Name', target: { kind: 'me', userId }, initialName: name.value } },
	{ id: 'image', label: 'Change Profile Image', icon: 'image', component: ImageActionScreen, props: { title: 'Change Profile Image', target: { kind: 'me', userId } } },
	{ id: 'logout', label: 'Logout', icon: 'log-out', danger: true, dangerText: 'Are you sure you want to log out?' },
]);

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
		<ActionMenu class="ms-auto" trigger-button-icon="settings" placement="top-right" :items="profileItems" @select="onSelect" @done="onDone" />
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
