import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import {
  Book,
  CheckCircle2,
  Layers3,
  Library,
  Loader2,
  LogOut,
  Plus,
  Star,
  StickyNote,
} from 'lucide-react'
import type { FormEvent, ReactNode } from 'react'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  type Collection,
  type LibraryItemWithSource,
  type Note,
  type Review,
  type Source,
  addLibraryItem,
  createBook,
  createCollection,
  createNote,
  createReview,
  deleteCollection,
  deleteLibraryItem,
  deleteNote,
  deleteReview,
  listCollections,
  listLibrary,
  listNotes,
  listReviews,
  listSources,
  updateLibraryItem,
} from '@/lib/api'
import { useAuth } from '@/lib/auth/AuthContext'

export const Route = createFileRoute('/dashboard')({
  component: DashboardPage,
})

type DashboardData = {
  sources: Source[]
  library: LibraryItemWithSource[]
  notes: Note[]
  reviews: Review[]
  collections: Collection[]
}

const emptyData: DashboardData = {
  sources: [],
  library: [],
  notes: [],
  reviews: [],
  collections: [],
}

function DashboardPage() {
  const navigate = useNavigate()
  const { isAuthenticated, isLoading, user, accessToken, logout } = useAuth()
  const [data, setData] = useState<DashboardData>(emptyData)
  const [loadingData, setLoadingData] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [bookTitle, setBookTitle] = useState('')
  const [bookAuthor, setBookAuthor] = useState('')
  const [bookISBN, setBookISBN] = useState('')
  const [savingBook, setSavingBook] = useState(false)
  const [selectedSourceID, setSelectedSourceID] = useState('')
  const [noteContent, setNoteContent] = useState('')
  const [notePublic, setNotePublic] = useState(false)
  const [reviewRating, setReviewRating] = useState('5')
  const [reviewContent, setReviewContent] = useState('')
  const [collectionName, setCollectionName] = useState('')
  const [collectionPublic, setCollectionPublic] = useState(false)
  const [savingActivity, setSavingActivity] = useState(false)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate({ to: '/login', search: { return_to: '/dashboard' } })
    }
  }, [isAuthenticated, isLoading, navigate])

  const loadDashboard = async () => {
    if (!accessToken) {
      return
    }
    setLoadingData(true)
    setError(null)
    try {
      const [sources, library, notes, reviews, collections] = await Promise.all([
        listSources(),
        listLibrary(accessToken),
        listNotes(accessToken),
        listReviews(accessToken),
        listCollections(accessToken),
      ])
      setData({ sources, library, notes, reviews, collections })
      if (!selectedSourceID && library.length > 0) {
        setSelectedSourceID(library[0].source_id)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load dashboard')
    } finally {
      setLoadingData(false)
    }
  }

  useEffect(() => {
    if (isAuthenticated && accessToken) {
      loadDashboard()
    }
  }, [isAuthenticated, accessToken])

  const handleLogout = async () => {
    await logout()
    navigate({ to: '/' })
  }

  const handleCreateBook = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!accessToken || !bookTitle.trim()) {
      return
    }

    setSavingBook(true)
    setError(null)
    try {
      const book = await createBook(accessToken, {
        title: bookTitle.trim(),
        isbn_13: bookISBN.trim() || undefined,
        contributors: bookAuthor.trim()
          ? [{ name: bookAuthor.trim(), role: 'author' }]
          : [],
      })
      await addLibraryItem(accessToken, book.source.id)
      setBookTitle('')
      setBookAuthor('')
      setBookISBN('')
      await loadDashboard()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create book')
    } finally {
      setSavingBook(false)
    }
  }

  const handleCreateNote = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!accessToken || !selectedSourceID || !noteContent.trim()) return

    setSavingActivity(true)
    setError(null)
    try {
      await createNote(accessToken, {
        source_id: selectedSourceID,
        content: noteContent.trim(),
        content_type: 'note',
        is_public: notePublic,
      })
      setNoteContent('')
      setNotePublic(false)
      await loadDashboard()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create note')
    } finally {
      setSavingActivity(false)
    }
  }

  const handleCreateReview = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!accessToken || !selectedSourceID) return

    setSavingActivity(true)
    setError(null)
    try {
      await createReview(accessToken, {
        source_id: selectedSourceID,
        rating: Number(reviewRating),
        content: reviewContent.trim() || undefined,
        is_public: true,
      })
      setReviewRating('5')
      setReviewContent('')
      await loadDashboard()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create review')
    } finally {
      setSavingActivity(false)
    }
  }

  const handleCreateCollection = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!accessToken || !selectedSourceID || !collectionName.trim()) return

    setSavingActivity(true)
    setError(null)
    try {
      await createCollection(accessToken, {
        name: collectionName.trim(),
        is_public: collectionPublic,
        source_ids: [selectedSourceID],
      })
      setCollectionName('')
      setCollectionPublic(false)
      await loadDashboard()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create collection')
    } finally {
      setSavingActivity(false)
    }
  }

  const runAction = async (action: () => Promise<unknown>, fallback: string) => {
    setError(null)
    try {
      await action()
      await loadDashboard()
    } catch (err) {
      setError(err instanceof Error ? err.message : fallback)
    }
  }

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </div>
    )
  }

  if (!isAuthenticated) {
    return null
  }

  const librarySourceIDs = new Set(data.library.map(item => item.source_id))
  const librarySources = data.library.map(item => ({ item, source: item.source }))
  return (
    <div className="min-h-screen bg-slate-50">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            <div className="rounded-lg bg-emerald-600 p-2">
              <Library className="h-5 w-5 text-white" />
            </div>
            <span className="text-xl font-bold text-slate-900">
              Bayt al Hikmah
            </span>
          </div>

          <div className="flex items-center gap-4">
            <span className="hidden text-sm text-slate-600 sm:inline">
              Welcome, {user.username || user.email}
            </span>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleLogout}
              className="text-slate-600 hover:text-red-600"
            >
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </Button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Dashboard</h1>
            <p className="mt-2 text-slate-600">
              Build and test your knowledge library with real platform data.
            </p>
          </div>
          <Link to="/settings">
            <Button variant="outline">Profile settings</Button>
          </Link>
        </div>

        {error && (
          <div className="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        <div className="mb-8 grid grid-cols-2 gap-4 lg:grid-cols-5">
          <StatCard icon={<Book />} label="Book sources" value={data.sources.length} />
          <StatCard icon={<Library />} label="Library" value={data.library.length} />
          <StatCard icon={<StickyNote />} label="Notes" value={data.notes.length} />
          <StatCard icon={<Star />} label="Reviews" value={data.reviews.length} />
          <StatCard
            icon={<Layers3 />}
            label="Collections"
            value={data.collections.length}
          />
        </div>

        <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_380px]">
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Library className="h-5 w-5 text-emerald-600" />
                  My Library
                </CardTitle>
                <CardDescription>
                  Your tracked sources. New books created here are added to your
                  library automatically.
                </CardDescription>
              </CardHeader>
              <CardContent>
                {loadingData ? (
                  <LoadingRow />
                ) : data.library.length === 0 ? (
                  <EmptyState message="No library items yet. Create a book to start testing the platform." />
                ) : (
                  <div className="space-y-3">
                    {data.library.slice(0, 8).map(item => {
                      const source = item.source
                      return (
                        <div
                          key={item.id}
                          className="rounded-lg border border-slate-200 bg-white p-4"
                        >
                          <div className="flex items-start justify-between gap-3">
                            <div>
                              <Link to="/sources/books/$id" params={{ id: item.source_id }} className="font-medium text-slate-900 hover:text-emerald-700">
                                {source?.title || item.source_id}
                              </Link>
                              <p className="mt-1 text-sm text-slate-500">
                                {item.status.replace('_', ' ')} · {item.visibility}
                              </p>
                            </div>
                            <div className="flex flex-col items-end gap-2">
                              <CheckCircle2 className="h-5 w-5 text-emerald-600" />
                              <select
                                className="h-8 rounded-md border border-slate-300 bg-white px-2 text-xs"
                                value={item.status}
                                onChange={event => {
                                  if (!accessToken) return
                                  runAction(
                                    () => updateLibraryItem(accessToken, item.id, { status: event.target.value }),
                                    'Failed to update status'
                                  )
                                }}
                              >
                                <option value="to_consume">To consume</option>
                                <option value="in_progress">In progress</option>
                                <option value="completed">Completed</option>
                                <option value="paused">Paused</option>
                                <option value="abandoned">Abandoned</option>
                              </select>
                              <div className="flex gap-2">
                                <Button
                                  variant="ghost"
                                  size="xs"
                                  onClick={() => {
                                    if (!accessToken) return
                                    runAction(
                                      () => updateLibraryItem(accessToken, item.id, { visibility: item.visibility === 'public' ? 'private' : 'public' }),
                                      'Failed to update visibility'
                                    )
                                  }}
                                >
                                  {item.visibility === 'public' ? 'Make private' : 'Make public'}
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="xs"
                                  className="text-red-600"
                                  onClick={() => {
                                    if (!accessToken) return
                                    runAction(() => deleteLibraryItem(accessToken, item.id), 'Failed to remove item')
                                  }}
                                >
                                  Remove
                                </Button>
                              </div>
                            </div>
                          </div>
                        </div>
                      )
                    })}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Available Book Sources</CardTitle>
                <CardDescription>
                  Public book records currently in the platform.
                </CardDescription>
              </CardHeader>
              <CardContent>
                {loadingData ? (
                  <LoadingRow />
                ) : data.sources.length === 0 ? (
                  <EmptyState message="No book sources yet." />
                ) : (
                  <div className="space-y-3">
                    {data.sources.slice(0, 8).map(source => (
                      <div
                        key={source.id}
                        className="flex flex-col gap-3 rounded-lg border border-slate-200 bg-white p-4 sm:flex-row sm:items-center sm:justify-between"
                      >
                        <div>
                          <p className="font-medium text-slate-900">
                            {source.title}
                          </p>
                          <p className="mt-1 text-sm text-slate-500">
                            {source.publisher || source.type}
                          </p>
                        </div>
                        {!librarySourceIDs.has(source.id) && (
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={async () => {
                              if (!accessToken) return
                              try {
                                await addLibraryItem(accessToken, source.id)
                                await loadDashboard()
                              } catch (err) {
                                setError(
                                  err instanceof Error
                                    ? err.message
                                    : 'Failed to add source'
                                )
                              }
                            }}
                          >
                            Add to library
                          </Button>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          <aside className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Plus className="h-5 w-5 text-emerald-600" />
                  Add a Book
                </CardTitle>
                <CardDescription>
                  Minimal book creation for MVP testing.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <form className="space-y-4" onSubmit={handleCreateBook}>
                  <Input
                    value={bookTitle}
                    onChange={event => setBookTitle(event.target.value)}
                    placeholder="Book title"
                    required
                  />
                  <Input
                    value={bookAuthor}
                    onChange={event => setBookAuthor(event.target.value)}
                    placeholder="Author"
                  />
                  <Input
                    value={bookISBN}
                    onChange={event => setBookISBN(event.target.value)}
                    placeholder="ISBN-13"
                  />
                  <Button
                    type="submit"
                    className="w-full"
                    disabled={savingBook || !bookTitle.trim()}
                  >
                    {savingBook ? 'Saving...' : 'Create and add'}
                  </Button>
                </form>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Work With a Source</CardTitle>
                <CardDescription>
                  Select a library item, then add notes, reviews, and collections.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <select
                  className="h-9 w-full rounded-md border border-slate-300 bg-white px-3 text-sm"
                  value={selectedSourceID}
                  onChange={event => setSelectedSourceID(event.target.value)}
                  disabled={librarySources.length === 0}
                >
                  {librarySources.length === 0 ? (
                    <option value="">No library sources yet</option>
                  ) : (
                    librarySources.map(({ item, source }) => (
                      <option key={item.id} value={item.source_id}>
                        {source?.title || item.source_id}
                      </option>
                    ))
                  )}
                </select>

                <form className="space-y-3" onSubmit={handleCreateNote}>
                  <textarea
                    className="min-h-24 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm"
                    value={noteContent}
                    onChange={event => setNoteContent(event.target.value)}
                    placeholder="Write a note..."
                  />
                  <label className="flex items-center gap-2 text-sm text-slate-600">
                    <input
                      type="checkbox"
                      checked={notePublic}
                      onChange={event => setNotePublic(event.target.checked)}
                    />
                    Public note
                  </label>
                  <Button
                    type="submit"
                    variant="outline"
                    className="w-full"
                    disabled={savingActivity || !selectedSourceID || !noteContent.trim()}
                  >
                    Add note
                  </Button>
                </form>

                <form className="space-y-3 border-t border-slate-200 pt-4" onSubmit={handleCreateReview}>
                  <div className="grid grid-cols-[90px_1fr] gap-3">
                    <select
                      className="h-9 rounded-md border border-slate-300 bg-white px-3 text-sm"
                      value={reviewRating}
                      onChange={event => setReviewRating(event.target.value)}
                    >
                      {[5, 4, 3, 2, 1].map(rating => (
                        <option key={rating} value={rating}>
                          {rating} star
                        </option>
                      ))}
                    </select>
                    <Input
                      value={reviewContent}
                      onChange={event => setReviewContent(event.target.value)}
                      placeholder="Short review"
                    />
                  </div>
                  <Button
                    type="submit"
                    variant="outline"
                    className="w-full"
                    disabled={savingActivity || !selectedSourceID}
                  >
                    Add public review
                  </Button>
                </form>

                <form className="space-y-3 border-t border-slate-200 pt-4" onSubmit={handleCreateCollection}>
                  <Input
                    value={collectionName}
                    onChange={event => setCollectionName(event.target.value)}
                    placeholder="Collection name"
                  />
                  <label className="flex items-center gap-2 text-sm text-slate-600">
                    <input
                      type="checkbox"
                      checked={collectionPublic}
                      onChange={event => setCollectionPublic(event.target.checked)}
                    />
                    Public collection
                  </label>
                  <Button
                    type="submit"
                    variant="outline"
                    className="w-full"
                    disabled={savingActivity || !selectedSourceID || !collectionName.trim()}
                  >
                    Create collection
                  </Button>
                </form>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Recent Notes</CardTitle>
              </CardHeader>
              <CardContent>
                {data.notes.length === 0 ? (
                  <EmptyState message="No notes yet." />
                ) : (
                  <div className="space-y-3">
                    {data.notes.slice(0, 5).map(note => (
                      <div
                        key={note.id}
                        className="rounded-lg border border-slate-200 bg-white p-3 text-sm text-slate-700"
                      >
                        <p>{note.content}</p>
                        <Button
                          variant="ghost"
                          size="xs"
                          className="mt-2 text-red-600"
                          onClick={() => {
                            if (!accessToken) return
                            runAction(() => deleteNote(accessToken, note.id), 'Failed to delete note')
                          }}
                        >
                          Delete
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Recent Reviews</CardTitle>
              </CardHeader>
              <CardContent>
                {data.reviews.length === 0 ? (
                  <EmptyState message="No reviews yet." />
                ) : (
                  <div className="space-y-3">
                    {data.reviews.slice(0, 5).map(review => (
                      <div key={review.id} className="rounded-lg border border-slate-200 bg-white p-3 text-sm text-slate-700">
                        <p className="font-medium">{review.rating}/5 stars</p>
                        {review.content && <p className="mt-1">{review.content}</p>}
                        <Button variant="ghost" size="xs" className="mt-2 text-red-600" onClick={() => {
                          if (!accessToken) return
                          runAction(() => deleteReview(accessToken, review.id), 'Failed to delete review')
                        }}>
                          Delete
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Collections</CardTitle>
              </CardHeader>
              <CardContent>
                {data.collections.length === 0 ? (
                  <EmptyState message="No collections yet." />
                ) : (
                  <div className="space-y-3">
                    {data.collections.slice(0, 5).map(collection => (
                      <div key={collection.id} className="rounded-lg border border-slate-200 bg-white p-3 text-sm text-slate-700">
                        <p className="font-medium">{collection.name}</p>
                        <p className="mt-1 text-slate-500">{collection.is_public ? 'Public' : 'Private'} · {collection.source_ids?.length || 0} sources</p>
                        <Button variant="ghost" size="xs" className="mt-2 text-red-600" onClick={() => {
                          if (!accessToken) return
                          runAction(() => deleteCollection(accessToken, collection.id), 'Failed to delete collection')
                        }}>
                          Delete
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </aside>
        </div>
      </main>
    </div>
  )
}

function StatCard({
  icon,
  label,
  value,
}: {
  icon: ReactNode
  label: string
  value: number
}) {
  return (
    <Card className="gap-2 py-4">
      <CardContent className="flex items-center gap-3 px-4">
        <div className="rounded-lg bg-emerald-50 p-2 text-emerald-700">
          {icon}
        </div>
        <div>
          <p className="text-2xl font-bold text-slate-900">{value}</p>
          <p className="text-xs text-slate-500">{label}</p>
        </div>
      </CardContent>
    </Card>
  )
}

function LoadingRow() {
  return (
    <div className="flex items-center gap-2 py-8 text-sm text-slate-500">
      <Loader2 className="h-4 w-4 animate-spin" />
      Loading platform data...
    </div>
  )
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="rounded-lg border border-dashed border-slate-300 bg-slate-50 p-6 text-center text-sm text-slate-500">
      {message}
    </div>
  )
}
