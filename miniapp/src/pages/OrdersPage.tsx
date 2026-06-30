import { useEffect, useState } from 'react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { BackButton } from '../components/BackButton'
import { Loader } from '../components/Loader'
import { getOrders } from '../api/client'
import type { Order } from '../types'
import type { Translations } from '../i18n'
import styles from './OrdersPage.module.css'

interface OrdersPageProps {
  t: Translations
}

type Filter = 'all' | 'paid' | 'pending' | 'cancelled'

const filters: Filter[] = ['all', 'paid', 'pending', 'cancelled']

export function OrdersPage({ t }: OrdersPageProps) {
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<Filter>('all')

  useEffect(() => {
    getOrders()
      .then(setOrders)
      .catch(() => setOrders([]))
      .finally(() => setLoading(false))
  }, [])

  const filtered = orders.filter((o) => (filter === 'all' ? true : o.status === filter))

  const statusLabel = (status: string) => {
    switch (status) {
      case 'paid':
        return t.paid
      case 'pending':
        return t.pending
      case 'cancelled':
        return t.cancelled
      case 'failed':
        return t.failed
      default:
        return status
    }
  }

  const statusClass = (status: string) => {
    switch (status) {
      case 'paid':
        return styles.statusPaid
      case 'pending':
        return styles.statusPending
      case 'cancelled':
      case 'failed':
        return styles.statusBad
      default:
        return ''
    }
  }

  return (
    <Page>
      <BackButton text={t.back} />

      <h2 className={styles.title}>{t.history}</h2>

      <div className={styles.filters}>
        {filters.map((f) => (
          <button
            key={f}
            type="button"
            className={filter === f ? styles.activeFilter : styles.filter}
            onClick={() => setFilter(f)}
          >
            {statusLabel(f)}
          </button>
        ))}
      </div>

      {loading ? (
        <Loader />
      ) : filtered.length === 0 ? (
        <Card className={styles.empty}>{t.history}</Card>
      ) : (
        <div className={styles.list}>
          {filtered.map((order) => (
            <Card key={order.id} className={styles.order}>
              <div className={styles.orderHeader}>
                <span className={styles.orderNumber}>#{order.order_number}</span>
                <span className={statusClass(order.status)}>{statusLabel(order.status)}</span>
              </div>
              <div className={styles.orderBody}>
                <span className={styles.orderType}>
                  {order.product_type === 'Stars' ? '⭐' : '💎'}{' '}
                  {order.product_type === 'Stars'
                    ? order.stars_amount ?? order.product_id
                    : `${order.product_id.replace(/\D/g, '')}${t.months}`}
                </span>
                <span className={styles.orderRecipient}>@{order.recipient_username}</span>
              </div>
              <div className={styles.orderFooter}>
                {order.amount.toFixed(2)} {order.currency}
              </div>
            </Card>
          ))}
        </div>
      )}
    </Page>
  )
}
