import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Library, ArrowRight } from 'lucide-react'

const ORY_HOSTED_URL = 'https://sleepy-swartz-8u1sjz0in0.projects.oryapis.com'

export const Route = createFileRoute('/recovery')({
  component: RecoveryPage,
})

function RecoveryPage() {
  const handleRecovery = () => {
    window.location.href = `${ORY_HOSTED_URL}/recovery`
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-b from-slate-50 to-slate-100 px-4">
      <div className="w-full max-w-md">
        <div className="mb-8 flex flex-col items-center">
          <div className="mb-4 rounded-xl bg-emerald-600 p-3">
            <Library className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-slate-900">Bayt al Hikmah</h1>
          <p className="mt-2 text-slate-600">Recover your account</p>
        </div>

        <Card>
          <CardHeader className="text-center">
            <CardTitle>Account Recovery</CardTitle>
            <CardDescription>
              Reset your password or recover access to your account
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Button onClick={handleRecovery} className="w-full" size="lg">
              Continue with Ory
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
