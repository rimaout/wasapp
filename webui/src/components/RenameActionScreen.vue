<script>
import ConfirmBar from './ui/ConfirmBar.vue';
import { updateName } from '../services/api.js';
import { getErrorMessage } from '../services/utils.js';
import { validateBaseName } from '../services/validators.js';

/**
 * RenameActionScreen component allows users to rename user or group by providing a new name.
 *
 * Used in:
 *  - ProfileBar.var to rename the current user.
	- ChatView.var to rename a group chat.
 *
 * Props:
 *  - title (String): The title displayed at the top of the rename screen.
 *  - initialName (String): The initial name of the target
 *  - target (Object): The target entity to be renamed, structure: { kind: 'me' | 'group', userId | chatId }.
 *
 * Emits:
 *  - 'done' - emitted after the name is successfully saved, with an object containing the action and new name.
 *  - 'cancel' - emitted when the user cancels the renaming process.
 */
export default {
	components: { ConfirmBar },
	props: {
		title:       { type: String, default: ''    },
		initialName: { type: String, default: ''    },
		target:      { type: Object, required: true },
	},
	/**
	 * Events:
	 *   done({ action:'rename', name }) - fired after the name is saved.
	 *   cancel - fired when the user cancels.
     *
	 * Note: action: 'rename' is neded to distinguish this event from other 'done' events in the parent component.
	 */
	emits: ['done', 'cancel'],

	data() {
		return {
			name: this.initialName || '',
			errormsg: null,
			busy: false,
		};
	},

	computed: {
		isValid() {
			// Ensure the name is not empty and has changed from the initial name
			const trimmed = this.name.trim();
			const current = (this.initialName || '').trim();
			return trimmed.length > 0 && trimmed !== current;
		},
	},

	methods: {
		// ---- ACTION HANDLERS ----
		// Confirm the renaming action, used when the user clicks the confirm button or presses enter
		async confirm() {
			if (this.busy || !this.isValid) return;

			// Validate the name and show the error
			const label = this.target.kind === 'group' ? 'Group name' : 'Username';
			const err = validateBaseName(this.name, label);
			if (err) {
				this.errormsg = err;
				return;
			}

			// Proceed to update the name
			this.busy = true;
			this.errormsg = null;

			// Call the API to update the name and emit the 'done' event on success
			try {
				const newName = await updateName(this.target, this.name.trim());
				this.$emit('done', { action: 'rename', name: newName });
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			} finally {
				this.busy = false;
			}
		},
	},
};
</script>

<template>
	<div class="action-screen rename">
		<!-- Display the title if provided -->
		<div v-if="title" class="action-title">{{ title }}</div>

		<!-- Input field for the new name  -->
		<input type="text" class="form-control name-input" v-model="name" maxlength="24" placeholder="Enter name" @keyup.enter="confirm" />

		<!-- Display error message if any -->
		<ErrorMsg v-if="errormsg" :msg="errormsg" />

		<!-- ConfirmBar component for confirming or canceling the action -->
		<ConfirmBar
			:cancel-disabled="busy"
			:confirm-disabled="busy || !isValid"
			confirm-text="Confirm"
			confirm-icon="check"
			@cancel="$emit('cancel')"
			@confirm="confirm"
		/>
	</div>
</template>

<style scoped>
.action-screen.rename {
	display: flex;
	flex-direction: column;
	gap: 12px;
	padding: 4px 8px 8px;
	width: 260px;
}

.action-title {
	color: var(--theme-fg);
	font-weight: 600;
	font-size: 1rem;
}

.name-input {
	background: var(--theme-bg-highlight);
	border-color: var(--theme-border);
	color: var(--theme-fg);
}
</style>
