const BACKEND_API_BASE =
  process.env.BACKEND_API_URL || "http://localhost:8081/api";

export class ApiError extends Error {
  status: number;
  data?: any;

  constructor(message: string, status: number, data?: any) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.data = data;
  }
}

/**
 * Server-side direct backend fetch (for Server Components and Route Handlers)
 */
export async function backendFetch<T = any>(
  endpoint: string,
  options: RequestInit & { token?: string } = {}
): Promise<T> {
  const { token, headers = {}, ...fetchOptions } = options;

  const normalizedEndpoint = endpoint.startsWith("/") ? endpoint : `/${endpoint}`;
  const url = `${BACKEND_API_BASE}${normalizedEndpoint}`;

  const requestHeaders: Record<string, string> = {
    "Content-Type": "application/json",
    ...(headers as Record<string, string>),
  };

  if (token) {
    requestHeaders["Authorization"] = `Bearer ${token}`;
  }

  try {
    const res = await fetch(url, {
      ...fetchOptions,
      headers: requestHeaders,
      cache: "no-store",
    });

    if (!res.ok) {
      let errorMessage = `API request failed with status ${res.status}`;
      try {
        const errorData = await res.json();
        errorMessage = errorData.error || errorData.message || errorMessage;
      } catch {
        const textData = await res.text();
        if (textData) errorMessage = textData;
      }
      throw new ApiError(errorMessage, res.status);
    }

    if (res.status === 204) {
      return {} as T;
    }

    return await res.json();
  } catch (error) {
    if (error instanceof ApiError) throw error;
    throw new ApiError((error as Error).message || "Network error", 500);
  }
}

/**
 * Client-side proxy fetch (calls /api/proxy which automatically forwards httpOnly cookies)
 */
export async function apiFetch<T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const normalizedEndpoint = endpoint.startsWith("/") ? endpoint : `/${endpoint}`;
  const url = `/api/proxy${normalizedEndpoint}`;

  const requestHeaders: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string> || {}),
  };

  try {
    const res = await fetch(url, {
      ...options,
      headers: requestHeaders,
    });

    if (!res.ok) {
      let errorMessage = `API request failed with status ${res.status}`;
      try {
        const errorData = await res.json();
        errorMessage = errorData.error || errorData.message || errorMessage;
      } catch {
        const textData = await res.text();
        if (textData) errorMessage = textData;
      }
      throw new ApiError(errorMessage, res.status);
    }

    if (res.status === 204) {
      return {} as T;
    }

    return await res.json();
  } catch (error) {
    if (error instanceof ApiError) throw error;
    throw new ApiError((error as Error).message || "Network error", 500);
  }
}
