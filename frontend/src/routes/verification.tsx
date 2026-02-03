import { useEffect, useState } from 'react'
import { Link, useSearch } from '@tanstack/react-router'
import type { VerificationFlow, UiText } from '@ory/client'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { createFileRoute } from '@tanstack/react-router'
import {
  createVerificationFlow,
  submitVerificationFlow,
  getCsrfToken,
} from '@/lib/auth/kratos'
import { Library, Loader2, CheckCircle } from 'lucide-react'

export const Route = createFileRoute('/verification')({
  component: VerificationPage,
})

function VerificationPage() {
  const search = useSearch({ from: '/verification' }) as { flow?: string }
  const [flow, setFlow] = useState<VerificationFlow | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)
  const [formData, setFormData] = useState<Record<string, string>>({})

  useEffect(() => {
    async function initFlow() {
      try {
        setIsLoading(true)
        setError(null)

        const result = await createVerificationFlow()

        if ('error' in result) {
          setError(result.error)
          return
        }

        setFlow(result.flow)

        // Initialize form data with CSRF token
        const csrfToken = getCsrfToken(result.flow.ui.nodes)
        if (csrfToken) {
          setFormData((prev) => ({ ...prev, csrf_token: csrfToken }))
        }
      } catch (err) {
        setError('Failed to initialize verification. Please try again.')
      } finally {
        setIsLoading(false)
      }
    }

    initFlow()
  }, [])

  const handleInputChange = (name: string, value: string) => {
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!flow) return

    setIsSubmitting(true)
    setError(null)
    setSuccessMessage(null)

    try {
      const result = await submitVerificationFlow(flow.id, formData)

      if ('error' in result) {
        setError(result.error)
        return
      }

      setFlow(result.flow)

      // Check for messages in the flow response
      if (result.flow.ui.messages && result.flow.ui.messages.length > 0) {
        const successMsg = result.flow.ui.messages.find((m: UiText) => m.type === 'info')
        const errorMsg = result.flow.ui.messages.find((m: UiText) => m.type === 'error')
        
        if (errorMsg) {
          setError(errorMsg.text)
          return
        }
        
        if (successMsg) {
          setSuccessMessage(successMsg.text)
        }
      }

      // Re-initialize CSRF token after submission
      const csrfToken = getCsrfToken(result.flow.ui.nodes)
      if (csrfToken) {
        setFormData({ csrf_token: csrfToken })
      }
    } catch (err) {
      setError('Failed to submit verification. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-slate-50">
        <Loader2 className="w-8 h-8 animate-spin text-emerald-600" />
      </div>
    )
  }

  if (!flow) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-slate-50">
        <Alert variant="destructive" className="max-w-md">
          <AlertDescription>{error || 'Failed to load verification form'}</AlertDescription>
        </Alert>
      </div>
    )
  }

  // Get input nodes from the flow
  const inputNodes = flow.ui.nodes.filter(
    (node) =>
      node.attributes &&
      'node_type' in node.attributes &&
      node.attributes.node_type === 'input' &&
      'name' in node.attributes &&
      node.attributes.name !== 'csrf_token'
  )

  return (
    <div className="flex items-center justify-center min-h-screen px-4 bg-gradient-to-b from-slate-50 to-slate-100">
      <div className="w-full max-w-md">
        <div className="flex flex-col items-center mb-8">
          <div className="p-3 mb-4 rounded-xl bg-emerald-600">
            <Library className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-slate-900">Bayt al Hikmah</h1>
          <p className="mt-2 text-slate-600">Verify your email address</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Verify Email</CardTitle>
            <CardDescription>
              Verify your email address to complete your account setup
            </CardDescription>
          </CardHeader>
          <CardContent>
            {error && (
              <Alert variant="destructive" className="mb-4">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            {successMessage && (
              <Alert className="mb-4 border-emerald-200 bg-emerald-50">
                <CheckCircle className="w-4 h-4 text-emerald-600" />
                <AlertDescription className="text-emerald-800">{successMessage}</AlertDescription>
              </Alert>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              {inputNodes.map((node) => {
                const attrs = node.attributes as {
                  name: string
                  type: string
                  label?: { text: string }
                  autocomplete?: string
                  required?: boolean
                }
                const label = node.meta?.label?.text || attrs.label?.text || attrs.name

                return (
                  <div key={attrs.name} className="space-y-2">
                    <Label htmlFor={attrs.name}>
                      {label}
                      {attrs.required && <span className="text-red-500"> *</span>}
                    </Label>
                    <Input
                      id={attrs.name}
                      name={attrs.name}
                      type={attrs.type === 'password' ? 'password' : 'text'}
                      autoComplete={attrs.autocomplete}
                      required={attrs.required}
                      value={formData[attrs.name] || ''}
                      onChange={(e) => handleInputChange(attrs.name, e.target.value)}
                      disabled={isSubmitting}
                    />
                  </div>
                )
              })}

              <Button
                type="submit"
                className="w-full bg-emerald-600 hover:bg-emerald-700"
                disabled={isSubmitting}
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Verifying...
                  </>
                ) : (
                  'Verify Email'
                )}
              </Button>
            </form>

            <div className="mt-6 space-y-2 text-center">
              <p className="text-sm text-slate-600">
                Already verified?{' '}
                <Link
                  to="/login"
                  className="font-medium text-emerald-600 hover:text-emerald-700"
                >
                  Sign in
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}