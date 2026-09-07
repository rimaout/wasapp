<script>
import BaseAvatar from './BaseAvatar.vue';

// Avatar for a user, loading from /users/{id}/avatar.
export default {
	components: { BaseAvatar },
	props: {
		userId:      { type: String, required: true }, // User ID to load the avatar for
		displayName: { type: String, required: true }, // Name to display when no image is available (used to generate the letter and color)
		size:        { type: Number, default: 48    }, // Size of the avatar in pixels (width and height)
		version:     { type: Number, default: 0     }, // Used to trigger avatar reload when the avatar is updated
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
