<script>
import ImagePicker from '../components/ui/ImagePicker.vue';
import MemberChips from '../components/ui/MemberChips.vue';
import IconButton from '../components/ui/IconButton.vue';
import ConfirmBar from '../components/ui/ConfirmBar.vue';
import { createGroup as createGroupRequest, updateAvatar } from '../services/api.js';
import { refreshChats } from '../composables/useChats.js';
import { navigateToChat } from '../services/chatNavigation.js';
import { getErrorMessage } from '../services/utils.js';

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

		onPickerError(message) {
			this.errormsg = message;
		},

		validate() {
			const n = this.name.trim();
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
		},

		async createGroup() {
			const err = this.validate();
			if (err) {
				this.errormsg = err;
				return;
			}
			if (this.creating) return;

			this.creating = true;
			this.errormsg = null;
			try {
				const chat = await createGroupRequest(this.name.trim(), this.members.map(m => m.id));

				if (this.imageFile) {
					await updateAvatar({ kind: 'group', chatId: chat.id }, this.imageFile);
				}

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
			<div class="section-box">
				<span class="section-title">Details</span>
				<ImagePicker v-model="imageFile" :size="120" @error="onPickerError" />
				<button v-if="imageFile" type="button" class="remove-image-btn" @click="imageFile = null">Remove image</button>
				<input type="text" class="form-control name-input" v-model="name" maxlength="24" placeholder="Group name" />
			</div>

			<div class="section-box">
				<div class="section-title-row">
					<span class="section-title">Members</span>
					<IconButton icon="plus" size="small" filled @click="$emit('add-members')" />
				</div>
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
