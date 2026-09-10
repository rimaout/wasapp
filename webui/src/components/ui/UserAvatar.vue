<script>
import BaseAvatar from './BaseAvatar.vue';

/**
 * UserAvatar — avatar rendering for a user, avatar image loading from /users/{id}/avatar.
 *
 * Props:
 *  - userId:      ID of the user to load the avatar for
 *  - displayName: Name to display when no image is available (used to generate the letter and color)
 *  - size:        Size of the avatar in pixels (width and height)
 *  - version:     Used to trigger avatar reload when the avatar is updated
 */
export default {
	components: { BaseAvatar },
	props: {
		userId:      { type: String, required: true },
		displayName: { type: String, required: true },
		size:        { type: Number, default: 48    },
		version:     { type: Number, default: 0     },
	},
	computed: {
		imageUrl() {
			return '/users/' + this.userId + '/avatar';
		},
	},
};
</script>

<template>
	<BaseAvatar :image-url="imageUrl" :display-name="displayName" :size="size" :is-group="false" :version="version" />
</template>
