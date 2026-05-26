import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import { Library, Loader2, Mail, Shield, User } from 'lucide-react'
import type { FormEvent } from 'react'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { type Profile, getProfile, updateProfile } from '@/lib/api'
import { useAuth } from '@/lib/auth/AuthContext'

export const Route = createFileRoute('/settings')({
  component: SettingsPage,
})

function SettingsPage() {
  const navigate = useNavigate()
  const { isAuthenticated, isLoading, user, accessToken } = useAuth()
  const [profile, setProfile] = useState<Profile | null>(null)
  const [displayName, setDisplayName] = useState('')
  const [bio, setBio] = useState('')
  const [publicProfile, setPublicProfile] = useState(false)
  const [loadingProfile, setLoadingProfile] = useState(true)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate({ to: '/login' })
    }
  }, [isAuthenticated, isLoading, navigate])

  useEffect(() => {
    if (!isAuthenticated || !accessToken) {
      return
    }
    let cancelled = false
    async function loadProfile() {
      setLoadingProfile(true)
      setError(null)
      try {
        const result = await getProfile(accessToken as string)
        if (cancelled) return
        setProfile(result)
        setDisplayName(result.display_name || '')
        setBio(result.bio || '')
        setPublicProfile(result.public_profile)
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Failed to load profile')
        }
      } finally {
        if (!cancelled) {
          setLoadingProfile(false)
        }
      }
    }
    loadProfile()
    return () => {
      cancelled = true
    }
  }, [isAuthenticated, accessToken])

  const handleSave = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!accessToken) return

    setSaving(true)
    setError(null)
    setMessage(null)
    try {
      const updated = await updateProfile(accessToken, {
        display_name: displayName.trim() || undefined,
        bio: bio.trim() || undefined,
        public_profile: publicProfile,
      })
      setProfile(updated)
      setMessage('Profile saved')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save profile')
    } finally {
      setSaving(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-slate-50 py-8">
      <div className="mx-auto max-w-2xl px-4">
        <div className="mb-8 flex items-center justify-between gap-3">
          <div className="flex items-center gap-3">
            <div className="rounded-xl bg-emerald-600 p-2">
              <Library className="h-6 w-6 text-white" />
            </div>
            <h1 className="text-2xl font-bold text-slate-900">
              Account Settings
            </h1>
          </div>
          <Link to="/dashboard">
            <Button variant="outline">Dashboard</Button>
          </Link>
        </div>

        {error && (
          <div className="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}
        {message && (
          <div className="mb-6 rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
            {message}
          </div>
        )}

        <Card>
          <CardContent className="space-y-6 pt-6">
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
                <User className="h-6 w-6 text-emerald-600" />
              </div>
              <div>
                <p className="font-medium text-slate-900">{user.username}</p>
                <p className="text-sm text-slate-500">Username</p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
                <Mail className="h-6 w-6 text-emerald-600" />
              </div>
              <div>
                <p className="font-medium text-slate-900">
                  {user.email || 'No email'}
                </p>
                <p className="text-sm text-slate-500">Email address</p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
                <Shield className="h-6 w-6 text-emerald-600" />
              </div>
              <div>
                <p className="font-medium text-slate-900">Password auth enabled</p>
                <p className="text-sm text-slate-500">
                  OAuth support can be added later.
                </p>
              </div>
            </div>

            <form className="space-y-4 border-t border-slate-200 pt-6" onSubmit={handleSave}>
              <div>
                <label className="mb-2 block text-sm font-medium text-slate-700">
                  Display name
                </label>
                <Input
                  value={displayName}
                  onChange={event => setDisplayName(event.target.value)}
                  placeholder="How your public profile should appear"
                  disabled={loadingProfile}
                />
              </div>
              <div>
                <label className="mb-2 block text-sm font-medium text-slate-700">
                  Bio
                </label>
                <textarea
                  className="min-h-28 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm"
                  value={bio}
                  onChange={event => setBio(event.target.value)}
                  placeholder="What are you reading, researching, or collecting?"
                  disabled={loadingProfile}
                />
              </div>
              <label className="flex items-start gap-3 rounded-lg border border-slate-200 bg-slate-50 p-4 text-sm text-slate-700">
                <input
                  type="checkbox"
                  className="mt-1"
                  checked={publicProfile}
                  onChange={event => setPublicProfile(event.target.checked)}
                  disabled={loadingProfile}
                />
                <span>
                  <span className="block font-medium text-slate-900">
                    Make my profile public
                  </span>
                  Public profiles can be opened at{' '}
                  <code>/users/{user.username}/profile</code> and can anchor public
                  notes, reviews, collections, and library items.
                </span>
              </label>
              {profile?.public_profile && user.username && (
                <p className="text-sm text-slate-500">
                  Public profile:{' '}
                  <Link
                    to="/users/$username/profile"
                    params={{ username: user.username }}
                    className="font-medium text-emerald-700 underline-offset-4 hover:underline"
                  >
                    /users/{user.username}/profile
                  </Link>
                </p>
              )}
              <Button type="submit" disabled={saving || loadingProfile}>
                {saving ? 'Saving...' : 'Save profile'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
