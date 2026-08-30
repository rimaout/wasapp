import { ref } from 'vue';
import axios from '../services/axios.js';
import { getUserId } from '../services/auth.js';
import { getErrorMessage } from '../services/utils.js';

// Shared reactive store for the users list — a single source of truth.
// Every component imports the same module-level refs, so the list is fetched
// once and shared; opening another picker re-fetches for freshness but
// concurrent calls share a single in-flight request.

const users = ref([]);
const loading = ref(false);
const errormsg = ref(null);

let inFlight = null;

export function fetchUsers() {
	// A fetch is already running: piggyback on it (dedupes concurrent mounts).
	if (inFlight) return inFlight;

	loading.value = true;
	errormsg.value = null;
	inFlight = (async () => {
		try {
			const response = await axios.get('/users');
			users.value = response.data.usersList;
		} catch (e) {
			errormsg.value = getErrorMessage(e);
		} finally {
			loading.value = false;
			inFlight = null;
		}
	})();
	return inFlight;
}

// Returns all users (excluding the logged-in user), optionally filtered by
// query and by a list of additional user IDs to exclude.
export function filteredUsers(query, excludeIds) {
	const selfId = getUserId();
	let list = users.value.filter(u => u.id !== selfId);
	if (excludeIds) {
		const s = new Set(excludeIds);
		list = list.filter(u => !s.has(u.id));
	}
	if (!query) return list;
	const q = query.toLowerCase();
	return list.filter(u => u.name.toLowerCase().includes(q));
}

export function useUsers() {
	return { users, loading, errormsg, fetchUsers, filteredUsers };
}