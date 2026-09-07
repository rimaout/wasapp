<script>
import BaseAvatar from './BaseAvatar.vue';

// Avatar for a chat (direct or group), loading from /chats/{id}/avatar.
export default {
	components: { BaseAvatar },
	props: {
		chatId:      { type: String,  required: true  }, // Chat ID to load the avatar for
		displayName: { type: String,  required: true  }, // Name to display when no image is available (used to generate the letter and color)
		size:        { type: Number,  default:  48    }, // Size of the avatar in pixels (width and height)
		isGroup:     { type: Boolean, default:  false }, // Whether this is a group chat (show group icon instead of letter)
		version:     { type: Number,  default:  0     }, // Used to trigger avatar reload when the avatar is updated
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
