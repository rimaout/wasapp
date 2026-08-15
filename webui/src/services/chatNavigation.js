// Collapse the sidebar on mobile (a no-op on larger screens).
export function collapseSidebar() {
	const sidebar = document.getElementById('sidebarMenu');
	if (sidebar && window.innerWidth < 768) {
		const bsCollapse = bootstrap.Collapse.getOrCreateInstance(sidebar);
		bsCollapse.hide();
	}
}

// Navigate to a chat, passing its name/group through query params
// (there is no GET /chats/{id} endpoint yet).
export function navigateToChat(router, chat) {
	router.push('/chats/' + chat.id +
		'?name=' + encodeURIComponent(chat.displayName) +
		'&group=' + (chat.isGroupChat ? '1' : '0'));
	collapseSidebar();
}
