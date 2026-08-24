import { getUserName, getUserId } from './auth.js';

const AVATAR_COLORS = [ //Tokyo Night color palette
	'#7aa2f7', '#bb9af7', '#9ece6a', '#e0af68', '#f7768e',
	'#7dcfff', '#ff9e64', '#c0caf5', '#565f89', '#414868',
];

// Generate a color based on the name string
export function getAvatarColor(name) {
	let hash = 0;
	for (let i = 0; i < name.length; i++) {
		hash = name.charCodeAt(i) + ((hash << 5) - hash);
	}
	return AVATAR_COLORS[Math.abs(hash) % AVATAR_COLORS.length];
}

// Get the first letter of the name, or '?' if name is empty
export function getAvatarLetter(name) {
	if (!name || name.length === 0) return '?';
	return name[0].toUpperCase();
}

// Extract a human-readable message from an axios error.
export function getErrorMessage(e) {
	return e.response?.data?.message || e.toString();
}

export function isMyMessage(senderName) {
	return senderName === getUserName();
}

export function isMyMessageById(senderId) {
	return senderId === getUserId();
}

export function formatPreviewTime(isoString) {
	let date = new Date(isoString);
	let now = new Date();
	let diffMins = Math.floor((now - date) / 60000);

	if (diffMins < 1) return 'now';
	if (diffMins < 60) return diffMins + 'm';

	let today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
	let msgDay = new Date(date.getFullYear(), date.getMonth(), date.getDate());
	let diffDays = Math.floor((today - msgDay) / 86400000);

	if (diffDays === 0) {
		return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
	}
	if (diffDays === 1) return 'Yesterday';

	if (date.getFullYear() === now.getFullYear()) {
		return date.toLocaleDateString([], { day: 'numeric', month: 'short' });
	}
	return date.toLocaleDateString([], { day: 'numeric', month: 'short', year: 'numeric' });
}

export function getStatusIcon(status) {
	return status === 'delivered' ? '✓' : '✓✓';
}

export function formatMessageTime(isoString) {
	let date = new Date(isoString);
	let now = new Date();
	let diffMins = Math.floor((now - date) / 60000);

	if (diffMins < 1) return 'now';
	if (diffMins < 60) return diffMins + 'm';

	return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

export function getStatusColor(status) {
	return status === 'read' ? 'text-primary' : 'text-muted';
}

// Returns a formatted string for the day of the message, e.g. "Today", "Yesterday", "March 5", "March 5, 2023"
export function formatDay(isoString) {
	let date = new Date(isoString);
	let now = new Date();
	let today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
	let msgDay = new Date(date.getFullYear(), date.getMonth(), date.getDate());
	let diffDays = Math.floor((today - msgDay) / 86400000);

	if (diffDays === 0) return 'Today';
	if (diffDays === 1) return 'Yesterday';
	if (date.getFullYear() === now.getFullYear()) {
		return date.toLocaleDateString([], { day: 'numeric', month: 'long' });
	}
	return date.toLocaleDateString([], { day: 'numeric', month: 'long', year: 'numeric' });
}

// Returns true if the two messages are on different days
export function isNewDay(prevMsg, currMsg) {
	let prev = new Date(prevMsg.sendTime);
	let curr = new Date(currMsg.sendTime);
	return prev.getFullYear() !== curr.getFullYear()
		|| prev.getMonth() !== curr.getMonth()
		|| prev.getDate() !== curr.getDate();
}

export function getChatSnippet(chat) {
	let lm = chat.lastMessage;
	if (lm.isInitMessage) {
		let who = isMyMessage(lm.senderName) ? 'You' : lm.senderName;
		return chat.isGroupChat
			? 'Group created by ' + who
			: who + ' started this chat';
	}
	if (lm.isJoinMessage) {
		let who = isMyMessage(lm.senderName) ? 'You' : lm.senderName;
		return who + ' joined the group';
	}
	if (lm.isLeaveMessage) {
		let who = isMyMessage(lm.senderName) ? 'You' : lm.senderName;
		return who + ' left the group';
	}
	if (lm.isForward) return 'Forwarded message';
	if (lm.isDeleted) return 'Message deleted';
	if (lm.content && lm.content.text) {
		let text = lm.content.text;
		return text.length > 50 ? text.substring(0, 47) + '...' : text;
	}
	if (lm.content && lm.content.msgImageId) return 'Photo';
	return '';
}
