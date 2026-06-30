import { createContext, useContext, useState, type ReactNode } from 'react'
import type { Profile, User } from '../types'

interface UserContextValue {
  user: User | null
  profile: Profile | null
  setUserData: (user: User, profile: Profile) => void
  setRefreshFn: (fn: () => Promise<void>) => void
  refresh: () => Promise<void>
}

const UserContext = createContext<UserContextValue | null>(null)

export function UserProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [profile, setProfile] = useState<Profile | null>(null)
  const [refresh, setRefresh] = useState<() => Promise<void>>(async () => {})

  const setUserData = (u: User, p: Profile) => {
    setUser(u)
    setProfile(p)
  }

  const setRefreshFn = (fn: () => Promise<void>) => {
    setRefresh(() => fn)
  }

  return (
    <UserContext.Provider value={{ user, profile, setUserData, setRefreshFn, refresh }}>
      {children}
    </UserContext.Provider>
  )
}

export function useUser() {
  const ctx = useContext(UserContext)
  if (!ctx) throw new Error('useUser must be used inside UserProvider')
  return ctx
}
