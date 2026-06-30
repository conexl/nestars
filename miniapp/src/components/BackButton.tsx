import { ChevronLeft } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import styles from './BackButton.module.css'

interface BackButtonProps {
  text: string
}

export function BackButton({ text }: BackButtonProps) {
  const navigate = useNavigate()

  return (
    <button type="button" className={styles.back} onClick={() => navigate(-1)}>
      <ChevronLeft size={20} />
      {text}
    </button>
  )
}
