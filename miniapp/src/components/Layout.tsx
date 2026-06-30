import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { Home, ListOrdered, User, Gift, Star, LogOut } from 'lucide-react'
import { clsx } from 'clsx'
import { useUser } from '../contexts/UserContext'
import { useTelegram } from '../hooks/useTelegram'
import { Header } from './Header'
import type { Translations } from '../i18n'
import styles from './Layout.module.css'

interface LayoutProps {
  t: Translations
}

export function Layout({ t }: LayoutProps) {
  const navigate = useNavigate()
  const { pathname } = useLocation()
  const { profile } = useUser()
  const { user: tgUser } = useTelegram()

  const mainLinks = [
    { to: '/', icon: Home, label: t.home },
    { to: '/orders', icon: ListOrdered, label: t.orders },
    { to: '/profile', icon: User, label: t.profile },
  ]

  const extraLinks = [
    { to: '/referral', icon: Gift, label: t.referral },
    { to: '/withdraw', icon: Star, label: t.withdraw },
  ]

  const isActive = (to: string) => pathname === to

  const NavLink = ({ to, icon: Icon, label }: { to: string; icon: typeof Home; label: string }) => (
    <button
      type="button"
      className={clsx(styles.navItem, isActive(to) && styles.navItemActive)}
      onClick={() => navigate(to)}
    >
      <Icon size={22} />
      <span>{label}</span>
    </button>
  )

  return (
    <div className={styles.layout}>
      <aside className={styles.sidebar}>
        <nav className={styles.sidebarNav}>
          {mainLinks.map((link) => (
            <NavLink key={link.to} {...link} />
          ))}
        </nav>

        <div className={styles.sidebarSection}>
          <span className={styles.sidebarSectionTitle}>{t.referral}</span>
          {extraLinks.map((link) => (
            <NavLink key={link.to} {...link} />
          ))}
        </div>

        <div className={styles.sidebarFooter}>
          <button
            type="button"
            className={styles.logout}
            onClick={() => navigate('/')}
          >
            <LogOut size={18} />
            {t.logout}
          </button>
        </div>
      </aside>

      <div className={styles.content}>
        {pathname !== '/profile' && <Header user={tgUser} profile={profile} />}
        <main className={styles.main}>
          <Outlet />
        </main>
      </div>

      <nav className={styles.bottomNav}>
        {mainLinks.map((link) => (
          <NavLink key={link.to} {...link} />
        ))}
      </nav>
    </div>
  )
}
