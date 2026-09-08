<script>
import SearchBar from './ui/SearchBar.vue';
import ChatAvatar from './ui/ChatAvatar.vue';
import UserAvatar from './ui/UserAvatar.vue';
import ConfirmBar from './ui/ConfirmBar.vue';
import { useChats } from '../composables/useChats.js';
import { useUsers } from '../composables/useUsers.js';
import { getErrorMessage } from '../services/utils.js';

/**
 * ForwardActionScreen — a chat/user picker used to forward a message. It is meant
 * to be supplied as the `component` of an ActionMenu item. Selecting a chat
 * forwards there directly; selecting a user creates (or reuses) a direct chat
 * with them first, then forwards.
 *
 * @property {Object} message - the message being forwarded (needs `id`, `chatId`).
 * @property {string} [title='Forward message'] - heading shown at the top.
 */
export default {
	components: { SearchBar, ChatAvatar, UserAvatar, ConfirmBar },
	props: {
		message: { type: Object, required: true             },
		title:   { type: String, default: 'Forward message' },
	},
	/**
	 * Events:
	 *   done({ action:'forward' }) - fired after the message is forwarded to all targets.
	 *   cancel - fired when the user closes the picker.
	 */
	emits: ['done', 'cancel'],

	// Shared chats/users state comes from the composables (module-level refs);
	// expose it here so the rest of the Options API can use it via `this`.
	setup(props, { emit }) {
		const { chats, refreshChats       } = useChats();
		const { fetchUsers, filteredUsers } = useUsers();
		return { chats, refreshChats, fetchUsers, filteredUsers };
	},

	data() {
		return {
			query: '',
			selected: [],
			sending: false,
			errormsg: null,
		};
	},

	computed: {
		filteredChats() {
			let q = this.query.toLowerCase();
			return this.chats.filter(c => {
				if (c.id === this.message.chatId) return false;
				if (q && !c.displayName.toLowerCase().includes(q)) return false;
				return true;
			});
		},
		visibleUsers() {
			return this.filteredUsers(this.query);
		},
	},

	mounted() {
		this.refreshChats().catch(() => {});
		this.fetchUsers();
	},

	methods: {
		isSelected(kind, id) {
			return this.selected.some(s => s.kind === kind && s.id === id);
		},

		toggleChat(chat) {
			let list = this.selected.slice();
			let idx = list.findIndex(s => s.kind === 'chat' && s.id === chat.id);
			if (idx === -1) {
				list.push({ kind: 'chat', id: chat.id, name: chat.displayName, chat });
			} else {
				list.splice(idx, 1);
			}
			this.selected = list;
		},

		toggleUser(user) {
			let list = this.selected.slice();
			let idx = list.findIndex(s => s.kind === 'user' && s.id === user.id);
			if (idx === -1) {
				list.push({ kind: 'user', id: user.id, name: user.name, user });
			} else {
				list.splice(idx, 1);
			}
			this.selected = list;
		},

		removeItem(item) {
			this.selected = this.selected.filter(s => !(s.kind === item.kind && s.id === item.id));
		},

		async resolveChatId(item) {
			if (item.kind === 'chat') return item.id;
			try {
				const res = await this.$axios.post('/users/' + item.id + '/chats');
				return res.data.id;
			} catch (e) {
				if (e.response && e.response.status === 409 && e.response.data && e.response.data.chatId) {
					return e.response.data.chatId;
				}
				throw e;
			}
		},

		async forward() {
			if (this.sending || this.selected.length === 0) return;
			this.sending = true;
			this.errormsg = null;
			try {
				let url = '/chats/' + this.message.chatId + '/messages/' + this.message.id + '/forwards';
				for (const item of this.selected) {
					const chatId = await this.resolveChatId(item);
					await this.$axios.post(url, { forwardTo: chatId });
				}
				this.refreshChats().catch(() => {});
				this.$emit('done', { action: 'forward' });
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
<div class="action-screen forward">
<div class="action-title">{{ title }}</div>
<SearchBar v-model="query" placeholder="Search chats or users..." compact light />
<div v-if="selected.length > 0" class="px-2 py-1">
<div class="forward-chips">
<span v-for="s in selected" :key="s.kind + s.id" class="forward-chip">
<ChatAvatar v-if="s.kind === 'chat'" :chatId="s.id" :displayName="s.name" :size="24" :isGroup="s.chat.isGroupChat" />
<UserAvatar v-else :userId="s.id" :displayName="s.name" :size="24" />
<span class="forward-chip-name">{{ s.name }}</span>
<button type="button" class="forward-chip-remove" @click="removeItem(s)">x</button>
</span>
</div>
</div>
<ErrorMsg v-if="errormsg" :msg="errormsg" />
<div class="forward-list fade-bottom">
<template v-if="filteredChats.length > 0">
<div class="section-header">Recent chats</div>
<button v-for="chat in filteredChats" :key="chat.id" type="button" class="list-group-item list-group-item-action d-flex align-items-center forward-row" :class="{ selected: isSelected('chat', chat.id) }" @click="toggleChat(chat)">
<ChatAvatar :chatId="chat.id" :displayName="chat.displayName" :size="40" :isGroup="chat.isGroupChat" class="me-3" />
<span class="forward-name fw-semibold text-truncate flex-grow-1">{{ chat.displayName }}</span>
<svg v-if="isSelected('chat', chat.id)" class="feather check-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
</button>
</template>
<template v-if="visibleUsers.length > 0">
<div class="section-header">Users</div>
<button v-for="user in visibleUsers" :key="user.id" type="button" class="list-group-item list-group-item-action d-flex align-items-center forward-row" :class="{ selected: isSelected('user', user.id) }" @click="toggleUser(user)">
<UserAvatar :userId="user.id" :displayName="user.name" :size="40" class="me-3" />
<span class="forward-name fw-semibold text-truncate flex-grow-1">{{ user.name }}</span>
<svg v-if="isSelected('user', user.id)" class="feather check-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
</button>
</template>
<div v-if="filteredChats.length === 0 && visibleUsers.length === 0" class="text-muted text-center py-5">
<p class="mb-2 fs-5">No results</p>
</div>
</div>
<div class="forward-footer px-1 pt-0 pb-1">
<ConfirmBar confirm-text="Forward" confirm-icon="corner-up-right" :cancel-disabled="sending" :confirm-disabled="sending || selected.length === 0" cancel-label="Cancel" @cancel="$emit('cancel')" @confirm="forward" />
</div>
</div>
</template>

<style scoped>
.action-screen.forward { width: 280px; padding: 3px 4px 2px; display: flex; flex-direction: column; gap: 4px; max-height: 420px; }
.action-screen.forward :deep(.px-3.pb-2) { margin-top: 6px; padding-bottom: 0; }
.action-title { color: var(--theme-fg); font-weight: 600; font-size: 1rem; }
.forward-chips { display: flex; flex-wrap: wrap; gap: 8px; }
.forward-chip { display: inline-flex; align-items: center; gap: 6px; background: var(--theme-overlay); border-radius: 999px; padding: 3px 8px 3px 3px; color: var(--theme-fg); }
.forward-chip-name { font-size: 0.85rem; max-width: 120px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.forward-chip-remove { background: none; border: none; padding: 0; cursor: pointer; color: var(--theme-fg-dark); line-height: 1; font-size: 1rem; }
.forward-chip-remove:hover { color: var(--theme-red); }
.forward-list { flex: 1 1 auto; min-height: 0; overflow-y: auto; scrollbar-width: none; -ms-overflow-style: none; padding-bottom: 20px; }
.forward-list::-webkit-scrollbar { display: none; }
.section-header { position: sticky; top: 0; z-index: 1; background: var(--theme-bg-light); color: var(--theme-fg-dark); font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px; padding: 3px 10px 1px; }
.forward-row { border: 0 !important; border-radius: 12px !important; margin: 1px 2px; width: auto; padding: 6px 10px; background-color: transparent !important; color: var(--theme-fg); }
.forward-row:hover { background-color: var(--theme-overlay) !important; }
.forward-row.selected { background-color: var(--theme-overlay) !important; }
.forward-name { color: var(--theme-fg); }
.check-icon { width: 20px; height: 20px; color: var(--theme-blue); flex-shrink: 0; }
.forward-footer { flex-shrink: 0; }
</style>
