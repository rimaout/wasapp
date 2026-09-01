import { ref } from 'vue';
import axios from '../services/axios.js';
import { getErrorMessage } from '../services/utils.js';

// Shared chat state. Call refreshChats() to sync across components.

const chats    = ref([]);
const errormsg = ref(null);
const version  = ref(0);

// Fetches the list of chats from the server and updates the shared state.
export async function refreshChats() {
	errormsg.value = null;
	try {
		const res = await axios.get('/chats');
		chats.value = res.data.chatsPreviewList;
		version.value++;
	} catch (e) {
		errormsg.value = getErrorMessage(e);
	}
}

// Returns the shared state and functions for managing chats.
export function useChats() {
	return { chats, errormsg, version, refreshChats };
}
