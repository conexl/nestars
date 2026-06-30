import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Star } from 'lucide-react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { Button } from '../components/Button'
import { Input } from '../components/Input'
import { BackButton } from '../components/BackButton'
import { useUser } from '../contexts/UserContext'
import { withdraw } from '../api/client'
import type { Translations } from '../i18n'
import styles from './WithdrawPage.module.css'

interface WithdrawPageProps {
  t: Translations
}

function parseUsername(input: string): string {
  const trimmed = input.trim()
  if (trimmed.startsWith('https://t.me/')) {
    return trimmed.replace('https://t.me/', '').replace('/', '')
  }
  return trimmed.replace(/^@/, '')
}

export function WithdrawPage({ t }: WithdrawPageProps) {
  const navigate = useNavigate()
  const { profile, refresh } = useUser()
  const [username, setUsername] = useState('')
  const [amount, setAmount] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleWithdraw = async () => {
    setError('')
    const parsed = parseUsername(username)
    const stars = parseInt(amount, 10)

    if (parsed.length < 2) {
      setError(t.invalidUsername)
      return
    }
    if (Number.isNaN(stars) || stars < 50 || stars > 1_000_000) {
      setError(t.invalidAmount)
      return
    }

    setLoading(true)
    try {
      await withdraw({ username: parsed, stars })
      await refresh()
      navigate('/profile')
    } catch (err) {
      setError(err instanceof Error ? err.message : t.error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Page>
      <BackButton text={t.back} />

      <div className={styles.header}>
        <div className={styles.icon}>
          <Star size={28} />
        </div>
        <h2 className={styles.title}>{t.withdraw}</h2>
      </div>

      <Card className={styles.balanceCard}>
        <span className={styles.balanceLabel}>{t.balance}</span>
        <span className={styles.balanceValue}>{profile?.bought_stars ?? 0} ⭐</span>
      </Card>

      <Input
        label={t.recipient}
        placeholder="@username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />

      <Input
        label={t.withdrawAmount}
        type="number"
        placeholder="1000"
        value={amount}
        onChange={(e) => setAmount(e.target.value)}
      />

      {error && <p className={styles.error}>{error}</p>}

      <Button fullWidth loading={loading} onClick={handleWithdraw}>
        {t.withdraw}
      </Button>
    </Page>
  )
}
