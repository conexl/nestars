export interface TelegramUser {
  id: number
  first_name: string
  last_name?: string
  username?: string
  language_code?: string
  photo_url?: string
}

export interface User {
  id: string
  telegram_id: number
  username?: string
  first_name?: string
  last_name?: string
  language_code?: string
  referral_code: string
  created_at: string
}

export interface Profile {
  id: string
  user_id: string
  telegram_id: number
  referral_code: string
  balance_usd: number
  balance_ton: number
  bought_stars: number
  bonus_stars: number
  referral_bonus_stars: number
  referral_count: number
}

export interface CatalogItem {
  id: string
  name: string
  value: number
  amount: number
  currency: string
}

export interface Order {
  id: string
  order_number: number
  user_id: string
  recipient_username: string
  product_type: 'Stars' | 'Premium'
  product_id: string
  stars_amount?: number
  amount: number
  currency: string
  status: 'pending' | 'paid' | 'cancelled' | 'failed'
  payment_method?: string
  payment_url?: string
  created_at: string
}

export interface CreateOrderInput {
  username: string
  product_type: 'Stars' | 'Premium'
  product_id: string
  stars?: number
  currency: string
}

export interface AuthResponse {
  user: User
  profile: Profile
}

export type PaymentMethod = 'Balance' | 'CryptoBot' | 'Yookassa' | 'Heleket'

export interface WithdrawInput {
  username: string
  stars: number
}
