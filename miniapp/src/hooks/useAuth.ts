import { useCallback, useEffect, useState } from 'react'
import { auth, getProfile } from '../api/client'
import type { Profile, User } from '../types'

interface AuthState {
  loading: boolean
  error: string | null
  user: User | null
  profile: Profile | null
}

const DISABLE_AUTH = import.meta.env.VITE_DISABLE_AUTH === 'true'

export function useAuth(rawInitData: string) {
  const [state, setState] = useState<AuthState>({
    loading: true,
    error: null,
    user: null,
    profile: null,
  })

  const refresh = useCallback(async () => {
    try {
      let data
      if (rawInitData || DISABLE_AUTH) {
        data = await auth(rawInitData)
      } else {
        data = await getProfile()
      }
      setState({ loading: false, error: null, user: data.user, profile: data.profile })
    } catch (err) {
      setState({
        loading: false,
        error: err instanceof Error ? err.message : 'auth error',
        user: null,
        profile: null,
      })
    }
  }, [rawInitData])

  useEffect(() => {
    refresh()
  }, [refresh])

  return { ...state, refresh }
}
