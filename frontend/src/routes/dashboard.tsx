import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import type { FormEvent } from 'react'
import { useEffect, useState } from 'react'
import {
  AddBookCard,
  CollectionsCard,
  type DashboardData,
  DashboardHeader,
  DashboardIntro,
  DashboardStats,
  LibraryPanel,
  RecentNotesCard,
  RecentReviewsCard,
  SourceActivityCard,
  SourcesPanel,
} from '@/components/dashboard/dashboard-panels'
import {
  addLibraryItem,
  createBook,
  createCollection,
  createNote,
  createReview,
  listCollections,
  listLibrary,
  listNotes,
  listReviews,
  listSources,
} from '@/lib/api'
import { useAuth } from '@/lib/auth/AuthContext'

export const Route = createFileRoute('/dashboard')({
  component: DashboardPage,
})

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
      <DashboardHeader
        displayName={user.username || user.email}
        onLogout={handleLogout}
      />

      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        <DashboardIntro />

        {error && (
          <div className="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        <DashboardStats data={data} />

        <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_380px]">
          <div className="space-y-6">
            <LibraryPanel
              accessToken={accessToken}
              items={data.library}
              loading={loadingData}
              runAction={runAction}
            />

            <SourcesPanel
              accessToken={accessToken}
              librarySourceIDs={librarySourceIDs}
              loading={loadingData}
              onAdd={loadDashboard}
              onError={setError}
              sources={data.sources}
            />
          </div>

          <aside className="space-y-6">
            <AddBookCard
              author={bookAuthor}
              isbn={bookISBN}
              saving={savingBook}
              title={bookTitle}
              onAuthorChange={setBookAuthor}
              onISBNChange={setBookISBN}
              onSubmit={handleCreateBook}
              onTitleChange={setBookTitle}
            />

            <SourceActivityCard
              collectionName={collectionName}
              collectionPublic={collectionPublic}
              librarySources={librarySources}
              noteContent={noteContent}
              notePublic={notePublic}
              reviewContent={reviewContent}
              reviewRating={reviewRating}
              saving={savingActivity}
              selectedSourceID={selectedSourceID}
              onCollectionNameChange={setCollectionName}
              onCollectionPublicChange={setCollectionPublic}
              onNoteContentChange={setNoteContent}
              onNotePublicChange={setNotePublic}
              onReviewContentChange={setReviewContent}
              onReviewRatingChange={setReviewRating}
              onSelectedSourceChange={setSelectedSourceID}
              onSubmitCollection={handleCreateCollection}
              onSubmitNote={handleCreateNote}
              onSubmitReview={handleCreateReview}
            />

            <RecentNotesCard
              accessToken={accessToken}
              notes={data.notes}
              runAction={runAction}
            />
            <RecentReviewsCard
              accessToken={accessToken}
              reviews={data.reviews}
              runAction={runAction}
            />
            <CollectionsCard
              accessToken={accessToken}
              collections={data.collections}
              runAction={runAction}
            />
          </aside>
        </div>
      </main>
    </div>
  )
}
