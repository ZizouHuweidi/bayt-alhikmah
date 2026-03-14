import { useEffect, useState } from 'react'
import { useSearch, createFileRoute } from '@tanstack/react-router'
import { ory } from '@/lib/auth/AuthContext'
import { Card, CardContent } from '@/components/ui/card'
import { Loader2, Library, CheckCircle, AlertCircle } from 'lucide-react'

export const Route = createFileRoute('/verification')({
  component: VerificationPage,
})

function VerificationPage() {
  const search = useSearch({ from: '/verification' }) as { flow?: string }
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading')
  const [message, setMessage] = useState('')

  useEffect(() => {
    const initFlow = async () => {
      try {
        const flowId = search.flow

        if (flowId) {
          const { data } = await ory.getVerificationFlow(flowId)
          if (data.state === 'passed_challenge') {
            setStatus('success')
            setMessage('Your email has been verified successfully!')
          } else {
            setStatus('success')
            setMessage('Verification in progress. Please check your email.')
          }
        } else {
          setStatus('error')
          setMessage('No verification flow found. Please use the link from your email.')
        }
      } catch (error) {
        console.error('Verification error:', error)
        setStatus('error')
        setMessage('Verification failed. Please try again or request a new verification email.')
      }
    }

    initFlow()
  }, [search.flow])

  if (status === 'loading') {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-b from-slate-50 to-slate-100 px-4">
      <div className="w-full max-w-md">
        <div className="mb-8 flex flex-col items-center">
          <div className="mb-4 rounded-xl bg-emerald-600 p-3">
            <Library className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-slate-900">Bayt al Hikmah</h1>
        </div>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col items-center text-center">
              {status === 'success' ? (
                <>
                  <CheckCircle className="h-16 w-16 text-emerald-500 mb-4" />
                  <h2 className="text-xl font-semibold text-slate-900 mb-2">Email Verified</h2>
                  <p className="text-slate-600">{message}</p>
                </>
              ) : (
                <>
                  <AlertCircle className="h-16 w-16 text-red-500 mb-4" />
                  <h2 className="text-xl font-semibold text-slate-900 mb-2">Verification Issue</h2>
                  <p className="text-slate-600">{message}</p>
                </>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
