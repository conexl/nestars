import { useEffect, useMemo } from 'react'
import {
  bindMiniAppCssVars,
  bindThemeParamsCssVars,
  bindViewportCssVars,
  expandViewport,
  initData,
  isMiniAppDark,
  setMiniAppHeaderColor,
  useSignal,
} from '@telegram-apps/sdk-react'
import type { TelegramUser } from '../types'

/**
 * Safely restores the Telegram init data state.
 * Outside Telegram the launch params are missing/invalid, so restore() throws.
 * We catch it and leave the signals empty – the app can still render a
 * graceful placeholder or run in dev mode.
 */
function restoreInitData(): void {
  try {
    initData.restore()
  } catch (error) {
    console.warn('Failed to restore Telegram init data:', error)
  }
}

export function useTelegram() {
  // Init data must be restored once before its signals become usable.
  useEffect(() => {
    restoreInitData()
  }, [])

  const rawInitData = useSignal(initData.raw)
  const tgUser = useSignal(initData.user) as TelegramUser | undefined
  const startParam = useSignal(initData.startParam)
  const isDark = useSignal(isMiniAppDark)

  useEffect(() => {
    try {
      expandViewport()
      bindMiniAppCssVars()
      bindThemeParamsCssVars()
      bindViewportCssVars()
      setMiniAppHeaderColor('bg_color')
    } catch (error) {
      console.warn('Error in Telegram CSS/theme init:', error)
    }
  }, [])

  return useMemo(
    () => ({
      rawInitData: rawInitData || '',
      user: tgUser,
      startParam: startParam || '',
      appearance: (isDark ? 'dark' : 'light') as 'dark' | 'light',
    }),
    [rawInitData, tgUser, startParam, isDark],
  )
}
