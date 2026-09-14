<script>
import UserPicker from '../components/UserPicker.vue';
import ConfirmBar from '../components/ui/ConfirmBar.vue';

/**
 * GroupMembersView component allows users to select members for creating a group chat.
 *
 * Used in: App.vue when the sidebar mode is set to 'group-members'.
 *
 * Props: members (Array) - the currently selected members for the group chat.
 *
 * Emits:
 *  - 'create' - emitted when the user confirms the creation of the group chat.
 *  - 'update:members' - emitted when the selected members are updated.
 *  - 'cancel' - emitted when the user cancels the group creation process.
 */

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
		<!-- UserPicker component allows selecting users for the group chat. -->
		<UserPicker :user-list="members" @update:user-list="v => $emit('update:members', v)" />
		<!-- Data Flow:
		  - :user-list="members" passes the current selected members to the UserPicker component.
		  - @update:user-list="v => $emit('update:members', v)" listens for updates from the UserPicker
				emits an event to update the selected members in the parent component. v is the new array of selected members.

		  Note: the members prop is updated in the parent component (App.vue) and not from this component directly.
		-->

		<!-- ConfirmBar (bottom) component provides a confirmation button to create the group chat -->
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
