import { useMemo, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { Wallet } from 'lucide-react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { Button } from '../components/Button'
import { BackButton } from '../components/BackButton'
import { useUser } from '../contexts/UserContext'
import { createOrder, payOrder } from '../api/client'
import type { PaymentMethod } from '../types'
import type { Translations } from '../i18n'
import styles from './CheckoutPage.module.css'

interface CheckoutState {
  productType: 'Stars' | 'Premium'
  productId: string
  stars: number
  name: string
  recipient: string
}

interface CheckoutPageProps {
  t: Translations
}

const methods: { key: PaymentMethod; label: string }[] = [
  { key: 'Balance', label: 'Balance' },
  { key: 'CryptoBot', label: 'Crypto Bot' },
  { key: 'Yookassa', label: 'Yookassa' },
  { key: 'Heleket', label: 'Heleket' },
]

export function CheckoutPage({ t }: CheckoutPageProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const { profile } = useUser()
  const state = location.state as CheckoutState | null

  const [method, setMethod] = useState<PaymentMethod>('CryptoBot')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const canPayBalance = useMemo(
    () => method !== 'Balance' || (profile?.balance_usd ?? 0) >= (state?.productType ? 0 : 1),
    [method, profile, state],
  )

  const handlePay = async () => {
    if (!state) return
    setLoading(true)
    setError('')
    try {
      const order = await createOrder({
        username: state.recipient,
        product_type: state.productType,
        product_id: state.productId,
        stars: state.stars,
        currency: 'RUB',
      })

      if (method === 'Balance' && (profile?.balance_usd ?? 0) < order.amount) {
        setError(t.insufficientBalance)
        setLoading(false)
        return
      }

      const payment = await payOrder(order.id, method)

      navigate('/payment-status', {
        state: {
          order,
          paymentUrl: payment.payment_url,
          method,
        },
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : t.error)
    } finally {
      setLoading(false)
    }
  }

  if (!state) {
    return (
      <Page>
        <BackButton text={t.back} />
        <p>{t.error}</p>
      </Page>
    )
  }

  return (
    <Page>
      <BackButton text={t.back} />

      <h2 className={styles.title}>{t.checkout}</h2>

      <Card className={styles.summary}>
        <div className={styles.row}>
          <span className={styles.label}>{state.productType === 'Stars' ? t.stars : t.premium}</span>
          <span className={styles.value}>{state.name}</span>
        </div>
        <div className={styles.row}>
          <span className={styles.label}>{t.recipient}</span>
          <span className={styles.value}>@{state.recipient}</span>
        </div>
      </Card>

      <div>
        <h3 className={styles.subtitle}>{t.paymentMethod}</h3>
        <div className={styles.methods}>
          {methods.map((m) => (
            <Card
              key={m.key}
              active={method === m.key}
              onClick={() => setMethod(m.key)}
              className={styles.method}
            >
              <span className={styles.methodName}>{m.label}</span>
              {m.key === 'Balance' && (
                <span className={styles.methodHint}>${profile?.balance_usd.toFixed(2)}</span>
              )}
            </Card>
          ))}
        </div>
      </div>

      {error && <p className={styles.error}>{error}</p>}

      <Button
        fullWidth
        loading={loading}
        disabled={loading || (method === 'Balance' && !canPayBalance)}
        onClick={handlePay}
      >
        <Wallet size={18} />
        {t.pay}
      </Button>
    </Page>
  )
}
