import { Link, createFileRoute } from '@tanstack/react-router'
import { BookOpen, Library, Loader2, UserRound } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { type Book, addLibraryItem, getBook } from '@/lib/api'
import { useAuth } from '@/lib/auth/AuthContext'

export const Route = createFileRoute('/sources/books/$id')({
  component: BookDetailPage,
})

function BookDetailPage() {
  const { id } = Route.useParams()
  const { isAuthenticated, accessToken } = useAuth()
  const [book, setBook] = useState<Book | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [added, setAdded] = useState(false)

  useEffect(() => {
    let cancelled = false
    async function loadBook() {
      setLoading(true)
      setError(null)
      try {
        const result = await getBook(id)
        if (!cancelled) setBook(result)
      } catch (err) {
        if (!cancelled) setError(err instanceof Error ? err.message : 'Failed to load book')
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    loadBook()
    return () => {
      cancelled = true
    }
  }, [id])

  const handleAdd = async () => {
    if (!accessToken) return
    setError(null)
    try {
      await addLibraryItem(accessToken, id)
      setAdded(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add to library')
    }
  }

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-50">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </main>
    )
  }

  if (!book) {
    return (
      <main className="min-h-screen bg-slate-50 px-4 py-16">
        <div className="mx-auto max-w-2xl rounded-2xl border bg-white p-8 text-center">
          <h1 className="text-2xl font-bold text-slate-900">Book not found</h1>
          <Link to="/dashboard"><Button className="mt-6">Back to dashboard</Button></Link>
        </div>
      </main>
    )
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-10">
      <div className="mx-auto max-w-5xl">
        <Link to="/dashboard" className="text-sm font-medium text-emerald-700 hover:underline">
          Back to dashboard
        </Link>
        {error && <div className="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div>}

        <section className="mt-6 rounded-3xl bg-white p-8 shadow-sm">
          <div className="flex flex-col gap-6 md:flex-row md:items-start md:justify-between">
            <div>
              <div className="mb-4 inline-flex items-center gap-2 rounded-full bg-emerald-50 px-3 py-1 text-sm font-medium text-emerald-700">
                <BookOpen className="h-4 w-4" /> Book
              </div>
              <h1 className="text-4xl font-bold text-slate-900">{book.source.title}</h1>
              {book.source.subtitle && <p className="mt-3 text-xl text-slate-600">{book.source.subtitle}</p>}
              {book.contributors && book.contributors.length > 0 && (
                <p className="mt-4 flex items-center gap-2 text-slate-600">
                  <UserRound className="h-4 w-4" />
                  {book.contributors.map(contributor => contributor.name).join(', ')}
                </p>
              )}
            </div>
            {isAuthenticated && (
              <Button onClick={handleAdd} disabled={added}>
                <Library className="h-4 w-4" /> {added ? 'Added' : 'Add to library'}
              </Button>
            )}
          </div>
        </section>

        <div className="mt-6 grid gap-6 md:grid-cols-2">
          <Card>
            <CardHeader><CardTitle>Metadata</CardTitle></CardHeader>
            <CardContent className="space-y-3 text-sm text-slate-700">
              <Metadata label="ISBN-13" value={book.metadata?.isbn_13 || book.source.isbn} />
              <Metadata label="Publisher" value={book.metadata?.publisher || book.source.publisher} />
              <Metadata label="Language" value={book.metadata?.language} />
              <Metadata label="Pages" value={book.metadata?.page_count?.toString()} />
            </CardContent>
          </Card>
          <Card>
            <CardHeader><CardTitle>Description</CardTitle></CardHeader>
            <CardContent>
              <p className="text-sm leading-6 text-slate-700">
                {book.source.description || 'No description yet.'}
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </main>
  )
}

function Metadata({ label, value }: { label: string; value?: string }) {
  return (
    <div className="flex justify-between gap-4 border-b border-slate-100 pb-2 last:border-0">
      <span className="text-slate-500">{label}</span>
      <span className="font-medium text-slate-900">{value || 'Not set'}</span>
    </div>
  )
}
