<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import axios from '../services/axios.js';
import MemberChips from './MemberChips.vue';
import UserPicker from './UserPicker.vue';
import ConfirmBar from './ConfirmBar.vue';
import { getErrorMessage } from '../services/utils.js';

// Button in the chat header showing the member count; opens a popup with the
// group member list and an "Add members" view (search + pick new members).
const props = defineProps({
	chatId: { type: String, required: true },
});
const emit = defineEmits(['add-members']);

const open = ref(false);
const view = ref('members'); // 'members' | 'add'
const members = ref([]);
const loading = ref(false);
const errormsg = ref(null);

const selected = ref([]);
const existingIds = computed(() => members.value.map(m => m.id));

const countLabel = computed(() => {
	if (loading.value) return 'Members';
	const n = members.value.length;
	return n === 1 ? '1 Member' : n + ' Members';
});

async function fetchMembers() {
	loading.value = true;
	errormsg.value = null;
	try {
		const res = await axios.get('/chats/' + props.chatId + '/members');
		members.value = res.data
			.filter(m => !m.leaveTime)
			.map(m => ({ id: m.userId, name: m.name }));
	} catch (e) {
		errormsg.value = getErrorMessage(e);
	}
	loading.value = false;
}

function toggle() {
	open.value = !open.value;
}

function close() {
	open.value = false;
	view.value = 'members';
}

function openAddView() {
	selected.value = [];
	view.value = 'add';
}

function confirmAdd() {
	emit('add-members', selected.value);
	view.value = 'members';
}

onMounted(fetchMembers);
watch(() => props.chatId, fetchMembers);
</script>

<template>
	<div class="members-popup">
		<button type="button" class="members-btn" @click="toggle">{{ countLabel }}</button>
		<template v-if="open">
			<div class="popup-backdrop" @click="close"></div>
			<div class="popup-panel">
				<template v-if="view === 'members'">
					<div class="panel-title">Group members</div>

					<ErrorMsg v-if="errormsg" :msg="errormsg" />
					<LoadingSpinner v-if="loading" />
					<MemberChips :members="members" :removable="false" light />

					<ConfirmBar confirm-text="Add members" confirm-icon="plus" cancel-label="Close" @cancel="close" @confirm="openAddView" />
				</template>

				<template v-else>
					<div class="panel-title">Add members</div>

					<UserPicker v-model="selected" :exclude-ids="existingIds" compact light class="user-picker-wrap" />

					<ConfirmBar confirm-text="Confirm" confirm-icon="check" cancel-label="Close" :confirm-disabled="selected.length === 0" @cancel="close" @confirm="confirmAdd" />
				</template>
			</div>
		</template>
	</div>
</template>

<style scoped>
.members-popup {
	position: relative;
	display: flex;
}

.members-btn {
	background: rgba(255, 255, 255, 0.08);
	border: none;
	color: var(--tn-fg-dark);
	font-size: 0.9rem;
	font-weight: 560;
	cursor: pointer;
	padding: 6px 14px;
	border-radius: 999px;
	flex-shrink: 0;
}

.members-btn:hover {
	background: rgba(255, 255, 255, 0.14);
	color: var(--tn-fg);
}

.popup-backdrop {
	position: fixed;
	inset: 0;
	z-index: 40;
}

.popup-panel {
	position: absolute;
	top: calc(100% + 6px);
	right: 0;
	width: 300px;
	max-height: 550px;
	background: var(--tn-bg-light);
	border: 1px solid var(--tn-border);
	border-radius: 12px;
	box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
	padding: 12px;
	z-index: 50;
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.user-picker-wrap {
	min-height: 0;
}

.panel-title {
	color: var(--tn-fg);
	font-weight: 600;
	font-size: 1rem;
}
</style>
