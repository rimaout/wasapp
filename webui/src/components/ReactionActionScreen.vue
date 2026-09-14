<script>
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
export default {
	props: {
		message: { type: Object, required: true },
	},
	/**
	 * Events:
	 *   done({ action:'react' }) - fired after the reaction is added/removed.
	 *   cancel - fired when the user closes the picker.
     *
     * Note: action:'react' is needed to distinguish this from other actions that may be emitted by the parent ActionMenu.
	 */
	emits: ['done', 'cancel'],

	data() {
		return {
			EMOJIS,
			sending: false,
			errormsg: null,
		};
	},

	computed: {
		myReactionId() {
			// Find the logged-in user's reaction to this message, if exists, and return its emojiId; otherwise return null.
			const me = getUserId();
			const found = (this.message.reactionsList || []).find(r => r.userId === me);
			return found ? found.emojiId : null;
		},
	},

	methods: {
		isActive(emojiId) {
			return this.myReactionId === emojiId;
		},

		// ---- ACTION HANDLERS ----
		async pick(emoji) {
			// Prevent multiple simultaneous requests
			if (this.sending) return;
			this.sending = true;
			this.errormsg = null;

			// Send the reaction request to the server: if the emoji is already active, remove it; otherwise, add it.
			try {
				const base = '/chats/' + this.message.chatId + '/messages/' + this.message.id + '/reactions';
				let res;
				if (this.isActive(emoji.id)) {
					res = await this.$axios.delete(base);
				} else {
					res = await this.$axios.post(base, { emojiId: emoji.id });
				}
				this.$emit('done', { action: 'react', message: res.data });
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			} finally {
				this.sending = false;
			}
		},
	},
};
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
				<span class="emoji-icon">{{ emoji.glyph }}</span>
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
	background-color: var(--theme-overlay);
}

.emoji-btn.active {
	background-color: rgba(122, 162, 247, 0.25);
}

.emoji-btn:disabled {
	opacity: 0.6;
	cursor: default;
}

.emoji-icon {
	font-size: 1.5rem;
	line-height: 1;
}
</style>
