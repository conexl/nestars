import type { AuthResponse, CatalogItem, CreateOrderInput, Order, PaymentMethod, Profile, User, WithdrawInput } from '../types'

const API_URL = import.meta.env.VITE_API_URL || ''

let initDataRaw = ''

export function setInitData(raw: string) {
  initDataRaw = raw
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(initDataRaw ? { Authorization: `tma ${initDataRaw}` } : {}),
  }

  if (options.headers) {
    Object.assign(headers, options.headers)
  }

  const response = await fetch(`${API_URL}/api/miniapp${path}`, {
    ...options,
    headers,
  })

  if (!response.ok) {
    const text = await response.text()
    throw new Error(text || `HTTP ${response.status}`)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return response.json() as Promise<T>
}

export async function auth(rawInitData: string): Promise<AuthResponse> {
  setInitData(rawInitData)
  return request('/auth', {
    method: 'POST',
    body: JSON.stringify({ init_data: rawInitData }),
  })
}

export async function getCatalogStars(): Promise<CatalogItem[]> {
  return request('/catalog/stars')
}

export async function getCatalogPremium(): Promise<CatalogItem[]> {
  return request('/catalog/premium')
}

export async function createOrder(input: CreateOrderInput): Promise<Order> {
  return request('/orders', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export async function payOrder(orderId: string, method: PaymentMethod): Promise<{ payment_url?: string }> {
  return request(`/orders/${orderId}/pay`, {
    method: 'POST',
    body: JSON.stringify({ method }),
  })
}

export async function getProfile(): Promise<{ user: User; profile: Profile }> {
  return request('/profile')
}

export async function getOrders(): Promise<Order[]> {
  return request('/orders')
}

export async function withdraw(input: WithdrawInput): Promise<void> {
  return request('/withdraw', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}
