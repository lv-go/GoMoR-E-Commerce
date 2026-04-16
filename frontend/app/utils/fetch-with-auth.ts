import { auth } from "../firebase-config";

const BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

export async function fetchWithAuth<T>(
  endpoint: string,
  { params = {}, ...options }: { params?: Record<string, any> } & RequestInit = {},
): Promise<T> {
  const user = auth.currentUser;
  const token = user ? await user.getIdToken() : null;

  const queryParams = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== "" && !(Array.isArray(value) && value.length === 0)) {
      if (Array.isArray(value)) {
        queryParams.append(key, value.join(","));
      } else {
        queryParams.append(key, String(value));
      }
    }
  }

  const headers = new Headers();
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  headers.set("Content-Type", "application/json");

  const response = await fetch(`${BASE_URL}${endpoint}?${queryParams}`, {
    headers,
    ...options,
  });

  if (!response.ok) {
    const errorData = await response.text();
    throw new Error(errorData || `HTTP error! Status: ${response.status}`);
  }

  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}
