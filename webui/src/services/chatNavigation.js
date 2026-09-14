// Collapse the sidebar on mobile (a no-op on larger screens).
export function collapseSidebar() {
	const sidebar = document.getElementById('sidebarMenu');
	if (sidebar && window.innerWidth < 768) {
		const bsCollapse = bootstrap.Collapse.getOrCreateInstance(sidebar);
		bsCollapse.hide();
	}
}

// navigateToChat - Navigates to the chat page for the given chat object, passing the chat name and group status as query parameters.
export function navigateToChat(router, chat) {
	router.push('/chats/' + chat.id +
		'?name=' + encodeURIComponent(chat.displayName) +
		'&group=' + (chat.isGroupChat ? '1' : '0'));

	// Collapse the sidebar on mobile devices after navigation
	collapseSidebar();
}
