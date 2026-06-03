import { create } from "zustand";
import { persist } from "zustand/middleware";
import { API_URL, getMe, type User } from "~/lib/api";

type AuthState = {
  accessToken: string | null;
  user: User;
  isAuthenticated: boolean;
  isLoading: boolean;
  setAccessToken: (token: string | null) => void;
  refreshSession: () => Promise<void>;
  login: (login: string, password: string) => Promise<void>;
  register: (email: string, username: string, password: string) => Promise<void>;
  logout: () => void;
};

const emptyUser: User = { id: "" };

export function accessTokenFromAuthResponse(data: unknown): string | null {
  if (!data || typeof data !== "object") {
    return null;
  }
  const tokens = (data as { tokens?: unknown }).tokens;
  if (!tokens || typeof tokens !== "object") {
    return null;
  }
  const accessToken = (tokens as { access_token?: unknown }).access_token;
  return typeof accessToken === "string" && accessToken ? accessToken : null;
}

function userFromResponse(data: unknown): User {
  const user = (data as { user?: User }).user;
  if (user?.id) {
    return { ...user, firstName: user.username };
  }
  return emptyUser;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      accessToken: null,
      user: emptyUser,
      isAuthenticated: false,
      isLoading: true,
      setAccessToken: (token) => set({ accessToken: token }),
      refreshSession: async () => {
        set({ isLoading: true });
        const tokenAtStart = get().accessToken;
        try {
          if (tokenAtStart) {
            const user = await getMe(tokenAtStart);
            set({
              accessToken: tokenAtStart,
              user: { ...user, firstName: user.username },
              isAuthenticated: true,
              isLoading: false,
            });
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
          const refreshedToken = accessTokenFromAuthResponse(data);
          if (!refreshedToken) {
            throw new Error("refresh response missing access token");
          }
          const user = await getMe(refreshedToken);
          set({
            accessToken: refreshedToken,
            user: { ...user, firstName: user.username },
            isAuthenticated: true,
            isLoading: false,
          });
        } catch {
          if (get().accessToken && get().accessToken !== tokenAtStart) {
            set({ isLoading: false });
            return;
          }
          set({
            accessToken: null,
            user: emptyUser,
            isAuthenticated: false,
            isLoading: false,
          });
        }
      },
      login: async (loginValue, password) => {
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
        const accessToken = accessTokenFromAuthResponse(data);
        if (!accessToken) {
          throw new Error("login response missing access token");
        }
        set({
          accessToken,
          user: userFromResponse(data),
          isAuthenticated: true,
          isLoading: false,
        });
      },
      register: async (email, username, password) => {
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
        const accessToken = accessTokenFromAuthResponse(data);
        if (!accessToken) {
          throw new Error("registration response missing access token");
        }
        set({
          accessToken,
          user: userFromResponse(data),
          isAuthenticated: true,
          isLoading: false,
        });
      },
      logout: () =>
        set({
          accessToken: null,
          user: emptyUser,
          isAuthenticated: false,
          isLoading: false,
        }),
    }),
    {
      name: "bayt-auth",
      partialize: (state) => ({ accessToken: state.accessToken }),
      onRehydrateStorage: () => (state) => {
        state?.refreshSession();
      },
    }
  )
);
