<script>
import { getAvatarColor, getAvatarLetter } from '../../services/utils.js';

// Shared avatar rendering: loads an image from `imageUrl`, falling back
// to a colored letter (or group icon) when no image is available.
export default {
	props: {
		imageUrl:    { type: String,  required: true }, // URL to the avatar image (can be a blob URL)
		displayName: { type: String,  required: true }, // Name to display when no image is available (used to generate the letter and color)
		size:        { type: Number,  default: 48    }, // Size of the avatar in pixels (width and height)
		isGroup:     { type: Boolean, default: false }, // Whether this is a group chat (show group icon instead of letter)
		version:     { type: Number,  default: 0     },	// Used to trigger avatar reload when the avatar is updated
	},
	data() {
		// Internal Data for the component
		return {
			imgSrc: null,
			showImage: false,
		};
	},
	computed: {
		// Intaernal Data that is computed based on props and data (re-computed when props/data change)
		style() {
			let px = this.size + 'px';
			return {
				width: px,
				height: px,
				minWidth: px,
				fontSize: (this.size * 0.45) + 'px',
				borderRadius: '50%',
			};
		},
		letter() {
			return getAvatarLetter(this.displayName);
		},
		bgColor() {
			return getAvatarColor(this.displayName);
		},
	},
	methods: {
		// Load the image from the server and create a blob URL for it.
		// NOTE: a blob URL is a local URL that points to a blob (binary data) in memory, it must be revoked when no longer needed.
		async loadImage() {
			try {
				// Request the image as a blob (binary data) from the server, and create a blob URL for it.
				let response = await this.$axios.get(this.imageUrl, { responseType: 'blob' });
				let url = URL.createObjectURL(response.data);

				// Release the previous blob URL before swapping to the new one
				if (this.imgSrc) URL.revokeObjectURL(this.imgSrc);
				this.imgSrc = url;
				this.showImage = true;
			} catch (e) {
				this.showImage = false;
			}
		},
	},
	watch: {
		// Reload when the image URL changes (switching to another chat) or when the version prop is bumped (avatar was re-uploaded).
		imageUrl: 'loadImage',
		version:  'loadImage',
	},
	mounted() {
		// Load the image when the component is mounted.
		this.loadImage();
	},
	beforeUnmount() {
		// Release the blob URL when the component is unmounted.
		if (this.imgSrc) URL.revokeObjectURL(this.imgSrc);
	},
};
</script>

<template>
	<div class="avatar-circle" :style="style">
		<!-- Show the image if it was successfully loaded -->
		<img v-if="showImage" :src="imgSrc" class="avatar-img" />

		<!-- Else, show the group icon if it's a group chat -->
		<span v-else-if="isGroup" class="avatar-group-icon" :style="{ backgroundColor: bgColor }">
			<i class="bi bi-people-fill"></i>
		</span>

		<!-- Else, show the letter avatar -->
		<span v-else class="avatar-letter" :style="{ backgroundColor: bgColor }">{{ letter }}</span>
	</div>
</template>

<style scoped>
.avatar-circle {
	overflow: hidden;
	flex-shrink: 0;
}

.avatar-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.avatar-letter {
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 600;
	color: var(--theme-bg-darker);
}

.avatar-group-icon {
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	color: var(--theme-bg-darker);
	font-size: 120%;
}
</style>
