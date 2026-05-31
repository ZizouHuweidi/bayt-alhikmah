import { PUBLIC_MAKTABA_API_URL } from '$env/static/public';

const API_URL = PUBLIC_MAKTABA_API_URL || 'http://localhost:8080';
const ACCESS_TOKEN_KEY = 'bayt_access_token';

export interface User {
	id: string;
	email?: string;
	username?: string;
	firstName?: string;
}

function createAuthState() {
	let isAuthenticated = $state(false);
	let isLoading = $state(true);
	let user = $state<User>({ id: '' });
	let accessToken = $state<string | null>(localStorage.getItem(ACCESS_TOKEN_KEY));

	function storeAccessToken(token: string | null) {
		accessToken = token;
		if (token) {
			localStorage.setItem(ACCESS_TOKEN_KEY, token);
		} else {
			localStorage.removeItem(ACCESS_TOKEN_KEY);
		}
	}

	async function fetchMe(token: string) {
		const response = await fetch(`${API_URL}/api/me`, {
			headers: { Authorization: `Bearer ${token}` },
			credentials: 'include',
		});
		if (!response.ok) {
			throw new Error('not authenticated');
		}
		const data = await response.json();
		user = {
			id: data.id,
			email: data.email,
			username: data.username,
			firstName: data.username,
		};
		isAuthenticated = true;
	}

	async function refreshSession() {
		isLoading = true;
		let tokenAtStart: string | null = null;
		try {
			tokenAtStart = localStorage.getItem(ACCESS_TOKEN_KEY);
			if (tokenAtStart) {
				await fetchMe(tokenAtStart);
				storeAccessToken(tokenAtStart);
				return;
			}

			const response = await fetch(`${API_URL}/auth/refresh`, {
				method: 'POST',
				credentials: 'include',
			});
			if (!response.ok) {
				throw new Error('refresh failed');
			}
			const data = await response.json();
			const refreshedToken = data.tokens?.access_token;
			if (!refreshedToken) {
				throw new Error('refresh response missing access token');
			}
			storeAccessToken(refreshedToken);
			await fetchMe(refreshedToken);
		} catch {
			const latestToken = localStorage.getItem(ACCESS_TOKEN_KEY);
			if (latestToken && latestToken !== tokenAtStart) {
				return;
			}
			isAuthenticated = false;
			user = { id: '' };
			storeAccessToken(null);
		} finally {
			isLoading = false;
		}
	}

	async function login(loginValue: string, password: string) {
		const response = await fetch(`${API_URL}/auth/login`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ login: loginValue, password }),
		});
		if (!response.ok) {
			throw new Error('invalid credentials');
		}
		const data = await response.json();
		storeAccessToken(data.tokens.access_token);
		user = {
			id: data.user.id,
			email: data.user.email,
			username: data.user.username,
			firstName: data.user.username,
		};
		isAuthenticated = true;
		isLoading = false;
	}

	async function register(email: string, username: string, password: string) {
		const response = await fetch(`${API_URL}/auth/register`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, username, password }),
		});
		if (!response.ok) {
			throw new Error('registration failed');
		}
		const data = await response.json();
		storeAccessToken(data.tokens.access_token);
		user = {
			id: data.user.id,
			email: data.user.email,
			username: data.user.username,
			firstName: data.user.username,
		};
		isAuthenticated = true;
		isLoading = false;
	}

	async function logout() {
		if (accessToken) {
			storeAccessToken(null);
		}
		try {
			await fetch(`${API_URL}/auth/logout`, {
				method: 'POST',
				credentials: 'include',
			});
		} catch {
			// Best-effort server-side logout
		}
		isAuthenticated = false;
		user = { id: '' };
	}

	return {
		get isAuthenticated() { return isAuthenticated; },
		get isLoading() { return isLoading; },
		get user() { return user; },
		get accessToken() { return accessToken; },
		login,
		register,
		logout,
		refreshSession,
	};
}

export const auth = createAuthState();
