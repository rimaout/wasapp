// Runs a function periodically and returns a cleanup function to stop it.
export function usePolling(fn, interval) {
	const id = setInterval(fn, interval);
	return () => clearInterval(id);
}
