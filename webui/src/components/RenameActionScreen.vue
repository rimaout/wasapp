<script>
import ConfirmBar from './ui/ConfirmBar.vue';
import { updateName } from '../services/api.js';
import { getErrorMessage } from '../services/utils.js';
import { validateBaseName } from '../services/validators.js';

/**
 * RenameActionScreen — an inline panel shown in place of the popup menu, used to
 * rename a profile or group chat. Supply it as the `component` of an ActionMenu item.
 *
 * @property {string} [title=''] - heading shown at the top of the panel.
 * @property {string} [initialName=''] - current name, pre-filled in the input.
 * @property {Object} target - what is being renamed; either { kind:'me', userId }
 *   or { kind:'group', chatId }.
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
			const trimmed = this.name.trim();
			const current = (this.initialName || '').trim();
			return trimmed.length > 0 && trimmed !== current;
		},
	},

	methods: {
		async confirm() {
			if (this.busy || !this.isValid) return;

			// Validate the name and show the error without disabling the button
			const label = this.target.kind === 'group' ? 'Group name' : 'Username';
			const err = validateBaseName(this.name, label);
			if (err) {
				this.errormsg = err;
				return;
			}

			this.busy = true;
			this.errormsg = null;
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
		<div v-if="title" class="action-title">{{ title }}</div>

		<input type="text" class="form-control name-input" v-model="name" maxlength="24" placeholder="Enter name" @keyup.enter="confirm" />

		<ErrorMsg v-if="errormsg" :msg="errormsg" />

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
