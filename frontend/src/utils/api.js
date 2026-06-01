const API = 'http://localhost:8000'

export function getToken() {
  return localStorage.getItem('token')
}

export async function apiFetch(endpoint, options = {}) {
  const token = getToken()
  const res = await fetch(`${API}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    }
  })

  if (res.status === 401) {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    window.location.href = '/login'
  }

  return res
}