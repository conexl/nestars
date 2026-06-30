import styles from './Loader.module.css'

interface LoaderProps {
  text?: string
}

export function Loader({ text }: LoaderProps) {
  return (
    <div className={styles.loader}>
      <span className={styles.spinner} />
      {text && <span className={styles.text}>{text}</span>}
    </div>
  )
}
