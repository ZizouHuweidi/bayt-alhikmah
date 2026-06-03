import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  addLibraryItem,
  apiRequest,
  createCollection,
  createNote,
  createReview,
  getPublicProfile,
  listPublicLibrary,
  updateProfile,
} from './api'

afterEach(() => {
  vi.restoreAllMocks()
})

describe('apiRequest', () => {
  it('returns JSON for successful responses', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(async () =>
        Response.json({ items: [1, 2, 3] }, { status: 200 })
      )
    )

    await expect(apiRequest('/test')).resolves.toEqual({ items: [1, 2, 3] })
  })

  it('uses backend error messages when available', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(async () => Response.json({ error: 'source not found' }, { status: 404 }))
    )

    await expect(apiRequest('/missing')).rejects.toThrow('source not found')
  })

  it('adds bearer authorization when an access token is provided', async () => {
    const fetchMock = vi.fn(async () => Response.json({ ok: true }, { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await apiRequest('/private', { accessToken: 'token-123' })

    const init = fetchMock.mock.calls[0]?.[1] as RequestInit
    expect((init.headers as Headers).get('Authorization')).toBe('Bearer token-123')
  })

  it('encodes usernames for public profile and library calls', async () => {
    const fetchMock = vi.fn(async () => Response.json({ ok: true }, { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await getPublicProfile('reader name')
    await listPublicLibrary('reader name')

    const urls = fetchMock.mock.calls.map(call => String(call[0]))
    expect(urls[0]).toContain('/users/reader%20name/profile')
    expect(urls[1]).toContain('/users/reader%20name/library')
  })
})

describe('MVP API helpers', () => {
  it('posts note payloads to the authenticated notes endpoint', async () => {
    const fetchMock = vi.fn(async () => Response.json({ id: 'note-1' }, { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await createNote('token', { source_id: 'source-1', content: 'Note', content_type: 'note' })

    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toContain('/api/notes')
    expect(init.method).toBe('POST')
    expect(JSON.parse(init.body as string)).toMatchObject({ source_id: 'source-1' })
  })

  it('posts review, collection, library, and profile payloads to expected endpoints', async () => {
    const fetchMock = vi.fn(async () => Response.json({ ok: true }, { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await createReview('token', { source_id: 'source-1', rating: 5, is_public: true })
    await createCollection('token', { name: 'List', source_ids: ['source-1'] })
    await addLibraryItem('token', 'source-1')
    await updateProfile('token', { display_name: 'Reader', public_profile: true })

    const urls = fetchMock.mock.calls.map(call => String(call[0]))
    expect(urls.some(url => url.endsWith('/api/reviews'))).toBe(true)
    expect(urls.some(url => url.endsWith('/api/collections'))).toBe(true)
    expect(urls.some(url => url.endsWith('/api/library/items'))).toBe(true)
    expect(urls.some(url => url.endsWith('/api/profile'))).toBe(true)
  })
})
