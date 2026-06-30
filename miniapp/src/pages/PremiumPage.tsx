import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Gem } from 'lucide-react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { Button } from '../components/Button'
import { BackButton } from '../components/BackButton'
import { Loader } from '../components/Loader'
import { getCatalogPremium } from '../api/client'
import type { CatalogItem } from '../types'
import type { Translations } from '../i18n'
import styles from './CatalogPage.module.css'

interface PremiumPageProps {
  t: Translations
}

export function PremiumPage({ t }: PremiumPageProps) {
  const navigate = useNavigate()
  const [items, setItems] = useState<CatalogItem[]>([])
  const [loading, setLoading] = useState(true)
  const [selected, setSelected] = useState<string | null>(null)

  useEffect(() => {
    getCatalogPremium()
      .then(setItems)
      .catch(() => setItems([]))
      .finally(() => setLoading(false))
  }, [])

  const selectedItem = useMemo(() => items.find((i) => i.id === selected), [items, selected])

  const handleContinue = () => {
    if (!selectedItem) return
    navigate('/recipient', {
      state: {
        productType: 'Premium',
        productId: selectedItem.id,
        stars: 0,
        name: selectedItem.name,
      },
    })
  }

  const badge = (index: number) => {
    if (index === items.length - 1) return 'best'
    if (index === 1) return 'value'
    return null
  }

  if (loading) return <Loader text={t.loading} />

  return (
    <Page>
      <BackButton text={t.back} />

      <div className={styles.header}>
        <div className={styles.icon} style={{ background: 'rgba(175, 81, 222, 0.12)', color: 'var(--purple)' }}>
          <Gem size={28} />
        </div>
        <h2 className={styles.title}>{t.buyPremium}</h2>
      </div>

      <div className={styles.grid}>
        {items.map((item, index) => (
          <Card
            key={item.id}
            active={selected === item.id}
            onClick={() => setSelected(item.id)}
            className={styles.plan}
          >
            <div className={styles.planInfo}>
              <span className={styles.planName}>{item.name}</span>
              <span className={styles.planPrice}>
                {item.amount.toFixed(0)} {item.currency}
              </span>
            </div>
            {badge(index) && (
              <span className={styles.badge}>
                {badge(index) === 'best' ? 'Best' : 'Value'}
              </span>
            )}
          </Card>
        ))}
      </div>

      <Button fullWidth disabled={!selectedItem} onClick={handleContinue}>
        {t.continue}
      </Button>
    </Page>
  )
}
