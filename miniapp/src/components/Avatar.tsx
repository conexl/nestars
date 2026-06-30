import { clsx } from 'clsx'
import styles from './Avatar.module.css'

interface AvatarProps {
  initials?: string
  src?: string
  alt?: string
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

export function Avatar({ initials, src, alt, size = 'md', className }: AvatarProps) {
  return (
    <div className={clsx(styles.avatar, styles[size], className)}>
      {src ? <img src={src} alt={alt} className={styles.img} /> : <span>{initials}</span>}
    </div>
  )
}
