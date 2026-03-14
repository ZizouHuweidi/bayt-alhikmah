import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from 'react'
import { Configuration, FrontendApi } from '@ory/client-fetch'
import type { Session, Identity } from '@ory/client-fetch'

const getOryClient = () => {
  console.log('Creating Ory client with URL:', import.meta.env.VITE_ORY_SDK_URL)
  const config = new Configuration({
    basePath: import.meta.env.VITE_ORY_SDK_URL || 'http://localhost:4000',
    baseOptions: {
      withCredentials: true,
    },
  })
  return new FrontendApi(config)
}

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
    const ory = getOryClient()
    try {
      const { data } = await ory.toSession()
      setSession(data)
      setIsAuthenticated(data.active === true)
      setIdentity(data.identity)
    } catch {
      setIsAuthenticated(false)
      setSession(null)
      setIdentity(null)
    } finally {
      setIsLoading(false)
    }
  }, [])

  const logout = useCallback(async () => {
    const ory = getOryClient()
    try {
      const { data: logoutFlow } = await ory.createBrowserLogoutFlow()

      if (logoutFlow.logout_url) {
        window.location.href = logoutFlow.logout_url
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      setIsAuthenticated(false)
      setSession(null)
      setIdentity(null)
    }
  }, [])

  useEffect(() => {
    refreshSession()
  }, [refreshSession])

  const traits = identity?.traits as
    | { email?: string; name?: { first?: string; last?: string } }
    | undefined

  const user = {
    id: identity?.id || '',
    email: traits?.email,
    firstName: traits?.name?.first,
    lastName: traits?.name?.last,
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

export const ory = getOryClient()
