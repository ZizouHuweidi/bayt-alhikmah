import { Link, createFileRoute } from '@tanstack/react-router'
import { BookOpen, Layers3, Library, Loader2, Star, StickyNote } from 'lucide-react'
import type { ReactNode } from 'react'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  type Collection,
  type LibraryItemWithSource,
  type Note,
  type Profile,
  type Review,
  getPublicProfile,
  listPublicCollectionsByUser,
  listPublicLibrary,
  listPublicNotesByUser,
  listPublicReviewsByUser,
} from '@/lib/api'

export const Route = createFileRoute('/users/$username/profile')({
  component: PublicProfilePage,
})

type PublicData = {
  profile: Profile | null
  library: LibraryItemWithSource[]
  notes: Note[]
  reviews: Review[]
  collections: Collection[]
}

function PublicProfilePage() {
  const { username } = Route.useParams()
  const [data, setData] = useState<PublicData>({
    profile: null,
    library: [],
    notes: [],
    reviews: [],
    collections: [],
  })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    async function loadProfile() {
      setLoading(true)
      setError(null)
      try {
        const profile = await getPublicProfile(username)
        const [library, notes, reviews, collections] = await Promise.all([
          listPublicLibrary(username),
          listPublicNotesByUser(profile.user_id),
          listPublicReviewsByUser(profile.user_id),
          listPublicCollectionsByUser(profile.user_id),
        ])
        if (!cancelled) {
          setData({ profile, library, notes, reviews, collections })
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Failed to load profile')
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }
    loadProfile()
    return () => {
      cancelled = true
    }
  }, [username])

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-50">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </main>
    )
  }

  if (error || !data.profile) {
    return (
      <main className="min-h-screen bg-slate-50 px-4 py-16">
        <div className="mx-auto max-w-2xl rounded-2xl border border-slate-200 bg-white p-8 text-center shadow-sm">
          <h1 className="text-2xl font-bold text-slate-900">Profile unavailable</h1>
          <p className="mt-3 text-slate-600">
            This profile is private or does not exist.
          </p>
          <Link to="/">
            <Button className="mt-6">Go home</Button>
          </Link>
        </div>
      </main>
    )
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-10">
      <div className="mx-auto max-w-6xl">
        <section className="mb-8 rounded-3xl bg-gradient-to-br from-emerald-700 to-teal-700 p-8 text-white shadow-lg">
          <div className="flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
            <div>
              <p className="text-sm font-medium uppercase tracking-wide text-emerald-100">
                Public Knowledge Profile
              </p>
              <h1 className="mt-3 text-4xl font-bold">
                {data.profile.display_name || data.profile.username || username}
              </h1>
              {data.profile.bio && (
                <p className="mt-4 max-w-2xl text-emerald-50">{data.profile.bio}</p>
              )}
            </div>
            <div className="rounded-2xl bg-white/10 px-4 py-3 text-sm text-emerald-50">
              @{data.profile.username || username}
            </div>
          </div>
        </section>

        <div className="mb-8 grid grid-cols-2 gap-4 md:grid-cols-4">
          <Stat icon={<Library />} label="Library" value={data.library.length} />
          <Stat icon={<StickyNote />} label="Notes" value={data.notes.length} />
          <Stat icon={<Star />} label="Reviews" value={data.reviews.length} />
          <Stat icon={<Layers3 />} label="Collections" value={data.collections.length} />
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          <PublicSection
            title="Public Library"
            description="Sources this reader has chosen to share."
            empty="No public library items yet."
            items={data.library.map(item => ({
              id: item.id,
              title: item.source?.title || item.source_id,
              meta: `${item.status.replace('_', ' ')} · ${item.visibility}`,
            }))}
          />
          <PublicSection
            title="Public Notes"
            description="Notes and reflections shared by this reader."
            empty="No public notes yet."
            items={data.notes.map(note => ({
              id: note.id,
              title: note.content,
              meta: note.content_type,
            }))}
          />
          <PublicSection
            title="Public Reviews"
            description="Ratings and reviews shared publicly."
            empty="No public reviews yet."
            items={data.reviews.map(review => ({
              id: review.id,
              title: review.content || `${review.rating}/5 stars`,
              meta: `${review.rating}/5 stars`,
            }))}
          />
          <PublicSection
            title="Public Collections"
            description="Curated lists from this profile."
            empty="No public collections yet."
            items={data.collections.map(collection => ({
              id: collection.id,
              title: collection.name,
              meta: `${collection.source_ids?.length || 0} sources`,
            }))}
          />
        </div>
      </div>
    </main>
  )
}

function Stat({
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
        <div className="rounded-lg bg-emerald-50 p-2 text-emerald-700">{icon}</div>
        <div>
          <p className="text-2xl font-bold text-slate-900">{value}</p>
          <p className="text-xs text-slate-500">{label}</p>
        </div>
      </CardContent>
    </Card>
  )
}

function PublicSection({
  title,
  description,
  empty,
  items,
}: {
  title: string
  description: string
  empty: string
  items: Array<{ id: string; title: string; meta: string }>
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <BookOpen className="h-5 w-5 text-emerald-600" />
          {title}
        </CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        {items.length === 0 ? (
          <div className="rounded-lg border border-dashed border-slate-300 bg-slate-50 p-6 text-center text-sm text-slate-500">
            {empty}
          </div>
        ) : (
          <div className="space-y-3">
            {items.slice(0, 8).map(item => (
              <div key={item.id} className="rounded-lg border border-slate-200 bg-white p-4">
                <p className="font-medium text-slate-900">{item.title}</p>
                <p className="mt-1 text-sm text-slate-500">{item.meta}</p>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
