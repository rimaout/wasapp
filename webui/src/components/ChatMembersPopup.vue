<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import axios from '../services/axios.js';
import MemberChips from './MemberChips.vue';
import { getErrorMessage } from '../services/utils.js';

// Button in the chat header showing the member count; opens a popup with the
// group member list and (for now UI-only) Leave / Add members buttons.
const props = defineProps({
	chatId: { type: String, required: true },
});
const emit = defineEmits(['add-members']);

const open = ref(false);
const members = ref([]);
const loading = ref(false);
const errormsg = ref(null);

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

onMounted(fetchMembers);
watch(() => props.chatId, fetchMembers);
</script>

<template>
	<div class="members-popup">
		<button type="button" class="members-btn" @click="toggle">{{ countLabel }}</button>
		<template v-if="open">
			<div class="popup-backdrop" @click="open = false"></div>
			<div class="popup-panel">
				<div class="panel-title">Group members</div>

				<ErrorMsg v-if="errormsg" :msg="errormsg" />
				<LoadingSpinner v-if="loading" />
				<MemberChips :members="members" :removable="false" />

				<div class="panel-footer">
					<button type="button" class="panel-btn close" @click="open = false" aria-label="Close" title="Close">
						<svg class="feather close-icon"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
					</button>
					<button type="button" class="panel-btn add" @click="emit('add-members')">Add members</button>
				</div>
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
	width: 280px;
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

.panel-title {
	color: var(--tn-fg);
	font-weight: 600;
	font-size: 1rem;
}

.panel-footer {
	display: flex;
	gap: 8px;
}

.panel-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	height: 40px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-weight: 600;
	font-size: 0.9rem;
	flex-shrink: 0;
}

.panel-btn.close {
	width: 40px;
	background: rgba(247, 118, 142, 0.12);
	color: var(--tn-red);
}

.panel-btn.close:hover {
	background: rgba(247, 118, 142, 0.22);
}

.panel-btn.add {
	flex: 1 1 auto;
	background: rgba(158, 206, 106, 0.15);
	color: var(--tn-green);
}

.panel-btn.add:hover {
	background: rgba(158, 206, 106, 0.25);
}

.close-icon {
	width: 22px;
	height: 22px;
}
</style>
