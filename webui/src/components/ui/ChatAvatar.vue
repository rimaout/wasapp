<script>
import BaseAvatar from './BaseAvatar.vue';

/**
 * ChatAvatar — avatar rendering for a chat (direct or group), avatar image
 * loading from /chats/{id}/avatar.
 *
 * Props:
 *  - chatId:      ID of the chat to load the avatar for
 *  - displayName: Name to display when no image is available (used to generate the letter and color)
 *  - size:        Size of the avatar in pixels (width and height)
 *  - isGroup:     Whether this is a group chat (show group icon instead of letter)
 *  - version:     Used to trigger avatar reload when the avatar is updated
 */
export default {
	components: { BaseAvatar },
	props: {
		chatId:      { type: String,  required: true  },
		displayName: { type: String,  required: true  },
		size:        { type: Number,  default:  48    },
		isGroup:     { type: Boolean, default:  false },
		version:     { type: Number,  default:  0     },
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
