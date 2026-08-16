<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import ImagePicker from '../components/ImagePicker.vue';
import { createGroup as createGroupRequest, updateAvatar } from '../services/api.js';
import { refreshChats } from '../composables/useChats.js';
import { navigateToChat } from '../services/chatNavigation.js';
import { getErrorMessage } from '../services/utils.js';

const props = defineProps({
	members: { type: Array, default: () => [] },
});
const emit = defineEmits(['done']);
const router = useRouter();

const name = ref('');
const imageFile = ref(null);
const errormsg = ref(null);
const creating = ref(false);

function onPickerError(message) {
	errormsg.value = message;
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
		const chat = await createGroupRequest(name.value.trim(), props.members.map(m => m.id));

		if (imageFile.value) {
			await updateAvatar({ kind: 'group', chatId: chat.id }, imageFile.value);
		}

		refreshChats().catch(() => {});
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
				<ImagePicker v-model="imageFile" radius="33%" @error="onPickerError" />
				<div>
					<div class="text-muted small">JPEG, PNG or WebP, up to 5 MB</div>
					<button v-if="imageFile" type="button" class="btn btn-sm btn-link remove-btn" @click="imageFile = null">Remove</button>
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

.remove-btn {
	padding: 0;
	color: var(--tn-red);
	text-decoration: none;
}

.footer {
	flex-shrink: 0;
}

.create-btn {
	font-weight: 600;
}
</style>
