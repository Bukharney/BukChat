export async function apiFetch(input: RequestInfo | URL, init: RequestInit = {}): Promise<Response> {
  const token = localStorage.getItem("token");

  const headers = new Headers(init.headers || {});
  if (!headers.has("Content-Type") && !(init.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }
  if (token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(input, {
    ...init,
    headers,
  });

  if (response.status === 401) {
    localStorage.removeItem("token");
    if (window.location.pathname !== "/login") {
      window.location.href = "/login";
    }
  }

  return response;
}

export function getWsUrl(path: string): string {
  const envWsUrl = import.meta.env.VITE_WS_URL;
  const cleanPath = path.startsWith("/") ? path : `/${path}`;

  if (envWsUrl) {
    const baseUrl = envWsUrl.replace(/\/$/, "");
    return `${baseUrl}${cleanPath}`;
  }

  // Fallback: derive dynamically from browser location (handles HTTPS/wss: and Nginx reverse proxying)
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const host = window.location.host;
  return `${protocol}//${host}${cleanPath}`;
}
