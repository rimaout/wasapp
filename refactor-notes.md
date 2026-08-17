# Refactoring Notes (frontend)

Pending cleanups identified after the group-members / popup / user-picker work.
To be done incrementally.

## 1. Extract overlay color variable (high value)

`rgba(255, 255, 255, 0.08)` is hardcoded **14 times** across 11 files as the de-facto
"hover / overlay" color. Introduce a CSS variable (e.g. `--tn-overlay`) in
`tokyo-night.css` and replace all usages.

Files:
- `components/IconButton.vue`
- `components/SearchBar.vue`
- `components/MemberChips.vue`
- `components/UserPicker.vue`
- `components/PopupMenu.vue`
- `components/PopupActionScreen.vue`
- `components/ChatMembersPopup.vue`
- `views/ChatView.vue`
- `views/GroupMembersView.vue`
- `views/GroupDetailsView.vue`

Also consider variables for the translucent red/green button backgrounds
(e.g. `--tn-red-bg`, `--tn-green-bg`).

## 2. Extract shared popup shell (high value)

`.popup-backdrop` + `.popup-panel` are duplicated between:
- `components/PopupMenu.vue`
- `components/ChatMembersPopup.vue`

Extract a generic `Popover` / `PopupPanel` component:
- trigger button
- fixed backdrop (click-to-close)
- absolutely-positioned panel (bg, border, radius, shadow, `direction: up/down`)
- slot for content

## 3. Extract shared action buttons (high value)

Green confirm/add/create button (`rgba(158,206,106,0.15)` bg + green text +
40px height + disabled-gray) duplicated in:
- `components/PopupActionScreen.vue`
- `components/ChatMembersPopup.vue`
- `views/GroupDetailsView.vue`
- `views/GroupMembersView.vue`

Red `x` close button duplicated in:
- `components/PopupActionScreen.vue`
- `components/ChatMembersPopup.vue`

Extract a shared `ActionButton` component (or shared CSS classes) for
`confirm` (green) and `close` (red) variants.

## 4. Backend spec/impl mismatch (correctness)

`GET /chats/{chatId}/members` returns a **bare array** (`[]Member`) in
`service/api/get-chat-members.go`, but `doc/api.yaml` documents
`{ membersList: [...] }`.

Decide and fix one side (probably align the spec to the implementation, or wrap
the response). The frontend currently reads the bare array directly.

## 5. useUsers singleton (minor)

`composables/useUsers.js` creates fresh refs per call, so each user list
independently fetches `GET /users`. Could be made module-level shared state
(like `useChats`) to avoid duplicate fetches.

## 6. (Optional) user-row duplication

`views/UserSearchView.vue` still has its own `.user-row` styles vs
`components/UserPicker.vue`. Different interaction (tap-to-create vs
multi-select), so extraction may not be worth it.
