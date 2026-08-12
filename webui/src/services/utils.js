import { getUserName } from './auth.js';

export function isMyMessage(senderName) {
	return senderName === getUserName();
}

export function formatTime(isoString) {
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

export function getStatusColor(status) {
	return status === 'read' ? 'text-primary' : 'text-muted';
}

export function getChatSnippet(chat) {
	let lm = chat.lastMessage;
	if (lm.isInitMessage) {
		return chat.isGroupChat
			? 'Group created by ' + lm.senderName
			: lm.senderName + ' started a chat';
	}
	if (lm.isDeleted) return 'Message deleted';
	if (lm.content && lm.content.text) {
		let text = lm.content.text;
		return text.length > 50 ? text.substring(0, 47) + '...' : text;
	}
	if (lm.content && lm.content.msgImageId) return 'Photo';
	return '';
}
