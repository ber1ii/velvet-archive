const BASE_URL = '/api/v1'; // Routed via our Vite dev proxy

interface RequestOptions extends Omit<RequestInit, 'body'> {
  body?: unknown;
}

async function handleResponse(response: Response) {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(
      errorData.error || `HTTP error! status: ${response.status}`,
    );
  }
  if (response.status === 204) return null;
  return response.json();
}

export const api = {
  get: async (endpoint: string, options: RequestOptions = {}) => {
    // Destructure body out so it is completely excluded from fetchOptions
    const { body: _, ...fetchOptions } = options;
    const token = localStorage.getItem("admin_key");
    const headers = new Headers(fetchOptions.headers);
    if (token) headers.set("Authorization", `Bearer ${token}`);

    const res = await fetch(`${BASE_URL}${endpoint}`, {
      ...fetchOptions, // Safe now! It completely lacks any 'body' field
      method: "GET",
      headers,
    });
    return handleResponse(res);
  },

  post: async (endpoint: string, options: RequestOptions = {}) => {
    const { body, ...fetchOptions } = options;
    const token = localStorage.getItem("admin_key");
    const headers = new Headers(fetchOptions.headers);
    headers.set("Content-Type", "application/json");
    if (token) headers.set("Authorization", `Bearer ${token}`);

    const res = await fetch(`${BASE_URL}${endpoint}`, {
      ...fetchOptions,
      method: "POST",
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });
    return handleResponse(res);
  },

  put: async (endpoint: string, options: RequestOptions = {}) => {
    const { body, ...fetchOptions } = options;
    const token = localStorage.getItem("admin_key");
    const headers = new Headers(fetchOptions.headers);
    headers.set("Content-Type", "application/json");
    if (token) headers.set("Authorization", `Bearer ${token}`);

    const res = await fetch(`${BASE_URL}${endpoint}`, {
      ...fetchOptions,
      method: "PUT",
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });
    return handleResponse(res);
  },
};