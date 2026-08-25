<script>
import { formatMessageTime, isMyMessageById } from '../services/utils.js';
import { emojiGlyph } from '../services/emojis.js';
import UserAvatar from './ui/UserAvatar.vue';
import MessageImage from './MessageImage.vue';
import ActionMenu from './ui/ActionMenu.vue';
import ForwardActionScreen from './ForwardActionScreen.vue';
import ReactionActionScreen from './ReactionActionScreen.vue';
import MessageStatusIcon from './MessageStatusIcon.vue';

export default {
	components: { UserAvatar, MessageImage, ActionMenu, ForwardActionScreen, ReactionActionScreen, MessageStatusIcon },
	emits: ['action'],
	props: {
		message: { type: Object, required: true },
		isGroup: { type: Boolean, default: false },
		showAvatar: { type: Boolean, default: false },
		previousMessage: { type: Object, default: null },
	},
	data() {
		return {
			hovered: false,
		};
	},
	computed: {
		isMine() {
			return isMyMessageById(this.message.sender.id);
		},
		showSenderHeader() {
			if (!this.isGroup || this.isMine) return false;
			return !this.sameSeries;
		},
		hasAvatar() {
			return this.showSenderHeader && this.showAvatar;
		},
		isForwarded() {
			return !!this.message.forwardedFrom;
		},
		forwardedEmpty() {
			if (!this.isForwarded) return false;
			const f = this.message.forwardedFrom;
			return !f.text && !f.msgImageId;
		},
		hasImage() {
			return this.isForwarded
				? !!this.message.forwardedFrom.msgImageId
				: !!this.message.content?.msgImageId;
		},
		hasCaption() {
			return this.isForwarded
				? !!this.message.forwardedFrom.text
				: !!this.message.content?.text;
		},
		overlayMeta() {
			return this.hasImage && !this.hasCaption;
		},
		hasRepliedTo() {
			return !!this.message.repliedTo;
		},
		repliedTo() {
			return this.message.repliedTo;
		},
		repliedToDeleted() {
			const r = this.repliedTo;
			return !!r && !r.senderName;
		},
		repliedToEmpty() {
			const r = this.repliedTo;
			return !!r && !!r.senderName && !r.text && !r.msgImageId;
		},
		repliedToHasImage() {
			const r = this.repliedTo;
			return !!(r && r.msgImageId);
		},
		repliedToText() {
			const r = this.repliedTo;
			return (r && r.text) || '';
		},
		repliedToName() {
			const r = this.repliedTo;
			if (!r) return '';
			if (!r.senderName) return 'Replied to a deleted message';
			if (isMyMessageById(r.senderId)) return 'Replied to you';
			return 'Replied to ' + r.senderName;
		},
		reactionGroups() {
			const groups = new Map();
			for (const r of this.message.reactionsList || []) {
				let g = groups.get(r.emojiId);
				if (!g) {
					g = { emojiId: r.emojiId, count: 0, includesMe: false };
					groups.set(r.emojiId, g);
				}
				g.count++;
				if (isMyMessageById(r.userId)) g.includesMe = true;
			}
			return [...groups.values()].sort((a, b) => a.emojiId - b.emojiId);
		},
		sameSeries() {
			// true when the previous message is a normal message from the same
			// sender within the time gap. Used both to hide the sender header and
			// to tighten the gap between consecutive same-sender messages.
			if (!this.isNormalMessage(this.message)) return false;
			let prev = this.previousMessage;
			if (!this.isNormalMessage(prev)) return false;
			if (prev.sender.id !== this.message.sender.id) return false;
			return !this.timeGapExceeded(prev.sendTime, this.message.sendTime);
		},
		isContinuation() {
			return this.sameSeries;
		},
		actionItems() {
			const items = [
				{ id: 'react', label: 'React', icon: 'smile', component: ReactionActionScreen, props: { message: this.message } },
				{ id: 'reply', label: 'Reply', icon: 'corner-up-left' },
				{ id: 'forward', label: 'Forward', icon: 'share', component: ForwardActionScreen, props: { message: this.message } },
			];
			if (this.isMine) {
				items.push({ id: 'delete', label: 'Delete', icon: 'trash-2', danger: true, dangerText: 'Are you sure you want to delete this message?' });
			}
			return items;
		},
	},
	methods: {
		formatMessageTime,
		emojiGlyph,

		isNormalMessage(m) {
			return !!m && !m.isInitMessage && !m.isJoinMessage && !m.isLeaveMessage && !m.isDeleted;
		},

		timeGapExceeded(a, b) {
			let diff = Math.abs(new Date(b) - new Date(a));
			return diff > 5 * 60 * 1000;
		},

		initMessageText() {
			let who = this.isMine ? 'You' : this.message.sender.name;
			return this.isGroup ? 'Group created by ' + who : who + ' started this chat';
		},

		systemMessageText(action) {
			let who = this.isMine ? 'You' : this.message.sender.name;
			return who + ' ' + action;
		},

		onMenuSelect(item) {
			this.$emit('action', { type: item.id, message: this.message });
		},

		onActionDone(payload) {
			if (payload && payload.action === 'react' && payload.message) {
				this.$emit('action', { type: 'react', message: this.message, updatedMessage: payload.message });
			}
		},
	},
};
</script>

