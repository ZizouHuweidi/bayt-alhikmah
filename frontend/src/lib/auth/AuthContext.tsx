import {
	createContext,
	useContext,
	useState,
	useEffect,
	useCallback,
	type ReactNode,
} from 'react'
import type { Identity, Session } from '@ory/client'
import {
	getSession,
	getIdentity,
	isAuthenticated as checkIsAuthenticated,
	logout as kratosLogout,
	getUserTraits,
} from '@/lib/auth/kratos'

interface AuthContextType {
	isAuthenticated: boolean
	isLoading: boolean
	session: Session | null
	identity: Identity | null
	user: {
		id: string
		email?: string
		firstName?: string
		lastName?: string
	}
	logout: () => Promise<void>
	refreshSession: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
	children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
	const [isAuthenticated, setIsAuthenticated] = useState(false)
	const [isLoading, setIsLoading] = useState(true)
	const [session, setSession] = useState<Session | null>(null)
	const [identity, setIdentity] = useState<Identity | null>(null)

	const refreshSession = useCallback(async () => {
		try {
			const session = await getSession()
			const authenticated = await checkIsAuthenticated()

			setSession(session)
			setIsAuthenticated(authenticated)

			if (session?.identity) {
				setIdentity(session.identity)
			} else {
				const id = await getIdentity()
				setIdentity(id)
			}
		} catch (error) {
			console.error('Failed to refresh session:', error)
			setIsAuthenticated(false)
			setSession(null)
			setIdentity(null)
		} finally {
			setIsLoading(false)
		}
	}, [])

	const logout = useCallback(async () => {
		const success = await kratosLogout()
		if (success) {
			setIsAuthenticated(false)
			setSession(null)
			setIdentity(null)
		}
	}, [])

	useEffect(() => {
		refreshSession()
	}, [refreshSession])

	const user = {
		id: identity?.id || '',
		...getUserTraits(identity),
	}

	const value: AuthContextType = {
		isAuthenticated,
		isLoading,
		session,
		identity,
		user,
		logout,
		refreshSession,
	}

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextType {
	const context = useContext(AuthContext)
	if (context === undefined) {
		throw new Error('useAuth must be used within an AuthProvider')
	}
	return context
}
