import { clsx } from 'clsx'
import type { ReactNode } from 'react'
import styles from './Page.module.css'

interface PageProps {
  children: ReactNode
  className?: string
}

export function Page({ children, className }: PageProps) {
  return <div className={clsx(styles.page, className)}>{children}</div>
}
