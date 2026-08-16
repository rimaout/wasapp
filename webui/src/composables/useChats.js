import { ref } from 'vue';
import axios from '../services/axios.js';
import { getErrorMessage } from '../services/utils.js';

// Shared reactive store for the chats list — a single source of truth.
// Every component imports the same module-level refs, so calling
// refreshChats() anywhere instantly updates the sidebar previews.

const chats = ref([]);
const errormsg = ref(null);
const version = ref(0);

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

export function useChats() {
	return { chats, errormsg, version, refreshChats };
}
