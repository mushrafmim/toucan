import { getAccessTokenSnapshot } from '@/shared/auth/access-token'
import { getActiveRoleSnapshot, TOUCAN_ROLE_HEADER } from '@/shared/auth/roles'

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export async function getJSON<T>(input: string): Promise<T> {
  const response = await fetch(input, {
    headers: requestHeaders(),
  })

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`
    try {
      const payload = (await response.json()) as { error?: string }
      if (payload.error) {
        message = payload.error
      }
    } catch {
      // Keep the fallback message when the response is not JSON.
    }
    throw new ApiError(message, response.status)
  }

  return response.json() as Promise<T>
}

export async function postJSON<T>(input: string, body: unknown): Promise<T> {
  const response = await fetch(input, {
    method: 'POST',
    headers: requestHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify(body),
  })

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`
    try {
      const payload = (await response.json()) as { error?: string }
      if (payload.error) {
        message = payload.error
      }
    } catch {
      // Keep the fallback message when the response is not JSON.
    }
    throw new ApiError(message, response.status)
  }

  if (response.status === 204) {
    return {} as T
  }

  return response.json() as Promise<T>
}
export async function deleteItem(input: string): Promise<void> {
  const response = await fetch(input, {
    method: 'DELETE',
    headers: requestHeaders(),
  })

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`
    try {
      const payload = (await response.json()) as { error?: string }
      if (payload.error) {
        message = payload.error
      }
    } catch {
      // Keep the fallback message when the response is not JSON.
    }
    throw new ApiError(message, response.status)
  }
}

function requestHeaders(extra?: HeadersInit): HeadersInit {


  const token = getAccessTokenSnapshot()

  return {
    Accept: 'application/json',
    [TOUCAN_ROLE_HEADER]: getActiveRoleSnapshot(),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...extra,
  }
}
