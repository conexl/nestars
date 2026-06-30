import { Sparkles, Gem, ChevronRight } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { Card } from '../components/Card'
import { Page } from '../components/Page'
import type { Translations } from '../i18n'
import styles from './HomePage.module.css'

interface HomePageProps {
  t: Translations
}

export function HomePage({ t }: HomePageProps) {
  const navigate = useNavigate()

  return (
    <Page>
      <div className={styles.actions}>
        <Card onClick={() => navigate('/stars')} className={styles.actionCard}>
          <div className={styles.actionIcon}>
            <Sparkles size={28} />
          </div>
          <div className={styles.actionContent}>
            <span className={styles.actionTitle}>{t.buyStars}</span>
            <span className={styles.actionDesc}>{t.minStars} — {t.maxStars}</span>
          </div>
          <ChevronRight size={20} className={styles.actionArrow} />
        </Card>
        <Card onClick={() => navigate('/premium')} className={styles.actionCard}>
          <div className={styles.actionIconPremium}>
            <Gem size={28} />
          </div>
          <div className={styles.actionContent}>
            <span className={styles.actionTitle}>{t.buyPremium}</span>
            <span className={styles.actionDesc}>3, 6, 12 {t.months}</span>
          </div>
          <ChevronRight size={20} className={styles.actionArrow} />
        </Card>
      </div>
    </Page>
  )
}
