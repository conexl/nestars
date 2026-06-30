import { useState } from 'react'
import { Copy, Check, Gift } from 'lucide-react'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import { Button } from '../components/Button'
import { BackButton } from '../components/BackButton'
import { useUser } from '../contexts/UserContext'
import type { Translations } from '../i18n'
import styles from './ReferralPage.module.css'

interface ReferralPageProps {
  t: Translations
}

export function ReferralPage({ t }: ReferralPageProps) {
  const { profile } = useUser()
  const [copied, setCopied] = useState(false)

  const botUsername = 'nestars_bot'
  const referralCode = profile?.referral_code || ''
  const link = `https://t.me/${botUsername}?start=${referralCode}`

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(link)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    } catch {
      // ignore
    }
  }

  return (
    <Page>
      <BackButton text={t.back} />

      <div className={styles.header}>
        <div className={styles.icon}>
          <Gift size={28} />
        </div>
        <h2 className={styles.title}>{t.referral}</h2>
      </div>

      <Card className={styles.linkCard}>
        <div className={styles.link}>{link}</div>
        <Button fullWidth variant="secondary" onClick={handleCopy}>
          {copied ? <Check size={18} /> : <Copy size={18} />}
          {copied ? 'OK' : t.copy}
        </Button>
      </Card>

      <div className={styles.stats}>
        <Card className={styles.stat}>
          <span className={styles.statValue}>{profile?.referral_bonus_stars ?? 0}</span>
          <span className={styles.statLabel}>{t.earned}</span>
        </Card>
        <Card className={styles.stat}>
          <span className={styles.statValue}>{profile?.referral_count ?? 0}</span>
          <span className={styles.statLabel}>{t.invited}</span>
        </Card>
        <Card className={styles.stat}>
          <span className={styles.statValue}>{profile?.bonus_stars ?? 0}</span>
          <span className={styles.statLabel}>{t.bonusBalance}</span>
        </Card>
      </div>
    </Page>
  )
}
