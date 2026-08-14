import { getUserName, getUserId } from './auth.js';

export function isMyMessage(senderName) {
	return senderName === getUserName();
}

export function isMyMessageById(senderId) {
	return senderId === getUserId();
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
	if (lm.isDeleted) return 'Message deleted';
	if (lm.content && lm.content.text) {
		let text = lm.content.text;
		return text.length > 50 ? text.substring(0, 47) + '...' : text;
	}
	if (lm.content && lm.content.msgImageId) return 'Photo';
	return '';
}
