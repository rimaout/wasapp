<script>
import MemberChips from './ui/MemberChips.vue';
import UserPicker from './UserPicker.vue';
import ConfirmBar from './ui/ConfirmBar.vue';
import { getErrorMessage } from '../services/utils.js';

// Button in the chat header showing the member count; opens a popup with the
// group member list and an "Add members" view (search + pick new members).
export default {
	components: { MemberChips, UserPicker, ConfirmBar },
	props: {
		chatId: { type: String, required: true },
	},

	data() {
		return {
			open: false,
			view: 'members', // 'members' | 'add'
			members: [],
			loading: false,
			errormsg: null,
			selected: [],
		};
	},

	computed: {
		existingIds() {
			return this.members.map(m => m.id);
		},
		countLabel() {
			if (this.loading) return 'Members';
			const n = this.members.length;
			return n === 1 ? '1 Member' : n + ' Members';
		},
	},

	watch: {
		// Reload when switching to another chat
		chatId: 'fetchMembers',
	},

	mounted() {
		this.fetchMembers();
	},

	methods: {
		async fetchMembers() {
			this.loading = true;
			this.errormsg = null;
			try {
				const res = await this.$axios.get('/chats/' + this.chatId + '/members');
				this.members = res.data.membersList
					.filter(m => !m.leaveTime)
					.map(m => ({ id: m.userId, name: m.name }));
			} catch (e) {
				this.errormsg = getErrorMessage(e);
			}
			this.loading = false;
		},

		toggle() {
			this.open = !this.open;
		},

		close() {
			this.open = false;
			this.view = 'members';
		},

		openAddView() {
			this.selected = [];
			this.view = 'add';
		},

		async confirmAdd() {
			this.errormsg = null;
			for (const user of this.selected) {
				try {
					await this.$axios.post('/chats/' + this.chatId + '/members', { userId: user.id });
				} catch (e) {
					this.errormsg = getErrorMessage(e);
					return;
				}
			}
			this.selected = [];
			await this.fetchMembers();
			this.close();
		},
	},
};
</script>

<template>
	<div class="members-popup">
		<button type="button" class="members-btn" @click="toggle">{{ countLabel }}</button>
		<template v-if="open">
			<div class="popup-backdrop" @click="close"></div>
			<div class="popup-panel">
				<template v-if="view === 'members'">
					<div class="panel-title">Group members</div>

					<ErrorMsg v-if="errormsg" :msg="errormsg" />
					<LoadingSpinner v-if="loading" />
					<MemberChips :members="members" :removable="false" light />

					<ConfirmBar confirm-text="Add members" confirm-icon="plus" cancel-label="Close" @cancel="close" @confirm="openAddView" />
				</template>

				<template v-else>
					<div class="panel-title">Add members</div>

					<UserPicker v-model="selected" :exclude-ids="existingIds" compact light class="user-picker-wrap" />

					<ConfirmBar class="confirm-bar-tight" confirm-text="Confirm" confirm-icon="check" cancel-label="Close" :confirm-disabled="selected.length === 0" @cancel="close" @confirm="confirmAdd" />
				</template>
			</div>
		</template>
	</div>
</template>

<style scoped>
.members-popup {
	position: relative;
	display: flex;
}

.members-btn {
	background: var(--theme-overlay);
	border: none;
	color: var(--theme-fg-dark);
	font-size: 0.9rem;
	font-weight: 560;
	cursor: pointer;
	padding: 6px 14px;
	border-radius: 999px;
	flex-shrink: 0;
}

.members-btn:hover {
	background: rgba(255, 255, 255, 0.14);
	color: var(--theme-fg);
}

.popup-backdrop {
	position: fixed;
	inset: 0;
	z-index: 40;
}

.popup-panel {
	position: absolute;
	top: calc(100% + 6px);
	right: 0;
	width: 300px;
	max-height: 550px;
	background: var(--theme-bg-light);
	border: 1px solid var(--theme-border);
	border-radius: 12px;
	box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
	padding: 8px 12px 12px;
	z-index: 50;
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.user-picker-wrap {
	min-height: 0;
}

.confirm-bar-tight {
	margin-top: -12px;
}

.panel-title {
	color: var(--theme-fg);
	font-weight: 600;
	font-size: 1rem;
}
</style>
