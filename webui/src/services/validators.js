// Shared input validators. These mirror the backend rules in
// service/api/helpers-validators.go so the client can give instant feedback.

// Validates a name against the BaseName rules used for usernames and group names:
// must be 3-24 characters, only alphanumeric + spaces/underscores/hyphens,
// and must contain at least one non-space character.
// @param {string} name  the value to validate
// @param {string} label the noun used in the error message (e.g. 'Username')
// @returns {string} '' when valid, otherwise a user-facing error message.
export function validateBaseName(name, label = 'Name') {
	const n = (name || '').trim();
	if (n.length < 3 || n.length > 24) {
		return `${label} must be 3-24 characters`;
	}
	if (!/^[A-Za-z0-9 _-]+$/.test(n)) {
		return `${label} can only contain letters, numbers, spaces, underscores and hyphens`;
	}
	if (!/\S/.test(n)) {
		return `${label} must contain at least one non-space character`;
	}
	return '';
}
