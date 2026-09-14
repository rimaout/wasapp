/**
 * Validates a base name according to the following rules:
 * - Must be between 3 and 24 characters long.
 * - Can only contain letters, numbers, spaces, underscores, and hyphens.
 * - Must contain at least one non-space character.
 *
 * Inputs:
 *  - name: The name to validate.
 *  - label: Optional. The label to use in the error message. Defaults to 'Name'.
 *
 * Returns:
 *  - An empty string if the name is valid.
 *  - An error message string if the name is invalid.
 */
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
