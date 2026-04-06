import type { ApiResponse } from '@/types/api'

export class HttpError extends Error {
  status: number
  details: string

  constructor(status: number, details: string) {
    super(details || `HTTP ${status}`)
    this.status = status
    this.details = details
  }
}

export async function requestJson<T>(input: string, init?: RequestInit): Promise<T> {
  const response = await fetch(input, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
    ...init,
  })

  const text = await response.text()
  let data: T & ApiResponse

  try {
    data = text ? (JSON.parse(text) as T & ApiResponse) : ({} as T & ApiResponse)
  } catch {
    const looksLikeHtml = text.trimStart().startsWith('<!doctype') || text.trimStart().startsWith('<html')
    if (looksLikeHtml) {
      throw new Error('当前前端已启动，但 `/api` 没有连到后端或代理未配置')
    }
    throw new Error('接口返回了无法解析的内容')
  }

  if (!response.ok) {
    const message = data.error || data.message || `请求失败 (${response.status})`
    throw new HttpError(response.status, message)
  }

  return data as T
}
