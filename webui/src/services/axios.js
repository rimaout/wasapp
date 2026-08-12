import axios from "axios";
import { getToken, clearAuth } from "./auth.js";

const instance = axios.create({
	baseURL: __API_URL__,
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
