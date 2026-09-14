import { ref } from 'vue';
import axios from '../services/axios.js';
import { getUserId } from '../services/auth.js';
import { getErrorMessage } from '../services/utils.js';

/**
 * Shared state for the list of users, loading status, and error messages.
 */

const users    = ref([]);
const loading  = ref(false);
const errormsg = ref(null);

let activeFetch = null; // Current active fetch request, if any. Used to avoid duplicate requests.

// Fetches the list of users from the server and updates the shared state.
export function fetchUsers() {
	// Return existing request if already fetching
	if (activeFetch) return activeFetch;

	loading.value = true;
	errormsg.value = null;
	activeFetch = (async () => {
		try {
			const response = await axios.get('/users');
			users.value = response.data.usersList;
		} catch (e) {
			errormsg.value = getErrorMessage(e);
		} finally {
			loading.value = false;
			activeFetch = null;
		}
	})();
	return activeFetch;
}

// Returns all users (excluding the logged-in user), optionally filtered by
// query and by a list of additional user IDs to exclude.
export function filteredUsers(query, excludeIds) {
	// Exclude the logged-in user from the list
	const selfId = getUserId();
	let list = users.value.filter(u => u.id !== selfId);

	// Exclude additional user IDs if provided
	if (excludeIds) {
		const s = new Set(excludeIds);
		list = list.filter(u => !s.has(u.id));
	}

	// Filter by query if provided
	if (!query) return list;
	const q = query.toLowerCase();
	return list.filter(u => u.name.toLowerCase().includes(q));
}

// Returns the shared state and functions for managing users.
export function useUsers() {
	return { users, loading, errormsg, fetchUsers, filteredUsers };
}
