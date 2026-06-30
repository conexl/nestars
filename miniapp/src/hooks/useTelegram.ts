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
  console.log('[useTelegram] before initData.restore()')
  try {
    initData.restore()
    console.log('[useTelegram] initData.restore() succeeded')
  } catch (error) {
    console.warn('[useTelegram] initData.restore() failed:', error)
  }
}

export function useTelegram() {
  console.log('[useTelegram] hook called')

  // Init data must be restored once before its signals become usable.
  useEffect(() => {
    console.log('[useTelegram] mount effect: calling restoreInitData()')
    restoreInitData()
    console.log('[useTelegram] mount effect: restoreInitData() returned')
  }, [])

  console.log('[useTelegram] before useSignal calls')
  const rawInitData = useSignal(initData.raw)
  console.log('[useTelegram] useSignal(initData.raw) =', rawInitData)
  const tgUser = useSignal(initData.user) as TelegramUser | undefined
  console.log('[useTelegram] useSignal(initData.user) =', tgUser)
  const startParam = useSignal(initData.startParam)
  console.log('[useTelegram] useSignal(initData.startParam) =', startParam)
  const isDark = useSignal(isMiniAppDark)
  console.log('[useTelegram] useSignal(isMiniAppDark) =', isDark)

  useEffect(() => {
    console.log('[useTelegram] theme effect: start')
    try {
      console.log('[useTelegram] calling expandViewport()')
      expandViewport()
      console.log('[useTelegram] calling bindMiniAppCssVars()')
      bindMiniAppCssVars()
      console.log('[useTelegram] calling bindThemeParamsCssVars()')
      bindThemeParamsCssVars()
      console.log('[useTelegram] calling bindViewportCssVars()')
      bindViewportCssVars()
      console.log('[useTelegram] calling setMiniAppHeaderColor()')
      setMiniAppHeaderColor('bg_color')
      console.log('[useTelegram] theme effect: done')
    } catch (error) {
      console.warn('[useTelegram] theme effect error:', error)
    }
  }, [])

  const result = useMemo(
    () => ({
      rawInitData: rawInitData || '',
      user: tgUser,
      startParam: startParam || '',
      appearance: (isDark ? 'dark' : 'light') as 'dark' | 'light',
    }),
    [rawInitData, tgUser, startParam, isDark],
  )

  console.log('[useTelegram] returning result:', result)
  return result
}
