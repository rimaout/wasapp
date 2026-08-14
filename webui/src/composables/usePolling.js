export function usePolling(fn, interval) {
	const id = setInterval(fn, interval);
	return () => clearInterval(id);
}
