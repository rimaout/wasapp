<script>
import UserPicker from '../components/UserPicker.vue';
import ConfirmBar from '../components/ui/ConfirmBar.vue';

export default {
	components: { UserPicker, ConfirmBar },
	props: {
		members: { type: Array, default: () => [] },
	},
	emits: ['create', 'update:members', 'cancel'],

	methods: {
		create() {
			if (this.members.length === 0) return;
			this.$emit('create');
		},
	},
};
</script>

<template>
	<div class="group-members-view">
		<UserPicker :model-value="members" @update:model-value="v => $emit('update:members', v)" />

		<div class="footer px-3 pt-0 pb-3">
			<ConfirmBar confirm-text="Continue" confirm-icon="arrow-right" icon-right :confirm-disabled="members.length === 0" @cancel="$emit('cancel')" @confirm="create" />
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
