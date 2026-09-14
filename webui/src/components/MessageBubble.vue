<script>
import { formatMessageTime, isMyMessageById } from '../services/utils.js';
import { getEmoji } from '../services/emojis.js';
import UserAvatar from './ui/UserAvatar.vue';
import MessageImage from './MessageImage.vue';
import ActionMenu from './ui/ActionMenu.vue';
import ForwardActionScreen from './ForwardActionScreen.vue';
import ReactionActionScreen from './ReactionActionScreen.vue';
import MessageStatusIcon from './MessageStatusIcon.vue';

/**
 * MessageBubble component - displays a single chat message with optional features like
 * sender header, avatar, image, caption, forwarded label, replied-to section, and reactions.
 *
 * Used in: ChatView.vue to render each message in the chat conversation.
 *
 * Props:
 *   - message:         The message object to display.
 *   - isGroup:         Boolean indicating if the chat is a group chat.
 *   - showAvatar:      Boolean indicating if the sender's avatar should be shown.
 *   - previousMessage: The previous message object for determining series continuity.
 *
 * Emits:
 *   - action: Emitted when an action is performed on the message (react, reply, forward, delete).
 */
export default {
	components: { UserAvatar, MessageImage, ActionMenu, MessageStatusIcon },
	props: {
		message:         { type: Object,  required: true },
		isGroup:         { type: Boolean, default: false },
		showAvatar:      { type: Boolean, default: false },
		previousMessage: { type: Object,  default: null  },
	},

	/**
	 * Events:
	 *   action({ type, message, updatedMessage }) - fired when an action is performed on the message.
	 *     - type: The type of action ('react', 'reply', 'forward', 'delete').
	 *     - message: The original message object.
	 *     - updatedMessage: (Optional) The updated message object after the action (example: after reacting).
	 */
	emits: ['action'],
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
			// Show the sender header (name + avatar) for group chats, unless it's a continuation of the same sender's messages
			if (!this.isGroup || this.isMine) return false;
			return !this.isContinuation;
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
		hasText() {
			// Determine if the message has a text to display
			return this.isForwarded
				? !!this.message.forwardedFrom.text
				: !!this.message.content?.text;
		},
		overlayTimeStatus() {
			// Determine if the message time and status should be overlaid on the image (if it has an image but no text)
			return this.hasImage && !this.hasText;
		},
		hasRepliedTo() {
			// Determine if this message is a reply to another message
			return !!this.message.repliedTo;
		},
		repliedTo() {
			// Get the replied-to message object, or null if it doesn't exist
			return this.message.repliedTo;
		},
		repliedToDeleted() {
			// Determine if the replied-to message was deleted (no sender name)
			const r = this.repliedTo;
			return !!r && !r.senderName;
		},
		repliedToEmpty() {
			// Determine if the replied-to message is empty (no text, no image)
			const r = this.repliedTo;
			return !!r && !!r.senderName && !r.text && !r.msgImageId;
		},
		repliedToHasImage() {
			// Determine if the replied-to message has an image
			const r = this.repliedTo;
			return !!(r && r.msgImageId);
		},
		repliedToText() {
			// Get the text of the replied-to message, or an empty string if it doesn't exist
			const r = this.repliedTo;
			return (r && r.text) || '';
		},
		repliedToName() {
			// Generate the "Replied to ..." label based on the replied-to message's sender
			const r = this.repliedTo;
			if (!r) return '';
			if (!r.senderName) return 'Replied to a deleted message';
			if (isMyMessageById(r.senderId)) return 'Replied to you';
			return 'Replied to ' + r.senderName;
		},

		reactionGroups() {
			const groups = new Map();
			this.message.reactionsList?.forEach(r => {
				let g = groups.get(r.emojiId); // Find existing group for this emoji
				// If no group exists, create a new one
				if (!g) {
					g = { emojiId: r.emojiId, count: 0, includesMe: false };
					groups.set(r.emojiId, g);
				}
				// Increment the count and check if the current user reacted with this emoji
				g.count++;
				if (isMyMessageById(r.userId)) g.includesMe = true;
			});
			// Return the groups as a sorted arrray of objects (emojiId, count, includesMe), sorted by emojiId
			return [...groups.values()].sort((a, b) => a.emojiId - b.emojiId);
		},
		isContinuation() {
			// Determine if this message is part of the same series as the previous message (same sender, within 5 minutes)
			if (!this.isNormalMessage(this.message)) return false;				// Only normal messages can be part of a series
			let prev = this.previousMessage;
			if (!this.isNormalMessage(prev)) return false;			            // Previous message must also be a normal message
			if (prev.sender.id !== this.message.sender.id) return false;        // Must be the same sender
			return !this.timeGapExceeded(prev.sendTime, this.message.sendTime); // Must be within 5 minutes
		},
		actionItems() {
			// Define the action items for the action menu (React, Reply, Forward, Delete)
			// Note: Delete is only available for the sender's own messages
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
		formatMessageTime, // Generate correct time string, ex: "now, 5m, 12.34"
		getEmoji,

		isNormalMessage(m) {
			return !!m && !m.isInitMessage && !m.isJoinMessage && !m.isLeaveMessage && !m.isDeleted;
		},

		timeGapExceeded(a, b) {
			let diff = Math.abs(new Date(b) - new Date(a));
			return diff > 5 * 60 * 1000; // 5 minutes in milliseconds
		},

		initMessageText() {
			let who = this.isMine ? 'You' : this.message.sender.name;
			return this.isGroup ? 'Group created by ' + who : who + ' started this chat';
		},

		systemMessageText(action) {
			let who = this.isMine ? 'You' : this.message.sender.name;
			return who + ' ' + action;
		},

		// ---- ACTION HANDLERS ----

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

		<!-- Display system messages for initialization, join, leave -->
		<div v-if="message.isInitMessage" class="message-init text-muted">{{ initMessageText() }}</div>
		<div v-else-if="message.isJoinMessage" class="message-init text-muted">{{ systemMessageText('joined the group') }}</div>
		<div v-else-if="message.isLeaveMessage" class="message-init text-muted">{{ systemMessageText('left the group') }}</div>

		<!-- Display deleted message-->
		<div v-else-if="message.isDeleted" class="message-deleted text-muted fst-italic">
			{{ isMine ? 'You deleted this message' : 'Message deleted' }}
		</div>

		<!-- Display "normal" messages (base, forwarded, replied-to) -->
		<div v-else class="message-row" @mouseenter="hovered = true" @mouseleave="hovered = false">

			<!-- Action menu for the sender's own messages (show on hover) -->
			<ActionMenu v-if="isMine" :items="actionItems" trigger-button-mode="hover" :trigger-button-visible="hovered" trigger-button-icon="more-vertical" trigger-button-size="small" trigger-button-filled placement="down-right" @select="onMenuSelect" @done="onActionDone" />

			<div class="message-bubble" :class="[isMine ? 'me' : 'other', { 'has-avatar': hasAvatar, 'has-image': hasImage }]">

				<!-- Sender header (avatar + name) for group chats -->
				<div v-if="showSenderHeader" class="message-sender">
					<UserAvatar v-if="showAvatar" :userId="message.sender.id" :displayName="message.sender.name" :size="18" />
					<span class="message-sender-name">{{ message.sender.name }}</span>
				</div>

				<!-- Is Reply Message -->
				<div v-if="hasRepliedTo" class="message-replied-to">

					<!-- Reply label and name -->
					<div class="replied-label">
						<svg class="feather replied-icon"><use href="/feather-sprite-v4.29.0.svg#corner-up-left"/></svg>
						<span class="replied-name">{{ repliedToName }}</span>
					</div>

					<!-- Reply content (image or text or deleted) -->
					<template v-if="!repliedToDeleted">
						<MessageImage v-if="repliedToHasImage" :chat-id="message.chatId" :image-id="message.repliedTo.msgImageId" thumbnail />
						<span v-if="repliedToText" class="replied-text">{{ repliedToText }}</span>
						<span v-else-if="repliedToEmpty" class="replied-text replied-muted">(empty message)</span>
					</template>
				</div>

				<!-- Is Forwarded Message -->
				<div v-if="isForwarded" class="message-forwarded">
					<svg class="feather forwarded-icon"><use href="/feather-sprite-v4.29.0.svg#corner-up-right"/></svg>
					<span>Forwarded</span>
				</div>
				<!-- Forwarded message deleted -->
				<div v-if="forwardedEmpty" class="message-text message-forwarded-deleted">
					The original message was deleted
				</div>
				<!-- Forwarded message image -->
				<div v-if="hasImage" class="message-image-wrap">
					<MessageImage :chat-id="message.chatId" :image-id="isForwarded ? message.forwardedFrom.msgImageId : message.content.msgImageId" />
					<div v-if="overlayTimeStatus" class="message-time-status message-time-status-overlay">
						<span class="message-time">{{ formatMessageTime(message.sendTime) }}</span>
						<MessageStatusIcon v-if="isMine" :status="message.status" class="message-check" />
					</div>
				</div>
				<!-- Forwarded message text -->
				<div v-if="hasText" class="message-text" :class="{ 'message-caption': hasImage }">
					{{ isForwarded ? message.forwardedFrom.text : message.content.text }}

					<!-- Message time and status -->
					<span v-if="!overlayTimeStatus" class="message-time-status message-time-status-inline">
						<span class="message-time">{{ formatMessageTime(message.sendTime) }}</span>
						<MessageStatusIcon v-if="isMine" :status="message.status" class="message-check" />
					</span>
				</div>
			</div>

			<!-- Action menu for other users' messages (show on hover) -->
			<ActionMenu v-if="!isMine" :items="actionItems" trigger-button-mode="hover" :trigger-button-visible="hovered" trigger-button-icon="more-vertical" trigger-button-size="small" trigger-button-filled placement="down-left" @select="onMenuSelect" @done="onActionDone" />
		</div>

		<!-- Reactions display -->
		<div v-if="reactionGroups.length > 0" class="message-reactions">
			<span v-for="g in reactionGroups" :key="g.emojiId" class="reaction-chip" :class="{ active: g.includesMe }">
				<span class="reaction-emoji">{{ getEmoji(g.emojiId) }}</span>
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
	background-color: var(--theme-blue);
	color: var(--theme-bg-darker);
	border-bottom-right-radius: 4px;
}

.message-bubble.me .message-forwarded {
	color: var(--theme-bg-darker);
}

.message-bubble.other {
	background-color: var(--theme-bg-highlight);
	color: var(--theme-fg);
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

.message-time-status-overlay {
	position: absolute;
	right: 6px;
	bottom: 6px;
	margin: 0;
	padding: 2px 7px;
	border-radius: 999px;
	background: rgba(0, 0, 0, 0.55);
}

.message-time-status-overlay .message-time {
	color: #fff;
	opacity: 0.9;
}

.message-forwarded {
	display: flex;
	align-items: center;
	gap: 4px;
	font-size: 0.75rem;
	font-weight: 600;
	color: var(--theme-fg-dark);
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
	color: var(--theme-cyan);
}

.message-text {
	font-size: 0.95rem;
	white-space: pre-wrap;
}

.message-time-status {
	display: flex;
	align-items: center;
	justify-content: flex-end;
	gap: 4px;
	margin-top: 2px;
}

.message-time-status-inline {
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

.message-check:not(.is-read) {
	opacity: 0.45;
}

.message-time-status-overlay .message-check {
	color: #fff;
}

.message-time-status-overlay .message-check:not(.is-read) {
	opacity: 0.7;
}

.message-time-status-overlay .message-check.is-read {
	color: var(--theme-cyan);
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
	background: var(--theme-bg-light);
	border: 1px solid var(--theme-border);
}

.reaction-chip.active {
	border-color: var(--theme-blue);
}

.reaction-emoji {
	font-size: 0.9rem;
	line-height: 1;
}

.reaction-count {
	color: var(--theme-fg-dark);
	font-weight: 600;
	font-size: 0.8rem;
}
</style>
