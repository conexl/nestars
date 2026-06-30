import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import './index.css'
import App from './App'

window.addEventListener('error', (event) => {
  console.error('[global error]', event.message, event.filename, event.lineno, event.colno, event.error)
})

window.addEventListener('unhandledrejection', (event) => {
  console.error('[global unhandledrejection]', event.reason)
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>,
)
