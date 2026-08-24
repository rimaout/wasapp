<script>
import BaseAvatar from './BaseAvatar.vue';

// Avatar for a chat (direct or group), loading from /chats/{id}/avatar.
export default {
	components: { BaseAvatar },
	props: {
		chatId: { type: String, required: true },
		displayName: { type: String, required: true },
		size: { type: Number, default: 48 },
		isGroup: { type: Boolean, default: false },
		version: { type: Number, default: 0 },
	},
	computed: {
		imageUrl() {
			return '/chats/' + this.chatId + '/avatar';
		},
	},
};
</script>

<template>
	<BaseAvatar :image-url="imageUrl" :display-name="displayName" :size="size" :is-group="isGroup" :version="version" />
</template>
