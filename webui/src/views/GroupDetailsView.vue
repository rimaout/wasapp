<script>
import ImagePicker from '../components/ui/ImagePicker.vue';
import MemberChips from '../components/ui/MemberChips.vue';
import IconButton from '../components/ui/IconButton.vue';
import ConfirmBar from '../components/ui/ConfirmBar.vue';
import { createGroup as createGroupRequest, updateAvatar } from '../services/api.js';
import { refreshChats } from '../composables/useChats.js';
import { navigateToChat } from '../services/chatNavigation.js';
import { getErrorMessage } from '../services/utils.js';
import { validateBaseName } from '../services/validators.js';

/**
 * GroupDetailsView component allows users to set the details for a new group chat, including the group name, image, and members.
 *
 * Used in: App.vue when the sidebar mode is set to 'group-details'.
 *
 * Props:
 *  - members (Array): The currently selected members for the group chat.
 *
 * Emits:
 *  - 'done' - emitted when the group creation process is completed successfully.
 *  - 'update:members' - emitted when the selected members are updated.
 *  - 'add-members' - emitted when the user wants to add more members to the group.
 *  - 'cancel' - emitted when the user cancels the group creation process.
 */
export default {
	components: { ImagePicker, MemberChips, IconButton, ConfirmBar },
	props: {
		members: { type: Array, default: () => [] },
	},
	emits: ['done', 'update:members', 'add-members', 'cancel'],
	data() {
		return {
			name: '',
			imageFile: null,
			errormsg: null,
			creating: false,
		};
	},
	methods: {
		removeMember(member) {
			this.$emit('update:members', this.members.filter(m => m.id !== member.id));
		},

		//Event handler for <ImagePicker>'s 'error' event.
		onPickerError(message) {
			this.errormsg = message;
		},

		validate() {
			return validateBaseName(this.name, 'Group name');
		},

		async createGroup() {
			// Validate the group name before proceeding
			const err = this.validate();
			if (err) {
				this.errormsg = err;
				return;
			}

			// Prevent multiple simultaneous group creation requests
			if (this.creating) return;
			this.creating = true;
			this.errormsg = null;

			try {
				// Create the group chat with the provided name and Members
				const chat = await createGroupRequest(this.name.trim(), this.members.map(m => m.id));

				// If an image file is selected, update the group's avatar
				if (this.imageFile) {
					await updateAvatar({ kind: 'group', chatId: chat.id }, this.imageFile);
				}

				// Refresh the chat list and navigate to the newly created group chat
				refreshChats().catch(() => {});
				navigateToChat(this.$router, chat);
				this.$emit('done');
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			} finally {
				this.creating = false;
			}
		},
	},
};
</script>

<template>
	<div class="group-details-view">
		<div class="body flex-grow-1 px-3 py-3">

			<!-- DETAILS SECTION -->
			<div class="section-box">
				<span class="section-title">Details</span>

				<!-- ImagePicker component allows users to select an image for the group avatar. -->
				<ImagePicker v-model:selected-image="imageFile" :size="120" @error="onPickerError" />
				<!--
				  @error="onPickerError" is Vue shorthand.
				  When ImagePicker calls $emit('error', 'File too large'), Vue automatically
				  passes that error text into onPickerError as the first argument.
				  (Equivalent to: @error="(msg) => onPickerError(msg)" or @error="onPickerError($event)")
				-->

				<!-- Remove image button appears only if an image is selected. -->
				<button v-if="imageFile" type="button" class="remove-image-btn" @click="imageFile = null">Remove image</button>

				<!-- Input field for the group name, bound to the 'name' data property. -->
				<input type="text" class="form-control name-input" v-model="name" maxlength="24" placeholder="Group name" />
			</div>

			<!-- MEMBERS SECTION -->
			<div class="section-box">
				<!-- Add Members Button -->
				<div class="section-title-row">
					<span class="section-title">Members</span>
					<IconButton icon="plus" size="small" filled @click="$emit('add-members')" />
				</div>

				<!-- Selected Members List -->
				<MemberChips :members="members" light @remove="removeMember" />
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="footer px-3 py-3">
			<ConfirmBar confirm-text="Create group" confirm-icon="check" :confirm-disabled="creating || !name.trim()" @cancel="$emit('cancel')" @confirm="createGroup" />
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
	background: var(--theme-bg-light);
	border: 1px solid var(--theme-border);
	border-radius: 12px;
	padding: 16px;
	margin-bottom: 16px;
}

.section-title {
	color: var(--theme-fg);
	font-weight: 600;
	font-size: 0.9rem;
}

.section-title-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.name-input {
	background: var(--theme-bg-highlight);
	border-color: var(--theme-border);
	color: var(--theme-fg);
}

.remove-image-btn {
	align-self: center;
	background: none;
	border: none;
	color: var(--theme-red);
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
