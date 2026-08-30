<script setup>
import { ref, computed } from 'vue';
import axios from '../services/axios.js';
import { EMOJIS } from '../services/emojis.js';
import { getUserId } from '../services/auth.js';
import { getErrorMessage } from '../services/utils.js';

/**
 * ReactionActionScreen — an inline emoji picker shown in place of the popup menu.
 * Picking an emoji adds (or replaces) the logged-in user's reaction; picking the
 * emoji they already used removes it. Supply it as the `component` of an
 * ActionMenu item.
 *
 * @property {Object} message - the message being reacted to (needs `id`, `chatId`, `reactionsList`).
 */
const props = defineProps({
	message: { type: Object, required: true },
});
/**
 * Events:
 *   done({ action:'react' }) - fired after the reaction is added/removed.
 *   cancel - fired when the user closes the picker.
 */
const emit = defineEmits(['done', 'cancel']);

const sending = ref(false);
const errormsg = ref(null);

const myReactionId = computed(() => {
	const me = getUserId();
	const found = (props.message.reactionsList || []).find(r => r.userId === me);
	return found ? found.emojiId : null;
});

function isActive(emojiId) {
	return myReactionId.value === emojiId;
}

async function pick(emoji) {
	if (sending.value) return;
	sending.value = true;
	errormsg.value = null;
	try {
		const base = '/chats/' + props.message.chatId + '/messages/' + props.message.id + '/reactions';
		let res;
		if (isActive(emoji.id)) {
			res = await axios.delete(base);
		} else {
			res = await axios.post(base, { emojiId: emoji.id });
		}
		emit('done', { action: 'react', message: res.data });
	} catch (e) {
		errormsg.value = getErrorMessage(e);
	} finally {
		sending.value = false;
	}
}
</script>

<template>
	<div class="action-screen react">
		<div class="emoji-grid">
			<button
				v-for="emoji in EMOJIS"
				:key="emoji.id"
				type="button"
				class="emoji-btn"
				:class="{ active: isActive(emoji.id) }"
				:disabled="sending"
				@click="pick(emoji)"
			>
				<span class="emoji-glyph">{{ emoji.glyph }}</span>
			</button>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg" />
	</div>
</template>

<style scoped>
.action-screen.react {
	display: flex;
	flex-direction: column;
	gap: 8px;
	padding: 4px 8px 8px;
	width: 240px;
}

.emoji-grid {
	display: grid;
	grid-template-columns: repeat(5, 1fr);
	gap: 6px;
}

.emoji-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	height: 40px;
	border: none;
	border-radius: 10px;
	background: none;
	cursor: pointer;
}

.emoji-btn:hover {
	background-color: var(--tn-overlay);
}

.emoji-btn.active {
	background-color: rgba(122, 162, 247, 0.25);
}

.emoji-btn:disabled {
	opacity: 0.6;
	cursor: default;
}

.emoji-glyph {
	font-size: 1.5rem;
	line-height: 1;
}
</style>
