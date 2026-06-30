import { useMemo } from 'react'
import { getDict, normalizeLang, type Lang, type Translations } from '../i18n'

export function useI18n(languageCode?: string): { lang: Lang; t: Translations } {
  const lang = useMemo(() => normalizeLang(languageCode), [languageCode])
  const t = useMemo(() => getDict(lang), [lang])
  return { lang, t }
}
