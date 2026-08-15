<script>
import { getAvatarColor, getAvatarLetter } from '../services/utils.js';

// Shared avatar rendering: loads an image from `imageUrl`, falling back
// to a colored letter (or group icon) when no image is available.
export default {
	props: {
		imageUrl: { type: String, required: true },
		displayName: { type: String, required: true },
		size: { type: Number, default: 48 },
		isGroup: { type: Boolean, default: false },
		version: { type: Number, default: 0 },	// Used to trigger avatar reload when the avatar is updated
	},
	data() {
		return {
			imgSrc: null,
			showImage: false,
		};
	},
	computed: {
		style() {
			let px = this.size + 'px';
			let radius = this.isGroup ? '33%' : '50%';
			return {
				width: px,
				height: px,
				minWidth: px,
				fontSize: (this.size * 0.45) + 'px',
				borderRadius: radius,
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
		async loadImage() {
			try {
				let response = await this.$axios.get(this.imageUrl, { responseType: 'blob' });
				this.imgSrc = URL.createObjectURL(response.data);
				this.showImage = true;
			} catch (e) {
				this.showImage = false;
			}
		},
	},
	watch: {
		version: {
			handler: 'loadImage',
			immediate: true,
		},
	},
	mounted() {
		this.loadImage();
	},
	beforeUnmount() {
		if (this.imgSrc) URL.revokeObjectURL(this.imgSrc);
	},
};
</script>

<template>
	<div class="avatar-circle" :style="style">
		<img v-if="showImage" :src="imgSrc" class="avatar-img" />
		<span v-else-if="isGroup" class="avatar-group-icon" :style="{ backgroundColor: bgColor }">
			<i class="bi bi-people-fill"></i>
		</span>
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
	color: var(--tn-bg-darker);
}

.avatar-group-icon {
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	color: var(--tn-bg-darker);
	font-size: 120%;
}
</style>
