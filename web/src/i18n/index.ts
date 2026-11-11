import { createI18n } from 'vue-i18n'
import en from './en.json'
import id from './id.json'

export type Locale = 'en' | 'id'

export function makeI18n(locale: Locale = 'en') {
  return createI18n({ legacy: false, locale, fallbackLocale: 'en', messages: { en, id } })
}

export function savedLocale(): Locale {
  const v = typeof localStorage !== 'undefined' ? localStorage.getItem('locale') : null
  return v === 'id' ? 'id' : 'en'
}
