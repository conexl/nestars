import { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { CheckCircle, Clock } from 'lucide-react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { Button } from '../components/Button'
import type { Order, PaymentMethod } from '../types'
import type { Translations } from '../i18n'
import styles from './PaymentStatusPage.module.css'

interface PaymentState {
  order: Order
  paymentUrl?: string
  method: PaymentMethod
}

interface PaymentStatusPageProps {
  t: Translations
}

export function PaymentStatusPage({ t }: PaymentStatusPageProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const state = location.state as PaymentState | null

  const [secondsLeft, setSecondsLeft] = useState(600)

  useEffect(() => {
    if (state?.method === 'Balance' || !state?.paymentUrl) return
    const timer = setInterval(() => {
      setSecondsLeft((s) => {
        if (s <= 1) clearInterval(timer)
        return Math.max(0, s - 1)
      })
    }, 1000)
    return () => clearInterval(timer)
  }, [state])

  if (!state) {
    return (
      <Page>
        <p>{t.error}</p>
        <Button onClick={() => navigate('/')}>{t.home}</Button>
      </Page>
    )
  }

  const { order, paymentUrl, method } = state
  const isBalance = method === 'Balance'

  const formatTime = (s: number) => {
    const m = Math.floor(s / 60)
    const sec = s % 60
    return `${m}:${sec.toString().padStart(2, '0')}`
  }

  return (
    <Page className={styles.page}>
      {isBalance ? (
        <>
          <CheckCircle className={styles.successIcon} size={72} />
          <h2 className={styles.title}>{t.paymentSuccess}</h2>
        </>
      ) : (
        <>
          <Clock className={styles.pendingIcon} size={72} />
          <h2 className={styles.title}>{t.invoiceValid}</h2>
          <p className={styles.timer}>{formatTime(secondsLeft)}</p>
        </>
      )}

      <Card className={styles.summary}>
        <div className={styles.row}>
          <span className={styles.label}>{t.orderNumber}</span>
          <span className={styles.value}>#{order.order_number}</span>
        </div>
        <div className={styles.row}>
          <span className={styles.label}>{t.amount}</span>
          <span className={styles.value}>
            {order.amount.toFixed(2)} {order.currency}
          </span>
        </div>
        <div className={styles.row}>
          <span className={styles.label}>{t.status}</span>
          <span className={styles.value}>{isBalance ? t.paid : t.pending}</span>
        </div>
      </Card>

      {!isBalance && paymentUrl && (
        <a href={paymentUrl} target="_blank" rel="noreferrer" className={styles.payLink}>
          <Button fullWidth>{t.payHere}</Button>
        </a>
      )}

      <Button variant="secondary" fullWidth onClick={() => navigate('/orders')}>
        {t.orders}
      </Button>
    </Page>
  )
}
