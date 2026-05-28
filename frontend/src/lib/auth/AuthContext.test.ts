import { describe, expect, it } from 'vitest'
import { accessTokenFromAuthResponse } from './AuthContext'

describe('accessTokenFromAuthResponse', () => {
  it('reads the backend nested access token shape', () => {
    expect(
      accessTokenFromAuthResponse({ tokens: { access_token: 'fresh-token' } })
    ).toBe('fresh-token')
  })

  it('rejects missing or legacy token shapes', () => {
    expect(accessTokenFromAuthResponse({ access_token: 'legacy-token' })).toBeNull()
    expect(accessTokenFromAuthResponse({ tokens: {} })).toBeNull()
    expect(accessTokenFromAuthResponse(null)).toBeNull()
  })
})
