<script>
import BaseAvatar from './BaseAvatar.vue';

// Avatar for a user, loading from /users/{id}/avatar.
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
