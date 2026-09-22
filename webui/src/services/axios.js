import axios from "axios";
import { getToken, clearAuth } from "./auth.js";

/**
 * Shared Axios instance for making API requests. Automatically adds the auth token to the Authorization
 * header if it exists, and handles 401 Unauthorized responses by clearing the auth and redirecting to the login page.
 */
const instance = axios.create({
	baseURL: `http://${window.location.hostname}:3000`,
	timeout: 1000 * 5
});

// Automatically add the token to the Authorization header of each request if it exists.
instance.interceptors.request.use(config => {
	const token = getToken();
	if (token) {
		config.headers.Authorization = 'Bearer ' + token;
	}
	return config;
});

// Automatically handle 401 Unauthorized responses by clearing the auth and redirecting to the login page.
instance.interceptors.response.use(
	response => response,
	error => {
		if (error.response && error.response.status === 401) {
			clearAuth();
			window.location.hash = '#/login';
		}
		return Promise.reject(error);
	}
);

export default instance;
