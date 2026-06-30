import { useNavigate } from 'react-router-dom'
import { Gift, Users, Star, Wallet } from 'lucide-react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { Avatar } from '../components/Avatar'
import { useUser } from '../contexts/UserContext'
import { useTelegram } from '../hooks/useTelegram'
import type { Translations } from '../i18n'
import styles from './ProfilePage.module.css'

interface ProfilePageProps {
  t: Translations
}

export function ProfilePage({ t }: ProfilePageProps) {
  const navigate = useNavigate()
  const { user, profile } = useUser()
  const { user: tgUser } = useTelegram()

  const initials = ((tgUser?.first_name?.[0] || '') + (tgUser?.last_name?.[0] || '')).toUpperCase()
    || tgUser?.first_name?.[0]?.toUpperCase()
    || '?'

  const menuItems = [
    { icon: Gift, label: t.referral, onClick: () => navigate('/referral') },
    { icon: Wallet, label: t.history, onClick: () => navigate('/orders') },
  ]

  return (
    <Page>
      <div className={styles.header}>
        <Avatar initials={initials} src={tgUser?.photo_url} size="lg" />
        <div className={styles.info}>
          <div className={styles.name}>{tgUser?.first_name || user?.username || 'User'}</div>
          <div className={styles.username}>@{tgUser?.username || user?.username}</div>
        </div>
      </div>

      <Card className={styles.balanceCard}>
        <div className={styles.balanceRow}>
          <Wallet size={20} />
          <span>{t.balance}</span>
        </div>
        <div className={styles.balanceValue}>${profile?.balance_usd.toFixed(2) ?? '0.00'}</div>
        <div className={styles.balanceTon}>{profile?.balance_ton.toFixed(2) ?? '0.00'} TON</div>
      </Card>

      <Card className={styles.statsCard}>
        <div className={styles.stat}>
          <Star size={20} className={styles.statIcon} />
          <span className={styles.statValue}>{profile?.bought_stars ?? 0}</span>
          <span className={styles.statLabel}>Stars</span>
        </div>
        <div className={styles.statDivider} />
        <div className={styles.stat}>
          <Users size={20} className={styles.statIcon} />
          <span className={styles.statValue}>{profile?.referral_count ?? 0}</span>
          <span className={styles.statLabel}>{t.invited}</span>
        </div>
      </Card>

      <div className={styles.menu}>
        {menuItems.map((item) => (
          <button key={item.label} type="button" className={styles.menuItem} onClick={item.onClick}>
            <div className={styles.menuLeft}>
              <item.icon size={20} />
              <span>{item.label}</span>
            </div>
          </button>
        ))}
      </div>
    </Page>
  )
}
