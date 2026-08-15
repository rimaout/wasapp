<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from '../services/axios.js';
import { navigateToChat } from '../services/chatNavigation.js';
import { getErrorMessage } from '../services/utils.js';

const props = defineProps({
	members: { type: Array, default: () => [] },
});
const emit = defineEmits(['done']);
const router = useRouter();

const name = ref('');
const imageFile = ref(null);
const imagePreview = ref(null);
const errormsg = ref(null);
const creating = ref(false);

function onFileChange(e) {
	const file = e.target.files && e.target.files[0];
	if (!file) return;

	if (file.size > 5 * 1024 * 1024) {
		errormsg.value = 'Image must be at most 5 MB';
		e.target.value = '';
		return;
	}

	imageFile.value = file;
	if (imagePreview.value) URL.revokeObjectURL(imagePreview.value);
	imagePreview.value = URL.createObjectURL(file);
}

function removeImage() {
	imageFile.value = null;
	if (imagePreview.value) URL.revokeObjectURL(imagePreview.value);
	imagePreview.value = null;
}

function validate() {
	const n = name.value.trim();
	if (n.length < 3 || n.length > 24) {
		return 'Group name must be 3-24 characters';
	}
	if (!/^[A-Za-z0-9 _-]+$/.test(n)) {
		return 'Group name can only contain letters, numbers, spaces, underscores and hyphens';
	}
	if (!/\S/.test(n)) {
		return 'Group name must contain at least one non-space character';
	}
	return '';
}

async function createGroup() {
	const err = validate();
	if (err) {
		errormsg.value = err;
		return;
	}
	if (creating.value) return;

	creating.value = true;
	errormsg.value = null;
	try {
		const response = await axios.post('/chats', {
			groupName: name.value.trim(),
			membersList: props.members.map(m => m.id),
		});
		const chat = response.data;

		if (imageFile.value) {
			const formData = new FormData();
			formData.append('binaryImage', imageFile.value);
			await axios.put('/chats/' + chat.id + '/avatar', formData);
		}

		navigateToChat(router, chat);
		emit('done');
	} catch (e) {
		errormsg.value = getErrorMessage(e);
	} finally {
		creating.value = false;
	}
}
</script>

<template>
	<div class="group-details-view">
		<div class="body flex-grow-1 px-3 py-3">
			<label class="form-label field-label">Group name</label>
			<input type="text" class="form-control name-input" v-model="name" maxlength="24" placeholder="Enter group name" />

			<label class="form-label field-label mt-4">Group image (optional)</label>
			<div class="d-flex align-items-center gap-3">
				<label class="image-picker">
					<img v-if="imagePreview" :src="imagePreview" class="image-preview" />
					<span v-else class="image-placeholder">
						<svg class="feather camera-icon"><use href="/feather-sprite-v4.29.0.svg#camera"/></svg>
					</span>
					<input type="file" accept="image/jpeg,image/png,image/webp" class="d-none" @change="onFileChange" />
				</label>
				<div>
					<div class="text-muted small">JPEG, PNG or WebP, up to 5 MB</div>
					<button v-if="imagePreview" type="button" class="btn btn-sm btn-link remove-btn" @click="removeImage">Remove</button>
				</div>
			</div>

			<p v-if="members.length > 0" class="text-muted small mt-3">
				{{ members.length }} member{{ members.length > 1 ? 's' : '' }} selected
			</p>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="footer px-3 py-3">
			<button type="button" class="btn btn-primary w-100 create-btn" :disabled="creating" @click="createGroup">
				Create group
			</button>
		</div>
	</div>
</template>

<style scoped>
.group-details-view {
	display: flex;
	flex-direction: column;
	height: 100%;
}

.body {
	overflow-y: auto;
	scrollbar-width: none;
	-ms-overflow-style: none;
}

.body::-webkit-scrollbar {
	display: none;
}

.field-label {
	color: var(--tn-fg);
	font-weight: 600;
	font-size: 0.85rem;
}

.name-input {
	background: var(--tn-bg-highlight);
	border-color: var(--tn-border);
	color: var(--tn-fg);
}

.image-picker {
	width: 96px;
	height: 96px;
	border-radius: 33%;
	overflow: hidden;
	cursor: pointer;
	flex-shrink: 0;
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
	background: var(--tn-bg-highlight);
	color: var(--tn-fg-dark);
}

.camera-icon {
	width: 28px;
	height: 28px;
}

.remove-btn {
	padding: 0;
	color: var(--tn-red);
	text-decoration: none;
}

.footer {
	border-top: 1px solid var(--tn-border);
	flex-shrink: 0;
}

.create-btn {
	font-weight: 600;
}
</style>
