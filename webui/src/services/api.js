import axios from './axios.js';

// Centralized API calls. `target` is `{ kind: 'me' }` for the logged-in user,
// or `{ kind: 'group', chatId }` for a group chat.

function avatarUrl(target) {
	return target.kind === 'me'
		? '/users/' + target.userId + '/avatar'
		: '/chats/' + target.chatId + '/avatar';
}

// Fetch an avatar as a Blob (throws on 404 when no avatar is set).
export async function fetchAvatar(target) {
	const res = await axios.get(avatarUrl(target), { responseType: 'blob' });
	return res.data;
}

// Upload or replace an avatar (multipart `binaryImage`).
export function updateAvatar(target, file) {
	const formData = new FormData();
	formData.append('binaryImage', file);
	return target.kind === 'me'
		? axios.put('/me/avatar', formData)
		: axios.put('/chats/' + target.chatId + '/avatar', formData);
}

// Remove the current avatar.
export function deleteAvatar(target) {
	return target.kind === 'me'
		? axios.delete('/me/avatar')
		: axios.delete('/chats/' + target.chatId + '/avatar');
}

// Rename and return the new name string.
export async function updateName(target, name) {
	if (target.kind === 'me') {
		const res = await axios.patch('/me/name', { userName: name });
		return res.data.name;
	}
	const res = await axios.patch('/chats/' + target.chatId + '/name', { groupName: name });
	return res.data.groupName;
}

// Create a group chat and return the chat object.
export async function createGroup(groupName, membersList) {
	const res = await axios.post('/chats', { groupName, membersList });
	return res.data;
}