<template>
	<div class="message-wrapper" :class="[isMine ? 'me' : 'other', { continuation: isContinuation }]">
		<div v-if="message.isInitMessage" class="message-init text-muted">{{ initMessageText() }}</div>

		<div v-else-if="message.isJoinMessage" class="message-init text-muted">{{ systemMessageText('joined the group') }}</div>

		<div v-else-if="message.isLeaveMessage" class="message-init text-muted">{{ systemMessageText('left the group') }}</div>

		<div v-else-if="message.isDeleted" class="message-deleted text-muted fst-italic">
			{{ isMine ? 'You deleted this message' : 'Message deleted' }}
		</div>

		<div v-else class="message-row" @mouseenter="hovered = true" @mouseleave="hovered = false">
			<ActionMenu v-if="isMine" :items="actionItems" trigger-button-mode="hover" :trigger-button-visible="hovered" trigger-button-icon="more-vertical" trigger-button-size="small" trigger-button-filled placement="down-right" @select="onMenuSelect" @done="onActionDone" />

			<div class="message-bubble" :class="[isMine ? 'me' : 'other', { 'has-avatar': hasAvatar, 'has-image': hasImage }]">
				<div v-if="showSenderHeader" class="message-sender">
					<UserAvatar v-if="showAvatar" :userId="message.sender.id" :displayName="message.sender.name" :size="18" />
					<span class="message-sender-name">{{ message.sender.name }}</span>
				</div>

				<div v-if="hasRepliedTo" class="message-replied-to">
					<div class="replied-label">
						<svg class="feather replied-icon"><use href="/feather-sprite-v4.29.0.svg#corner-up-left"/></svg>
						<span class="replied-name">{{ repliedToName }}</span>
					</div>
					<template v-if="!repliedToDeleted">
						<MessageImage v-if="repliedToHasImage" :chat-id="message.chatId" :image-id="message.repliedTo.msgImageId" thumbnail />
						<span v-if="repliedToText" class="replied-text">{{ repliedToText }}</span>
						<span v-else-if="repliedToEmpty" class="replied-text replied-muted">(empty message)</span>
					</template>
				</div>

				<div v-if="isForwarded" class="message-forwarded">
					<svg class="feather forwarded-icon"><use href="/feather-sprite-v4.29.0.svg#corner-up-right"/></svg>
					<span>Forwarded</span>
				</div>

				<div v-if="forwardedEmpty" class="message-text message-forwarded-deleted">
					The original message was deleted
				</div>

				<div v-if="hasImage" class="message-image-wrap">
					<MessageImage :chat-id="message.chatId" :image-id="isForwarded ? message.forwardedFrom.msgImageId : message.content.msgImageId" />
					<div v-if="overlayMeta" class="message-meta message-meta-overlay">
						<span class="message-time">{{ formatMessageTime(message.sendTime) }}</span>
						<MessageStatusIcon v-if="isMine" :status="message.status" class="message-check" />
					</div>
				</div>

				<div v-if="hasCaption" class="message-text" :class="{ 'message-caption': hasImage }">
					{{ isForwarded ? message.forwardedFrom.text : message.content.text }}
					<span v-if="!overlayMeta" class="message-meta message-meta-inline">
						<span class="message-time">{{ formatMessageTime(message.sendTime) }}</span>
						<MessageStatusIcon v-if="isMine" :status="message.status" class="message-check" />
					</span>
				</div>
			</div>

			<ActionMenu v-if="!isMine" :items="actionItems" trigger-button-mode="hover" :trigger-button-visible="hovered" trigger-button-icon="more-vertical" trigger-button-size="small" trigger-button-filled placement="down-left" @select="onMenuSelect" @done="onActionDone" />

		</div>

		<div v-if="reactionGroups.length > 0" class="message-reactions">
			<span v-for="g in reactionGroups" :key="g.emojiId" class="reaction-chip" :class="{ active: g.includesMe }">
				<span class="reaction-glyph">{{ emojiGlyph(g.emojiId) }}</span>
				<span class="reaction-count">{{ g.count }}</span>
			</span>
		</div>
	</div>
</template>

<style scoped>
.message-wrapper {
	display: flex;
	flex-direction: column;
	align-items: flex-start;
}

