import {
	createContext,
	type ReactNode,
	useCallback,
	useContext,
	useEffect,
	useState,
} from "react";

const API_URL = import.meta.env.VITE_MAKTABA_API_URL || "http://localhost:8080";
const ACCESS_TOKEN_KEY = "bayt_access_token";

interface User {
	id: string;
	email?: string;
	username?: string;
	firstName?: string;
	lastName?: string;
}

interface AuthContextType {
	isAuthenticated: boolean;
	isLoading: boolean;
	session: null;
	identity: null;
	user: User;
	accessToken: string | null;
	login: (login: string, password: string) => Promise<void>;
	register: (
		email: string,
		username: string,
		password: string,
	) => Promise<void>;
	logout: () => Promise<void>;
	refreshSession: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
	const [isAuthenticated, setIsAuthenticated] = useState(false);
	const [isLoading, setIsLoading] = useState(true);
	const [user, setUser] = useState<User>({ id: "" });
	const [accessToken, setAccessToken] = useState<string | null>(() =>
		typeof window === "undefined"
			? null
			: localStorage.getItem(ACCESS_TOKEN_KEY),
	);

	const storeAccessToken = useCallback((token: string | null) => {
		setAccessToken(token);
		if (typeof window === "undefined") {
			return;
		}
		if (token) {
			localStorage.setItem(ACCESS_TOKEN_KEY, token);
		} else {
			localStorage.removeItem(ACCESS_TOKEN_KEY);
		}
	}, []);

	const fetchMe = useCallback(async (token: string) => {
		const response = await fetch(`${API_URL}/api/me`, {
			headers: { Authorization: `Bearer ${token}` },
			credentials: "include",
		});
		if (!response.ok) {
			throw new Error("not authenticated");
		}
		const data = await response.json();
		setUser({
			id: data.id,
			email: data.email,
			username: data.username,
			firstName: data.username,
		});
		setIsAuthenticated(true);
	}, []);

	const refreshSession = useCallback(async () => {
		setIsLoading(true);
		try {
			const currentToken =
				typeof window === "undefined"
					? null
					: localStorage.getItem(ACCESS_TOKEN_KEY);
			if (currentToken) {
				await fetchMe(currentToken);
				storeAccessToken(currentToken);
				return;
			}

			const response = await fetch(`${API_URL}/auth/refresh`, {
				method: "POST",
				credentials: "include",
			});
			if (!response.ok) {
				throw new Error("refresh failed");
			}
			const data = await response.json();
			storeAccessToken(data.access_token);
			await fetchMe(data.access_token);
		} catch {
			setIsAuthenticated(false);
			setUser({ id: "" });
			storeAccessToken(null);
		} finally {
			setIsLoading(false);
		}
	}, [fetchMe, storeAccessToken]);

	const login = useCallback(
		async (loginValue: string, password: string) => {
			const response = await fetch(`${API_URL}/auth/login`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				credentials: "include",
				body: JSON.stringify({ login: loginValue, password }),
			});
			if (!response.ok) {
				throw new Error("invalid credentials");
			}
			const data = await response.json();
			storeAccessToken(data.tokens.access_token);
			setUser({
				id: data.user.id,
				email: data.user.email,
				username: data.user.username,
				firstName: data.user.username,
			});
			setIsAuthenticated(true);
		},
		[storeAccessToken],
	);

	const register = useCallback(
		async (email: string, username: string, password: string) => {
			const response = await fetch(`${API_URL}/auth/register`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				credentials: "include",
				body: JSON.stringify({ email, username, password }),
			});
			if (!response.ok) {
				throw new Error("registration failed");
			}
			const data = await response.json();
			storeAccessToken(data.tokens.access_token);
			setUser({
				id: data.user.id,
				email: data.user.email,
				username: data.user.username,
				firstName: data.user.username,
			});
			setIsAuthenticated(true);
		},
		[storeAccessToken],
	);

	const logout = useCallback(async () => {
		storeAccessToken(null);
		setIsAuthenticated(false);
		setUser({ id: "" });
	}, [storeAccessToken]);

	useEffect(() => {
		refreshSession();
	}, [refreshSession]);

	return (
		<AuthContext.Provider
			value={{
				isAuthenticated,
				isLoading,
				session: null,
				identity: null,
				user,
				accessToken,
				login,
				register,
				logout,
				refreshSession,
			}}
		>
			{children}
		</AuthContext.Provider>
	);
}

export function useAuth(): AuthContextType {
	const context = useContext(AuthContext);
	if (context === undefined) {
		throw new Error("useAuth must be used within an AuthProvider");
	}
	return context;
}
