import type { AuthResponse, CatalogItem, CreateOrderInput, Order, PaymentMethod, Profile, User, WithdrawInput } from '../types'

const API_URL = import.meta.env.VITE_API_URL || ''

let initDataRaw = ''

export function setInitData(raw: string) {
  initDataRaw = raw
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_URL}/api/miniapp${path}`
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(initDataRaw ? { Authorization: `tma ${initDataRaw}` } : {}),
  }

  if (options.headers) {
    Object.assign(headers, options.headers)
  }

  console.log('[api/client] request:', { url, method: options.method || 'GET', headers })

  const response = await fetch(url, {
    ...options,
    headers,
  })

  console.log('[api/client] response:', { url, status: response.status, statusText: response.statusText, contentType: response.headers.get('content-type') })

  if (!response.ok) {
    const text = await response.text()
    console.error('[api/client] non-ok response body:', text.slice(0, 1000))
    throw new Error(text || `HTTP ${response.status}`)
  }

  if (response.status === 204) {
    return undefined as T
  }

  const responseText = await response.text()
  console.log('[api/client] response body:', responseText.slice(0, 1000))

  try {
    return JSON.parse(responseText) as T
  } catch (parseError) {
    console.error('[api/client] JSON.parse failed:', parseError)
    console.error('[api/client] raw response was:', responseText.slice(0, 1000))
    throw new Error(`Invalid JSON response from ${url}: ${responseText.slice(0, 200)}`)
  }
}

export async function auth(rawInitData: string): Promise<AuthResponse> {
  console.log('[api/client] auth called with rawInitData length:', rawInitData.length)
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
