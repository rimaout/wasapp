<script setup>
import { ref, computed } from 'vue';
import ConfirmBar from './ui/ConfirmBar.vue';
import { updateName } from '../services/api.js';
import { getErrorMessage } from '../services/utils.js';

/**
 * RenameActionScreen — an inline panel shown in place of the popup menu, used to
 * rename a profile or group chat. Supply it as the `component` of an ActionMenu item.
 *
 * @property {string} [title=''] - heading shown at the top of the panel.
 * @property {string} [initialName=''] - current name, pre-filled in the input.
 * @property {Object} target - what is being renamed; either { kind:'me', userId }
 *   or { kind:'group', chatId }.
 */
const props = defineProps({
	title: { type: String, default: '' },
	initialName: { type: String, default: '' },
	target: { type: Object, required: true },
});
/**
 * Events:
 *   done({ action:'rename', name }) - fired after the name is saved.
 *   cancel - fired when the user cancels.
 */
const emit = defineEmits(['done', 'cancel']);

const name = ref(props.initialName || '');
const errormsg = ref(null);
const busy = ref(false);

const isValid = computed(() => {
	const trimmed = name.value.trim();
	const current = (props.initialName || '').trim();
	return trimmed.length > 0 && trimmed !== current;
});

async function confirm() {
	if (busy.value || !isValid.value) return;
	busy.value = true;
	errormsg.value = null;
	try {
		const newName = await updateName(props.target, name.value.trim());
		emit('done', { action: 'rename', name: newName });
	} catch (e) {
		errormsg.value = getErrorMessage(e);
	} finally {
		busy.value = false;
	}
}
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
			@cancel="emit('cancel')"
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
