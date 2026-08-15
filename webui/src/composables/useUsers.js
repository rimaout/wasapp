import { ref } from 'vue';
import axios from '../services/axios.js';
import { getUserId } from '../services/auth.js';
import { getErrorMessage } from '../services/utils.js';

export function useUsers() {
	const users = ref([]);
	const loading = ref(true);
	const errormsg = ref(null);

	async function fetchUsers() {
		loading.value = true;
		errormsg.value = null;
		try {
			const response = await axios.get('/users');
			users.value = response.data.usersList;
		} catch (e) {
			errormsg.value = getErrorMessage(e);
		}
		loading.value = false;
	}

	// Returns all users (excluding the logged-in user), optionally filtered by query.
	function filteredUsers(query) {
		const selfId = getUserId();
		let list = users.value.filter(u => u.id !== selfId);
		if (!query) return list;
		const q = query.toLowerCase();
		return list.filter(u => u.name.toLowerCase().includes(q));
	}

	return { users, loading, errormsg, fetchUsers, filteredUsers };
}
