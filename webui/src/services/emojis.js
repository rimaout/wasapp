// Shared list of emoji reactions. The index (id) is the `emojiId` used by the
// reactions API, as defined in the OpenAPI spec (0-9).
export const EMOJIS = [
	{ id: 0, glyph: '👍' },
	{ id: 1, glyph: '❤️' },
	{ id: 2, glyph: '😂' },
	{ id: 3, glyph: '😮' },
	{ id: 4, glyph: '😢' },
	{ id: 5, glyph: '🤔' },
	{ id: 6, glyph: '👏' },
	{ id: 7, glyph: '🙌' },
	{ id: 8, glyph: '😡' },
	{ id: 9, glyph: '🤯' },
];

export function getEmoji(emojiId) {
	const e = EMOJIS.find(e => e.id === emojiId);
	return e ? e.glyph : '';
}