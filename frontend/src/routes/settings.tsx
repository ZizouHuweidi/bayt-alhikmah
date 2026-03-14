import { useEffect, useState } from 'react'
import { useNavigate, createFileRoute } from '@tanstack/react-router'
import { ory } from '@/lib/auth/AuthContext'
import { useAuth } from '@/lib/auth/AuthContext'
import { Card, CardContent } from '@/components/ui/card'
import { Loader2, Library, User, Mail, Shield } from 'lucide-react'

const ORY_HOSTED_URL = 'https://sleepy-swartz-8u1sjz0in0.projects.oryapis.com'

export const Route = createFileRoute('/settings')({
  component: SettingsPage,
})

function SettingsPage() {
  const navigate = useNavigate()
  const { isAuthenticated, isLoading: authLoading, identity } = useAuth()
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate({ to: '/login' })
      return
    }

    if (isAuthenticated) {
      setIsLoading(false)
    }
  }, [authLoading, isAuthenticated, navigate])

  const handleSettings = () => {
    window.location.href = `${ORY_HOSTED_URL}/settings`
  }

  if (authLoading || isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </div>
    )
  }

  const traits = identity?.traits as
    | { email?: string; name?: { first?: string; last?: string } }
    | undefined
  const email = traits?.email
  const name = traits?.name

  return (
    <div className="min-h-screen bg-slate-50 py-8">
      <div className="mx-auto max-w-2xl px-4">
        <div className="mb-8 flex items-center gap-3">
          <div className="rounded-xl bg-emerald-600 p-2">
            <Library className="h-6 w-6 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900">Account Settings</h1>
        </div>

        <Card>
          <CardContent className="pt-6 space-y-6">
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
                <User className="h-6 w-6 text-emerald-600" />
              </div>
              <div>
                <p className="font-medium text-slate-900">
                  {name?.first} {name?.last}
                </p>
                <p className="text-sm text-slate-500">Your profile</p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
                <Mail className="h-6 w-6 text-emerald-600" />
              </div>
              <div>
                <p className="font-medium text-slate-900">{email || 'No email'}</p>
                <p className="text-sm text-slate-500">Email address</p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
                <Shield className="h-6 w-6 text-emerald-600" />
              </div>
              <div>
                <p className="font-medium text-slate-900">Security Settings</p>
                <p className="text-sm text-slate-500">Password, 2FA, and more</p>
              </div>
            </div>

            <div className="pt-4">
              <button
                type="button"
                onClick={handleSettings}
                className="inline-flex items-center justify-center rounded-md bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-700 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2"
              >
                Manage Account on Ory
              </button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
