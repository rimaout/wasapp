const TOKEN_KEY = 'wasapp_token'
const USER_ID_KEY = 'wasapp_userId'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function getUserId() {
  return localStorage.getItem(USER_ID_KEY)
}

export function setAuth(token, userId) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_ID_KEY, userId)
}

export function clearAuth() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_ID_KEY)
}
