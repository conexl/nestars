import { useEffect, useMemo } from 'react'
import {
  bindMiniAppCssVars,
  bindThemeParamsCssVars,
  bindViewportCssVars,
  expandViewport,
  isMiniAppDark,
  setMiniAppHeaderColor,
  useLaunchParams,
  useRawInitData,
  useSignal,
} from '@telegram-apps/sdk-react'
import type { TelegramUser } from '../types'

function useSafeLaunchParams() {
  try {
    const params = useLaunchParams(false)
    console.log('LaunchParams:', params)
    return params
  } catch (error) {
    console.warn('Error getting launch params:', error)
    return undefined
  }
}

function useSafeRawInitData() {
  try {
    const data = useRawInitData()
    console.log('RawInitData from hook:', data)
    return data
  } catch (error) {
    console.warn('Error getting raw init data:', error)
    return undefined
  }
}

export function useTelegram() {
  const lp = useSafeLaunchParams()
  const rawInitData = useSafeRawInitData()
  const isDark = useSignal(isMiniAppDark)

  const user = lp?.tgWebAppData?.user as TelegramUser | undefined

  useEffect(() => {
    console.log('Telegram hook mounted')
    console.log('LP data:', lp)
    console.log('Raw init data:', rawInitData)
    console.log('User:', user)
    
    try {
      expandViewport()
      bindMiniAppCssVars()
      bindThemeParamsCssVars()
      bindViewportCssVars()
      setMiniAppHeaderColor('bg_color')
    } catch (error) {
      console.warn('Error in Telegram init:', error)
    }
  }, [lp, rawInitData, user])

  return useMemo(
    () => ({
      rawInitData: rawInitData || '',
      user,
      startParam: lp?.tgWebAppStartParam || '',
      appearance: isDark ? 'dark' : 'light' as 'dark' | 'light',
    }),
    [rawInitData, user, lp?.tgWebAppStartParam, isDark],
  )
}
