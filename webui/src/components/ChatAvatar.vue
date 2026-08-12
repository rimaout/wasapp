<script>
const AVATAR_COLORS = [ //Tokyo Night color palette
	'#7aa2f7', '#bb9af7', '#9ece6a', '#e0af68', '#f7768e',
	'#7dcfff', '#ff9e64', '#c0caf5', '#565f89', '#414868',
];

// Generate a color based on the name string
function getAvatarColor(name) {
	let hash = 0;
	for (let i = 0; i < name.length; i++) {
		hash = name.charCodeAt(i) + ((hash << 5) - hash);
	}
	return AVATAR_COLORS[Math.abs(hash) % AVATAR_COLORS.length];
}

// Get the first letter of the name, or '?' if name is empty
function getAvatarLetter(name) {

	if (!name || name.length === 0) return '?';

	return name[0].toUpperCase();
}

// Vue component for displaying a chat avatar
export default {

	// Input Chat Data Needed
	props: {
		chatId: { type: String, required: true },
		displayName: { type: String, required: true },
		size: { type: Number, default: 48 },
		isGroup: { type: Boolean, default: false },
		version: { type: Number, default: 0 },	// Used to trigger avatar reload when the avatar is updated
	},

	// Component Data
	data() {
		return {
			imgSrc: null,
			showImage: false,
		};
	},

	// Computed Properties
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

	// Component Methods
	methods: {
		async loadImage() {
			// Load the avatar image from the server
			try {
				let response = await this.$axios.get('/chats/' + this.chatId + '/avatar', {
					responseType: 'blob',
				});
				this.imgSrc = URL.createObjectURL(response.data);
				this.showImage = true;
			} catch (e) {
				this.showImage = false;
			}
		},
	},
	watch: {
		version: {
			// If the version prop changes, reload the avatar image
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
