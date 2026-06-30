import { forwardRef } from 'react'
import { clsx } from 'clsx'
import styles from './Input.module.css'

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className, ...props }, ref) => {
    return (
      <label className={styles.wrapper}>
        {label && <span className={styles.label}>{label}</span>}
        <input
          ref={ref}
          className={clsx(styles.input, error && styles.error, className)}
          {...props}
        />
        {error && <span className={styles.errorText}>{error}</span>}
      </label>
    )
  },
)
