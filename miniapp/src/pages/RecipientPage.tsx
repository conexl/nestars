import { useMemo, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { User } from 'lucide-react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { Button } from '../components/Button'
import { Input } from '../components/Input'
import { BackButton } from '../components/BackButton'
import { useTelegram } from '../hooks/useTelegram'
import type { Translations } from '../i18n'
import styles from './RecipientPage.module.css'

interface LocationState {
  productType: 'Stars' | 'Premium'
  productId: string
  stars: number
  name: string
}

interface RecipientPageProps {
  t: Translations
}

function parseUsername(input: string): string {
  const trimmed = input.trim()
  if (trimmed.startsWith('https://t.me/')) {
    return trimmed.replace('https://t.me/', '').replace('/', '')
  }
  return trimmed.replace(/^@/, '')
}

export function RecipientPage({ t }: RecipientPageProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const { user } = useTelegram()
  const state = location.state as LocationState | null

  const [mode, setMode] = useState<'self' | 'friend'>('self')
  const [username, setUsername] = useState('')
  const [error, setError] = useState('')

  const selfUsername = useMemo(() => user?.username, [user])

  const canContinue = useMemo(() => {
    if (mode === 'self') return !!selfUsername
    const parsed = parseUsername(username)
    return parsed.length >= 2
  }, [mode, selfUsername, username])

  const handleContinue = () => {
    setError('')
    const recipient = mode === 'self' ? selfUsername : parseUsername(username)
    if (!recipient) {
      setError(t.invalidUsername)
      return
    }
    if (!state) {
      navigate('/')
      return
    }
    navigate('/checkout', {
      state: {
        ...state,
        recipient,
      },
    })
  }

  return (
    <Page>
      <BackButton text={t.back} />

      <div className={styles.header}>
        <div className={styles.icon}>
          <User size={28} />
        </div>
        <h2 className={styles.title}>{t.recipient}</h2>
      </div>

      <div className={styles.tabs}>
        <button
          type="button"
          className={mode === 'self' ? styles.activeTab : styles.tab}
          onClick={() => setMode('self')}
        >
          {t.forMyself}
        </button>
        <button
          type="button"
          className={mode === 'friend' ? styles.activeTab : styles.tab}
          onClick={() => setMode('friend')}
        >
          {t.forFriend}
        </button>
      </div>

      {mode === 'self' ? (
        <Card className={styles.selfCard}>
          <span className={styles.selfLabel}>@{selfUsername || '—'}</span>
          {!selfUsername && <span className={styles.selfHint}>{t.invalidUsername}</span>}
        </Card>
      ) : (
        <Input
          label={t.enterUsername}
          placeholder="@username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          error={error}
        />
      )}

      <Button fullWidth disabled={!canContinue} onClick={handleContinue}>
        {t.continue}
      </Button>
    </Page>
  )
}
