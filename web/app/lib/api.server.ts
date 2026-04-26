import axios, { type AxiosRequestConfig } from "axios"
import { ok, err, validationErr, type Result } from "~/lib/result"

const API_URL = process.env.API_URL ?? "http://localhost:3000"

type SearchParams = Record<string, string | undefined>

function buildUrl(path: string, params?: SearchParams): string {
  const url = new URL(path, API_URL)
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (v !== undefined) url.searchParams.set(k, v)
    }
  }
  return url.toString()
}

function authHeader(token?: string): Record<string, string> {
  return token ? { Authorization: `Bearer ${token}` } : {}
}

function normalizeError(error: unknown): Result<never> {
  if (!axios.isAxiosError(error) || !error.response) return err("UNKNOWN")

  const data = error.response.data as Record<string, unknown> | undefined
  const code = typeof data?.code === "string" ? data.code : null

  if (code === "VALIDATION_ERROR" && Array.isArray(data?.errors)) {
    return validationErr(data.errors as { field: string; tag: string }[])
  }

  return err(code ?? "UNKNOWN")
}

async function request<T>(config: AxiosRequestConfig): Promise<Result<T>> {
  try {
    const res = await axios<T>(config)
    return ok(res.data)
  } catch (e) {
    return normalizeError(e)
  }
}

export function apiClient(tenantId?: string) {
  const baseHeaders = {
    "Content-Type": "application/json",
    ...(tenantId ? { "X-Tenant-ID": tenantId } : {}),
  }

  return {
    get<T>(path: string, params?: SearchParams, token?: string): Promise<Result<T>> {
      return request<T>({ method: "GET", url: buildUrl(path, params), headers: { ...baseHeaders, ...authHeader(token) } })
    },

    post<T>(path: string, body: unknown, token?: string): Promise<Result<T>> {
      return request<T>({ method: "POST", url: buildUrl(path), data: body, headers: { ...baseHeaders, ...authHeader(token) } })
    },

    patch<T>(path: string, body: unknown, token?: string): Promise<Result<T>> {
      return request<T>({ method: "PATCH", url: buildUrl(path), data: body, headers: { ...baseHeaders, ...authHeader(token) } })
    },
  }
}
