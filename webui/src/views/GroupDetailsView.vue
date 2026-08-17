<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import ImagePicker from '../components/ImagePicker.vue';
import MemberChips from '../components/MemberChips.vue';
import IconButton from '../components/IconButton.vue';
import ConfirmBar from '../components/ConfirmBar.vue';
import { createGroup as createGroupRequest, updateAvatar } from '../services/api.js';
import { refreshChats } from '../composables/useChats.js';
import { navigateToChat } from '../services/chatNavigation.js';
import { getErrorMessage } from '../services/utils.js';

const props = defineProps({
	members: { type: Array, default: () => [] },
});
const emit = defineEmits(['done', 'update:members', 'add-members', 'cancel']);

function removeMember(member) {
	emit('update:members', props.members.filter(m => m.id !== member.id));
}
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
			<div class="section-box">
				<span class="section-title">Info</span>
				<ImagePicker v-model="imageFile" :size="120" @error="onPickerError" />
				<button v-if="imageFile" type="button" class="remove-image-btn" @click="imageFile = null">Remove image</button>
				<input type="text" class="form-control name-input" v-model="name" maxlength="24" placeholder="Group name" />
			</div>

			<div class="section-box">
				<div class="section-title-row">
					<span class="section-title">Members</span>
					<IconButton icon="plus" size="small" label="Add members" @click="emit('add-members')" />
				</div>
				<MemberChips :members="members" @remove="removeMember" />
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="footer px-3 py-3">
			<ConfirmBar confirm-text="Create group" confirm-icon="check" :confirm-disabled="creating || !name.trim()" @cancel="emit('cancel')" @confirm="createGroup" />
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
	display: flex;
	flex-direction: column;
	overflow-y: auto;
	scrollbar-width: none;
	-ms-overflow-style: none;
}

.body::-webkit-scrollbar {
	display: none;
}

.section-box {
	display: flex;
	flex-direction: column;
	gap: 12px;
	background: var(--tn-bg-light);
	border: 1px solid var(--tn-border);
	border-radius: 12px;
	padding: 16px;
	margin-bottom: 16px;
}

.section-title {
	color: var(--tn-fg);
	font-weight: 600;
	font-size: 0.9rem;
}

.section-title-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.name-input {
	background: var(--tn-bg-highlight);
	border-color: var(--tn-border);
	color: var(--tn-fg);
}

.remove-image-btn {
	align-self: center;
	background: none;
	border: none;
	color: var(--tn-red);
	font-size: 0.85rem;
	cursor: pointer;
	padding: 0;
}

.remove-image-btn:hover {
	text-decoration: underline;
}

.footer {
	flex-shrink: 0;
}
</style>
