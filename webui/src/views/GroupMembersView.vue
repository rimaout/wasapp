<script setup>
import UserPicker from '../components/UserPicker.vue';
import ConfirmBar from '../components/ConfirmBar.vue';

const props = defineProps({
	members: { type: Array, default: () => [] },
});
const emit = defineEmits(['create', 'update:members', 'cancel']);

function create() {
	if (props.members.length === 0) return;
	emit('create');
}
</script>

<template>
	<div class="group-members-view">
		<UserPicker :model-value="members" @update:model-value="v => emit('update:members', v)" />

		<div class="footer px-3 py-3">
			<ConfirmBar confirm-text="Continue" confirm-icon="arrow-right" icon-right :confirm-disabled="members.length === 0" @cancel="emit('cancel')" @confirm="create" />
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
</style>
