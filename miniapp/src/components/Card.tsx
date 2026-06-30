import { clsx } from 'clsx'
import type { ReactNode } from 'react'
import styles from './Card.module.css'

interface CardProps {
  children: ReactNode
  className?: string
  onClick?: () => void
  active?: boolean
  glow?: boolean
}

export function Card({ children, className, onClick, active, glow }: CardProps) {
  return (
    <div
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      onClick={onClick}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') onClick?.()
      }}
      className={clsx(
        styles.card,
        active && styles.active,
        onClick && styles.clickable,
        glow && styles.glow,
        className,
      )}
    >
      {children}
    </div>
  )
}
