import axios from "axios";
import { getToken } from "./auth.js";

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

export default instance;
