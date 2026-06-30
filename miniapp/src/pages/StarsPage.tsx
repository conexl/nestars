import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Sparkles } from 'lucide-react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { Button } from '../components/Button'
import { Input } from '../components/Input'
import { BackButton } from '../components/BackButton'
import { Loader } from '../components/Loader'
import { getCatalogStars } from '../api/client'
import type { CatalogItem } from '../types'
import type { Translations } from '../i18n'
import styles from './CatalogPage.module.css'

interface StarsPageProps {
  t: Translations
}

export function StarsPage({ t }: StarsPageProps) {
  const navigate = useNavigate()
  const [items, setItems] = useState<CatalogItem[]>([])
  const [loading, setLoading] = useState(true)
  const [selected, setSelected] = useState<string | null>(null)
  const [custom, setCustom] = useState('')

  useEffect(() => {
    getCatalogStars()
      .then(setItems)
      .catch(() => setItems([]))
      .finally(() => setLoading(false))
  }, [])

  const selectedValue = useMemo(() => {
    if (selected) {
      const item = items.find((i) => i.id === selected)
      return item ? { productId: item.id, stars: 0 } : null
    }
    const value = parseInt(custom, 10)
    if (!Number.isNaN(value) && value >= 50 && value <= 1_000_000) {
      return { productId: '', stars: value }
    }
    return null
  }, [selected, custom, items])

  const handleContinue = () => {
    if (!selectedValue) return
    navigate('/recipient', {
      state: {
        productType: 'Stars',
        productId: selectedValue.productId,
        stars: selectedValue.stars,
        name: selected ? items.find((i) => i.id === selected)?.name : custom,
      },
    })
  }

  if (loading) return <Loader text={t.loading} />

  return (
    <Page>
      <BackButton text={t.back} />

      <div className={styles.header}>
        <div className={styles.icon}>
          <Sparkles size={28} />
        </div>
        <h2 className={styles.title}>{t.buyStars}</h2>
      </div>

      <div className={styles.grid}>
        {items.map((item) => (
          <Card
            key={item.id}
            active={selected === item.id}
            onClick={() => {
              setSelected(item.id)
              setCustom('')
            }}
            className={styles.item}
          >
            <span className={styles.itemName}>{item.name}</span>
            <span className={styles.itemPrice}>
              {item.amount.toFixed(0)} {item.currency}
            </span>
          </Card>
        ))}
      </div>

      <Card className={styles.customCard}>
        <Input
          label={t.customAmount}
          type="number"
          placeholder="1000"
          value={custom}
          onChange={(e) => {
            setCustom(e.target.value)
            setSelected(null)
          }}
        />
        <p className={styles.hint}>
          {t.minStars} — {t.maxStars}
        </p>
      </Card>

      <Button fullWidth disabled={!selectedValue} onClick={handleContinue}>
        {t.continue}
      </Button>
    </Page>
  )
}
