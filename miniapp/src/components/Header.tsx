import { useMemo } from 'react'
import { Wallet } from 'lucide-react'
import { Avatar } from './Avatar'
import type { Profile, TelegramUser } from '../types'
import styles from './Header.module.css'

interface HeaderProps {
  user?: TelegramUser
  profile?: Profile | null
}

export function Header({ user, profile }: HeaderProps) {
  const initials = useMemo(() => {
    if (!user) return '?'
    return (user.first_name?.[0] || '') + (user.last_name?.[0] || '') || user.first_name?.[0] || '?'
  }, [user])

  return (
    <header className={styles.header}>
      <div className={styles.brand}>
        <span className={styles.brandName}>neStars</span>
      </div>

      <div className={styles.right}>
        <div className={styles.balance}>
          <Wallet size={16} />
          <span>${profile?.balance_usd.toFixed(2) ?? '0.00'}</span>
        </div>
        <Avatar initials={initials} src={user?.photo_url} />
      </div>
    </header>
  )
}
