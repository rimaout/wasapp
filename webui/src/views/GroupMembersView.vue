<script setup>
import UserPicker from '../components/UserPicker.vue';

const props = defineProps({
	members: { type: Array, default: () => [] },
});
const emit = defineEmits(['create', 'update:members']);

function create() {
	if (props.members.length === 0) return;
	emit('create');
}
</script>

<template>
	<div class="group-members-view">
		<UserPicker :model-value="members" @update:model-value="v => emit('update:members', v)" />

		<div class="footer px-3 py-3">
			<button type="button" class="create-btn" :disabled="members.length === 0" @click="create">
				Continue
			</button>
		</div>
	</div>
</template>

<style scoped>
.group-members-view {
	display: flex;
	flex-direction: column;
	height: 100%;
}

.footer {
	flex-shrink: 0;
}

.create-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 100%;
	height: 40px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	background: rgba(158, 206, 106, 0.15);
	color: var(--tn-green);
	font-weight: 600;
}

.create-btn:hover:not(:disabled) {
	background: rgba(158, 206, 106, 0.25);
}

.create-btn:disabled {
	background: rgba(255, 255, 255, 0.08);
	color: var(--tn-fg-dark);
	cursor: default;
}
</style>