.message-wrapper.me {
	align-items: flex-end;
}

.message-wrapper.continuation {
	margin-top: -8px;
}

.message-row {
	display: flex;
	align-items: center;
	gap: 6px;
	width: 100%;
	position: relative;
}

.message-wrapper.other .message-row {
	justify-content: flex-start;
}

.message-wrapper.me .message-row {
	justify-content: flex-end;
}

.message-init {
	font-size: 0.8rem;
	text-align: center;
	width: 100%;
	padding: 4px 0;
}

.message-deleted {
	font-size: 0.85rem;
	padding: 4px 0;
}

.message-bubble {
	max-width: 70%;
	padding: 8px 12px;
	border-radius: 14px;
	word-wrap: break-word;
}

.message-bubble.me {
	background-color: var(--tn-blue);
	color: var(--tn-bg-darker);
	border-bottom-right-radius: 4px;
}

.message-bubble.me .message-forwarded {
	color: var(--tn-bg-darker);
}

.message-bubble.other {
	background-color: var(--tn-bg-highlight);
	color: var(--tn-fg);
	border-bottom-left-radius: 4px;
}

.message-bubble.has-avatar {
	padding-left: 8px;
}

.message-bubble.has-avatar .message-sender {
	margin-bottom: 6px;
}

.message-bubble.has-image {
	padding: 3px;
	overflow: hidden;
}

.message-bubble.has-image .message-sender {
	padding: 4px 6px 2px;
	margin-bottom: 0;
}

.message-image-wrap {
	position: relative;
	border-radius: 10px;
	overflow: hidden;
}

.message-caption {
	padding: 6px 8px 4px;
}

.message-meta-overlay {
	position: absolute;
	right: 6px;
	bottom: 6px;
	margin: 0;
	padding: 2px 7px;
	border-radius: 999px;
	background: rgba(0, 0, 0, 0.55);
}

.message-meta-overlay .message-time {
	color: #fff;
	opacity: 0.9;
}

.message-forwarded {
	display: flex;
	align-items: center;
	gap: 4px;
	font-size: 0.75rem;
	font-weight: 600;
	color: var(--tn-fg-dark);
	margin-bottom: 3px;
}

.message-forwarded .forwarded-icon {
	width: 14px;
	height: 14px;
}

.message-forwarded-deleted {
	font-size: 0.85rem;
	opacity: 0.6;
	font-style: italic;
}

.message-replied-to {
	display: flex;
	flex-direction: column;
	align-items: flex-start;
	gap: 4px;
	margin: 4px 0 6px;
	padding: 6px 6px;
	border-radius: 8px;
	background: rgba(0, 0, 0, 0.12);
	min-width: 160px;
}

.replied-label {
	display: flex;
	align-items: center;
	gap: 4px;
}

.replied-icon {
	width: 14px;
	height: 14px;
}

.replied-name {
	font-size: 0.8rem;
	font-weight: 600;
	color: inherit;
}

.replied-text {
	font-size: 0.8rem;
	opacity: 0.85;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 100%;
}

.replied-muted {
	opacity: 0.5;
	font-style: italic;
}

.message-sender {
	display: flex;
	align-items: center;
	gap: 6px;
	margin-bottom: 2px;
}

.message-sender-name {
	font-size: 0.75rem;
	font-weight: 600;
	color: var(--tn-cyan);
}

.message-text {
	font-size: 0.95rem;
	white-space: pre-wrap;
}

.message-meta {
	display: flex;
	align-items: center;
	justify-content: flex-end;
	gap: 4px;
	margin-top: 2px;
}

.message-meta-inline {
	float: right;
	display: inline-flex;
	margin: 0.4em 0 2px 8px;
}

.message-time {
	font-size: 0.7rem;
	opacity: 0.7;
}

.message-check {
	color: inherit;
}

.message-check.is-received {
	opacity: 0.7;
}

.message-meta-overlay .message-check {
	color: #fff;
}

.message-meta-overlay .message-check.is-received {
	opacity: 0.7;
}

.message-reactions {
	display: flex;
	flex-wrap: wrap;
	gap: 4px;
	margin-top: 4px;
	padding: 0 6px;
}

.reaction-chip {
	display: inline-flex;
	align-items: center;
	gap: 4px;
	padding: 2px 8px;
	border-radius: 999px;
	background: var(--tn-bg-light);
	border: 1px solid var(--tn-border);
}

.reaction-chip.active {
	border-color: var(--tn-blue);
}

.reaction-glyph {
	font-size: 0.9rem;
	line-height: 1;
}

.reaction-count {
	color: var(--tn-fg-dark);
	font-weight: 600;
	font-size: 0.8rem;
}
</style>
