import { describe, expect, it } from 'vitest'
import { normalizeDashboardData } from './dashboard'

describe('normalizeDashboardData', () => {
  it('turns null dashboard lists into empty arrays', () => {
    expect(
      normalizeDashboardData({
        sources: null as never,
        library: null as never,
        notes: null as never,
        reviews: null as never,
        collections: null as never,
      })
    ).toEqual({
      sources: [],
      library: [],
      notes: [],
      reviews: [],
      collections: [],
    })
  })
})
