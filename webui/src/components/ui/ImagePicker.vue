<script>

/**
 * ImagePicker — a Vue 3 component for selecting and previewing an image file.
 *
 * Props:
 *  - selectedImage: The selected File object (v-model binding).
 *  - previewUrl:    Initial/fallback image URL when no local File object is selected.
 *  - radius:        CSS border-radius value controlling container shape (e.g., '50%' or '12px').
 *  - size:          Width and height of the component in pixels.
 *
 * Emits:
 *  - update:selectedImage: Emitted when a new file is selected, passing the File object to the parent.
 *  - error:                Emitted when validation fails (e.g., file size exceeds limit), passing an error message to the parent.
 */
export default {
	props: {
		selectedImage: { type: File,   default: null  },
		previewUrl:    { type: String, default: null  },
		radius:        { type: String, default: '50%' },
		size:          { type: Number, default: 96    },
	},

	emits: ['update:selectedImage', 'error'],

	data() {
		return {
			// Temporary browser blob/object URL created for local image preview
			objectUrl: null,
		};
	},

	computed: {
		// Calculates camera icon size proportional to overall component dimensions
		iconSize() {
			return Math.round(this.size * 28 / 96); //Note: 28 is the default icon size for 96px component size
		},
		// Display source priority: newly selected local file > initial preview URL > fallback null
		displaySrc() {
			return this.objectUrl || this.previewUrl || null;
		},
	},

	watch: {
		// Automatically manages blob URLs: removes previous object URL and generates a new preview URL
		selectedImage(file) {
			if (this.objectUrl) URL.revokeObjectURL(this.objectUrl);  // Cleanup previous object URL
			this.objectUrl = file ? URL.createObjectURL(file) : null; // Generate new object URL for the selected file or reset to null if no file is selected
		},
	},

	methods: {
		// Handles file input selection: validates max file size (5MB) and updates v-model
		onFileChange(e) {
			// Extract the first selected file from the input event
			let file = null;
			if (e.target.files) {
				file = e.target.files[0];
			}
			if (!file) return;

			// Validate image size limit (5 MB)
			if (file.size > 5 * 1024 * 1024) {
				this.$emit('error', 'Image must be at most 5 MB');
				e.target.value = ''; // Reset input selection
				return;
			}

			// Emit the selected file to the parent component
			this.$emit('update:selectedImage', file);
		},
	},

	beforeUnmount() {
		// Memory cleanup: revokes generated object URL when component unmounts
		if (this.objectUrl) URL.revokeObjectURL(this.objectUrl);
	},
};
</script>

<template>
    <!-- Clickable label wrapper that triggers the hidden native file input -->
    <label class="image-picker" :style="{ borderRadius: radius, width: size + 'px', height: size + 'px' }">

        <!-- Rendered image preview when a local blob URL or previewUrl is available -->
        <img v-if="displaySrc" :src="displaySrc" class="image-preview" />

        <!-- Default camera placeholder icon displayed when no image source is present -->
        <span v-else class="image-placeholder">
            <svg class="feather camera-icon" :style="{ width: iconSize + 'px', height: iconSize + 'px' }">
                <use href="/feather-sprite-v4.29.0.svg#camera"/>
            </svg>
        </span>

        <!-- Hidden file input constrained to allowed image MIME types -->
        <input type="file" accept="image/jpeg,image/png,image/webp" class="d-none" @change="onFileChange" />

    </label>
</template>

<style scoped>
.image-picker {
    overflow: hidden;
    cursor: pointer;
    flex-shrink: 0;
    align-self: center;
    display: flex;
    align-items: center;
    justify-content: center;
}

.image-preview {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.image-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--theme-bg-highlight);
    color: var(--theme-fg-dark);
}
</style>
