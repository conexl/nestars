import { useEffect } from 'react'
import { Navigate, Route, Routes } from 'react-router-dom'
import { useTelegram } from './hooks/useTelegram'
import { useAuth } from './hooks/useAuth'
import { useI18n } from './hooks/useI18n'
import { UserProvider, useUser } from './contexts/UserContext'
import { Layout } from './components/Layout'
import { Loader } from './components/Loader'
import { HomePage } from './pages/HomePage'
import { StarsPage } from './pages/StarsPage'
import { PremiumPage } from './pages/PremiumPage'
import { RecipientPage } from './pages/RecipientPage'
import { CheckoutPage } from './pages/CheckoutPage'
import { PaymentStatusPage } from './pages/PaymentStatusPage'
import { OrdersPage } from './pages/OrdersPage'
import { ProfilePage } from './pages/ProfilePage'
import { ReferralPage } from './pages/ReferralPage'
import { WithdrawPage } from './pages/WithdrawPage'
import styles from './App.module.css'

function AppContent() {
  const { rawInitData, user: tgUser } = useTelegram()
  const { user, profile, loading, error, refresh } = useAuth(rawInitData)
  const { setUserData, setRefreshFn } = useUser()
  const { t } = useI18n(tgUser?.language_code)

  useEffect(() => {
    if (user && profile) {
      setUserData(user, profile)
    }
  }, [user, profile, setUserData])

  useEffect(() => {
    setRefreshFn(refresh)
  }, [refresh, setRefreshFn])

  const content = (() => {
    if (loading) {
      return (
        <div className={styles.centered}>
          <Loader text={t.loading} />
        </div>
      )
    }

    if (error) {
      return (
        <div className={styles.centered}>
          <p className={styles.error}>{t.error}: {error}</p>
        </div>
      )
    }

    return (
      <Routes>
        <Route element={<Layout t={t} />}>
          <Route path="/" element={<HomePage t={t} />} />
          <Route path="/stars" element={<StarsPage t={t} />} />
          <Route path="/premium" element={<PremiumPage t={t} />} />
          <Route path="/recipient" element={<RecipientPage t={t} />} />
          <Route path="/checkout" element={<CheckoutPage t={t} />} />
          <Route path="/payment-status" element={<PaymentStatusPage t={t} />} />
          <Route path="/orders" element={<OrdersPage t={t} />} />
          <Route path="/profile" element={<ProfilePage t={t} />} />
          <Route path="/referral" element={<ReferralPage t={t} />} />
          <Route path="/withdraw" element={<WithdrawPage t={t} />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    )
  })()

  return <div className={styles.app}>{content}</div>
}

function App() {
  return (
    <UserProvider>
      <AppContent />
    </UserProvider>
  )
}

export default App
